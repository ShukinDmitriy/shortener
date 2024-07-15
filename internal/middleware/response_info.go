package middleware

import (
	"github.com/ShukinDmitriy/shortener/internal/logger"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// ResponseInfo middleware for logging response
func ResponseInfo(applicationLogger logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			next(c)

			applicationLogger.Info("HTTP response",
				zap.String("status", strconv.FormatInt(int64(c.Response().Status), 10)),
				zap.String("size", strconv.FormatInt(c.Response().Size, 10)),
			)

			return nil
		}
	}
}
