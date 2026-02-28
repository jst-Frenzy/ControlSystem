package orderService

import "gorm.io/gorm"

type OrderPostgresRep interface {
	AddToCart(CartItem) (int, error)
	RemoveFromCart(int, string) error
	GetCart(int) ([]CartItem, float64, error)
}

type orderPostgresRep struct {
	db *gorm.DB
}

func NewOrderPostgresRep(db *gorm.DB) OrderPostgresRep {
	return &orderPostgresRep{db: db}
}

func (r *orderPostgresRep) AddToCart(CartItem) (int, error) {
	return 0, nil
}

func (r *orderPostgresRep) RemoveFromCart(int, string) error {
	return nil
}

func (r *orderPostgresRep) GetCart(int) ([]CartItem, float64, error) {
	return []CartItem{}, 0, nil
}
