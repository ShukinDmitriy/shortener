package models

import (
	"context"

	"github.com/ShukinDmitriy/shortener/internal/environments"
)

// URLRepository repository interface for working with URL
type URLRepository interface {
	Initialize(configuration environments.Configuration) error

	Get(shortKey string) (Event, bool)

	Save(ctx context.Context, events []*Event) error

	Delete(ctx context.Context, events []DeleteRequestBatch) error

	GetShortKeyByOriginalURL(originalURL string) (string, bool)

	GetEventsByUserID(ctx context.Context, userID string) []*Event

	GetStats(ctx context.Context) (countUser int, countURL int, err error)
}
