package orderService

type CartItem struct {
	Id        int     `json:"id"`
	CartID    int     `json:"cart_id"`
	Name      string  `json:"name"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
