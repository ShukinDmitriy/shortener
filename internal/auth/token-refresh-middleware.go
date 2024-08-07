package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// TokenRefreshMiddleware for refresh user token
func TokenRefreshMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		u, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return next(c)
		}

		claims, ok := u.Claims.(*Claims)
		if !ok {
			return next(c)
		}

		user := &User{
			ID: claims.ID,
		}

		if time.Until(time.Unix(claims.ExpiresAt.Unix(), 0)) < 15*time.Minute {
			rc, err := c.Cookie(refreshTokenCookieName)
			if err == nil && rc != nil {
				tkn, err := jwt.ParseWithClaims(rc.Value, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(GetRefreshJWTSecret()), nil
				})
				if err != nil {
					if errors.Is(err, jwt.ErrSignatureInvalid) {
						c.Response().Writer.WriteHeader(http.StatusUnauthorized)
					}
				}

				if tkn != nil && tkn.Valid {
					_ = GenerateTokensAndSetCookies(c, user)
				}
			}
		}

		return next(c)
	}
}
