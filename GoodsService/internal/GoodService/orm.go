package GoodService

type Item struct {
	ID          string `bson:"_id"`
	Name        string `bson:"name"`
	Description string `bson:"description"`
	Quantity    int    `bson:"quantity"`
	SellerID    string `bson:"seller_id"`
}

type Seller struct {
	Id     string `bson:"_id"`
	userID int    `bson:"user_id"`
	Name   string `bson:"name"`
}
