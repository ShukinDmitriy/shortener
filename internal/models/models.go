package models

import (
	"bufio"
	"os"
)

type CreateRequest struct {
	URL string `json:"url"`
}

type CreateRequestBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type CreateResponse struct {
	Result string `json:"result"`
}

type CreateResponseBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type GetUserURLsResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
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
