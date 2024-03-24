package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Get("user") == nil {
			return next(c)
		}
		u := c.Get("user").(*jwt.Token)

		claims := u.Claims.(*Claims)

		if time.Unix(claims.ExpiresAt.Unix(), 0).Sub(time.Now()) < 15*time.Minute {
			rc, err := c.Cookie(refreshTokenCookieName)
			if err == nil && rc != nil {
				tkn, err := jwt.ParseWithClaims(rc.Value, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(GetRefreshJWTSecret()), nil
				})
				if err != nil {
					if err == jwt.ErrSignatureInvalid {
						c.Response().Writer.WriteHeader(http.StatusUnauthorized)
					}
				}

				if tkn != nil && tkn.Valid {
					_ = GenerateTokensAndSetCookies(&User{
						ID: claims.ID,
					}, c)
				}
			}
		}

		return next(c)
	}
}
