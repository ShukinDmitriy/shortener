package logger_test

import (
	"github.com/ShukinDmitriy/shortener/internal/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestInitializeSuccess проверяет успешную инициализацию логгера.
func TestInitializeSuccess(t *testing.T) {
	assert.NoError(t, logger.Initialize("info"))
}
