package gRPC

import (
	"context"
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/jst-Frenzy/ControlSystem/AuthService/internal/AuthService"
	mock_AuthService "github.com/jst-Frenzy/ControlSystem/AuthService/internal/mocks"
	gen "github.com/jst-Frenzy/ControlSystem/protobuf/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestServer_ValidateToken(t *testing.T) {
	type mockBehavior func(s *mock_AuthService.MockAuthService, accessToken string)

	testTable := []struct {
		name                          string
		inputValidateTokenRequest     *gen.ValidateTokenRequest
		inputAccessToken              string
		mockBehavior                  mockBehavior
		expectedValidateTokenResponse *gen.ValidateTokenResponse
		expectedError                 error
	}{
		{
			name:                      "OK",
			inputValidateTokenRequest: &gen.ValidateTokenRequest{AccessToken: "access_token"},
			inputAccessToken:          "access_token",
			mockBehavior: func(s *mock_AuthService.MockAuthService, accessToken string) {
				s.EXPECT().ParseToken(accessToken).Return(AuthService.InfoFromToken{
					ID:       1,
					Role:     "user",
					UserName: "test name",
				}, nil)
			},
			expectedValidateTokenResponse: &gen.ValidateTokenResponse{
				Valid:    true,
				UserId:   "1",
				Role:     "user",
				UserName: "test name",
			},
			expectedError: nil,
		},
		{
			name:                          "EmptyAccessToken",
			inputValidateTokenRequest:     &gen.ValidateTokenRequest{AccessToken: ""},
			inputAccessToken:              "",
			mockBehavior:                  func(s *mock_AuthService.MockAuthService, accessToken string) {},
			expectedValidateTokenResponse: nil,
			expectedError:                 status.Error(codes.InvalidArgument, "access_token is required"),
		},
		{
			name:                      "Fail Parse Token",
			inputValidateTokenRequest: &gen.ValidateTokenRequest{AccessToken: "access_token"},
			inputAccessToken:          "access_token",
			mockBehavior: func(s *mock_AuthService.MockAuthService, accessToken string) {
				s.EXPECT().ParseToken(accessToken).Return(AuthService.InfoFromToken{}, errors.New("fail parse"))
			},
			expectedValidateTokenResponse: &gen.ValidateTokenResponse{
				Valid:  false,
				UserId: "",
				Role:   "",
			},
			expectedError: errors.New("can't parse token"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			servMock := mock_AuthService.NewMockAuthService(c)
			testCase.mockBehavior(servMock, testCase.inputAccessToken)

			gRPCServ := NewGRPCServer(Deps{
				Logger:      nil,
				AuthService: servMock,
			})

			resp, err := gRPCServ.ValidateToken(context.Background(), &gen.ValidateTokenRequest{AccessToken: testCase.inputAccessToken})

			assert.Equal(t, resp, testCase.expectedValidateTokenResponse)
			assert.Equal(t, err, testCase.expectedError)
		})
	}
}
