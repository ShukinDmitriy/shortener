package logger_test

import (
	"testing"

	"github.com/ShukinDmitriy/shortener/internal/logger"
	"github.com/stretchr/testify/assert"
)

// TestInitializeSuccess проверяет успешную инициализацию логгера.
func TestInitializeSuccess(t *testing.T) {
	assert.NoError(t, logger.Initialize("info"))
}
