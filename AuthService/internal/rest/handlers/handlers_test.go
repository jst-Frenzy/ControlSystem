package handlers

import (
	"AuthService/internal/AuthService"
	mock_AuthService "AuthService/internal/AuthService/mocks"
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"net/http/httptest"
	"testing"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_AuthService.MockAuthService, user AuthService.UserSignUp)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            AuthService.UserSignUp
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"user_name": "test", "email": "test@test.com", "password":"qwerty"}`,
			inputUser: AuthService.UserSignUp{
				UserName: "test",
				Email:    "test@test.com",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_AuthService.MockAuthService, u AuthService.UserSignUp) {
				s.EXPECT().SignUp(u).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Empty fields",
			inputBody:            `{"user_name": "test", "password":"qwerty"}`,
			mockBehavior:         func(s *mock_AuthService.MockAuthService, u AuthService.UserSignUp) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"Message":"invalid input body"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"user_name": "test", "email": "test@test.com", "password":"qwerty"}`,
			inputUser: AuthService.UserSignUp{
				UserName: "test",
				Email:    "test@test.com",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_AuthService.MockAuthService, u AuthService.UserSignUp) {
				s.EXPECT().SignUp(u).Return(0, errors.New("ServiceFailure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"Message":"ServiceFailure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			authServ := mock_AuthService.NewMockAuthService(c)
			testCase.mockBehavior(authServ, testCase.inputUser)

			handler := NewAuthHandler(authServ)

			//Test Server
			r := gin.New()
			r.POST("/sign-up", handler.SignUp)

			//Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(testCase.inputBody))
			req.Header.Set("Content-Type", "application/json")

			//Perform Request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_SignIn(t *testing.T) {
	type mockBehavior func(s *mock_AuthService.MockAuthService, u AuthService.UserSignIn)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            AuthService.UserSignIn
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email":"test@test.com", "password":"qwerty"}`,
			inputUser: AuthService.UserSignIn{
				Email:    "test@test.com",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_AuthService.MockAuthService, u AuthService.UserSignIn) {
				s.EXPECT().SignIn(u).Return(AuthService.Tokens{
					AccessToken:  "abc",
					RefreshToken: "def",
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"access token":"abc","refresh token":"def"}`,
		},
		{
			name:                 "Empty Fields",
			inputBody:            `"email":"test@test.com", "password":"qwerty"}`,
			mockBehavior:         func(s *mock_AuthService.MockAuthService, u AuthService.UserSignIn) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"Message":"invalid input body"}`,
		},
		{
			name:      "OK",
			inputBody: `{"email":"test@test.com", "password":"qwerty"}`,
			inputUser: AuthService.UserSignIn{
				Email:    "test@test.com",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_AuthService.MockAuthService, u AuthService.UserSignIn) {
				s.EXPECT().SignIn(u).Return(AuthService.Tokens{}, errors.New("ServiceFailure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"Message":"ServiceFailure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			authServ := mock_AuthService.NewMockAuthService(c)
			testCase.mockBehavior(authServ, testCase.inputUser)

			handler := NewAuthHandler(authServ)

			r := gin.New()
			r.POST("/sign-in", handler.SignIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(testCase.inputBody))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
