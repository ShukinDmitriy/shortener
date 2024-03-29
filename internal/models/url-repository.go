package models

import (
	"context"
)

type URLRepository interface {
	Initialize() error

	Get(shortKey string) (Event, bool)

	Save(ctx context.Context, events []*Event) error

	Delete(ctx context.Context, events []DeleteRequestBatch) error

	GetShortKeyByOriginalURL(originalURL string) (string, bool)

	GetEventsByUserID(ctx context.Context, userID string) []*Event
}
