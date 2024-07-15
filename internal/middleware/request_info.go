// Package middleware custom middlewares
package middleware

import (
	"github.com/ShukinDmitriy/shortener/internal/logger"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// RequestInfo middleware for logging request
func RequestInfo(applicationLogger logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestStart := time.Now()

			if err := next(c); err != nil {
				applicationLogger.Error("error process HTTP request", zap.String("err", err.Error()))

				c.Error(err)
			}

			duration := time.Since(requestStart)
			applicationLogger.Info("got incoming HTTP request",
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.String("duration", duration.String()),
			)

			return nil
		}
	}
}
