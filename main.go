package main

import (
	"fmt"
	"net/http"
	"time"

	"example.com/mod/controller"
	"example.com/mod/db"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func main() {
	engine := gin.New()

	database := db.DataBase{}

	router := controller.Router(engine, &database)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", "localhost", "8080"),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		return server.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Error().Msg(fmt.Sprintf("Error: %s", err))
	}
}
