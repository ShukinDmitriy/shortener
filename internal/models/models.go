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
