package handlers

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/jst-Frenzy/ControlSystem/GoodsService/internal/GoodService"
	mock "github.com/jst-Frenzy/ControlSystem/GoodsService/internal/mocks"
	"net/http/httptest"
	"testing"
)

func TestHandler_addItem(t *testing.T) {
	type mockBehavior func(s *mock.MockGoodService, i GoodService.Item, u GoodService.UserCtx)

	testTable := []struct {
		name                 string
		inputBody            string
		inputItem            GoodService.Item
		inputUserCtx         GoodService.UserCtx
		mockBehavior         mockBehavior
		userRole             string
		userName             string
		userID               int
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name": "testName", "description": "test description", "quantity": 1}`,
			inputItem: GoodService.Item{
				Name:        "testName",
				Description: "test description",
				Quantity:    1,
			},
			inputUserCtx: GoodService.UserCtx{
				ID:   1,
				Name: "testName",
			},
			mockBehavior: func(s *mock.MockGoodService, i GoodService.Item, u GoodService.UserCtx) {
				s.EXPECT().AddItem(i, u).Return("itemID", nil)
			},
			userRole:             "seller",
			userName:             "testName",
			userID:               1,
			expectedStatusCode:   201,
			expectedResponseBody: `{"id":"itemID"}`,
		},
		{
			name:                 "Incorrect Role",
			mockBehavior:         func(s *mock.MockGoodService, i GoodService.Item, u GoodService.UserCtx) {},
			userRole:             "user",
			expectedStatusCode:   400,
			expectedResponseBody: `{"Message":"not enough rights"}`,
		},
		{
			name:                 "Empty Fields",
			inputBody:            `{"name": "testName", "quantity": 1}`,
			mockBehavior:         func(s *mock.MockGoodService, i GoodService.Item, u GoodService.UserCtx) {},
			userRole:             "seller",
			expectedStatusCode:   400,
			expectedResponseBody: `{"Message":"invalid input body"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"name": "testName", "description": "test description", "quantity": 1}`,
			inputItem: GoodService.Item{
				Name:        "testName",
				Description: "test description",
				Quantity:    1,
			},
			inputUserCtx: GoodService.UserCtx{
				ID:   1,
				Name: "testName",
			},
			mockBehavior: func(s *mock.MockGoodService, i GoodService.Item, u GoodService.UserCtx) {
				s.EXPECT().AddItem(i, u).Return("", errors.New("service failure"))
			},
			userRole:             "seller",
			userName:             "testName",
			userID:               1,
			expectedStatusCode:   500,
			expectedResponseBody: `{"Message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			goodService := mock.NewMockGoodService(c)
			authClient := mock.NewMockAuthClient(c)
			testCase.mockBehavior(goodService, testCase.inputItem, testCase.inputUserCtx)

			handler := NewGoodsHandlers(goodService, authClient)

			r := gin.New()
			r.POST("/item", func(ctx *gin.Context) {
				ctx.Set("userID", testCase.userID)
				ctx.Set("userRole", testCase.userRole)
				ctx.Set("userName", testCase.userName)
				handler.AddItem(ctx)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/item", bytes.NewBufferString(testCase.inputBody))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getGoods(t *testing.T) {
	type mockBehavior func(s *mock.MockGoodService)
	testTable := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mock.MockGoodService) {
				s.EXPECT().GetGoods().Return([]GoodService.Item{
					{
						Name:        "apple",
						Description: "apple description",
					},
					{
						Name:        "orange",
						Description: "orange description",
					},
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"apple":"apple description","orange":"orange description"}`,
		},
		{
			name: "Server Error",
			mockBehavior: func(s *mock.MockGoodService) {
				s.EXPECT().GetGoods().Return([]GoodService.Item{}, errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"Message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			goodService := mock.NewMockGoodService(c)
			testCase.mockBehavior(goodService)

			authClient := mock.NewMockAuthClient(c)

			handler := NewGoodsHandlers(goodService, authClient)

			r := gin.New()
			r.GET("/catalog", handler.GetGoods)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/catalog", nil)
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_deleteItem(t *testing.T) {
	type mockBehavior func(s *mock.MockGoodService, itemID string, userID int)

	testTable := []struct {
		name                 string
		userRole             string
		userID               int
		userName             string
		itemID               string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "OK",
			userRole: "seller",
			userID:   1,
			userName: "testName",
			itemID:   "123",
			mockBehavior: func(s *mock.MockGoodService, itemID string, userID int) {
				s.EXPECT().DeleteItem(itemID, userID).Return(nil)
			},
			expectedStatusCode:   204,
			expectedResponseBody: ``,
		},
		{
			name:                 "Incorrect Role",
			userRole:             "user",
			itemID:               "123",
			mockBehavior:         func(s *mock.MockGoodService, itemID string, userID int) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"Message":"not enough rights"}`,
		},
		{
			name:     "Server error",
			userRole: "seller",
			userID:   1,
			userName: "testName",
			itemID:   "123",
			mockBehavior: func(s *mock.MockGoodService, itemID string, userID int) {
				s.EXPECT().DeleteItem(itemID, userID).Return(errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"Message":"service failure"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			goodService := mock.NewMockGoodService(c)
			testCase.mockBehavior(goodService, testCase.itemID, testCase.userID)

			authClient := mock.NewMockAuthClient(c)

			handler := NewGoodsHandlers(goodService, authClient)

			r := gin.New()
			r.DELETE("/item/:id", func(ctx *gin.Context) {
				ctx.Set("userID", testCase.userID)
				ctx.Set("userRole", testCase.userRole)
				ctx.Set("userName", testCase.userName)
				handler.DeleteItem(ctx)
			})

			url := "/item/" + testCase.itemID
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", url, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}

}

func TestHandler_updateItem(t *testing.T) {
	type mockBehavior func(s *mock.MockGoodService, i GoodService.Item, userID int)

	testTable := []struct {
		name                 string
		userRole             string
		userName             string
		userID               int
		inputItem            GoodService.Item
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "OK",
			userRole: "seller",
			userName: "testName",
			userID:   100,
			inputItem: GoodService.Item{
				ID:          "123",
				Name:        "apple",
				Description: "new description",
				Quantity:    10,
			},
			inputBody: `{"_id":"123","name":"apple","description":"new description","quantity": 10}`,
			mockBehavior: func(s *mock.MockGoodService, i GoodService.Item, userID int) {
				s.EXPECT().UpdateItem(i, userID).Return(GoodService.Item{
					ID:          "123",
					Name:        "apple",
					Description: "new description",
					Quantity:    10,
					SellerID:    "1",
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"_id":"123","name":"apple","description":"new description","quantity":10,"sellerID":"1"}`,
		},
		{
			name:                 "Incorrect role",
			userRole:             "user",
			mockBehavior:         func(s *mock.MockGoodService, i GoodService.Item, userID int) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"Message":"not enough rights"}`,
		},
		{
			name:                 "Empty Fields",
			userRole:             "seller",
			inputBody:            `{"_id":"123","name":"apple","quantity": 10}`,
			mockBehavior:         func(s *mock.MockGoodService, i GoodService.Item, userID int) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"Message":"invalid input body"}`,
		},
		{
			name:     "Server Error",
			userRole: "seller",
			inputItem: GoodService.Item{
				ID:          "123",
				Name:        "apple",
				Description: "new description",
				Quantity:    10,
			},
			inputBody: `{"_id":"123","name":"apple","description":"new description","quantity": 10}`,
			mockBehavior: func(s *mock.MockGoodService, i GoodService.Item, userID int) {
				s.EXPECT().UpdateItem(i, userID).Return(GoodService.Item{}, errors.New("server failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"Message":"server failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			goodService := mock.NewMockGoodService(c)
			testCase.mockBehavior(goodService, testCase.inputItem, testCase.userID)

			authClient := mock.NewMockAuthClient(c)

			handler := NewGoodsHandlers(goodService, authClient)

			r := gin.New()
			r.PUT("/item", func(ctx *gin.Context) {
				ctx.Set("userID", testCase.userID)
				ctx.Set("userRole", testCase.userRole)
				ctx.Set("userName", testCase.userName)
				handler.UpdateItem(ctx)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/item", bytes.NewBufferString(testCase.inputBody))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
