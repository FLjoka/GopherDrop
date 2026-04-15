package handler

import (
	"apps/go-api/internal/core/usecase"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArchivoHandler struct {
	service *usecase.ArchivoService
}

func NewArchivoHandler(service *usecase.ArchivoService) *ArchivoHandler {
	return &ArchivoHandler{service: service}
}

func (h *ArchivoHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file"})
		return
	}
	dst := fmt.Sprintf("./uploads/%s", file.Filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "file uploaded successfully"})
}
