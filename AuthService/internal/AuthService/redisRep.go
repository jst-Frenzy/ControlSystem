package AuthService

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=redisRep.go -destination=mocks/mockRedis.go

type AuthRedisRepo interface {
	AddUserWithEmail(User) error
	GetUserWithEmail(string) (User, error)
	AddUserWithRefreshToken(User, string) error
	GetUserWithRefreshToken(string) (User, error)
	EditRoleWithRefreshToken(string, string) error
}

type authRedisRepo struct {
	db  *redis.Client
	ctx context.Context
}

func NewAuthRedisRepo(db *redis.Client) AuthRedisRepo {
	return &authRedisRepo{
		db:  db,
		ctx: context.Background(),
	}
}

func (r *authRedisRepo) AddUserWithEmail(u User) error {
	key := u.Email

	jsonData, err := json.Marshal(u)
	if err != nil {
		return err
	}

	return r.db.Set(r.ctx, key, jsonData, userEmailTTL).Err()
}

func (r *authRedisRepo) AddUserWithRefreshToken(u User, refreshTokenHash string) error {
	key := refreshTokenHash

	jsonData, err := json.Marshal(u)
	if err != nil {
		return err
	}

	return r.db.Set(r.ctx, key, jsonData, userRefreshTokenTTL).Err()
}

func (r *authRedisRepo) GetUserWithEmail(email string) (User, error) {
	var user User
	jsonData, err := r.db.Get(r.ctx, email).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return user, errors.New("user not found in cache")
		}
		return user, err
	}

	if errUnmarshal := json.Unmarshal([]byte(jsonData), &user); errUnmarshal != nil {
		return user, err
	}

	return user, nil
}

func (r *authRedisRepo) GetUserWithRefreshToken(refreshTokenHash string) (User, error) {
	var user User
	jsonData, err := r.db.Get(r.ctx, refreshTokenHash).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return user, errors.New("user not found in cache")
		}
		return user, err
	}

	if errUnmarshal := json.Unmarshal([]byte(jsonData), &user); errUnmarshal != nil {
		return user, err
	}

	return user, nil
}

func (r *authRedisRepo) EditRoleWithRefreshToken(refreshTokenHash, newRole string) error {
	u, err := r.GetUserWithRefreshToken(refreshTokenHash)
	if err != nil {
		return err
	}
	u.Role = newRole

	jsonUser, err := json.Marshal(u)
	if err != nil {
		return err
	}

	return r.db.Set(r.ctx, refreshTokenHash, jsonUser, userRefreshTokenTTL).Err()
}
