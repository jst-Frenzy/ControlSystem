package server

import (
	"context"
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/jst-Frenzy/ControlSystem/GoodsService/internal/GoodService"
	mock "github.com/jst-Frenzy/ControlSystem/GoodsService/internal/mocks"
	gen "github.com/jst-Frenzy/ControlSystem/protobuf/gen/goods"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestServer_GetItemQuantityAndPrice(t *testing.T) {
	type mockBehavior func(s *mock.MockGoodService, itemID string)

	testTable := []struct {
		name                     string
		inputItemIdRequest       *gen.ItemQuantityAndPriceRequest
		inputItemId              string
		mockBehavior             mockBehavior
		expectedItemInfoResponse *gen.ItemQuantityAndPriceResponse
		expectedError            error
	}{
		{
			name:               "OK",
			inputItemIdRequest: &gen.ItemQuantityAndPriceRequest{ItemId: "Object id"},
			inputItemId:        "Object id",
			mockBehavior: func(s *mock.MockGoodService, itemID string) {
				s.EXPECT().GetItemInfoForCart(itemID).Return(GoodService.ItemInfoForCart{
					Quantity: 6,
					Price:    15,
				}, nil)
			},
			expectedItemInfoResponse: &gen.ItemQuantityAndPriceResponse{
				Valid:    true,
				Quantity: "6",
				Price:    "15.00",
			},
			expectedError: nil,
		},
		{
			name:                     "Empty item id",
			inputItemIdRequest:       &gen.ItemQuantityAndPriceRequest{ItemId: ""},
			inputItemId:              "",
			mockBehavior:             func(s *mock.MockGoodService, itemID string) {},
			expectedItemInfoResponse: nil,
			expectedError:            status.Errorf(codes.InvalidArgument, "item id is required"),
		},
		{
			name:               "Fail get info",
			inputItemIdRequest: &gen.ItemQuantityAndPriceRequest{ItemId: "Object id"},
			inputItemId:        "Object id",
			mockBehavior: func(s *mock.MockGoodService, itemID string) {
				s.EXPECT().GetItemInfoForCart(itemID).Return(GoodService.ItemInfoForCart{}, errors.New("can't get info"))
			},
			expectedItemInfoResponse: &gen.ItemQuantityAndPriceResponse{},
			expectedError:            errors.New("can't get info"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			servMock := mock.NewMockGoodService(c)
			testCase.mockBehavior(servMock, testCase.inputItemId)

			gRPCServ := NewGRPCServer(Deps{
				GoodsService: servMock,
				Logger:       nil,
			})

			resp, err := gRPCServ.GetItemQuantityAndPrice(context.Background(), &gen.ItemQuantityAndPriceRequest{ItemId: testCase.inputItemId})

			assert.Equal(t, resp, testCase.expectedItemInfoResponse)
			assert.Equal(t, err, testCase.expectedError)
		})
	}
}
