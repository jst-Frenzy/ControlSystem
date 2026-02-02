package AuthService

import "time"

type User struct {
	ID       int
	UserName string

	Email        string
	PasswordHash string

	Role string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type RefreshToken struct {
	ID     int
	UserID int

	TokenHash string

	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}
