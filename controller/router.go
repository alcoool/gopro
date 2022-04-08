package controller

import (
	"log"
	"net/http"

	"example.com/mod/db"
	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine, database db.DataBaseInterface) http.Handler {

	router.Use(gin.Recovery())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "Route not found!")
	})

	err := database.Connect()
	if err != nil {
		log.Fatal("Unable to connect to db")
	}

	ctrl := &Controller{DB: &(database)}

	router.GET("/seed", ctrl.Seed)
	router.POST("/login", ctrl.Login)
	router.POST("/logout", ctrl.Logout)
	router.POST("/add-to-cart", ctrl.AddToCart)
	router.POST("/checkout", ctrl.Checkout)

	return router
}

