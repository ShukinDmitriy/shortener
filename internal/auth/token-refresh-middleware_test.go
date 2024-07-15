package auth_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ShukinDmitriy/shortener/internal/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestTokenRefreshMiddleware(t *testing.T) {
	longAccessExpirationTime := time.Now().Add(20 * time.Minute)
	accessExpirationTime := time.Now().Add(5 * time.Minute)
	refreshExpirationTime := time.Now().Add(1 * time.Hour)
	longAccessClaims := &auth.Claims{
		ID: "123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(longAccessExpirationTime),
		},
	}
	accessClaims := &auth.Claims{
		ID: "123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpirationTime),
		},
	}
	refreshClaims := &auth.Claims{
		ID: "123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
		},
	}

	longAccessToken := jwt.NewWithClaims(auth.GetSigningMethod(), longAccessClaims)
	longAccessTokenString, _ := longAccessToken.SignedString([]byte(auth.GetJWTSecret()))
	accessToken := jwt.NewWithClaims(auth.GetSigningMethod(), accessClaims)
	accessTokenString, _ := accessToken.SignedString([]byte(auth.GetJWTSecret()))
	refreshToken := jwt.NewWithClaims(auth.GetSigningMethod(), refreshClaims)
	refreshTokenString, _ := refreshToken.SignedString([]byte(auth.GetRefreshJWTSecret()))

	type args struct {
		targetPath   string
		accessToken  string
		refreshToken string
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
				targetPath: "/skip",
			},
			want: want{
				statusCode: 200,
				hasCookie:  false,
			},
		},
		{
			name: "positive test #2",
			args: args{
				targetPath: "/",
			},
			want: want{
				statusCode: 200,
				hasCookie:  true,
			},
		},
		{
			name: "positive test #3",
			args: args{
				targetPath:  "/",
				accessToken: longAccessTokenString,
			},
			want: want{
				statusCode: 200,
				hasCookie:  false,
			},
		},
		{
			name: "positive test #4",
			args: args{
				targetPath:  "/",
				accessToken: accessTokenString,
				// refreshToken incorrect
				refreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjQxNGU1NWIzLWQyNWYtNDYyZC1hN2NjLTY4MTQ4OTM0ODhkOCIsImV4cCI6MTc1MjU2ODM1NH0.esFYBIffRE2xcWTjzZMKVy4ExICqKFzezzGRxSopVv8",
			},
			want: want{
				statusCode: 200,
				hasCookie:  false,
			},
		},
		{
			name: "positive test #5",
			args: args{
				targetPath:   "/",
				accessToken:  accessTokenString,
				refreshToken: refreshTokenString,
			},
			want: want{
				statusCode: 200,
				hasCookie:  true,
			},
		},
	}
	e := echo.New()
	e.Use(auth.CreateTokenWithConfig(auth.CreateTokenConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "/skip")
		},
	}))
	e.Use(auth.TokenRefreshMiddleware)
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, nil)
	})
	e.GET("/skip", func(c echo.Context) error {
		return c.JSON(http.StatusOK, nil)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.args.targetPath, nil)
			cookieString := []string{}
			if tt.args.accessToken != "" {
				cookieString = append(cookieString, "access-token="+tt.args.accessToken)
			}
			if tt.args.refreshToken != "" {
				cookieString = append(cookieString, "refresh-token="+tt.args.refreshToken)
			}
			if len(cookieString) > 0 {
				req.Header.Set("Cookie", strings.Join(cookieString, "; "))
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
