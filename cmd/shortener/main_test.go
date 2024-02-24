package main

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
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
				response:    "empty url",
				contentType: "text/plain; charset=utf-8",
			},
			body: "",
		},
	}

	var shortener = &URLShortener{
		urls: make(map[string]string),
	}

	e := echo.New()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := shortener.HandleShorten(c)

			// Assertions
			if err != nil {
				res, ok := err.(*echo.HTTPError)

				require.NotNil(t, ok)
				assert.Equal(t, test.want.code, res.Code)

				resBody := res.Message
				assert.Contains(t, resBody, test.want.response)
			} else {
				res := rec.Result()

				assert.Equal(t, test.want.code, res.StatusCode)

				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)

				require.NoError(t, err)
				assert.Contains(t, string(resBody), test.want.response)
				assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func TestURLShortener_HandleCreateShorten(t *testing.T) {
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
				contentType: "application/json; charset=UTF-8",
			},
			body: `{"url":"https://yandex.ru"}`,
		},
		{
			name: "negative test #1",
			want: want{
				code:        400,
				response:    "empty url",
				contentType: "application/json; charset=UTF-8",
			},
			body: `{"url":""}`,
		},
		{
			name: "negative test #2",
			want: want{
				code:        400,
				response:    "empty url",
				contentType: "application/json; charset=UTF-8",
			},
			body: `{"test":"test"}`,
		},
		{
			name: "negative test #3",
			want: want{
				code:        500,
				response:    "invalid JSON",
				contentType: "application/json; charset=UTF-8",
			},
			body: `{"test" "test"}`,
		},
	}

	var shortener = &URLShortener{
		urls: make(map[string]string),
	}

	e := echo.New()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(test.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := shortener.HandleCreateShorten(c)

			// Assertions
			if err != nil {
				res, ok := err.(*echo.HTTPError)

				require.NotNil(t, ok)
				assert.Equal(t, test.want.code, res.Code)

				resBody := res.Message
				assert.Contains(t, resBody, test.want.response)
			} else {
				res := rec.Result()
				defer res.Body.Close()

				body := rec.Body.Bytes()
				var data struct {
					Result string `json:"result"`
				}
				json.Unmarshal(body, &data)

				assert.Equal(t, test.want.code, res.StatusCode)
				require.NoError(t, err)
				assert.Contains(t, data.Result, test.want.response)
				assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			}
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
				response:    "https://yandex.ru",
				contentType: "",
			},
			preRequest: "https://yandex.ru",
		},
		{
			name: "negative test #1",
			want: want{
				code:        400,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #2",
			want: want{
				code:        404,
				response:    "URL not found",
				contentType: "text/plain; charset=utf-8",
			},
			preRequest: "https://yandex.ru",
			target:     "/incorrectUrl",
		},
	}

	var shortener = &URLShortener{
		urls: make(map[string]string),
	}

	e := echo.New()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var shortURL []byte

			if test.preRequest != "" {
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.preRequest))

				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				shortener.HandleShorten(c)

				res := rec.Result()

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

			req := httptest.NewRequest(http.MethodGet, target, nil)
			// создаём новый Recorder
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if len(target) > 6 {
				c.SetPath("/:id")
				c.SetParamNames("id")
				c.SetParamValues(target[len(target)-6:])
			}

			err := shortener.HandleRedirect(c)

			// Assertions
			if err != nil {
				res, ok := err.(*echo.HTTPError)

				require.NotNil(t, ok)
				assert.Equal(t, test.want.code, res.Code)

				resBody := res.Message
				assert.Contains(t, resBody, test.want.response)
			} else {
				res := rec.Result()
				defer res.Body.Close()

				assert.Equal(t, test.want.code, res.StatusCode)

				require.NoError(t, err)
				assert.Contains(t, res.Header.Get("Location"), test.want.response)
				assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}

}
