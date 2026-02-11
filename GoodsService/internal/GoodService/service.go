package GoodService

type GoodService interface {
	GetGoods() ([]Item, error)

	AddItem(Item) (string, error)
	DeleteItem(string) error
	UpdateItem(Item) (Item, error)
}

type goodService struct {
	repo GoodsMongoRepo
}

func NewGoodService(repo GoodsMongoRepo) GoodService {
	return &goodService{repo: repo}
}

func (s *goodService) GetGoods() ([]Item, error) {
	return s.repo.GetGoods()
}

func (s *goodService) AddItem(i Item) (string, error) {
	return s.repo.CreateItem(i)
}

func (s *goodService) DeleteItem(itemID string) error {
	return s.repo.DeleteItem(itemID)
}

func (s *goodService) UpdateItem(i Item) (Item, error) {
	return s.repo.UpdateItem(i)
}
