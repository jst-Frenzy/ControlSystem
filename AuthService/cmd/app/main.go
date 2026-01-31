package main

import (
	"AuthService/internal/handlers"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", handlers.Handler)

	err := router.Run(":8080")
	if err != nil {
		fmt.Println("Error starting server")
	}
}
