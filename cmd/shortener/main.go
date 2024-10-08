// The package starts the application Shortener
package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ShukinDmitriy/shortener/internal/app"
	"github.com/ShukinDmitriy/shortener/internal/auth"
	"github.com/ShukinDmitriy/shortener/internal/environments"
	"github.com/ShukinDmitriy/shortener/internal/logger"
	internalMiddleware "github.com/ShukinDmitriy/shortener/internal/middleware"
	"github.com/ShukinDmitriy/shortener/internal/models"
	pb "github.com/ShukinDmitriy/shortener/proto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var (
	// buildVersion build version
	buildVersion = "N/A"
	// buildDate build date
	buildDate = "N/A"
	// buildCommit build commit hash
	buildCommit = "N/A"
)

func urlRepositoryFactory(configuration environments.Configuration) (models.URLRepository, error) {
	var repository models.URLRepository

	if configuration.DatabaseDSN != "" {
		repository = &models.PGURLRepository{}
	} else {
		repository = &models.MemoryURLRepository{}
	}

	err := repository.Initialize(configuration)
	if err != nil {
		return nil, err
	}

	return repository, nil
}

func main() {
	log.Printf("Build version: %v\n", buildVersion)
	log.Printf("Build date: %v\n", buildDate)
	log.Printf("Build commit: %v\n", buildCommit)

	// Профилирование
	runProf()

	configuration := environments.ParseFlags()

	if err := logger.Initialize(configuration.LogLevel); err != nil {
		fmt.Println(err)
		return
	}

	repository, err := urlRepositoryFactory(configuration)
	if err != nil {
		fmt.Println(err)
		return
	}

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	conn, err := pgx.Connect(context.Background(), configuration.DatabaseDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	} else {
		defer conn.Close(context.Background())
	}

	authService := auth.NewAuthService()
	_, subnet, err := net.ParseCIDR(configuration.TrustedSubnet)
	if err != nil {
		fmt.Println(err)
	}
	shortener := app.NewURLShortener(repository, conn, authService, subnet)

	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	// Routing
	e.GET("/:id", shortener.HandleRedirect)
	e.POST("/", shortener.HandleShorten)
	e.POST("/api/shorten", shortener.HandleCreateShorten)
	e.POST("/api/shorten/batch", shortener.HandleCreateShortenBatch)
	e.GET("/ping", shortener.HandlePing)
	e.GET("/api/user/urls", shortener.HandleUserURLGet)
	e.DELETE("/api/user/urls", shortener.HandleUserURLDelete)
	e.GET("/api/internal/stats", shortener.HandleGetStats)

	//-------------------
	// middleware
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

	// auth
	e.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &auth.Claims{}
		},
		Skipper: func(c echo.Context) bool {
			return !strings.Contains(c.Path(), "/api/user/")
		},
		SigningKey:    []byte(auth.GetJWTSecret()),
		SigningMethod: auth.GetSigningMethod().Alg(),
		TokenLookup:   "cookie:access-token", // "<source>:<name>"
		ErrorHandler:  auth.JWTErrorChecker,
	}))

	e.Use(auth.CreateTokenWithConfig(auth.CreateTokenConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "/api/user/")
		},
	}))
	e.Use(auth.TokenRefreshMiddleware)

	// Start gRPC
	var grpcServer *grpc.Server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", configuration.GrpcPort))
		if err != nil {
			log.Printf("gRPC server failed to listen: %v", err.Error())
			return err
		}
		grpcServer = grpc.NewServer()
		shortenerGRPC := app.NewURLShortenerGRPC(repository, conn, subnet)
		pb.RegisterURLServer(grpcServer, shortenerGRPC)
		log.Printf("grpc server listening at %v", listener.Addr())
		return grpcServer.Serve(listener)
	})
	// Start server
	go func() {
		if configuration.EnableHTTPS {
			if err := e.StartTLS(configuration.RunAddr, "ssl/localhost.crt", "ssl/device.key"); err != nil && !errors.Is(err, http.ErrServerClosed) {
				e.Logger.Fatal("shutting down the server")
			}
		} else {
			if err := e.Start(configuration.RunAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
				e.Logger.Fatal("shutting down the server")
			}
		}

		zap.L().Info("Running server", zap.String("address", configuration.RunAddr))
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	<-ctx.Done()

	// Запускаем остановку
	shutdownChan := shortener.Shutdown()
	<-shutdownChan

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	if grpcServer != nil {
		grpcServer.GracefulStop()
	}
}
