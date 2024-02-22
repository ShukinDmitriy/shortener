package middleware

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"strconv"
	"sync"
)

type (
	ResponseInfo struct {
		mutex  sync.RWMutex
		logger zap.Logger
	}
)

func NewResponseInfo(l *zap.Logger) *ResponseInfo {
	return &ResponseInfo{
		logger: *l,
	}
}

func (ri *ResponseInfo) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		next(c)

		ri.logger.Info("HTTP response",
			zap.String("status", strconv.FormatInt(int64(c.Response().Status), 10)),
			zap.String("size", strconv.FormatInt(c.Response().Size, 10)),
		)

		return nil
	}
}
