package usecase

import "apps/go-api/internal/core/domain"

type ArchivoService struct {
	repo domain.ArchivoRepository
}

func NewArchivoService(repo domain.ArchivoRepository) *ArchivoService {
	return &ArchivoService{repo: repo}
}

func (s *ArchivoService) Save(archivo *domain.Archivo) error {
	return s.repo.Save(archivo)
}

func (s *ArchivoService) FindAll() ([]domain.Archivo, error) {
	return s.repo.FindAll()
}

func (s *ArchivoService) FindByID(id int) (*domain.Archivo, error) {
	return s.repo.FindByID(id)
}

func (s *ArchivoService) Update(archivo *domain.Archivo) error {
	return s.repo.Update(archivo)
}

func (s *ArchivoService) Delete(id int) error {
	return s.repo.Delete(id)
}
