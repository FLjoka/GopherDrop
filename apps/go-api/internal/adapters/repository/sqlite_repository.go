package repository

import (
	"apps/go-api/internal/core/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLFileMetadataRepository struct {
	db *gorm.DB
}

func NewSQLFileMetadataRepository(db *gorm.DB) *SQLFileMetadataRepository {
	return &SQLFileMetadataRepository{db: db}
}

func (r *SQLFileMetadataRepository) Save(archivo *domain.FileMetadata) error {
	return r.db.Create(archivo).Error
}

func (r *SQLFileMetadataRepository) FindAll() ([]domain.FileMetadata, error) {
	var files []domain.FileMetadata
	result := r.db.Find(&files)
	if result.Error != nil {
		return nil, result.Error
	}
	return files, nil
}

func (r *SQLFileMetadataRepository) FindByName(name string) (*domain.FileMetadata, error) {
	var file domain.FileMetadata
	if err := r.db.Where("original_name = ?", name).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *SQLFileMetadataRepository) FindByID(id int) (*domain.FileMetadata, error) {
	var file domain.FileMetadata
	if err := r.db.First(&file, id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *SQLFileMetadataRepository) Update(archivo *domain.FileMetadata) error {
	return r.db.Save(archivo).Error
}

func (r *SQLFileMetadataRepository) Delete(id int) error {
	return r.db.Delete(&domain.FileMetadata{}, id).Error
}

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("data/db.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&domain.FileMetadata{})

	return db, nil
}
