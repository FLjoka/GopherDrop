package handler

import (
	"apps/go-api/internal/core/domain"
	"apps/go-api/internal/core/usecase"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArchivoHandler struct {
	service *usecase.FileService
}

func NewArchivoHandler(service *usecase.FileService) *ArchivoHandler {
	return &ArchivoHandler{service: service}
}

// @Summary Subir archivo
// @Description Sube un archivo
// @Tags archivos
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Archivo a subir"
// @Success 200 {object} domain.FileMetadata
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /upload [post]
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
	metadara := domain.FileMetadata{
		OriginalName: file.Filename,
		StoredName:   file.Filename,
		FilePath:     dst,
		Size:         file.Size,
		ContentType:  file.Header.Get("Content-Type"),
	}
	if err := h.service.Save(&metadara); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file metadata"})
		return
	}
	c.JSON(http.StatusOK, metadara)
}

// @Summary Listar archivos
// @Description Obtiene una lista de todos los archivos
// @Tags archivos
// @Produce json
// @Success 200 {array} domain.FileMetadata
// @Router /files [get]
func (h *ArchivoHandler) GetFiles(c *gin.Context) {
	files, err := h.service.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los archivos"})
		return
	}
	c.JSON(http.StatusOK, files)
}

// @Summary Descargar archivo
// @Description Descarga un archivo por su ID
// @Tags archivos
// @Produce octet-stream
// @Param id path int true "ID del archivo"
// @Success 200 {file} file
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /download/{id} [get]
func (h *ArchivoHandler) Download(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de archivo inválido"})
		return
	}

	file, err := h.service.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Archivo no encontrado"})
		return
	}
	c.File(file.FilePath)
}

// @Summary Obtener miniatura
// @Description Devuelve la miniatura de una imagen por su ID
// @Tags archivos
// @Produce image/jpeg
// @Param id path int true "ID del archivo"
// @Success 200 {file} file
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /thumbnails/{id} [get]
func (h *ArchivoHandler) GetThumbnail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de archivo inválido"})
		return
	}

	file, err := h.service.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Archivo no encontrado"})
		return
	}

	if file.ThumbnailPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Miniatura no disponible para este archivo"})
		return
	}

	c.File(file.ThumbnailPath)
}

// @Summary Borrar archivo
// @Description Elimina el registro y los archivos físicos asociados
// @Tags archivos
// @Param id path int true "ID del archivo"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /files/{id} [delete]
func (h *ArchivoHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de archivo inválido"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo borrar el archivo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Archivo borrado correctamente"})
}
