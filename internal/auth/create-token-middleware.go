package auth

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CreateTokenConfig struct
type CreateTokenConfig struct {
	Skipper middleware.Skipper
}

// CreateTokenWithConfig create middleware with config
func CreateTokenWithConfig(config CreateTokenConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			accessTokenCookie, err := c.Request().Cookie(GetAccessTokenCookieName())

			if err == nil && accessTokenCookie != nil {
				claims := &Claims{}
				token, err1 := jwt.ParseWithClaims(accessTokenCookie.Value, claims,
					func(t *jwt.Token) (interface{}, error) {
						if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
							return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
						}
						return []byte(GetJWTSecret()), nil
					})
				if err1 != nil {
					return echo.NewHTTPError(http.StatusUnauthorized, "Token is incorrect")
				}

				if !token.Valid {
					return echo.NewHTTPError(http.StatusUnauthorized, "Token is incorrect")
				}

				c.Set("user", token)
			}

			if c.Get("user") != nil {
				return next(c)
			}

			storedUser := LoadTestUser()
			err = GenerateTokensAndSetCookies(c, storedUser)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token is incorrect")
			}

			return next(c)
		}
	}
}

// DefaultSkipper skipper function
func DefaultSkipper(echo.Context) bool {
	return false
}
