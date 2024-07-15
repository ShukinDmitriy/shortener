package auth_test

import (
	"github.com/ShukinDmitriy/shortener/internal/auth"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateTokenWithConfig(t *testing.T) {
	type args struct {
		targetPath  string
		accessToken string
	}
	type want struct {
		statusCode int
		hasCookie  bool
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
				statusCode: 200,
				hasCookie:  true,
			},
		},
		{
			name: "positive test #2",
			args: args{
				targetPath: "/skip",
			},
			want: want{
				statusCode: 200,
			},
		},
		{
			name: "positive test #3",
			args: args{
				targetPath:  "/",
				accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjQxNGU1NWIzLWQyNWYtNDYyZC1hN2NjLTY4MTQ4OTM0ODhkOCIsImV4cCI6NTMyMTAyMzIyMn0.hArR7cfGcdRExHwzeA4_S1fMIZwW3tKfERt2jHgIYrY",
			},
			want: want{
				statusCode: 200,
			},
		},
		{
			name: "negative test #1",
			args: args{
				targetPath:  "/",
				accessToken: "invalidToken",
			},
			want: want{
				statusCode: 401,
			},
		},
		{
			name: "negative test #2",
			args: args{
				targetPath:  "/",
				accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjQxNGU1NWIzLWQyNWYtNDYyZC1hN2NjLTY4MTQ4OTM0ODhkOCIsImV4cCI6MTcyMTExMjM4OH0.IRdlAZia-AXDJeqAWszGpL7M2MGQ1g0M5s71BfiFVbk",
			},
			want: want{
				statusCode: 401,
			},
		},
	}
	e := echo.New()
	e.Use(auth.CreateTokenWithConfig(auth.CreateTokenConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "/skip")
		},
	}))
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, nil)
	})
	e.GET("/skip", func(c echo.Context) error {
		return c.JSON(http.StatusOK, nil)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, tt.args.targetPath, nil)
			if tt.args.accessToken != "" {
				req.Header.Set("Cookie", "access-token="+tt.args.accessToken)
			}

			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			// Проверяем, что ответ сервера OK
			assert.Equal(t, tt.want.statusCode, res.StatusCode)

			headers := res.Header.Values("Set-Cookie")
			hasAccessCookie := false
			for _, h := range headers {
				tmp := strings.Split(h, "=")
				if tmp[0] == auth.GetAccessTokenCookieName() {
					hasAccessCookie = true
					break
				}
			}
			assert.Equal(t, tt.want.hasCookie, hasAccessCookie)
		})
	}
}
