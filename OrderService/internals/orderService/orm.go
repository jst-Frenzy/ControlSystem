package orderService

type CartItem struct {
	Id        int `json:"id"`
	CartID    int `json:"cart_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
