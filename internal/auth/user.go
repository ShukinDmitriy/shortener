package auth

import (
	"github.com/google/uuid"
)

// User struct
type User struct {
	ID string `json:"id"`
}

// LoadTestUser create user for tests
func LoadTestUser() *User {
	return &User{
		ID: uuid.New().String(),
	}
}
