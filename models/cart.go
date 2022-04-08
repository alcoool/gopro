package models

type Cart struct {
	Id 		    uint   `json:"id"`
	UserId 	    int    `json:"user_id"`
	ProductId 	int    `json:"product_id"`
	Quantity 	uint64 `json:"quantity"`
}
