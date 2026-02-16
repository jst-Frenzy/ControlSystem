package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	mock "github.com/jst-Frenzy/ControlSystem/GoodsService/internal/mocks"
	proto "github.com/jst-Frenzy/ControlSystem/protobuf/gen"
	"net/http/httptest"
	"testing"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(authClient *mock.MockAuthClient, token string)

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
			mockBehavior: func(authClient *mock.MockAuthClient, token string) {
				authClient.EXPECT().ValidateToken(gomock.Any(), token).Return(&proto.ValidateTokenResponse{
					Valid:    true,
					UserId:   "1",
					Role:     "user",
					UserName: "testName",
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"role":"user","name":"testName"}`,
		},
		{
			name:                 "No header",
			mockBehavior:         func(authClient *mock.MockAuthClient, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"Message":"empty auth header"}`,
		},
		{
			name:                 "Invalid header",
			headerName:           "Authorization",
			headerValue:          "Bearer token and some more parts",
			mockBehavior:         func(authClient *mock.MockAuthClient, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"Message":"invalid auth header"}`,
		},
		{
			name:                 "Invalid Bearer word",
			headerName:           "Authorization",
			headerValue:          "Berr token",
			token:                "token",
			mockBehavior:         func(authClient *mock.MockAuthClient, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"Message":"invalid auth header"}`,
		},
		{
			name:                 "Empty token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "",
			mockBehavior:         func(authClient *mock.MockAuthClient, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"Message":"token is empty"}`,
		},
		{
			name:        "Auth client error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(authClient *mock.MockAuthClient, token string) {
				authClient.EXPECT().ValidateToken(gomock.Any(), token).Return(&proto.ValidateTokenResponse{},
					errors.New("auth client failure"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"Message":"auth client failure"}`,
		},
		{
			name:        "Auth client error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(authClient *mock.MockAuthClient, token string) {
				authClient.EXPECT().ValidateToken(gomock.Any(), token).Return(&proto.ValidateTokenResponse{
					Valid:    false,
					UserId:   "",
					Role:     "",
					UserName: "",
				}, nil)
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"Message":"invalid token"}`,
		},
		{
			name:        "Auth client error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(authClient *mock.MockAuthClient, token string) {
				authClient.EXPECT().ValidateToken(gomock.Any(), token).Return(&proto.ValidateTokenResponse{
					Valid:    true,
					UserId:   "abc",
					Role:     "",
					UserName: "",
				}, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"Message":"invalid user id"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			authClient := mock.NewMockAuthClient(c)
			testCase.mockBehavior(authClient, testCase.token)

			goodService := mock.NewMockGoodService(c)

			handler := NewGoodsHandlers(goodService, authClient)

			r := gin.New()
			r.POST("/protected", handler.UserIdentity, func(ctx *gin.Context) {
				id, _ := ctx.Get("userID")
				role, _ := ctx.Get("userRole")
				name, _ := ctx.Get("userName")
				ctx.String(200, fmt.Sprintf(`{"id":%d,"role":"%v","name":"%v"}`, id, role, name))
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/protected", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
