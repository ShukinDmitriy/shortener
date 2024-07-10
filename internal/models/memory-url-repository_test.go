package models_test

import (
	"testing"

	"github.com/ShukinDmitriy/shortener/internal/models"
)

func BenchmarkMemoryURLRepository_Initialize(b *testing.B) {
	repository := &models.MemoryURLRepository{}

	b.Run("initialize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = repository.Initialize()
		}
	})
}
