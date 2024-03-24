package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type CreateTokenConfig struct {
	Skipper middleware.Skipper
}

func CreateTokenWithConfig(config CreateTokenConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			if c.Get("user") != nil {
				return next(c)
			}

			storedUser := LoadTestUser()
			err := GenerateTokensAndSetCookies(storedUser, c)

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token is incorrect")
			}

			SetUser(storedUser)

			return next(c)
		}
	}
}

func DefaultSkipper(echo.Context) bool {
	return false
}
