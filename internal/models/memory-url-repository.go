package models

import "github.com/ShukinDmitriy/shortener/internal/environments"

type MemoryURLRepository struct {
	DBConsumer *Consumer
	DBProducer *Producer
	urls       map[string]string
}

func (r *MemoryURLRepository) Initialize() error {
	r.urls = make(map[string]string)

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
		r.urls[event.ShortKey] = event.OriginalURL
	}
}

func (r *MemoryURLRepository) Get(shortKey string) (string, bool) {
	// Поиск в памяти
	var originalURL string
	var found = false

	originalURL, found = r.urls[shortKey]

	return originalURL, found
}

func (r *MemoryURLRepository) Save(shortKey string, originalURL string) {
	// Хранение в памяти
	r.urls[shortKey] = originalURL

	if r.DBProducer == nil {
		return
	}

	// Хранение в файле
	r.DBProducer.WriteEvent(&Event{
		ShortKey:    shortKey,
		OriginalURL: originalURL,
	})
}
