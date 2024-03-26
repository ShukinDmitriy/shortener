package models

import (
	"context"
)

type URLRepository interface {
	Initialize() error

	Get(shortKey string) (string, bool)

	Save(ctx context.Context, events []*Event) error

	GetShortKeyByOriginalURL(originalURL string) (string, bool)

	GetEventsByUserID(ctx context.Context, userID string) []*Event
}
