package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	mock_AuthService "github.com/jst-Frenzy/ControlSystem/AuthService/internal/AuthService/mocks"
	"net/http/httptest"
	"testing"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(s *mock_AuthService.MockAuthService, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_AuthService.MockAuthService, token string) {
				s.EXPECT().ParseToken(token).Return(1, "user", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1, "role":"user"}`,
		},
		{
			name:                 "No Header",
			headerName:           "",
			mockBehavior:         func(s *mock_AuthService.MockAuthService, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"Message":"empty auth header"}`,
		},
		{
			name:                 "Invalid Bearer",
			headerName:           "Authorization",
			headerValue:          "Beer token",
			mockBehavior:         func(s *mock_AuthService.MockAuthService, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"Message":"invalid auth header"}`,
		},
		{
			name:                 "Invalid Token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			mockBehavior:         func(s *mock_AuthService.MockAuthService, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"Message":"token is empty"}`,
		},
		{
			name:        "Service Failure",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_AuthService.MockAuthService, token string) {
				s.EXPECT().ParseToken(token).Return(0, "", errors.New("failed to parse token"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"Message":"failed to parse token"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			authServ := mock_AuthService.NewMockAuthService(c)
			testCase.mockBehavior(authServ, testCase.token)

			handler := NewAuthHandler(authServ)

			//Test Server
			r := gin.New()
			r.POST("/protected", handler.UserIdentity, func(ctx *gin.Context) {
				id, _ := ctx.Get(userIDCtx)
				role, _ := ctx.Get(userRoleCtx)
				ctx.String(200, fmt.Sprintf(`{"id":%d, "role":"%v"}`, id.(int), role))
			})

			//Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/protected", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			//Make Request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
