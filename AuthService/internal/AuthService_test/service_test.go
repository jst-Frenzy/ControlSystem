package AuthService_test

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/jst-Frenzy/ControlSystem/AuthService/internal/AuthService"
	mockauthservice "github.com/jst-Frenzy/ControlSystem/AuthService/internal/AuthService/mocks"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestService_signUp(t *testing.T) {
	type mockBehavior func(r *mockauthservice.MockAuthPostgresRepo, redisMock *mockauthservice.MockAuthRedisRepo, u AuthService.User)

	testTable := []struct {
		name          string
		inputUser     AuthService.UserSignUp
		repoUser      AuthService.User
		mockBehavior  mockBehavior
		expectedID    int
		expectedError error
	}{
		{
			name: "OK",
			inputUser: AuthService.UserSignUp{
				UserName: "test",
				Email:    "test@test.com",
				Password: "qwerty",
			},
			repoUser: AuthService.User{
				UserName:     "test",
				Email:        "test@test.com",
				PasswordHash: "",
			},
			mockBehavior: func(r *mockauthservice.MockAuthPostgresRepo, redisMock *mockauthservice.MockAuthRedisRepo, u AuthService.User) {
				expectedUser := AuthService.User{ID: 1, UserName: "test", Email: "test@test.com"}
				r.EXPECT().CreateUser(gomock.Any()).Return(expectedUser, nil)
				redisMock.EXPECT().AddUserWithEmail(gomock.Any()).Return(nil).AnyTimes()
			},
			expectedID:    1,
			expectedError: nil,
		},
		{
			name: "DB Error",
			inputUser: AuthService.UserSignUp{
				UserName: "test",
				Email:    "test@test.com",
				Password: "qwerty",
			},
			repoUser: AuthService.User{
				UserName:     "test",
				Email:        "test@test.com",
				PasswordHash: "",
			},
			mockBehavior: func(r *mockauthservice.MockAuthPostgresRepo, redisMock *mockauthservice.MockAuthRedisRepo, u AuthService.User) {
				r.EXPECT().CreateUser(gomock.Any()).Return(AuthService.User{}, errors.New("DB Failure"))
			},
			expectedID:    0,
			expectedError: errors.New("DB Failure"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			authPostgresRepo := mockauthservice.NewMockAuthPostgresRepo(c)

			tokenManager := mockauthservice.NewMockTokenManager(c)

			authRedisRepo := mockauthservice.NewMockAuthRedisRepo(c)

			testCase.mockBehavior(authPostgresRepo, authRedisRepo, testCase.repoUser)

			service := AuthService.NewAuthService(authPostgresRepo, authRedisRepo, tokenManager)

			id, err := service.SignUp(testCase.inputUser)

			assert.Equal(t, id, testCase.expectedID)
			assert.Equal(t, err, testCase.expectedError)
		})
	}
}

