package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ShukinDmitriy/shortener/internal/logger"
	internalMiddleware "github.com/ShukinDmitriy/shortener/internal/middleware"
	"github.com/ShukinDmitriy/shortener/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

type URLShortener struct {
	urls map[string]string
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rng.Intn(len(charset))]
	}
	return string(shortKey)
}

func saveShortKey(us *URLShortener, shortKey string, originalURL string) {
	// Хранение в памяти
	us.urls[shortKey] = originalURL

	if models.DBProducer == nil {
		return
	}

	// Хранение в файле
	models.DBProducer.WriteEvent(&models.Event{
		ShortKey:    shortKey,
		OriginalURL: originalURL,
	})
}

func getOriginalURL(us *URLShortener, shortKey string) (string, bool) {
	// Поиск в памяти
	var originalURL string
	var event *models.Event
	var found = false
	var err error

	originalURL, found = us.urls[shortKey]

	if models.DBConsumer == nil {
		return originalURL, found
	}

	// Поиск в файле
	defer models.DBConsumer.Close()

	for {
		event, err = models.DBConsumer.ReadEvent()

		if err != nil {
			return originalURL, found
		}

		// Сохраняем значение в память, т.к. повторно файл не вычитывается
		us.urls[event.ShortKey] = event.OriginalURL

		if event.ShortKey == shortKey {
			return event.OriginalURL, true
		}
	}
}

func prepareFullURL(shortKey string, ctx echo.Context) string {
	var host string

	if flagBaseAddr != "" {
		host = flagBaseAddr
	} else {
		host = "http://" + ctx.Request().Host
	}

	return host + "/" + shortKey
}

func (us *URLShortener) HandleShorten(ctx echo.Context) error {
	originalURL, err := io.ReadAll(ctx.Request().Body)

	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "can't read body. internal error")
	}

	if string(originalURL) == "" {
		err := "empty url"
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Generate a unique shortened key for the original URL
	shortKey := generateShortKey()
	saveShortKey(us, shortKey, string(originalURL))
	result := prepareFullURL(shortKey, ctx)

	ctx.Response().Header().Set("Content-Type", "text/plain; charset=utf-8")

	return ctx.String(http.StatusCreated, result)
}

func (us *URLShortener) HandleCreateShorten(ctx echo.Context) error {
	// десериализуем запрос в структуру модели
	zap.L().Debug("decoding request")
	var req models.CreateRequest
	dec := json.NewDecoder(ctx.Request().Body)
	if err := dec.Decode(&req); err != nil {
		zap.L().Debug("cannot decode request JSON body", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid JSON")
	}

	// проверяем, что пришёл запрос понятного типа
	if string(req.URL) == "" {
		err := "empty url"
		ctx.Logger().Error(err)
		zap.L().Debug("unsupported request url", zap.String("url", req.URL))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Generate a unique shortened key for the original URL
	shortKey := generateShortKey()
	saveShortKey(us, shortKey, req.URL)

	// заполняем модель ответа
	resp := models.CreateResponse{
		Result: prepareFullURL(shortKey, ctx),
	}

	return ctx.JSON(http.StatusCreated, resp)
}

func (us *URLShortener) HandleRedirect(ctx echo.Context) error {
	shortKey := ctx.Param("id")

	if shortKey == "" {
		ctx.Logger().Error("empty id")
		return echo.NewHTTPError(http.StatusBadRequest, "")
	}

	// Retrieve the original URL from the `urls` map using the shortened key
	originalURL, found := getOriginalURL(us, shortKey)
	if !found {
		err := "URL not found"
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return ctx.Redirect(http.StatusTemporaryRedirect, originalURL)
}

var shortener = &URLShortener{
	urls: make(map[string]string),
}

func main() {
	parseFlags()

	if err := logger.Initialize(flagLogLevel); err != nil {
		return
	}

	if err := models.Initialize(flagFileStoragePath); err != nil {
		fmt.Println(err)
		return
	}

	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	//-------------------
	// Custom middleware
	//-------------------
	// ResponseInfo
	resInfo := internalMiddleware.NewResponseInfo(zap.L())
	e.Use(resInfo.Process)

	// RequestInfo
	reqInfo := internalMiddleware.NewRequestInfo(zap.L())
	e.Use(reqInfo.Process)

	// gzip Отдавать сжатый ответ клиенту, который поддерживает обработку
	// сжатых ответов (с HTTP-заголовком Accept-Encoding)
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			skipByAcceptEncodingHeader := true
			skipByContentTypeHeader := true

			acceptEncodingRaw := c.Request().Header.Get("Accept-Encoding")
			acceptEncodingValues := strings.Split(acceptEncodingRaw, ",")

			for _, value := range acceptEncodingValues {
				parts := strings.Split(value, ";")
				format := strings.TrimSpace(parts[0])

				if format == "gzip" {
					skipByAcceptEncodingHeader = false
					break
				}
			}

			contentTypeRaw := c.Request().Header.Get("Content-Type")
			contentTypeValues := strings.Split(contentTypeRaw, ",")

			for _, value := range contentTypeValues {
				if value == "application/json" || value == "text/html" {
					skipByContentTypeHeader = false
					break
				}
			}

			return skipByAcceptEncodingHeader && skipByContentTypeHeader
		},
	}))

	// decompress
	e.Use(middleware.DecompressWithConfig(middleware.DecompressConfig{}))

	e.GET("/:id", shortener.HandleRedirect)
	e.POST("/", shortener.HandleShorten)
	e.POST("/api/shorten", shortener.HandleCreateShorten)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(flagRunAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server")
		}

		zap.L().Info("Running server", zap.String("address", flagRunAddr))
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
