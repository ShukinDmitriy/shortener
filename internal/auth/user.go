package auth

import (
	"github.com/google/uuid"
)

type User struct {
	ID string `json:"id"`
}

func LoadTestUser() *User {
	return &User{
		ID: uuid.New().String(),
	}
}
