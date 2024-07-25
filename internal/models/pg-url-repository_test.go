package models_test

import (
	"context"
	"os"
	"testing"

	"github.com/ShukinDmitriy/shortener/internal/environments"
	"github.com/ShukinDmitriy/shortener/internal/models"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestPGURLRepository_Initialize(t *testing.T) {
	// Будем скипать тест если нет переменных в test.env
	godotenv.Load("../../test.env")
	databaseDSN := os.Getenv("DATABASE_DSN")
	if databaseDSN == "" {
		t.Skip("Skipping testing")
	}

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
				dsn: databaseDSN,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configuration := environments.Configuration{}
			if tt.args.dsn != "" {
				configuration.DatabaseDSN = tt.args.dsn
			}

			repository := &models.PGURLRepository{}
			assert.NoError(t, repository.Initialize(configuration))
		})
	}
}

func TestPGURLRepository_CRUD(t *testing.T) {
	// Будем скипать тест если нет переменных в test.env
	godotenv.Load("../../test.env")
	databaseDSN := os.Getenv("DATABASE_DSN")
	if databaseDSN == "" {
		t.Skip("Skipping testing")
	}

	type args struct {
		dsn    string
		events []models.Event
	}
	type want struct {
		countUsers int
		countURLs  int
	}
	eventShortKey := models.GenerateShortKey()
	repeatShortKey := models.GenerateShortKey()
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive test #1",
			args: args{
				dsn: databaseDSN,
				events: []models.Event{
					{
						OriginalURL: "https://" + repeatShortKey + ".com",
						ShortKey:    repeatShortKey,
						UserID:      repeatShortKey,
					},
					{
						OriginalURL: "https://" + repeatShortKey + ".com",
						ShortKey:    models.GenerateShortKey(),
						UserID:      models.GenerateShortKey(),
					},
					{
						OriginalURL: "https://" + eventShortKey + ".com",
						ShortKey:    eventShortKey,
						UserID:      eventShortKey,
					},
				},
			},
			want: want{
				countUsers: 2,
				countURLs:  2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configuration := environments.Configuration{}
			if tt.args.dsn != "" {
				configuration.DatabaseDSN = tt.args.dsn
			}

			repository := &models.PGURLRepository{}
			assert.NoError(t, repository.Initialize(configuration))

			originalURLs := make(map[string]bool, 3)

			countUsersBefore, countURLsBefore, err := repository.GetStats(context.TODO())
			if err != nil {
				t.Error(err)
			}

			for _, event := range tt.args.events {
				// notExpectedFound ожидается ли в бд такое событие
				_, notExpectedFound := originalURLs[event.OriginalURL]
				err := repository.Save(context.TODO(), []*models.Event{
					{
						ShortKey:      event.ShortKey,
						OriginalURL:   event.OriginalURL,
						CorrelationID: "testCorrelationID",
						UserID:        event.UserID,
						DeletedFlag:   false,
					},
				})
				originalURLs[event.OriginalURL] = true

				if !notExpectedFound {
					assert.NoError(t, err)
				}

				getEvent, found := repository.Get(event.ShortKey)

				assert.Equal(t, found, !notExpectedFound)
				if !notExpectedFound {
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
			}

			countUsers, countURLs, err := repository.GetStats(context.TODO())
			assert.Equal(t, tt.want.countUsers, countUsers-countUsersBefore)
			assert.Equal(t, tt.want.countURLs, countURLs-countURLsBefore)
			assert.Nil(t, err)
		})
	}
}
