package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ShukinDmitriy/shortener/internal/environments"
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

func prepareFullURL(ctx echo.Context, shortKey string) string {
	var host string

	if environments.FlagBaseAddr != "" {
		host = environments.FlagBaseAddr
	} else {
		host = "http://" + ctx.Request().Host
	}

	return host + "/" + shortKey
}

func urlRepositoryFactory() (models.URLRepository, error) {
	var repository models.URLRepository

	if environments.FlagDatabaseDSN != "" {
		repository = &models.PGURLRepository{}
	} else {
		repository = &models.MemoryURLRepository{}
	}

	err := repository.Initialize()
	if err != nil {
		return nil, err
	}

	return repository, nil
}

func main() {
	environments.ParseFlags()

	if err := logger.Initialize(environments.FlagLogLevel); err != nil {
		return
	}

	repository, err := urlRepositoryFactory()

	if err != nil {
		fmt.Println(err)
		return
	}

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	conn, err := pgx.Connect(context.Background(), environments.FlagDatabaseDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	} else {
		defer conn.Close(context.Background())
	}

	var shortener = newURLShortener(repository, conn)

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
	e.POST("/api/shorten/batch", shortener.HandleCreateShortenBatch)
	e.GET("/ping", shortener.HandlePing)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(environments.FlagRunAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server")
		}

		zap.L().Info("Running server", zap.String("address", environments.FlagRunAddr))
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
