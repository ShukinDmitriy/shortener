package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ShukinDmitriy/shortener/internal/logger"
	internalMiddleware "github.com/ShukinDmitriy/shortener/internal/middleware"
	"github.com/ShukinDmitriy/shortener/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

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

func initMapFromDB(us *URLShortener) {
	var event *models.Event
	var err error

	if models.DBConsumer == nil {
		return
	}

	defer models.DBConsumer.Close()

	for {
		event, err = models.DBConsumer.ReadEvent()

		if event == nil || err != nil {
			return
		}

		// Сохраняем значение в память, т.к. повторно файл не вычитывается
		us.urls[event.ShortKey] = event.OriginalURL
	}

}

func getOriginalURL(us *URLShortener, shortKey string) (string, bool) {
	// Поиск в памяти
	var originalURL string
	var found = false

	originalURL, found = us.urls[shortKey]

	return originalURL, found
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

func main() {
	parseFlags()

	if err := logger.Initialize(flagLogLevel); err != nil {
		return
	}

	if err := models.Initialize(flagFileStoragePath); err != nil {
		fmt.Println(err)
		return
	}

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	conn, err := pgx.Connect(context.Background(), flagDatabaseDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var shortener = newURLShortener(make(map[string]string), conn)

	initMapFromDB(shortener)

	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	//-------------------
	// Custom middleware
	//-------------------
	// ResponseInfo
	e.Use(internalMiddleware.ResponseInfo(zap.L()))

	// RequestInfo
	e.Use(internalMiddleware.RequestInfo(zap.L()))

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
	e.GET("/ping", shortener.HandlePing)

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
