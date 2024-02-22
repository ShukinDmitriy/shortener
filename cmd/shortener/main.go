package main

import (
	"context"
	"errors"
	"github.com/ShukinDmitriy/shortener/internal/logger"
	"github.com/ShukinDmitriy/shortener/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
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
	us.urls[shortKey] = string(originalURL)

	ctx.Response().Header().Set("Content-Type", "text/plain; charset=utf-8")

	var host string

	if flagBaseAddr != "" {
		host = flagBaseAddr
	} else {
		host = "http://" + ctx.Request().Host
	}

	return ctx.String(http.StatusCreated, host+"/"+shortKey)
}

func (us *URLShortener) HandleRedirect(ctx echo.Context) error {
	shortKey := ctx.Param("id")

	if shortKey == "" {
		ctx.Logger().Error("empty id")
		return echo.NewHTTPError(http.StatusBadRequest, "")
	}

	// Retrieve the original URL from the `urls` map using the shortened key
	originalURL, found := us.urls[shortKey]
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

	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	//-------------------
	// Custom middleware
	//-------------------
	// ResponseInfo
	resInfo := middleware.NewResponseInfo(logger.Log)
	e.Use(resInfo.Process)

	// RequestInfo
	reqInfo := middleware.NewRequestInfo(logger.Log)
	e.Use(reqInfo.Process)

	e.GET("/:id", shortener.HandleRedirect)
	e.POST("/", shortener.HandleShorten)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(flagRunAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server")
		}

		logger.Log.Info("Running server", zap.String("address", flagRunAddr))
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
