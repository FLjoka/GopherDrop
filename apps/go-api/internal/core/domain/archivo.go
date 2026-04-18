package domain

import "gorm.io/gorm"

type FileMetadata struct {
	gorm.Model          // Agrega ID, CreatedAt, UpdatedAt, DeletedAt automáticamente
	OriginalName string `json:"original_name"`
	StoredName   string `json:"stored_name"`
	FilePath     string `json:"file_path"`
	Size         int64  `json:"size"`
	ContentType  string `json:"content_type"`
	ThumbnailPath string `json:"thumbnail_path"`
}

type FileMetadataRepository interface {
	Save(archivo *FileMetadata) error
	FindAll() ([]FileMetadata, error)
	FindByID(id int) (*FileMetadata, error)
	FindByName(name string) (*FileMetadata, error)
	Update(archivo *FileMetadata) error
	Delete(id int) error
}
