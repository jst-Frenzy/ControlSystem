package AuthService

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"time"
)

//go:generate mockgen -source=tokenManager.go -destination=../mocks/mockManager.go

type TokenManager interface {
	NewJWT(user User, ttl time.Duration) (string, error)
	NewRefreshToken() (string, error)
	Parse(accessToken string) (InfoFromToken, error)
}

type manager struct {
	signingKey string
}

func NewManager(signingKey string) TokenManager {
	return &manager{signingKey: signingKey}
}

func (m *manager) NewJWT(user User, ttl time.Duration) (string, error) {
	claims := &CustomClaims{
		Role:     user.Role,
		UserName: user.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.signingKey))
}

func (m *manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (m *manager) Parse(accessToken string) (InfoFromToken, error) {
	token, err := jwt.ParseWithClaims(accessToken, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return InfoFromToken{}, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		var userID int
		_, errScan := fmt.Sscanf(claims.Subject, "%d", &userID)
		if errScan != nil {
			return InfoFromToken{}, errors.New("invalid user ID in token")
		}
		return InfoFromToken{
			ID:       userID,
			Role:     claims.Role,
			UserName: claims.UserName,
		}, nil
	}

	return InfoFromToken{}, errors.New("invalid token")
}
