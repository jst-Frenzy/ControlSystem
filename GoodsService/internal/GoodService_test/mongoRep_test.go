package GoodService_test

import (
	"github.com/jst-Frenzy/ControlSystem/GoodsService/internal/GoodService"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func TestMongoRep_createItem(t *testing.T) {
	type mockBehavior func(m *mtest.T)

	testTable := []struct {
		name         string
		inputItem    GoodService.Item
		mockBehavior mockBehavior
		expectedID   string
		wantErr      bool
	}{
		{
			name: "OK",
			inputItem: GoodService.Item{
				ID:          "1",
				Name:        "apple",
				Description: "tasty apple",
				Quantity:    1,
				SellerID:    "100",
			},
			mockBehavior: func(m *mtest.T) {
				response := bson.D{
					{Key: "ok", Value: 1},
					{Key: "n", Value: 1},
					{Key: "insertedId", Value: "1"},
				}
				m.AddMockResponses(response)
			},
			expectedID: "1",
			wantErr:    false,
		},
		{
			name: "Error insert",
			inputItem: GoodService.Item{
				ID:          "1",
				Name:        "apple",
				Description: "tasty apple",
				Quantity:    1,
				SellerID:    "100",
			},
			mockBehavior: func(m *mtest.T) {
				response := bson.D{
					{Key: "ok", Value: 0},
					{Key: "n", Value: 0},
				}
				m.AddMockResponses(response)
			},
			expectedID: "",
			wantErr:    true,
		},
	}

	for _, testCase := range testTable {
		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

		mt.Run(testCase.name, func(mt *mtest.T) {

			mongoRep := GoodService.NewGoodsMongoRepo(mt.Client)

			testCase.mockBehavior(mt)

			res, err := mongoRep.CreateItem(testCase.inputItem)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.expectedID, res)
		})
	}
}

