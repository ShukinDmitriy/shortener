package models

import (
	"bufio"
	"os"
)

type CreateRequest struct {
	URL string `json:"url"`
}

type CreateResponse struct {
	Result string `json:"result"`
}

type Event struct {
	ShortKey      string `json:"short_key"`
	OriginalURL   string `json:"original_url"`
	CorrelationId string `json:"correlation_id"`
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
