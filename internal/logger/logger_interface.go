package logger

import "go.uber.org/zap"

// Logger interface for logger
type Logger interface {
	// Info level log
	Info(msg string, fields ...zap.Field)
	// Warn level log
	Warn(msg string, fields ...zap.Field)
	// Error level log
	Error(msg string, fields ...zap.Field)
	// Fatal level log
	Fatal(msg string, fields ...zap.Field)
}
