package GoodService

type GoodService interface {
	GetGoods() ([]Item, error)
	AddItem(Item) (int, error)
	DeleteItem(int) error
	UpdateItem(Item) (Item, error)
}

type goodService struct {
	repo GoodsMongoRepo
}

func NewGoodService(repo GoodsMongoRepo) GoodService {
	return &goodService{repo: repo}
}

func (gs *goodService) GetGoods() ([]Item, error) {
	return gs.repo.GetGoods()
}

func (gs *goodService) AddItem(i Item) (int, error) {
	return gs.repo.CreateItem(i)
}

func (gs *goodService) DeleteItem(id int) error {
	return gs.repo.DeleteItem(id)
}

func (gs *goodService) UpdateItem(i Item) (Item, error) {
	return gs.repo.UpdateItem(i)
}
