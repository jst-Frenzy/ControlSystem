package orderService

import (
	"context"
	"errors"
	"github.com/jst-Frenzy/ControlSystem/OrderService/internals/gRPC/client"
	"strconv"
)

type OrderService interface {
	AddToCart(CartItem) (int, error)
	RemoveFromCart(int, string) error
	GetCart(int, context.Context) ([]CartItem, float64, error)
}

type orderService struct {
	repo        OrderPostgresRep
	goodsClient client.GoodsClient
}

func NewOrderService(repo OrderPostgresRep) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) AddToCart(i CartItem) (int, error) {
	return s.repo.AddToCart(i)
}

func (s *orderService) RemoveFromCart(cartID int, itemID string) error {
	return s.repo.RemoveFromCart(cartID, itemID)
}

func (s *orderService) GetCart(cartID int, ctx context.Context) ([]CartItem, float64, error) {
	cart, err := s.repo.GetCart(cartID)

	if err != nil {
		return nil, 0, err
	}

	var totalPrice float64

	for _, item := range cart {
		resp, errGet := s.goodsClient.GetItemQuantityAndPrice(ctx, item.ProductID)
		if errGet != nil {
			return nil, 0, errGet
		}
		if !resp.Valid {
			return nil, 0, errors.New("can't get info about item")
		}
		price, _ := strconv.ParseFloat(resp.Price, 64)
		quantity, _ := strconv.Atoi(resp.Quantity)
		if price != item.Price {
			item.Price = price
		}
		if quantity != item.Quantity {
			item.Quantity = quantity
		}
	}

	return cart, totalPrice, nil
}
