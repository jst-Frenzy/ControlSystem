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
	type mockBehavior func(s *mock_AuthService.MockAuthPostgresRepo)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           AuthService.UserSignUp
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"user_name": "test", "email": "test@test.com", "password":"qwerty"}`,
			inputUser: AuthService.UserSignUp{
				UserName: "test",
				Email:    "test@test.com",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_AuthService.MockAuthPostgresRepo) {
				s.EXPECT().CreateUser(gomock.Any()).Return(1, nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:                "Empty fields",
			inputBody:           `{"user_name": "test", "password":"qwerty"}`,
			mockBehavior:        func(s *mock_AuthService.MockAuthPostgresRepo) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"Message":"invalid input body"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"user_name": "test", "email": "test@test.com", "password":"qwerty"}`,
			inputUser: AuthService.UserSignUp{
				UserName: "test",
				Email:    "test@test.com",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_AuthService.MockAuthPostgresRepo) {
				s.EXPECT().CreateUser(gomock.Any()).Return(0, errors.New("ServiceFailure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"Message":"ServiceFailure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			authRepo := mock_AuthService.NewMockAuthPostgresRepo(c)
			testCase.mockBehavior(authRepo)

			manager := mock_AuthService.NewMockTokenManager(c)

			services := AuthService.NewAuthService(authRepo, manager)
			handler := NewAuthHandler(services)

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
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
