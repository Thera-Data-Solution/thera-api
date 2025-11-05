package services

import (
	"fmt"
	"thera-api/models"
	"thera-api/repositories"
)

type GalleryService struct {
	GalleryRepo *repositories.GalleryRepository
}

func (s *GalleryService) GetAllGallery(tenant string) ([]models.Gallery, error) {
	fmt.Println("testing disini2")
	return s.GalleryRepo.FindAll(tenant)
}
