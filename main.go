package main

import (
	"log"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

func main() {
	config := newConfig()

	db := newDB(config)
	err := db.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	service := todoService{db: db}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/todos", service.GetTodos)
	router.GET("/todos/:id", service.GetTodo)
	router.POST("/todos", service.AddTodo)
	router.PUT("/todos/:id", service.UpdateTodo)
	router.DELETE("/todos/:id", service.DeleteTodo)

	log.Fatal(router.Run(config.Addr))
}
