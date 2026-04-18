package usecase

import (
	"apps/go-api/internal/core/domain"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileService struct {
	repo              domain.FileMetadataRepository
	thumbnailProvider domain.ThumbnailProvider
}

func NewFileService(repo domain.FileMetadataRepository, thumbProvider domain.ThumbnailProvider) *FileService {
	return &FileService{
		repo:              repo,
		thumbnailProvider: thumbProvider,
	}
}

func (s *FileService) Save(archivo *domain.FileMetadata) error {
	// Check if it's an image
	lowerContentType := strings.ToLower(archivo.ContentType)
	if strings.HasPrefix(lowerContentType, "image/") {
		// Ensure thumbnails directory exists
		thumbDir := "data/thumbnails"
		if err := os.MkdirAll(thumbDir, 0755); err != nil {
			return fmt.Errorf("failed to create thumbnails directory: %w", err)
		}

		// Generate thumbnail path
		fileName := filepath.Base(archivo.FilePath)
		thumbName := "thumb_" + fileName
		thumbPath := filepath.Join(thumbDir, thumbName)

		// Generate thumbnail: 200x200
		if err := s.thumbnailProvider.GenerateThumbnail(archivo.FilePath, thumbPath, 200, 200); err != nil {
			// We log the error but allow the file to be saved without thumbnail
			fmt.Printf("Warning: failed to generate thumbnail for %s: %v\n", archivo.OriginalName, err)
		} else {
			archivo.ThumbnailPath = thumbPath
		}
	}

	return s.repo.Save(archivo)
}

func (s *FileService) FindAll() ([]domain.FileMetadata, error) {
	return s.repo.FindAll()
}

func (s *FileService) FindByID(id int) (*domain.FileMetadata, error) {
	return s.repo.FindByID(id)
}

func (s *FileService) FindByName(name string) (*domain.FileMetadata, error) {
	return s.repo.FindByName(name)
}

func (s *FileService) Update(archivo *domain.FileMetadata) error {
	return s.repo.Update(archivo)
}

func (s *FileService) Delete(id int) error {
	// 1. Get metadata to find paths
	archivo, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// 2. Delete original file
	if archivo.FilePath != "" {
		_ = os.Remove(archivo.FilePath) // ignore error if file not found
	}

	// 3. Delete thumbnail if exists
	if archivo.ThumbnailPath != "" {
		_ = os.Remove(archivo.ThumbnailPath)
	}

	// 4. Delete from repository
	return s.repo.Delete(id)
}