func TestService_signIn(t *testing.T) {
	type repoBehavior func(r *mockauthservice.MockAuthPostgresRepo, redisMock *mockauthservice.MockAuthRedisRepo, email string)
	type tokenBehavior func(t *mockauthservice.MockTokenManager)

	testTable := []struct {
		name           string
		inputUser      AuthService.UserSignIn
		inputEmail     string
		repoBehavior   repoBehavior
		tokenBehavior  tokenBehavior
		expectedTokens AuthService.Tokens
		expectedError  error
	}{
		{
			name: "OK",
			inputUser: AuthService.UserSignIn{
				Email:    "test@test.com",
				Password: "qwerty",
			},
			inputEmail: "test@test.com",

			repoBehavior: func(r *mockauthservice.MockAuthPostgresRepo, redisMock *mockauthservice.MockAuthRedisRepo, email string) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("qwerty"), bcrypt.DefaultCost)

				redisMock.EXPECT().GetUserWithEmail(email).Return(AuthService.User{}, errors.New("cache miss"))
				redisMock.EXPECT().AddUserWithRefreshToken(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				r.EXPECT().GetUser(email).Return(AuthService.User{
					ID:           1,
					UserName:     "test",
					Email:        "test@test.com",
					PasswordHash: string(hashedPassword),
				}, nil)

				redisMock.EXPECT().AddUserWithEmail(gomock.Any()).Return(nil).AnyTimes()

				r.EXPECT().SaveRefreshToken(gomock.Any()).Return(nil)
			},
			tokenBehavior: func(t *mockauthservice.MockTokenManager) {
				t.EXPECT().NewJWT(gomock.Any(), gomock.Any()).Return("access_token", nil)
				t.EXPECT().NewRefreshToken().Return("refresh_token", nil)
			},
			expectedTokens: AuthService.Tokens{
				AccessToken:  "access_token",
				RefreshToken: "refresh_token",
			},
			expectedError: nil,
		},
		{
			name: "Fail Get User From Repo",
			inputUser: AuthService.UserSignIn{
				Email:    "test@test.com",
				Password: "qwerty",
			},
			inputEmail: "test@test.com",
			repoBehavior: func(r *mockauthservice.MockAuthPostgresRepo, redisMock *mockauthservice.MockAuthRedisRepo, email string) {
				redisMock.EXPECT().GetUserWithEmail(email).Return(AuthService.User{}, errors.New("cache miss"))
				r.EXPECT().GetUser(email).Return(AuthService.User{}, errors.New("fail Get User"))
			},
			tokenBehavior:  func(t *mockauthservice.MockTokenManager) {},
			expectedTokens: AuthService.Tokens{},
			expectedError:  errors.New("fail Get User"),
		},
		{
			name: "OK with cache hit",
			inputUser: AuthService.UserSignIn{
				Email:    "test@test.com",
				Password: "qwerty",
			},
			inputEmail: "test@test.com",
			repoBehavior: func(r *mockauthservice.MockAuthPostgresRepo, redisMock *mockauthservice.MockAuthRedisRepo, email string) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("qwerty"), bcrypt.DefaultCost)

				redisMock.EXPECT().GetUserWithEmail(email).Return(AuthService.User{
					ID:           1,
					UserName:     "test",
					Email:        "test@test.com",
					PasswordHash: string(hashedPassword),
				}, nil)

				r.EXPECT().SaveRefreshToken(gomock.Any()).Return(nil).AnyTimes()

				redisMock.EXPECT().AddUserWithRefreshToken(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			tokenBehavior: func(t *mockauthservice.MockTokenManager) {
				t.EXPECT().NewJWT(gomock.Any(), gomock.Any()).Return("access_token", nil)
				t.EXPECT().NewRefreshToken().Return("refresh_token", nil)
			},
			expectedTokens: AuthService.Tokens{
				AccessToken:  "access_token",
				RefreshToken: "refresh_token",
			},
			expectedError: nil,
		},
		{
			name: "Fail Generate JWT Tokens",
			inputUser: AuthService.UserSignIn{
				Email:    "test@test.com",
				Password: "qwerty",
			},
			inputEmail: "test@test.com",

			repoBehavior: func(r *mockauthservice.MockAuthPostgresRepo, redisMock *mockauthservice.MockAuthRedisRepo, email string) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("qwerty"), bcrypt.DefaultCost)

				redisMock.EXPECT().GetUserWithEmail(email).Return(AuthService.User{}, errors.New("cache miss"))

				r.EXPECT().GetUser(email).Return(AuthService.User{
					ID:           1,
					UserName:     "test",
					Email:        "test@test.com",
					PasswordHash: string(hashedPassword),
				}, nil)

				redisMock.EXPECT().AddUserWithEmail(gomock.Any()).Return(nil)
			},
			tokenBehavior: func(t *mockauthservice.MockTokenManager) {
				t.EXPECT().NewJWT(gomock.Any(), gomock.Any()).Return("", errors.New("error generate JWT"))
			},
			expectedTokens: AuthService.Tokens{
				AccessToken:  "",
				RefreshToken: "",
			},
			expectedError: errors.New("error generate JWT"),
		},
		{
			name: "Fail Generate Refresh Tokens",
			inputUser: AuthService.UserSignIn{
				Email:    "test@test.com",
				Password: "qwerty",
			},
			inputEmail: "test@test.com",

			repoBehavior: func(r *mockauthservice.MockAuthPostgresRepo, redisMock *mockauthservice.MockAuthRedisRepo, email string) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("qwerty"), bcrypt.DefaultCost)

				redisMock.EXPECT().GetUserWithEmail(email).Return(AuthService.User{}, errors.New("cache miss"))

				r.EXPECT().GetUser(email).Return(AuthService.User{
					ID:           1,
					UserName:     "test",
					Email:        "test@test.com",
					PasswordHash: string(hashedPassword),
				}, nil)

				redisMock.EXPECT().AddUserWithEmail(gomock.Any()).Return(nil)
			},
			tokenBehavior: func(t *mockauthservice.MockTokenManager) {
				t.EXPECT().NewJWT(gomock.Any(), gomock.Any()).Return("access_token", nil)
				t.EXPECT().NewRefreshToken().Return("", errors.New("error generate Refresh"))
			},
			expectedTokens: AuthService.Tokens{
				AccessToken:  "",
				RefreshToken: "",
			},
			expectedError: errors.New("error generate Refresh"),
		},
		{
			name: "Fail Save Refresh Token",
			inputUser: AuthService.UserSignIn{
				Email:    "test@test.com",
				Password: "qwerty",
			},
			inputEmail: "test@test.com",

			repoBehavior: func(r *mockauthservice.MockAuthPostgresRepo, redisMock *mockauthservice.MockAuthRedisRepo, email string) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("qwerty"), bcrypt.DefaultCost)

				redisMock.EXPECT().GetUserWithEmail(email).Return(AuthService.User{}, errors.New("cache miss"))

				r.EXPECT().GetUser(email).Return(AuthService.User{
					ID:           1,
					UserName:     "test",
					Email:        "test@test.com",
					PasswordHash: string(hashedPassword),
				}, nil)

				redisMock.EXPECT().AddUserWithEmail(gomock.Any()).Return(errors.New("error save user")).AnyTimes()

				r.EXPECT().SaveRefreshToken(gomock.Any()).Return(errors.New("error save refresh token"))
			},
			tokenBehavior: func(t *mockauthservice.MockTokenManager) {
				t.EXPECT().NewJWT(gomock.Any(), gomock.Any()).Return("access_token", nil)
				t.EXPECT().NewRefreshToken().Return("refresh_token", nil)
			},
			expectedTokens: AuthService.Tokens{
				AccessToken:  "",
				RefreshToken: "",
			},
			expectedError: errors.New("error save refresh token"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			authPostgresRepo := mockauthservice.NewMockAuthPostgresRepo(c)

			tokenManager := mockauthservice.NewMockTokenManager(c)
			testCase.tokenBehavior(tokenManager)

			authRedisRepo := mockauthservice.NewMockAuthRedisRepo(c)

			testCase.repoBehavior(authPostgresRepo, authRedisRepo, testCase.inputEmail)

			service := AuthService.NewAuthService(authPostgresRepo, authRedisRepo, tokenManager)

			tokens, err := service.SignIn(testCase.inputUser)

			assert.Equal(t, err, testCase.expectedError)
			assert.Equal(t, tokens, testCase.expectedTokens)
		})
	}
}

