package GoodService

import (
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	GoodService "github.com/jst-Frenzy/ControlSystem/GoodsService/internal/GoodService"
	mock "github.com/jst-Frenzy/ControlSystem/GoodsService/internal/mocks"
	"testing"
)

func TestService_addItem(t *testing.T) {
	type mockBehavior func(r *mock.MockGoodsMongoRepo, i GoodService.Item, user GoodService.UserCtx)

	testTable := []struct {
		name          string
		inputItem     GoodService.Item
		inputUser     GoodService.UserCtx
		mockBehavior  mockBehavior
		expectedId    string
		expectedError error
	}{
		{
			name: "OK, user exists",
			inputItem: GoodService.Item{
				Name:        "apple",
				Description: "tasty apple",
				Quantity:    1,
			},
			inputUser: GoodService.UserCtx{
				ID:   1,
				Name: "test name",
			},
			mockBehavior: func(r *mock.MockGoodsMongoRepo, i GoodService.Item, user GoodService.UserCtx) {
				r.EXPECT().GetSellerIDByUserID(user.ID).Return("100", nil)
				expectedItem := i
				expectedItem.SellerID = "100"
				r.EXPECT().CreateItem(expectedItem).Return("itemID", nil)
			},
			expectedId:    "itemID",
			expectedError: nil,
		},
		{
			name: "OK, user not exists",
			inputItem: GoodService.Item{
				Name:        "apple",
				Description: "tasty apple",
				Quantity:    1,
			},
			inputUser: GoodService.UserCtx{
				ID:   1,
				Name: "test name",
			},
			mockBehavior: func(r *mock.MockGoodsMongoRepo, i GoodService.Item, user GoodService.UserCtx) {
				r.EXPECT().GetSellerIDByUserID(user.ID).Return("", errors.New("seller not found"))
				r.EXPECT().CreateSeller(user.ID, user.Name).Return("100", nil)
				expectedItem := i
				expectedItem.SellerID = "100"
				r.EXPECT().CreateItem(expectedItem).Return("itemID", nil)
			},
			expectedId:    "itemID",
			expectedError: nil,
		},
		{
			name: "Error create seller",
			inputUser: GoodService.UserCtx{
				ID:   1,
				Name: "test name",
			},
			mockBehavior: func(r *mock.MockGoodsMongoRepo, i GoodService.Item, user GoodService.UserCtx) {
				r.EXPECT().GetSellerIDByUserID(user.ID).Return("", errors.New("seller not found"))
				r.EXPECT().CreateSeller(user.ID, user.Name).Return("", errors.New("error create seller"))
			},
			expectedId:    "",
			expectedError: errors.New("error create seller"),
		},
		{
			name: "Error looking seller",
			inputUser: GoodService.UserCtx{
				ID:   1,
				Name: "test name",
			},
			mockBehavior: func(r *mock.MockGoodsMongoRepo, i GoodService.Item, user GoodService.UserCtx) {
				r.EXPECT().GetSellerIDByUserID(user.ID).Return("", errors.New("random error"))
			},
			expectedId:    "",
			expectedError: errors.New("random error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mongoRep := mock.NewMockGoodsMongoRepo(c)
			testCase.mockBehavior(mongoRep, testCase.inputItem, testCase.inputUser)

			serv := GoodService.NewGoodService(mongoRep)
			id, err := serv.AddItem(testCase.inputItem, testCase.inputUser)

			assert.Equal(t, id, testCase.expectedId)
			assert.Equal(t, err, testCase.expectedError)
		})
	}
}

func TestService_deleteItem(t *testing.T) {
	type mockBehavior func(r *mock.MockGoodsMongoRepo, itemID string, userID int)

	testTable := []struct {
		name          string
		itemID        string
		userID        int
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name:   "OK",
			itemID: "itemID",
			userID: 1,
			mockBehavior: func(r *mock.MockGoodsMongoRepo, itemID string, userID int) {
				r.EXPECT().GetSellerIDByUserID(userID).Return("100", nil)
				r.EXPECT().DeleteItem(itemID, "100").Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "Error looking seller",
			userID: 1,
			mockBehavior: func(r *mock.MockGoodsMongoRepo, itemID string, userID int) {
				r.EXPECT().GetSellerIDByUserID(userID).Return("", errors.New("error looking seller"))
			},
			expectedError: errors.New("error looking seller"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mongoRep := mock.NewMockGoodsMongoRepo(c)
			testCase.mockBehavior(mongoRep, testCase.itemID, testCase.userID)

			serv := GoodService.NewGoodService(mongoRep)

			err := serv.DeleteItem(testCase.itemID, testCase.userID)

			assert.Equal(t, err, testCase.expectedError)
		})
	}
}

