package models

import (
	"context"

	"github.com/ShukinDmitriy/shortener/internal/environments"
)

// MemoryURLRepository repository for working with a memory
type MemoryURLRepository struct {
	DBConsumer *Consumer
	DBProducer *Producer
	urls       map[string]Event
}

// Initialize repository
func (r *MemoryURLRepository) Initialize(configuration environments.Configuration) error {
	r.urls = make(map[string]Event)

	filename := configuration.FileStoragePath
	if filename == "" {
		return nil
	}

	var err error

	r.DBProducer, err = NewProducer(filename)
	if err != nil {
		return err
	}

	r.DBConsumer, err = NewConsumer(filename)
	if err != nil {
		return err
	}

	var event *Event

	defer r.DBConsumer.Close()

	for {
		event, err = r.DBConsumer.ReadEvent()

		if event == nil || err != nil {
			return err
		}

		// Сохраняем значение в память, т.к. повторно файл не вычитывается
		r.urls[event.ShortKey] = *event
	}
}

// Get event by short key
func (r *MemoryURLRepository) Get(shortKey string) (Event, bool) {
	// Поиск в памяти
	var event Event
	found := false

	event, found = r.urls[shortKey]

	return event, found
}

// Save batch save events
func (r *MemoryURLRepository) Save(_ context.Context, events []*Event) error {
	for _, event := range events {
		shortKey, found := r.GetShortKeyByOriginalURL(event.OriginalURL)
		if found {
			event.ShortKey = shortKey
			continue
		}

		// Хранение в памяти
		r.urls[event.ShortKey] = *event

		if r.DBProducer == nil {
			continue
		}

		// Хранение в файле
		r.DBProducer.WriteEvent(event)
	}

	return nil
}

// Delete batch delete event
func (r *MemoryURLRepository) Delete(_ context.Context, events []DeleteRequestBatch) error {
	for _, deleteEvent := range events {
		for _, shortKey := range deleteEvent.ShortKeys {
			event := r.urls[shortKey]

			if event.DeletedFlag || event.UserID != deleteEvent.UserID {
				continue
			}

			event.DeletedFlag = true

			// Хранение в памяти
			r.urls[event.ShortKey] = event

			if r.DBProducer == nil {
				continue
			}

			// Хранение в файле
			r.DBProducer.WriteEvent(&event)
		}
	}

	return nil
}

// GetShortKeyByOriginalURL get short link from full link
func (r *MemoryURLRepository) GetShortKeyByOriginalURL(originalURL string) (string, bool) {
	for _, event := range r.urls {
		if event.OriginalURL == originalURL && !event.DeletedFlag {
			return event.ShortKey, true
		}
	}

	return "", false
}

// GetEventsByUserID get events by user ID
func (r *MemoryURLRepository) GetEventsByUserID(_ context.Context, userID string) []*Event {
	var events []*Event
	for _, event := range r.urls {
		if event.UserID == userID && !event.DeletedFlag {
			events = append(events, &event)
		}
	}

	return events
}