func TestMongoRep_deleteItem(t *testing.T) {
	type mockBehavior func(m *mtest.T)

	testTable := []struct {
		name          string
		inputItemID   string
		inputSellerID string
		mockBehavior  mockBehavior
		wantErr       bool
	}{
		{
			name:          "OK",
			inputItemID:   "507f1f77bcf86cd799439011",
			inputSellerID: "100",
			mockBehavior: func(m *mtest.T) {
				response := bson.D{
					{Key: "ok", Value: 1},
					{Key: "n", Value: 1},
				}
				m.AddMockResponses(response)
			},
			wantErr: false,
		},
		{
			name:          "Cant parse itemId",
			inputItemID:   "1",
			inputSellerID: "100",
			mockBehavior:  func(m *mtest.T) {},
			wantErr:       true,
		},
		{
			name:          "DB error",
			inputItemID:   "507f1f77bcf86cd799439011",
			inputSellerID: "100",
			mockBehavior: func(m *mtest.T) {
				response := bson.D{
					{Key: "ok", Value: 0},
					{Key: "n", Value: 0},
				}
				m.AddMockResponses(response)
			},
			wantErr: true,
		},
		{
			name:          "Item not in DB",
			inputItemID:   "507f1f77bcf86cd799439011",
			inputSellerID: "100",
			mockBehavior: func(m *mtest.T) {
				response := bson.D{
					{Key: "ok", Value: 1},
					{Key: "n", Value: 0},
				}
				m.AddMockResponses(response)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

		mt.Run(testCase.name, func(mt *mtest.T) {
			mongoRep := GoodService.NewGoodsMongoRepo(mt.Client)

			testCase.mockBehavior(mt)

			err := mongoRep.DeleteItem(testCase.inputItemID, testCase.inputSellerID)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMongoRep_updateItem(t *testing.T) {
	type mockBehavior func(m *mtest.T)

	testTable := []struct {
		name         string
		inputItem    GoodService.Item
		mockBehavior mockBehavior
		expectedItem GoodService.Item
		wantErr      bool
	}{
		{
			name: "OK",
			inputItem: GoodService.Item{
				ID:          "507f1f77bcf86cd799439011",
				Name:        "apple",
				Description: "tasty apple",
				Quantity:    1,
				SellerID:    "100",
			},
			mockBehavior: func(m *mtest.T) {
				objectID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
				response := bson.D{
					{Key: "ok", Value: 1},
					{Key: "value", Value: bson.D{
						{Key: "_id", Value: objectID},
						{Key: "name", Value: "apple"},
						{Key: "description", Value: "tasty apple"},
						{Key: "quantity", Value: 1},
						{Key: "seller_id", Value: "100"},
					}},
				}
				m.AddMockResponses(response)
			},
			expectedItem: GoodService.Item{
				ID:          "507f1f77bcf86cd799439011",
				Name:        "apple",
				Description: "tasty apple",
				Quantity:    1,
				SellerID:    "100",
			},
			wantErr: false,
		},
		{
			name: "Cant parse itemId",
			inputItem: GoodService.Item{
				ID:          "1",
				Name:        "apple",
				Description: "tasty apple",
				Quantity:    1,
				SellerID:    "100",
			},
			mockBehavior: func(m *mtest.T) {},
			expectedItem: GoodService.Item{},
			wantErr:      true,
		},
		{
			name: "Item not found",
			inputItem: GoodService.Item{
				ID:          "507f1f77bcf86cd799439011",
				Name:        "apple",
				Description: "tasty apple",
				Quantity:    1,
				SellerID:    "100",
			},
			mockBehavior: func(m *mtest.T) {
				response := bson.D{
					{Key: "ok", Value: 1},
					{Key: "value", Value: nil},
				}
				m.AddMockResponses(response)
			},
			expectedItem: GoodService.Item{},
			wantErr:      true,
		},
		{
			name: "Decode error",
			inputItem: GoodService.Item{
				ID:          "507f1f77bcf86cd799439011",
				Name:        "apple",
				Description: "tasty apple",
				Quantity:    1,
				SellerID:    "100",
			},
			mockBehavior: func(m *mtest.T) {
				m.AddMockResponses(bson.D{
					{Key: "ok", Value: 1},
					{Key: "value", Value: nil},
					{Key: "errmsg", Value: "decode error"},
				})
			},
			expectedItem: GoodService.Item{},
			wantErr:      true,
		},
	}

	for _, testCase := range testTable {
		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		mt.Run(testCase.name, func(mt *mtest.T) {
			mongoRep := GoodService.NewGoodsMongoRepo(mt.Client)

			testCase.mockBehavior(mt)

			newItem, err := mongoRep.UpdateItem(testCase.inputItem)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.expectedItem, newItem)
		})
	}
}

func TestMongoRep_getGoods(t *testing.T) {
	type mockBehavior func(m *mtest.T)

	testTable := []struct {
		name           string
		mockBehavior   mockBehavior
		expectedAnswer []GoodService.Item
		wantErr        bool
	}{
		{
			name: "OK",
			mockBehavior: func(m *mtest.T) {
				objectID1, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
				objectID2, _ := primitive.ObjectIDFromHex("607f1f77bcf86cd799439012")
				response := mtest.CreateCursorResponse(2, "test.items", mtest.FirstBatch,
					bson.D{
						{Key: "_id", Value: objectID1},
						{Key: "name", Value: "apple"},
						{Key: "description", Value: "tasty apple"},
						{Key: "quantity", Value: 1},
						{Key: "seller_id", Value: "100"},
					},
					bson.D{
						{Key: "_id", Value: objectID2},
						{Key: "name", Value: "orange"},
						{Key: "description", Value: "tasty orange"},
						{Key: "quantity", Value: 10},
						{Key: "seller_id", Value: "250"},
					})
				m.AddMockResponses(response)

				secondResponse := mtest.CreateCursorResponse(
					0, "test.items", mtest.NextBatch,
				)
				m.AddMockResponses(secondResponse)
			},
			expectedAnswer: []GoodService.Item{
				{
					ID:          "507f1f77bcf86cd799439011",
					Name:        "apple",
					Description: "tasty apple",
					Quantity:    1,
					SellerID:    "100",
				},
				{
					ID:          "607f1f77bcf86cd799439012",
					Name:        "orange",
					Description: "tasty orange",
					Quantity:    10,
					SellerID:    "250",
				},
			},
			wantErr: false,
		},
		{
			name: "Error find",
			mockBehavior: func(m *mtest.T) {
				m.AddMockResponses(bson.D{
					{Key: "ok", Value: 1},
					{Key: "value", Value: nil},
					{Key: "errmsg", Value: "mongo error"},
				})
			},
			expectedAnswer: nil,
			wantErr:        true,
		},
		{
			name: "Error decode",
			mockBehavior: func(m *mtest.T) {
				objectID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
				response := mtest.CreateCursorResponse(
					1, "test.items", mtest.FirstBatch,
					bson.D{
						{Key: "_id", Value: objectID},
						{Key: "name", Value: 12345},
						{Key: "description", Value: true},
						{Key: "quantity", Value: "120"},
						{Key: "seller_id", Value: 999},
					},
				)
				m.AddMockResponses(response)
			},
			expectedAnswer: nil,
			wantErr:        true,
		},
	}

	for _, testCase := range testTable {
		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		mt.Run(testCase.name, func(mt *mtest.T) {
			mongoRep := GoodService.NewGoodsMongoRepo(mt.Client)

			testCase.mockBehavior(mt)

			items, err := mongoRep.GetGoods()

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.expectedAnswer, items)
		})
	}
}

func TestMongoRep_getQuantity(t *testing.T) {
	type mockBehavior func(m *mtest.T)

	testTable := []struct {
		name          string
		inputItemID   string
		mockBehavior  mockBehavior
		expectedCount int
		wantErr       bool
	}{
		{
			name:        "OK",
			inputItemID: "507f1f77bcf86cd799439011",
			mockBehavior: func(m *mtest.T) {
				objectID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
				response := mtest.CreateCursorResponse(
					1, "test.items", mtest.FirstBatch,
					bson.D{
						{Key: "_id", Value: objectID},
						{Key: "name", Value: "apple"},
						{Key: "description", Value: "tasty apple"},
						{Key: "quantity", Value: 12},
						{Key: "seller_id", Value: "100"},
					},
				)
				m.AddMockResponses(response)
			},
			expectedCount: 12,
			wantErr:       false,
		},
		{
			name:          "Cant parse itemId",
			inputItemID:   "1",
			mockBehavior:  func(m *mtest.T) {},
			expectedCount: 0,
			wantErr:       true,
		},
		{
			name:        "Item not found",
			inputItemID: "507f1f77bcf86cd799439011",
			mockBehavior: func(m *mtest.T) {
				response := bson.D{
					{Key: "ok", Value: 1},
					{Key: "value", Value: nil},
				}
				m.AddMockResponses(response)
			},
			expectedCount: 0,
			wantErr:       true,
		},
		{
			name:        "Decode error",
			inputItemID: "507f1f77bcf86cd799439011",
			mockBehavior: func(m *mtest.T) {
				response := bson.D{
					{Key: "ok", Value: 1},
					{Key: "value", Value: nil},
					{Key: "errmsg", Value: "decode error"},
				}
				m.AddMockResponses(response)
			},
			expectedCount: 0,
			wantErr:       true,
		},
	}

	for _, testCase := range testTable {
		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		mt.Run(testCase.name, func(mt *mtest.T) {
			mongoRep := GoodService.NewGoodsMongoRepo(mt.Client)

			testCase.mockBehavior(mt)

			cnt, err := mongoRep.GetQuantity(testCase.inputItemID)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.expectedCount, cnt)
		})
	}
}

func TestMongoRep_getSellerIDByUserID(t *testing.T) {
	type mockBehavior func(m *mtest.T)

	testTable := []struct {
		name             string
		inputUserID      int
		mockBehavior     mockBehavior
		expectedSellerID string
		wantErr          bool
	}{
		{
			name:        "OK",
			inputUserID: 20,
			mockBehavior: func(m *mtest.T) {
				response := mtest.CreateCursorResponse(
					1, "test.items", mtest.FirstBatch,
					bson.D{
						{Key: "_id", Value: "15"},
						{Key: "user_id", Value: 20},
						{Key: "name", Value: "test name"},
					},
				)
				m.AddMockResponses(response)
			},
			expectedSellerID: "15",
			wantErr:          false,
		},
		{
			name:        "Item not found",
			inputUserID: 20,
			mockBehavior: func(m *mtest.T) {
				response := bson.D{
					{Key: "ok", Value: 1},
					{Key: "value", Value: nil},
				}
				m.AddMockResponses(response)
			},
			expectedSellerID: "",
			wantErr:          true,
		},
		{
			name:        "Decode error",
			inputUserID: 20,
			mockBehavior: func(m *mtest.T) {
				response := bson.D{
					{Key: "ok", Value: 1},
					{Key: "value", Value: nil},
					{Key: "errmsg", Value: "decode error"},
				}
				m.AddMockResponses(response)
			},
			expectedSellerID: "",
			wantErr:          true,
		},
	}

	for _, testCase := range testTable {
		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		mt.Run(testCase.name, func(mt *mtest.T) {
			mongoRep := GoodService.NewGoodsMongoRepo(mt.Client)

			testCase.mockBehavior(mt)

			sellerID, err := mongoRep.GetSellerIDByUserID(testCase.inputUserID)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.expectedSellerID, sellerID)
		})
	}
}

func TestMongoRep_createSeller(t *testing.T) {
	type mockBehavior func(m *mtest.T)

	testTable := []struct {
		name          string
		inputUserID   int
		inputUserName string
		mockBehavior  mockBehavior
		expectedLenID int
		wantErr       bool
	}{
		{
			name:          "OK",
			inputUserID:   15,
			inputUserName: "test name",
			mockBehavior: func(m *mtest.T) {
				response := bson.D{
					{Key: "ok", Value: 1},
					{Key: "n", Value: 1},
					{Key: "insertedId", Value: primitive.NewObjectID()},
				}
				m.AddMockResponses(response)
			},
			wantErr:       false,
			expectedLenID: 24,
		},
		{
			name:          "Error insert",
			inputUserID:   15,
			inputUserName: "test name",
			mockBehavior: func(m *mtest.T) {
				response := bson.D{
					{Key: "ok", Value: 0},
					{Key: "n", Value: 0},
				}
				m.AddMockResponses(response)
			},
			wantErr:       true,
			expectedLenID: 0,
		},
	}

	for _, testCase := range testTable {
		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		mt.Run(testCase.name, func(mt *mtest.T) {
			mongoRep := GoodService.NewGoodsMongoRepo(mt.Client)

			testCase.mockBehavior(mt)

			sellerID, err := mongoRep.CreateSeller(testCase.inputUserID, testCase.inputUserName)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Len(t, sellerID, testCase.expectedLenID)
		})
	}
}

func TestMongoRep_getItemByID(t *testing.T) {
	type mockBehavior func(m *mtest.T)

	testTable := []struct {
		name         string
		inputItemId  string
		mockBehavior mockBehavior
		expectedItem GoodService.Item
		wantErr      bool
	}{
		{
			name:        "OK",
			inputItemId: "507f1f77bcf86cd799439011",
			mockBehavior: func(m *mtest.T) {
				objectID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
				response := mtest.CreateCursorResponse(1, "item.test", mtest.FirstBatch,
					bson.D{
						{Key: "_id", Value: objectID},
						{Key: "name", Value: "apple"},
						{Key: "description", Value: "tasty apple"},
						{Key: "quantity", Value: 15},
						{Key: "seller_id", Value: "100"},
					})
				m.AddMockResponses(response)
			},
			expectedItem: GoodService.Item{
				ID:          "507f1f77bcf86cd799439011",
				Name:        "apple",
				Description: "tasty apple",
				Quantity:    15,
				SellerID:    "100",
			},
			wantErr: false,
		},
		{
			name:         "Cant parse itemId",
			inputItemId:  "1",
			mockBehavior: func(m *mtest.T) {},
			expectedItem: GoodService.Item{},
			wantErr:      true,
		},
		{
			name:        "Item not found",
			inputItemId: "507f1f77bcf86cd799439011",
			mockBehavior: func(m *mtest.T) {
				response := bson.D{
					{Key: "ok", Value: 1},
					{Key: "value", Value: nil},
				}
				m.AddMockResponses(response)
			},
			expectedItem: GoodService.Item{},
			wantErr:      true,
		},
		{
			name:        "Decode error",
			inputItemId: "507f1f77bcf86cd799439011",
			mockBehavior: func(m *mtest.T) {
				response := bson.D{
					{Key: "ok", Value: 1},
					{Key: "value", Value: nil},
					{Key: "errmsg", Value: "decode error"},
				}
				m.AddMockResponses(response)
			},
			expectedItem: GoodService.Item{},
			wantErr:      true,
		},
	}

	for _, testCase := range testTable {
		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		mt.Run(testCase.name, func(mt *mtest.T) {
			mongoRepo := GoodService.NewGoodsMongoRepo(mt.Client)

			testCase.mockBehavior(mt)

			item, err := mongoRepo.GetItemByID(testCase.inputItemId)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, item, testCase.expectedItem)
		})
	}
}
