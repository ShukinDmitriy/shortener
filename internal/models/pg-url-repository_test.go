package models_test

import (
	"context"
	"github.com/ShukinDmitriy/shortener/internal/environments"
	"github.com/ShukinDmitriy/shortener/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPGURLRepository_Initialize(t *testing.T) {
	// Будем скипать тест через 10 секунд, т.к. в github нет запуска бд
	finalTest := make(chan interface{})
	go func() {
		select {
		case <-finalTest:
			return
		case <-time.After(10 * time.Second):
			t.Skip("Skipping testing in CI environment")
		}
	}()

	type args struct {
		dsn string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test #1",
			args: args{
				dsn: "postgresql://postgres:postgres@192.168.160.11:5432/praktikum?sslmode=disable",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.dsn != "" {
				environments.FlagDatabaseDSN = tt.args.dsn
			}

			repository := &models.PGURLRepository{}
			assert.NoError(t, repository.Initialize())
		})
	}

	finalTest <- struct{}{}
}

func TestPGURLRepository_CRUD(t *testing.T) {
	// Будем скипать тест через 10 секунд, т.к. в github нет запуска бд
	finalTest := make(chan interface{})
	go func() {
		select {
		case <-finalTest:
			return
		case <-time.After(10 * time.Second):
			t.Skip("Skipping testing in CI environment")
		}
	}()

	type args struct {
		dsn    string
		events []models.Event
	}
	repeatShortKey := models.GenerateShortKey()
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test #1",
			args: args{
				dsn: "postgresql://postgres:postgres@192.168.160.11:5432/praktikum?sslmode=disable",
				events: []models.Event{
					{
						OriginalURL: "https://example.com",
						ShortKey:    repeatShortKey,
						UserID:      "1",
					},
					{
						OriginalURL: "https://example.com",
						ShortKey:    repeatShortKey,
						UserID:      "2",
					},
					{
						OriginalURL: "https://example1.com",
						ShortKey:    models.GenerateShortKey(),
						UserID:      "3",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.dsn != "" {
				environments.FlagDatabaseDSN = tt.args.dsn
			}

			repository := &models.PGURLRepository{}
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

				shortKey, found := repository.GetShortKeyByOriginalURL(event.OriginalURL)
				assert.True(t, found)
				assert.Equal(t, shortKey, event.ShortKey)

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

	finalTest <- struct{}{}
}
