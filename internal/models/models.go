// Package models some models and repositories for application
package models

import (
	"bufio"
	"os"
)

// CreateRequest request to link creation
type CreateRequest struct {
	URL string `json:"url"`
}

// CreateRequestBatch request to link creation
type CreateRequestBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// DeleteRequestBatch request to link delete
type DeleteRequestBatch struct {
	ShortKeys []string `json:"short_keys"`
	UserID    string   `json:"user_id"`
}

// CreateResponse response to link creation
type CreateResponse struct {
	Result string `json:"result"`
}

// CreateResponseBatch response to link creation
type CreateResponseBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// GetUserURLsResponse response to receiving user links
type GetUserURLsResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// Consumer for read file
type Consumer struct {
	file *os.File
	// заменяем Reader на Scanner
	scanner *bufio.Scanner
}

// Producer for write file
type Producer struct {
	file *os.File
	// добавляем Writer в Producer
	writer *bufio.Writer
}
