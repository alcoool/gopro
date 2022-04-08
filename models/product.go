package models

type Product struct {
	Id 		uint   `json:"id"`
	Name 	string `json:"name"`
	Price 	uint   `json:"price"`
	Stock 	uint64   `json:"stock"`
}
