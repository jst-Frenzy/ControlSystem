package AuthService

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	refreshTTL = 720 * time.Hour
	accessTTL  = 120 * time.Minute
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type CustomClaims struct {
	Role string
	jwt.StandardClaims
}

type RefreshTokenRequest struct {
	RefreshToken string `binding:"required"`
}

type UserSignUp struct {
	UserName string `json:"user_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserSignIn struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
