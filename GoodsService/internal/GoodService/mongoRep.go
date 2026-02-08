package GoodService

type GoodsMongoRepo interface {
	GetGoods() ([]Item, error)

	CreateItem(item Item) (int, error)
	DeleteItem(int) error
	UpdateItem(Item) (Item, error)
}

type goodsMongoRepo struct {
	db interface{}
}

func NewGoodsMongoRepo(db interface{}) GoodsMongoRepo {
	return &goodsMongoRepo{db: db}
}

func (gr *goodsMongoRepo) GetGoods() ([]Item, error) {
	return []Item{}, nil
}

func (gr *goodsMongoRepo) CreateItem(item Item) (int, error) {
	return 0, nil
}

func (gr *goodsMongoRepo) DeleteItem(int) error {
	return nil
}
func (gr *goodsMongoRepo) UpdateItem(Item) (Item, error) {
	return Item{}, nil
}
