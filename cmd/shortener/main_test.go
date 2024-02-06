package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestURLShortener_HandleShorten(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		want want
		body string
	}{
		{
			name: "positive test #1",
			want: want{
				code:        201,
				response:    "http://example.com/",
				contentType: "text/plain; charset=utf-8",
			},
			body: "https://yandex.ru",
		},
		{
			name: "negative test #1",
			want: want{
				code:        400,
				response:    "empty url\n",
				contentType: "text/plain; charset=utf-8",
			},
			body: "",
		},
	}

	var shortener = &URLShortener{
		urls: make(map[string]string),
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
			// создаём новый Recorder
			w := httptest.NewRecorder()
			shortener.HandleShorten(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Contains(t, string(resBody), test.want.response)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestURLShortener_HandleRedirect(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name       string
		want       want
		preRequest string
		target     string
	}{
		{
			name: "positive test #1",
			want: want{
				code:        307,
				response:    "<a href=\"https://yandex.ru\">Temporary Redirect</a>.\n\n",
				contentType: "text/html; charset=utf-8",
			},
			preRequest: "https://yandex.ru",
		},
		{
			name: "negative test #1",
			want: want{
				code:        400,
				response:    "\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #2",
			want: want{
				code:        400,
				response:    "\n",
				contentType: "text/plain; charset=utf-8",
			},
			preRequest: "https://yandex.ru",
			target:     "/incorrectUrl",
		},
	}

	var shortener = &URLShortener{
		urls: make(map[string]string),
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var shortURL []byte

			if test.preRequest != "" {
				request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.preRequest))
				// создаём новый Recorder
				w := httptest.NewRecorder()
				shortener.HandleShorten(w, request)

				res := w.Result()
				// получаем и проверяем тело запроса
				defer res.Body.Close()
				shortURL, _ = io.ReadAll(res.Body)
			}

			target := string(shortURL)
			if target == "" {
				target = "/"
			}
			if test.target != "" {
				target = test.target
			}
			request := httptest.NewRequest(http.MethodGet, target, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			shortener.HandleRedirect(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}

}
