package main

import (
	"apps/go-api/internal/adapters/handler"
	"apps/go-api/internal/core/usecase"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	service := &usecase.ArchivoService{}
	h := handler.NewArchivoHandler(service)
	fmt.Println("Handler initialized:", h)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8080")
}
