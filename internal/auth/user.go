package auth

type User struct {
	ID string `json:"id"`
}

func LoadTestUser() *User {
	return &User{
		ID: "testUserId",
	}
}
