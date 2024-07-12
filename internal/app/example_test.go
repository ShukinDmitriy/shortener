package app_test

import (
	"encoding/json"
	"fmt"
	"github.com/ShukinDmitriy/shortener/internal/app"
	"github.com/ShukinDmitriy/shortener/internal/auth"
	"github.com/ShukinDmitriy/shortener/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strings"
)

func Example() {
	type bodyItem struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}
	body := []bodyItem{
		{
			CorrelationID: "847b5414-7f41-4363-be2a-e316fbfc2b33",
			OriginalURL:   "https://practicum.yandex.ru",
		},
		{
			CorrelationID: "022d3f81-2fb5-4fda-bb19-e89bad595b09",
			OriginalURL:   "https://yandex.ru",
		},
		{
			CorrelationID: "847b5414-7f41-4363-be2a-e316fbfc2b33",
			OriginalURL:   "https://music.yandex.ru",
		},
	}

	repository := &models.MemoryURLRepository{}
	err := repository.Initialize()
	if err != nil {
		panic(err)
	}
	authService := auth.NewAuthService()
	shortener := app.NewURLShortener(repository, nil, authService)

	e := echo.New()
	stringBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(string(stringBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = shortener.HandleCreateShortenBatch(c)
	if err != nil {
		panic(err)
	}

	res := rec.Result()
	defer res.Body.Close()

	responseBody := rec.Body.Bytes()
	var data []bodyItem
	json.Unmarshal(responseBody, &data)

	fmt.Println(len(data))

	// Output:
	// 3
}
