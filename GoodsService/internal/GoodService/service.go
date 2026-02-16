package GoodService

import (
	"errors"
)

//go:generate mockgen -source=service.go -destination=../mocks/mockServ.go -package=mocks

type GoodService interface {
	GetGoods() ([]Item, error)

	AddItem(Item, UserCtx) (string, error)
	DeleteItem(string, int) error
	UpdateItem(Item, int) (Item, error)
}

type goodService struct {
	repo GoodsMongoRep
}

func NewGoodService(repo GoodsMongoRep) GoodService {
	return &goodService{repo: repo}
}

func (s *goodService) GetGoods() ([]Item, error) {
	return s.repo.GetGoods()
}

func (s *goodService) AddItem(i Item, seller UserCtx) (string, error) {
	var sellerID string
	var err error
	sellerID, err = s.repo.GetSellerIDByUserID(seller.ID)
	if err != nil {
		if err.Error() == "seller not found" {
			sellerID, err = s.repo.CreateSeller(seller.ID, seller.Name)
			if err != nil {
				return "", errors.New("error create seller")
			}
		} else {
			return "", err
		}
	}
	i.SellerID = sellerID
	return s.repo.CreateItem(i)
}

func (s *goodService) DeleteItem(itemID string, userID int) error {
	sellerID, err := s.repo.GetSellerIDByUserID(userID)

	if err != nil {
		return err
	}
	return s.repo.DeleteItem(itemID, sellerID)
}

func (s *goodService) UpdateItem(i Item, userID int) (Item, error) {
	sellerID, err := s.repo.GetSellerIDByUserID(userID)
	if err != nil {
		return Item{}, err
	}

	oldItem, err := s.repo.GetItemByID(i.ID)
	if err != nil {
		return Item{}, err
	}

	if oldItem.SellerID != sellerID {
		return Item{}, errors.New("it's not your item")
	}
	i.SellerID = sellerID
	return s.repo.UpdateItem(i)
}
