package GoodService

type Item struct {
	ID          string `json:"_id" bson:"_id,omitempty"`
	Name        string `json:"name" bson:"name" binding:"required"`
	Description string `json:"description" bson:"description" binding:"required"`
	Quantity    int    `json:"quantity" bson:"quantity" binding:"required"`
	SellerID    string `json:"sellerID" bson:"seller_id"`
}

type Seller struct {
	Id     string `bson:"_id,omitempty"`
	UserID int    `bson:"user_id" binding:"required"`
	Name   string `bson:"name" binding:"required"`
}
