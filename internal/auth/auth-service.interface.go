package auth

import "github.com/labstack/echo/v4"

type AuthServiceInterface interface {
	GetUserID(c echo.Context) string
}
