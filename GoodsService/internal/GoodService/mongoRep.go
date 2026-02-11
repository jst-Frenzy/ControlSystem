package GoodService

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GoodsMongoRepo interface {
	GetGoods() ([]Item, error)

	GetQuantity(string) (int, error)

	CreateItem(Item) (string, error)
	DeleteItem(string) error
	UpdateItem(Item) (Item, error)
}

type goodsMongoRepo struct {
	itemCollection   *mongo.Collection
	sellerCollection *mongo.Collection
	ctx              context.Context
}

func NewGoodsMongoRepo(client *mongo.Client) GoodsMongoRepo {
	db := client.Database("GoodsInfo")
	return &goodsMongoRepo{
		itemCollection:   db.Collection("goods"),
		sellerCollection: db.Collection("sellers"),
		ctx:              context.Background(),
	}
}

func (r *goodsMongoRepo) GetGoods() ([]Item, error) {
	resp, errFind := r.itemCollection.Find(r.ctx, bson.M{})
	if errFind != nil {
		return nil, errFind
	}

	var allItems []Item
	if err := resp.All(r.ctx, &allItems); err != nil {
		return nil, err
	}

	return allItems, nil
}

func (r *goodsMongoRepo) CreateItem(item Item) (string, error) {
	res, err := r.itemCollection.InsertOne(r.ctx, item)

	if err != nil {
		return "", err
	}

	return res.InsertedID.(string), err
}

func (r *goodsMongoRepo) DeleteItem(itemID string) error {
	filter := bson.D{{"_id", itemID}}
	_, err := r.itemCollection.DeleteOne(r.ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (r *goodsMongoRepo) UpdateItem(item Item) (Item, error) {
	filter := bson.D{{"_id", item.ID}}
	res := r.itemCollection.FindOneAndReplace(r.ctx, filter, item, options.FindOneAndReplace().SetReturnDocument(options.After))

	var newItem Item
	err := res.Decode(&newItem)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Item{}, errors.New("item not found")
		}
		return Item{}, err
	}

	return newItem, nil
}

func (r *goodsMongoRepo) GetQuantity(itemID string) (int, error) {
	filter := bson.D{{"_id", itemID}}
	res := r.itemCollection.FindOne(r.ctx, filter)

	var i Item
	if err := res.Decode(&i); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, errors.New("item not found")
		}
		return 0, err
	}

	return i.Quantity, nil
}