func TestService_RefreshTokens(t *testing.T) {
	type repoBehavior func(r *mockauthservice.MockAuthPostgresRepo, redisMock *mockauthservice.MockAuthRedisRepo, refreshToken string)
	type tokenBehavior func(t *mockauthservice.MockTokenManager)

	testTable := []struct {
		name              string
		inputRefreshToken string
		repoBehavior      repoBehavior
		tokenBehavior     tokenBehavior
		expectedToken     string
		expectedError     error
	}{
		{
			name:              "OK",
			inputRefreshToken: "refresh_token",
			repoBehavior: func(r *mockauthservice.MockAuthPostgresRepo, redisMock *mockauthservice.MockAuthRedisRepo, refreshToken string) {
				redisMock.EXPECT().GetUserWithRefreshToken(refreshToken).Return(AuthService.User{}, errors.New("cache miss"))

				r.EXPECT().GetUserByRefreshToken(refreshToken).Return(AuthService.User{
					ID:       1,
					UserName: "test",
					Email:    "test@test.com",
					Role:     "user",
				}, nil)

				redisMock.EXPECT().AddUserWithRefreshToken(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
			tokenBehavior: func(t *mockauthservice.MockTokenManager) {
				t.EXPECT().NewJWT(gomock.Any(), gomock.Any()).Return("access_token", nil)
			},
			expectedToken: "access_token",
			expectedError: nil,
		},
		{
			name:              "OK with cache hit",
			inputRefreshToken: "refresh_token",
			repoBehavior: func(r *mockauthservice.MockAuthPostgresRepo, redisMock *mockauthservice.MockAuthRedisRepo, refreshTokenHash string) {
				redisMock.EXPECT().GetUserWithRefreshToken(refreshTokenHash).
					Return(AuthService.User{
						ID:       1,
						UserName: "test",
						Email:    "test@test.com",
						Role:     "user",
					}, nil)

			},
			tokenBehavior: func(t *mockauthservice.MockTokenManager) {
				t.EXPECT().NewJWT(gomock.Any(), gomock.Any()).Return("access_token", nil)
			},
			expectedToken: "access_token",
			expectedError: nil,
		},
		{
			name:              "Fail Get User From Repo",
			inputRefreshToken: "refresh_token",
			repoBehavior: func(r *mockauthservice.MockAuthPostgresRepo, redisMock *mockauthservice.MockAuthRedisRepo, refreshToken string) {
				redisMock.EXPECT().GetUserWithRefreshToken(refreshToken).Return(AuthService.User{}, errors.New("cache miss"))
				r.EXPECT().GetUserByRefreshToken(refreshToken).Return(AuthService.User{}, errors.New("error get user"))
			},
			tokenBehavior: func(t *mockauthservice.MockTokenManager) {},
			expectedToken: "",
			expectedError: errors.New("error get user"),
		},
		{
			name:              "Fail Create New JWT",
			inputRefreshToken: "refresh_token",
			repoBehavior: func(r *mockauthservice.MockAuthPostgresRepo, redisMock *mockauthservice.MockAuthRedisRepo, refreshToken string) {
				redisMock.EXPECT().GetUserWithRefreshToken(refreshToken).Return(AuthService.User{}, errors.New("cache miss"))

				r.EXPECT().GetUserByRefreshToken(refreshToken).Return(AuthService.User{
					ID:       1,
					UserName: "test",
					Email:    "test@test.com",
					Role:     "user",
				}, nil)
			},
			tokenBehavior: func(t *mockauthservice.MockTokenManager) {
				t.EXPECT().NewJWT(gomock.Any(), gomock.Any()).Return("", errors.New("error create new jwt"))
			},
			expectedToken: "",
			expectedError: errors.New("error create new jwt"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			authPostgresRepo := mockauthservice.NewMockAuthPostgresRepo(c)
			hash := sha256.Sum256([]byte(testCase.inputRefreshToken))

			tokenManager := mockauthservice.NewMockTokenManager(c)
			testCase.tokenBehavior(tokenManager)

			authRedisRepo := mockauthservice.NewMockAuthRedisRepo(c)
			testCase.repoBehavior(authPostgresRepo, authRedisRepo, hex.EncodeToString(hash[:]))

			service := AuthService.NewAuthService(authPostgresRepo, authRedisRepo, tokenManager)

			token, err := service.RefreshTokens(testCase.inputRefreshToken)

			assert.Equal(t, token, testCase.expectedToken)
			assert.Equal(t, err, testCase.expectedError)
		})
	}
}
