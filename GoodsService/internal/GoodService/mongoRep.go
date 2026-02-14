package GoodService

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GoodsMongoRepo interface {
	GetGoods() ([]Item, error)

	GetQuantity(string) (int, error)

	CreateItem(Item) (string, error)
	DeleteItem(string, string) error
	UpdateItem(Item) (Item, error)

	GetSellerIDByUserID(int) (string, error)
	GetItemByID(string) (Item, error)
	CreateSeller(int, string) (string, error)
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

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("cant convert id to ObjectID")
	}

	return id.Hex(), nil
}

func (r *goodsMongoRepo) DeleteItem(itemID string, sellerID string) error {
	objectID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", objectID}, {"seller_id", sellerID}}
	result, err := r.itemCollection.DeleteOne(r.ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("item not found")
	}

	return nil
}

func (r *goodsMongoRepo) UpdateItem(item Item) (Item, error) {
	objectID, err := primitive.ObjectIDFromHex(item.ID)
	if err != nil {
		return Item{}, err
	}

	filter := bson.D{{"_id", objectID}}
	update := bson.D{
		{"$set", bson.D{
			{"name", item.Name},
			{"description", item.Description},
			{"quantity", item.Quantity},
			{"seller_id", item.SellerID},
		}},
	}
	res := r.itemCollection.FindOneAndUpdate(r.ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))

	var newItem Item
	errDecode := res.Decode(&newItem)
	if errDecode != nil {
		if errors.Is(errDecode, mongo.ErrNoDocuments) {
			return Item{}, errors.New("item not found")
		}
		return Item{}, err
	}

	return newItem, nil
}

func (r *goodsMongoRepo) GetQuantity(itemID string) (int, error) {
	objectID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return 0, err
	}

	filter := bson.D{{"_id", objectID}}
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

func (r *goodsMongoRepo) GetSellerIDByUserID(id int) (string, error) {
	filter := bson.D{{"user_id", id}}
	res := r.sellerCollection.FindOne(r.ctx, filter)

	var s Seller
	if err := res.Decode(&s); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", errors.New("item not found")
		}
		return "", err
	}
	return s.Id, nil
}

func (r *goodsMongoRepo) CreateSeller(userID int, name string) (string, error) {
	var s = Seller{
		UserID: userID,
		Name:   name,
	}
	res, err := r.sellerCollection.InsertOne(r.ctx, s)
	if err != nil {
		return "", err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("cant convert id to ObjectID")
	}

	return id.Hex(), nil
}

func (r *goodsMongoRepo) GetItemByID(id string) (Item, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Item{}, err
	}

	filter := bson.D{{"_id", objectID}}
	res := r.itemCollection.FindOne(r.ctx, filter)

	var i Item
	if errDecode := res.Decode(&i); errDecode != nil {
		if errors.Is(errDecode, mongo.ErrNoDocuments) {
			return Item{}, errors.New("item not found")
		}
		return Item{}, errDecode
	}

	return i, nil
}