func TestService_updateItem(t *testing.T) {
	type mockBehavior func(r *mock.MockGoodsMongoRepo, i GoodService.Item, userID int)

	testTable := []struct {
		name          string
		inputItem     GoodService.Item
		inputUserID   int
		mockBehavior  mockBehavior
		expectedItem  GoodService.Item
		expectedError error
	}{
		{
			name: "OK",
			inputItem: GoodService.Item{
				ID:          "itemID",
				Name:        "apple",
				Description: "tasty apple",
				Quantity:    1,
			},
			inputUserID: 1,
			mockBehavior: func(r *mock.MockGoodsMongoRepo, i GoodService.Item, userID int) {
				r.EXPECT().GetSellerIDByUserID(userID).Return("100", nil)
				r.EXPECT().GetItemByID(i.ID).Return(GoodService.Item{
					ID:          "itemID",
					Name:        "apple",
					Description: "not tasty apple",
					Quantity:    2,
					SellerID:    "100",
				}, nil)
				expectedItem := i
				expectedItem.SellerID = "100"
				r.EXPECT().UpdateItem(expectedItem).Return(GoodService.Item{
					ID:          "itemID",
					Name:        "apple",
					Description: "tasty apple",
					Quantity:    1,
					SellerID:    "100",
				}, nil)
			},
			expectedItem: GoodService.Item{
				ID:          "itemID",
				Name:        "apple",
				Description: "tasty apple",
				Quantity:    1,
				SellerID:    "100",
			},
			expectedError: nil,
		},
		{
			name:        "Error looking user",
			inputUserID: 1,
			mockBehavior: func(r *mock.MockGoodsMongoRepo, i GoodService.Item, userID int) {
				r.EXPECT().GetSellerIDByUserID(userID).Return("", errors.New("error looking user"))
			},
			expectedError: errors.New("error looking user"),
		},
		{
			name: "Error get item by id",
			inputItem: GoodService.Item{
				ID:          "itemID",
				Name:        "apple",
				Description: "tasty apple",
				Quantity:    1,
			},
			inputUserID: 1,
			mockBehavior: func(r *mock.MockGoodsMongoRepo, i GoodService.Item, userID int) {
				r.EXPECT().GetSellerIDByUserID(userID).Return("100", nil)
				r.EXPECT().GetItemByID(i.ID).Return(GoodService.Item{}, errors.New("error get item by id"))

			},
			expectedItem:  GoodService.Item{},
			expectedError: errors.New("error get item by id"),
		},
		{
			name: "Not equal seller id",
			inputItem: GoodService.Item{
				ID:          "itemID",
				Name:        "apple",
				Description: "tasty apple",
				Quantity:    1,
			},
			inputUserID: 1,
			mockBehavior: func(r *mock.MockGoodsMongoRepo, i GoodService.Item, userID int) {
				r.EXPECT().GetSellerIDByUserID(userID).Return("100", nil)
				r.EXPECT().GetItemByID(i.ID).Return(GoodService.Item{
					ID:          "itemID",
					Name:        "apple",
					Description: "not tasty apple",
					Quantity:    2,
					SellerID:    "999",
				}, nil)
			},
			expectedItem:  GoodService.Item{},
			expectedError: errors.New("it's not your item"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mongoRep := mock.NewMockGoodsMongoRepo(c)
			testCase.mockBehavior(mongoRep, testCase.inputItem, testCase.inputUserID)

			serv := GoodService.NewGoodService(mongoRep)

			i, err := serv.UpdateItem(testCase.inputItem, testCase.inputUserID)

			assert.Equal(t, testCase.expectedItem, i)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestService_getGoods(t *testing.T) {
	type mockBehavior func(r *mock.MockGoodsMongoRepo)
	testTable := []struct {
		name          string
		mockBehavior  mockBehavior
		expectedSlice []GoodService.Item
		expectedError error
	}{
		{
			name: "OK",
			mockBehavior: func(r *mock.MockGoodsMongoRepo) {
				r.EXPECT().GetGoods().Return([]GoodService.Item{
					{
						ID:          "1",
						Name:        "test name 1",
						Description: "test description 1",
						Quantity:    1,
						SellerID:    "1",
					},
					{
						ID:          "2",
						Name:        "test name 2",
						Description: "test description 2",
						Quantity:    2,
						SellerID:    "2",
					},
					{
						ID:          "3",
						Name:        "test name 3",
						Description: "test description 3",
						Quantity:    3,
						SellerID:    "3",
					},
				}, nil)
			},
			expectedSlice: []GoodService.Item{
				{
					ID:          "1",
					Name:        "test name 1",
					Description: "test description 1",
					Quantity:    1,
					SellerID:    "1",
				},
				{
					ID:          "2",
					Name:        "test name 2",
					Description: "test description 2",
					Quantity:    2,
					SellerID:    "2",
				},
				{
					ID:          "3",
					Name:        "test name 3",
					Description: "test description 3",
					Quantity:    3,
					SellerID:    "3",
				},
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mongoRep := mock.NewMockGoodsMongoRepo(c)
			testCase.mockBehavior(mongoRep)

			serv := GoodService.NewGoodService(mongoRep)

			items, err := serv.GetGoods()

			assert.Equal(t, testCase.expectedSlice, items)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
