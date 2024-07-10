package models

import (
	"github.com/ShukinDmitriy/shortener/internal/environments"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateShortKey(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "success",
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateShortKey(); len(got) != tt.want {
				t.Errorf("GenerateShortKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkGenerateShortKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateShortKey()
	}
}

func TestPrepareFullURL(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	type args struct {
		ctx      echo.Context
		shortKey string
		baseAddr string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: args{
				ctx:      c,
				shortKey: "test",
				baseAddr: "https://test.com",
			},
			want: "https://test.com/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.baseAddr != "" {
				environments.FlagBaseAddr = tt.args.baseAddr
			}

			if got := PrepareFullURL(tt.args.ctx, tt.args.shortKey); got != tt.want {
				t.Errorf("PrepareFullURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
