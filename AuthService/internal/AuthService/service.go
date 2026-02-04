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

func (as *authService) SignUp(u UserSignUp) (int, error) {
	passwordHash, err := as.makePasswordHash(u.Password)
	if err != nil {
		return 0, err
	}

	var user = User{
		UserName:     u.UserName,
		Email:        u.Email,
		PasswordHash: passwordHash,
	}

	createdUser, err := as.repoPostgres.CreateUser(user)
	if err != nil {
		return 0, err
	}

	go func() {
		if errCache := as.repoRedis.AddUserWithEmail(createdUser); errCache != nil {
			logrus.Warn("can't save user to cache")
		}
	}()

	return createdUser.ID, nil
}

func (as *authService) SignIn(u UserSignIn) (Tokens, error) {
	var user User
	var err error
	user, err = as.repoRedis.GetUserWithEmail(u.Email)
	if err != nil {
		logrus.Warn(err)

		user, err = as.repoPostgres.GetUser(u.Email)
		if err != nil {
			return Tokens{}, err
		}

		go func() {
			if errCache := as.repoRedis.AddUserWithEmail(user); errCache != nil {
				logrus.Warn("can't save user to cache")
			}
		}()
	}

	if !as.checkPassword(u.Password, user.PasswordHash) {
		return Tokens{}, errors.New("password is wrong")
	}

	return as.generateTokensPair(user)
}

func (as *authService) RefreshTokens(refreshToken string) (string, error) {
	hashRefreshToken := as.makeHash(refreshToken)

	var user User
	var err error

	user, err = as.repoRedis.GetUserWithRefreshToken(hashRefreshToken)
	if err != nil {
		logrus.Warn(err)

		user, err = as.repoPostgres.GetUserByRefreshToken(hashRefreshToken)
		if err != nil {
			return "", err
		}

		go func() {
			if errCache := as.repoRedis.AddUserWithRefreshToken(user, hashRefreshToken); errCache != nil {
				logrus.Warn("can't save user to cache")
			}
		}()
	}

	accessToken, errAccess := as.tokenManager.NewJWT(user, accessTTL)
	if errAccess != nil {
		return "", errAccess
	}

	return accessToken, nil
}

func (as *authService) generateTokensPair(user User) (Tokens, error) {
	accessToken, errAccess := as.tokenManager.NewJWT(user, accessTTL)
	if errAccess != nil {
		return Tokens{}, errAccess
	}

	refreshToken, errRefresh := as.tokenManager.NewRefreshToken()
	if errRefresh != nil {
		return Tokens{}, errRefresh
	}

	hashRefreshToken := as.makeHash(refreshToken)

	refreshExpiresAt := time.Now().Add(refreshTTL)

	refreshTokenStruct := RefreshToken{
		UserID:    user.ID,
		TokenHash: hashRefreshToken,
		ExpiresAt: refreshExpiresAt,
	}

	errSave := as.repoPostgres.SaveRefreshToken(refreshTokenStruct)
	if errSave != nil {
		return Tokens{}, errSave
	}

	go func() {
		if errCache := as.repoRedis.AddUserWithRefreshToken(user, hashRefreshToken); errCache != nil {
			logrus.Warn("can't save user to cache")
		}
	}()

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (as *authService) ParseToken(accessToken string) (int, string, error) {
	return as.tokenManager.Parse(accessToken)
}

func (as *authService) makeHash(str string) string {
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])
}

func (as *authService) makePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (as *authService) checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
