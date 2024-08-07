package models

import (
	"bufio"
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// NewProducer create producer
func NewProducer(filename string) (*Producer, error) {
	dirPaths := strings.Split(filename, string(filepath.Separator))

	if len(dirPaths) > 1 {
		directory := path.Join(dirPaths[0 : len(dirPaths)-1]...)

		// создаем директорию, если она не существует
		if err := os.MkdirAll(directory, 0o755); err != nil {
			return nil, err
		}
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file: file,
		// создаём новый Writer
		writer: bufio.NewWriter(file),
	}, nil
}

// WriteEvent to the file
func (p *Producer) WriteEvent(event interface{}) error {
	data, err := json.Marshal(event)
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

// NewConsumer create consumer
func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0o666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file: file,
		// создаём новый scanner
		scanner: bufio.NewScanner(file),
	}, nil
}

// ReadEvent from file
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

// Close the file
func (c *Consumer) Close() error {
	return c.file.Close()
}
