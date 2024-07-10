package models

import (
	"math/rand"
	"time"

	"github.com/ShukinDmitriy/shortener/internal/environments"
	"github.com/labstack/echo/v4"
)

// Event short link generation event structure
type Event struct {
	ShortKey      string `json:"short_key,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
	UserID        string `json:"user_id,omitempty"`
	DeletedFlag   bool   `json:"is_deleted,omitempty"`
}

// GenerateShortKey generate random string
func GenerateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rng.Intn(len(charset))]
	}
	return string(shortKey)
}

// PrepareFullURL prepare full link
func PrepareFullURL(ctx echo.Context, shortKey string) string {
	var host string

	if environments.FlagBaseAddr != "" {
		host = environments.FlagBaseAddr
	} else {
		host = "http://" + ctx.Request().Host
	}

	return host + "/" + shortKey
}
