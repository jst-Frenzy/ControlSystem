package AuthService

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
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
	repo         AuthPostgresRepo
	tokenManager TokenManager
}

func NewAuthService(repo AuthPostgresRepo, tokenManager TokenManager) AuthService {
	return &authService{
		repo:         repo,
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
	return as.repo.CreateUser(user)
}

func (as *authService) SignIn(u UserSignIn) (Tokens, error) {
	user, errGet := as.repo.GetUser(u.Email)
	if errGet != nil {
		return Tokens{}, errGet
	}

	if !as.checkPassword(u.Password, user.PasswordHash) {
		return Tokens{}, errors.New("password is wrong")
	}

	return as.generateTokensPair(user)
}

func (as *authService) RefreshTokens(refreshToken string) (string, error) {
	hashRefreshToken := as.makeHash(refreshToken)
	user, errGet := as.repo.GetUserByRefreshToken(hashRefreshToken)
	if errGet != nil {
		return "", errGet
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

	errSave := as.repo.SaveRefreshToken(user.ID, hashRefreshToken, refreshExpiresAt)
	if errSave != nil {
		return Tokens{}, errSave
	}

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
