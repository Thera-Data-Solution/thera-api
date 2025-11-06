package services

import (
	"thera-api/models"
	"thera-api/repositories"
	"time"
)

type GalleryService struct {
	GalleryRepo *repositories.GalleryRepository
}

func (s *GalleryService) CreateGallery(
	title *string,
	description *string,
	image *string,
	createdAt time.Time,
	tenantId *string,
) (*models.Gallery, error) {
	gallery := &models.Gallery{
		Title:       title,
		Description: description,
		ImageUrl:    *image,
		CreatedAt:   createdAt,
		TenantId:    tenantId,
	}
	if err := s.GalleryRepo.Create(gallery); err != nil {
		return nil, err
	}
	return gallery, nil
}

func (s *GalleryService) UpdateGallery(
	id string,
	title *string,
	description *string,
	image *string,
	createdAt time.Time,
	tenantId string,
) (*models.Gallery, error) {
	gallery, err := s.GalleryRepo.FindByID(id, tenantId)

	if err != nil {
		return nil, err
	}

	if title != nil {
		gallery.Title = title
	}
	if description != nil {
		gallery.Description = description
	}
	if image != nil {
		gallery.ImageUrl = *image
	}

	if err := s.GalleryRepo.Update(gallery); err != nil {
		return nil, err
	}

	return gallery, nil
}

func (s *GalleryService) DeleteGallery(id string, tenantId string) error {
	return s.GalleryRepo.Delete(id, tenantId)
}

func (s *GalleryService) GetAllGallery(tenant string) ([]models.Gallery, error) {
	return s.GalleryRepo.FindAll(tenant)
}

func (s *GalleryService) GetGalleryByID(id string, tenant string) (*models.Gallery, error) {
	return s.GalleryRepo.FindByID(id, tenant)
}

func (s *GalleryService) GetGalleryByIDAndTenant(id string, tenant string) (*models.Gallery, error) {
	return s.GalleryRepo.FindByIDAndTenant(id, tenant)
}
