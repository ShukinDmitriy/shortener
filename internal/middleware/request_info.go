package middleware

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"sync"
	"time"
)

type (
	RequestInfo struct {
		requestStart time.Time
		mutex        sync.RWMutex
		logger       zap.Logger
	}
)

func NewRequestInfo(l *zap.Logger) *RequestInfo {
	return &RequestInfo{
		logger: *l,
	}
}

func (ri *RequestInfo) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ri.mutex.Lock()
		defer ri.mutex.Unlock()

		ri.requestStart = time.Now()

		if err := next(c); err != nil {
			ri.logger.Error("error process HTTP request", zap.String("err", err.Error()))

			c.Error(err)
		}

		duration := time.Since(ri.requestStart)
		ri.logger.Info("got incoming HTTP request",
			zap.String("method", c.Request().Method),
			zap.String("path", c.Request().URL.Path),
			zap.String("duration", duration.String()),
		)

		return nil
	}
}
