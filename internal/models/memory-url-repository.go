package models

import (
	"context"
	"github.com/ShukinDmitriy/shortener/internal/environments"
)

type MemoryURLRepository struct {
	DBConsumer *Consumer
	DBProducer *Producer
	urls       map[string]Event
}

func (r *MemoryURLRepository) Initialize() error {
	r.urls = make(map[string]Event)

	filename := environments.FlagFileStoragePath
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

func (r *MemoryURLRepository) Get(shortKey string) (string, bool) {
	// Поиск в памяти
	var event Event
	var found = false

	event, found = r.urls[shortKey]

	return event.OriginalURL, found
}

func (r *MemoryURLRepository) Save(ctx context.Context, events []*Event) error {
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

func (r *MemoryURLRepository) GetShortKeyByOriginalURL(originalURL string) (string, bool) {
	for _, event := range r.urls {
		if event.OriginalURL == originalURL {
			return event.ShortKey, true
		}
	}

	return "", false
}

func (r *MemoryURLRepository) GetEventsByUserID(ctx context.Context, userID string) []*Event {
	var events []*Event
	for _, event := range r.urls {
		if event.UserID == userID {
			events = append(events, &event)
		}
	}

	return events
}
