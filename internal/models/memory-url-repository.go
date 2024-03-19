package models

import "github.com/ShukinDmitriy/shortener/internal/environments"

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

func (r *MemoryURLRepository) Save(events []Event) error {
	for _, event := range events {
		// Хранение в памяти
		r.urls[event.ShortKey] = event

		if r.DBProducer == nil {
			continue
		}

		// Хранение в файле
		r.DBProducer.WriteEvent(&event)
	}

	return nil
}

func (r *MemoryURLRepository) SaveBatch(events []Event) error {
	return r.Save(events)
}
