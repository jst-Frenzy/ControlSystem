package orderService

type OrderService interface {
	AddToCart(CartItem) (int, error)
	RemoveFromCart(int, string) error
	GetCart(int) ([]CartItem, float64, error)
}

type orderService struct {
	repo OrderPostgresRep
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

func (s *orderService) GetCart(cartID int) ([]CartItem, float64, error) {
	//update info about items grpc call to goods
	return s.repo.GetCart(cartID)
}
