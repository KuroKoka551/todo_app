package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()
	service := todoService{db: db}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/todos", service.GetTodos)
	router.GET("/todos/:id", service.GetTodo)
	router.POST("/todos", service.AddTodo)
	router.PUT("/todos/:id", service.UpdateTodo)
	router.DELETE("/todos/:id", service.DeleteTodo)

	srv := http.Server{Addr: config.Addr, Handler: router}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("listen error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Starting server, press CTRL+C to stop")
	<-quit
	log.Println("Shutting down server...")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("Server shutdown error: %s\n", err)
	}

	log.Println("Server exited")
}
