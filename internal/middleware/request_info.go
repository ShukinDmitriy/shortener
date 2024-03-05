package middleware

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"time"
)

func RequestInfo(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestStart := time.Now()

			if err := next(c); err != nil {
				logger.Error("error process HTTP request", zap.String("err", err.Error()))

				c.Error(err)
			}

			duration := time.Since(requestStart)
			logger.Info("got incoming HTTP request",
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.String("duration", duration.String()),
			)

			return nil
		}
	}
}
