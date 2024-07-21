package models_test

import (
	"context"
	"github.com/ShukinDmitriy/shortener/internal/environments"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/ShukinDmitriy/shortener/internal/models"
)

func BenchmarkMemoryURLRepository_Initialize(b *testing.B) {
	repository := &models.MemoryURLRepository{}

	b.Run("initialize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = repository.Initialize()
		}
	})
}

func TestMemoryURLRepository_Initialize(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test #1",
			args: args{
				filename: "./testdata/events.json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.filename != "" {
				environments.FlagFileStoragePath = tt.args.filename
			}

			repository := &models.MemoryURLRepository{}
			assert.NoError(t, repository.Initialize())
			if tt.args.filename != "" {
				assert.FileExists(t, tt.args.filename)
				_ = os.Remove(tt.args.filename)
			}
		})
	}
}

func TestMemoryURLRepository_CRUD(t *testing.T) {
	type args struct {
		filename string
		events   []models.Event
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test #1",
			args: args{
				filename: "./events.json",
				events: []models.Event{
					{
						OriginalURL: "https://example.com",
						ShortKey:    "short1",
						UserID:      "1",
					},
					{
						OriginalURL: "https://example.com",
						ShortKey:    "short1",
						UserID:      "2",
					},
					{
						OriginalURL: "https://example1.com",
						ShortKey:    "short2",
						UserID:      "3",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.filename != "" {
				environments.FlagFileStoragePath = tt.args.filename
				defer os.Remove(tt.args.filename)
			}

			repository := &models.MemoryURLRepository{}
			assert.NoError(t, repository.Initialize())

			for _, event := range tt.args.events {
				assert.NoError(t, repository.Save(context.TODO(), []*models.Event{
					{
						ShortKey:      event.ShortKey,
						OriginalURL:   event.OriginalURL,
						CorrelationID: "testCorrelationID",
						UserID:        event.UserID,
						DeletedFlag:   false,
					},
				}))

				getEvent, found := repository.Get(event.ShortKey)
				assert.True(t, found)
				assert.Equal(t, event.OriginalURL, getEvent.OriginalURL)

				userEvents := repository.GetEventsByUserID(context.TODO(), event.UserID)
				assert.Len(t, userEvents, 1)

				assert.NoError(t, repository.Delete(context.TODO(), []models.DeleteRequestBatch{
					{
						ShortKeys: []string{event.ShortKey},
						UserID:    event.UserID,
					},
				}))

				getEvent, found = repository.Get(event.ShortKey)
				assert.True(t, found)
				assert.True(t, getEvent.DeletedFlag)
			}
		})
	}
}
