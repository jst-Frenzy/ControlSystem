package orderService

import (
	"fmt"
	"gorm.io/gorm"
)

type OrderPostgresRep interface {
	AddToCart(CartItem) (int, error)
	RemoveFromCart(int, string) error
	GetCart(int) ([]CartItem, error)
}

type orderPostgresRep struct {
	db *gorm.DB
}

func NewOrderPostgresRep(db *gorm.DB) OrderPostgresRep {
	return &orderPostgresRep{db: db}
}

func (r *orderPostgresRep) AddToCart(i CartItem) (int, error) {
	if err := r.db.Table("carts").Create(&i).Error; err != nil {
		return 0, err
	}
	return i.Id, nil
}

func (r *orderPostgresRep) RemoveFromCart(cartID int, productID string) error {
	res := r.db.Table("carts").Delete(&CartItem{}, "cart_id = ? AND product_id = ?", cartID, productID)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("item with product id: %s in cart with id: %d don't exists", productID, cartID)
	}

	return nil
}

func (r *orderPostgresRep) GetCart(cartID int) ([]CartItem, error) {
	var items []CartItem
	if err := r.db.Table("carts").Where("cart_id = ?", cartID).Order("created_at desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
