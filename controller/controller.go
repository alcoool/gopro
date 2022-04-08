package controller

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"example.com/mod/db"
	"example.com/mod/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "this should be a secret"

var mutex sync.Mutex

type Controller struct {
	DB *db.DataBaseInterface
}

func (ctrl *Controller) Seed(c *gin.Context) {
	err := (*ctrl.DB).CreateDb()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not create db",
		})
		return
	}
	err = (*ctrl.DB).Seed()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not seed db",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
	return
}

func (ctrl *Controller) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := (*ctrl.DB).GetUser("email", email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not login",
		})
		return
	}

	if user.Id == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "not allowed",
		})
		return
	}

	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "not allowed",
		})
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not login",
		})
		return
	}

	c.SetCookie("jwt", token, 600, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
	return
}

func (ctrl *Controller) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
	return
}

func (ctrl *Controller) AddToCart(c *gin.Context) {
	user, err := ctrl.checkUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	productId, err := strconv.Atoi(c.PostForm("productId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad product id",
		})
		return
	}

	quantity, err := strconv.ParseUint(c.PostForm("quantity"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad quantity",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad product id",
		})
		return
	}

	product, err := (*ctrl.DB).GetProduct(productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not get product from db",
		})
		return
	}

	if product.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "product not found",
		})
		return
	}

	if product.Stock < quantity {
		c.JSON(http.StatusFound, gin.H{
			"message": "not enough stock",
		})
		return
	}

	err = (*ctrl.DB).AddToCart(user, product, quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not add to cart",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
	return
}

func (ctrl *Controller) Checkout(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()

	user, err := ctrl.checkUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	cartInfo, err := (*ctrl.DB).GetCartInfo(int(user.Id))

	cartTotals := make(map[int]int)

	for _, cart := range cartInfo {
		if _, ok := cartTotals[cart.ProductId]; ok {
			cartTotals[cart.ProductId] = cartTotals[cart.ProductId] + int(cart.Quantity)
		} else {
			cartTotals[cart.ProductId] = int(cart.Quantity)
		}
	}

	err = (*ctrl.DB).Checkout(user, cartTotals)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
	return
}

func (ctrl *Controller) checkUser(c *gin.Context) (models.User, error) {
	user := models.User{}

	cookie, err := c.Cookie("jwt")

	if err != nil {
		return user, err
	}

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		return user, err
	}

	claims := token.Claims.(*jwt.StandardClaims)

	user, err = (*ctrl.DB).GetUser("id", claims.Issuer)

	if err != nil {
		return user, err
	}

	return user, nil
}
