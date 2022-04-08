package db

import (
	"database/sql"
	"errors"
	"fmt"

	"example.com/mod/models"
	"example.com/mod/payment"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type DataBaseInterface interface {
	Connect() error
	CreateDb() error
	Seed() error
	GetUser(string, string) (models.User, error)
	GetProduct(int) (models.Product, error)
	AddToCart(models.User, models.Product, uint64) error
	GetCartInfo(int) ([]models.Cart, error)
	Checkout(models.User, map[int]int) error
}

type DataBase struct {
	connection *sql.DB
}

func (db *DataBase) Connect() error {
	connection, err := sql.Open("sqlite3", "./shop.db")

	if err != nil {
		return err
	}

	db.connection = connection

	return nil
}

func (db *DataBase) CreateDb() error {
	err := db.createTable("users", "id INTEGER PRIMARY KEY, name TEXT, password TEXT, email TEXT, token TEXT")
	if err != nil {
		return err
	}
	err = db.createTable("products", "id INTEGER PRIMARY KEY, name TEXT, price INTEGER, stock INTEGER")
	if err != nil {
		return err
	}
	err = db.createTable("cart", "id INTEGER PRIMARY KEY, user_id INTEGER, product_id INTEGER, quantity INTEGER")
	if err != nil {
		return err
	}
	return nil
}

func (db *DataBase) Seed() error {
	err := db.truncateTable("users")
	if err != nil {
		return err
	}
	err = db.truncateTable("products")
	if err != nil {
		return err
	}

	statement, err := db.connection.Prepare("INSERT INTO users (name, password, email) VALUES (?, ?, ?), (?, ?, ?)")
	if err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte("1234"), 14)

	_, err = statement.Exec("Jon Doe", password, "jon.doe@company.com", "Jon Doe 2", password, "jon.doe2@company.com")
	if err != nil {
		return err
	}

	statement, err = db.connection.Prepare("INSERT INTO products (name, price, stock) VALUES (?, ?, ?), (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec("Product 1", 100, 2, "Product 2", 200, 3)
	if err != nil {
		return err
	}

	return nil
}

func (db *DataBase) createTable(tableName string, fields string) error {
	statement, err := db.connection.Prepare(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, fields))
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (db *DataBase) truncateTable(tableName string) error {
	statement, err := db.connection.Prepare(fmt.Sprintf("DELETE FROM %s", tableName))
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (db *DataBase) GetUser(field string, value string) (models.User, error) {
	user := models.User{}

	rows, err := db.connection.Query(fmt.Sprintf("SELECT * FROM users WHERE %s = \"%s\" LIMIT 1", field, value))

	if err != nil {
		return user, err
	}

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Password, &user.Email)

		if err != nil {
			return user, err
		}
	}

	err = rows.Close()

	if err != nil {
		return user, err
	}

	return user, nil
}

func (db *DataBase) GetProduct(productId int) (models.Product, error) {
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

func (db *DataBase) AddToCart(user models.User, product models.Product, quantity uint64) error {
	statement, err := db.connection.Prepare("INSERT INTO cart (user_id, product_id, quantity) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = statement.Exec(user.Id, product.Id, quantity)
	if err != nil {
		return err
	}

	return nil
}

func (db *DataBase) GetCartInfo(userId int) ([]models.Cart, error) {
	var cartInfo []models.Cart

	rows, err := db.connection.Query(fmt.Sprintf("SELECT * FROM cart WHERE user_id = %d", userId))

	if err != nil {
		return cartInfo, err
	}

	cart := models.Cart{}

	for rows.Next() {
		err = rows.Scan(&cart.Id, &cart.UserId, &cart.ProductId, &cart.Quantity)

		if err != nil {
			return cartInfo, err
		}

		cartInfo = append(cartInfo, cart)
	}

	err = rows.Close()

	if err != nil {
		return cartInfo, err
	}

	return cartInfo, nil
}

func (db *DataBase) Checkout(user models.User, cartTotals map[int]int) error {
	tx, err := db.connection.Begin()

	if err != nil {
		return err
	}

	cart := models.Cart{}
	totalPrice := 0

	for productId, total := range cartTotals {
		row := tx.QueryRow("SELECT * FROM products WHERE id = ? AND stock >= ?", productId, total)
		err = row.Scan(&cart.Id, &cart.UserId, &cart.ProductId, &cart.Quantity)

		if cart.Id == 0 {
			err = tx.Rollback()
			if err != nil {
				return err
			}
			return errors.New(fmt.Sprintf("Product %d has insuficent stock", productId))
		}

		_, err = tx.Exec("UPDATE products SET stock = stock - ? WHERE id = ?", total, productId)

		if err != nil {
			err = tx.Rollback()
			if err != nil {
				return err
			}
			return err
		}

		product, err := db.GetProduct(productId)
		if err != nil {
			return err
		}

		totalPrice = totalPrice + int(product.Price)
	}

	err = payment.Do(totalPrice)

	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return err
		}

		return err
	} else {
		_, txErr := tx.Exec("DELETE FROM cart WHERE user_id = ?", user.Id)
		if txErr != nil {
			return err
		}

		txErr = tx.Commit()
		if txErr != nil {
			return txErr
		}
	}

	return nil
}
