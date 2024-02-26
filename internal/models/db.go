package models

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Event struct {
	ShortKey    string `json:"shortKey"`
	OriginalURL string `json:"originalURL"`
}

type Consumer struct {
	file *os.File
	// заменяем Reader на Scanner
	scanner *bufio.Scanner
}

type Producer struct {
	file *os.File
	// добавляем Writer в Producer
	writer *bufio.Writer
}

func NewProducer(filename string) (*Producer, error) {
	dirPaths := strings.Split(filename, string(filepath.Separator))

	if len(dirPaths) > 1 {
		directory := strings.Join(dirPaths[0:len(dirPaths)-1], "/")

		// создаем директорию, если она не существует
		if err := os.MkdirAll(directory, 0755); err != nil {
			return nil, err
		}
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file: file,
		// создаём новый Writer
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *Producer) WriteEvent(event *Event) error {
	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	// записываем событие в буфер
	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	// добавляем перенос строки
	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	// записываем буфер в файл
	return p.writer.Flush()
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file: file,
		// создаём новый scanner
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *Consumer) ReadEvent() (*Event, error) {
	// одиночное сканирование до следующей строки
	if !c.scanner.Scan() {
		return nil, c.scanner.Err()
	}
	// читаем данные из scanner
	data := c.scanner.Bytes()

	event := Event{}
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

var DBConsumer *Consumer
var DBProducer *Producer

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize(filename string) error {
	if filename == "" {
		return nil
	}

	var err error

	DBProducer, err = NewProducer(filename)
	if err != nil {
		return err
	}

	DBConsumer, err = NewConsumer(filename)
	if err != nil {
		return err
	}

	return nil
}
