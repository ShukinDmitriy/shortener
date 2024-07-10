package middleware

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func ResponseInfo(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			next(c)

			logger.Info("HTTP response",
				zap.String("status", strconv.FormatInt(int64(c.Response().Status), 10)),
				zap.String("size", strconv.FormatInt(c.Response().Size, 10)),
			)

			return nil
		}
	}
}
