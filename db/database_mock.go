package db

import (
	"database/sql"
	"errors"
	"fmt"

	"example.com/mod/models"
)

type DataBaseMock struct {
	connection *sql.DB
}

func (db *DataBaseMock) Connect() error {
	return nil
}

func (db *DataBaseMock) CreateDb() error {
	return nil
}

func (db *DataBaseMock) Seed() error {
	return nil
}

func (db *DataBaseMock) createTable (tableName string, fields string) error {
	return nil
}

func (db *DataBaseMock) truncateTable (tableName string) error {
	return nil
}

func (db *DataBaseMock) GetUser(field string, value string) (models.User, error) {
	user := models.User{}

	if value == "error" {
		return user, errors.New("an error")
	}

	if value == "jon.doe@company.com" {
		user.Id = 1
		user.Name = "Jon Doe"
		user.Password = []byte(`$2a$14$bkJ6sCgPTrvBV5naw5CxNO1ukyb/iSp5wI0AqmgF0GOOmwr.yrAYG`)
		user.Email = "jon.doe@company.com"
	}

	return user, nil
}

func (db *DataBaseMock) GetProduct(productId int) (models.Product, error) {
	product := models.Product{}

	rows, err := db.connection.Query(fmt.Sprintf("SELECT * FROM products WHERE id = %d", productId))

	if err != nil {
		return product, err
	}

	for rows.Next() {
		err = rows.Scan(&product.Id, &product.Name, &product.Price, &product.Stock)

		if err != nil {
			return product, err
		}
	}

	err = rows.Close()

	if err != nil {
		return product, err
	}

	return product, nil
}

func (db *DataBaseMock) AddToCart(user models.User, product models.Product, quantity uint64) error {
	return nil
}

func (db *DataBaseMock) GetCartInfo(userId int) ([]models.Cart, error) {
	var cartInfo []models.Cart

	return cartInfo, nil
}

func (db *DataBaseMock) Checkout(user models.User, cartTotals map[int]int) error {
	return nil
}