package GoodService_test

import (
	"github.com/jst-Frenzy/ControlSystem/GoodsService/internal/GoodService"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func TestMongoRep_createItem(t *testing.T) {
	type mockBehavior func(m *mtest.T)

	testTable := []struct {
		name          string
		inputItem     GoodService.Item
		mockBehavior  mockBehavior
		expectedError bool
		textError     string
	}{
		{
			name: "OK",
			inputItem: GoodService.Item{
				Name:        "apple",
				Description: "tasty apple",
				Quantity:    1,
				SellerID:    "100",
			},
			mockBehavior: func(m *mtest.T) {
				m.AddMockResponses(mtest.CreateSuccessResponse())
			},
			expectedError: false,
		},
		{
			name: "Error insert",
			inputItem: GoodService.Item{
				Name:        "apple",
				Description: "tasty apple",
				Quantity:    1,
				SellerID:    "100",
			},
			mockBehavior: func(m *mtest.T) {
				m.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Message: "Error insert",
				}))
			},
			expectedError: true,
			textError:     "Error insert",
		},
	}

	for _, testCase := range testTable {
		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

		mt.Run(testCase.name, func(mt *mtest.T) {

			mongoRep := GoodService.NewGoodsMongoRepo(mt.Client)

			testCase.mockBehavior(mt)

			res, err := mongoRep.CreateItem(testCase.inputItem)

			if testCase.expectedError {
				assert.Error(t, err)
				assert.Empty(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
			}
		})
	}
}

func TestMongoRep_deleteItem(t *testing.T) {

}

func TestMongoRep_updateItem(t *testing.T) {

}

func TestMongoRep_getGoods(t *testing.T) {

}

func TestMongoRep_getQuantity(t *testing.T) {

}

func TestMongoRep_getSellerIDyUserID(t *testing.T) {

}

func TestMongoRep_createSeller(t *testing.T) {

}

func TestMongoRep_getItemByID(t *testing.T) {

}
