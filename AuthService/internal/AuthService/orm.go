package AuthService

import "time"

const (
	userEmailTTL        = 15 * time.Minute
	userRefreshTokenTTL = 45 * time.Minute
)

type User struct {
	ID       int    `redis:"ID"`
	UserName string `redis:"UserName"`

	Email        string `redis:"Email"`
	PasswordHash string `redis:"PasswordHash"`

	Role string `redis:"Role"`

	CreatedAt time.Time `redis:"CreatedAt"`
	UpdatedAt time.Time `redis:"UpdatedAt"`
}

type RefreshToken struct {
	ID     int
	UserID int

	TokenHash string

	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}
