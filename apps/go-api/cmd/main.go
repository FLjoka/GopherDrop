package main

import (
	"apps/go-api/internal/adapters/handler"
	"apps/go-api/internal/adapters/infrastructure"
	"apps/go-api/internal/adapters/repository"
	"apps/go-api/internal/core/usecase"

	_ "apps/go-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title GopherDrop API
// @version 1.0
// @description API para el manejo de archivos en GopherDrop.
// @host localhost:8080
// @BasePath /

func main() {
	db, err := repository.InitDB()
	if err != nil {
		panic("failed to connect database")
	}

	// Infrastructure Adapters
	repo := repository.NewSQLFileMetadataRepository(db)
	thumbProvider := infrastructure.NewImagingProvider()

	// Use Cases
	service := usecase.NewFileService(repo, thumbProvider)

	// Handlers
	h := handler.NewArchivoHandler(service)

	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/upload", h.Upload)
	r.GET("/files", h.GetFiles)
	r.GET("/download/:id", h.Download)
	r.GET("/thumbnails/:id", h.GetThumbnail)
	r.DELETE("/files/:id", h.Delete)
	r.Run(":8080")
}
