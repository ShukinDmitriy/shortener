package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ShukinDmitriy/shortener/internal/middleware"
	logger2 "github.com/ShukinDmitriy/shortener/mocks/internal_/logger"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

func TestRequestInfo(t *testing.T) {
	type args struct {
		targetPath string
	}
	type want struct {
		infoCallCount int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive test #1",
			args: args{
				targetPath: "/",
			},
			want: want{
				infoCallCount: 1,
			},
		},
	}

	for _, tt := range tests {
		e := echo.New()
		mockLogger := new(logger2.Logger)
		mockLogger.EXPECT().Info(
			"got incoming HTTP request",
			mock.AnythingOfType("zapcore.Field"),
			mock.AnythingOfType("zapcore.Field"),
			mock.AnythingOfType("zapcore.Field"),
		).Return()
		e.Use(middleware.RequestInfo(mockLogger))
		e.GET("/", func(c echo.Context) error {
			return c.JSON(http.StatusOK, nil)
		})

		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.args.targetPath, nil)

			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			mockLogger.AssertNumberOfCalls(t, "Info", tt.want.infoCallCount)
		})
	}
}
