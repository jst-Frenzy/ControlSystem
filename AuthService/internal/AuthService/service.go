package AuthService

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks/mockServ.go

type AuthService interface {
	SignUp(u UserSignUp) (int, error)
	SignIn(u UserSignIn) (Tokens, error)
	RefreshTokens(refreshToken string) (string, error)
	ParseToken(accessToken string) (int, string, error)
	ChangeRole(UserSignIn, int, string) error
}

type authService struct {
	repoPostgres AuthPostgresRepo
	repoRedis    AuthRedisRepo
	tokenManager TokenManager
}

func NewAuthService(repoPostgres AuthPostgresRepo, repoRedis AuthRedisRepo, tokenManager TokenManager) AuthService {
	return &authService{
		repoPostgres: repoPostgres,
		repoRedis:    repoRedis,
		tokenManager: tokenManager,
	}
}

func (s *authService) SignUp(u UserSignUp) (int, error) {
	passwordHash, err := s.makePasswordHash(u.Password)
	if err != nil {
		return 0, err
	}

	var user = User{
		UserName:     u.UserName,
		Email:        u.Email,
		PasswordHash: passwordHash,
		Role:         "user",
	}

	createdUser, err := s.repoPostgres.CreateUser(user)
	if err != nil {
		return 0, err
	}

	go func() {
		if errCache := s.repoRedis.AddUserWithEmail(createdUser); errCache != nil {
			logrus.Warn("can't save user to cache")
		}
	}()

	return createdUser.ID, nil
}

func (s *authService) SignIn(u UserSignIn) (Tokens, error) {
	var user User
	var err error
	user, err = s.repoRedis.GetUserWithEmail(u.Email)
	if err != nil {
		logrus.Warn(err)

		user, err = s.repoPostgres.GetUser(u.Email)
		if err != nil {
			return Tokens{}, err
		}

		go func() {
			if errCache := s.repoRedis.AddUserWithEmail(user); errCache != nil {
				logrus.Warn("can't save user to cache")
			}
		}()
	}

	if !s.checkPassword(u.Password, user.PasswordHash) {
		return Tokens{}, errors.New("password is wrong")
	}

	return s.generateTokensPair(user)
}

func (s *authService) RefreshTokens(refreshToken string) (string, error) {
	hashRefreshToken := s.makeHash(refreshToken)

	var user User
	var err error

	user, err = s.repoRedis.GetUserWithRefreshToken(hashRefreshToken)
	if err != nil {
		logrus.Warn(err)

		user, err = s.repoPostgres.GetUserByRefreshToken(hashRefreshToken)
		if err != nil {
			return "", err
		}

		go func() {
			if errCache := s.repoRedis.AddUserWithRefreshToken(user, hashRefreshToken); errCache != nil {
				logrus.Warn("can't save user to cache")
			}
		}()
	}

	accessToken, errAccess := s.tokenManager.NewJWT(user, accessTTL)
	if errAccess != nil {
		return "", errAccess
	}

	return accessToken, nil
}

func (s *authService) ChangeRole(user UserSignIn, id int, newRole string) error {
	if user.Email == "admin" && user.Password == "admin" {
		return s.repoPostgres.ChangeRole(id, newRole)
	}
	return errors.New("not enough rights")
}

func (s *authService) generateTokensPair(user User) (Tokens, error) {
	accessToken, errAccess := s.tokenManager.NewJWT(user, accessTTL)
	if errAccess != nil {
		return Tokens{}, errAccess
	}

	refreshToken, errRefresh := s.tokenManager.NewRefreshToken()
	if errRefresh != nil {
		return Tokens{}, errRefresh
	}

	hashRefreshToken := s.makeHash(refreshToken)

	refreshExpiresAt := time.Now().Add(refreshTTL)

	refreshTokenStruct := RefreshToken{
		UserID:    user.ID,
		TokenHash: hashRefreshToken,
		ExpiresAt: refreshExpiresAt,
	}

	errSave := s.repoPostgres.SaveRefreshToken(refreshTokenStruct)
	if errSave != nil {
		return Tokens{}, errSave
	}

	go func() {
		if errCache := s.repoRedis.AddUserWithRefreshToken(user, hashRefreshToken); errCache != nil {
			logrus.Warn("can't save user to cache")
		}
	}()

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) ParseToken(accessToken string) (int, string, error) {
	return s.tokenManager.Parse(accessToken)
}

func (s *authService) makeHash(str string) string {
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])
}

func (s *authService) makePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *authService) checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
