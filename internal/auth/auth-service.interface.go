package auth

import "github.com/labstack/echo/v4"

// AuthServiceInterface interface for auth service
type AuthServiceInterface interface {
	// GetUserID get user
	GetUserID(c echo.Context) string
}
