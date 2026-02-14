package GoodService

type Item struct {
	ID          string `json:"_id" bson:"_id,omitempty"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Quantity    int    `json:"quantity" bson:"quantity"`
	SellerID    string `json:"sellerID" bson:"seller_id"`
}

type Seller struct {
	Id     string `bson:"_id,omitempty"`
	UserID int    `bson:"user_id"`
	Name   string `bson:"name"`
}
