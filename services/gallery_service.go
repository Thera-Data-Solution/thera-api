package services

import (
	"math"
	"thera-api/dto"
	"thera-api/logger"
	"thera-api/models"
	"thera-api/repositories"
	"time"

	"go.uber.org/zap"
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

func (s *GalleryService) GetAllGallery(tenantID string, page, pageSize int) (*dto.GalleryPaginationResponse, error) {
	gallery, total, err := s.GalleryRepo.FindAllWithPagination(tenantID, page, pageSize)
	if err != nil {
		logger.Log.Error("Failed to fetch all gallery with pagination", zap.String("tenantId", tenantID), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.Error(err))
		return nil, err
	}

	galleryResponses := make([]dto.GalleryResponse, len(gallery))
	for i, gallery := range gallery {
		galleryResponses[i] = dto.GalleryResponse{
			ID:          gallery.ID,
			Title:       gallery.Title,
			Description: gallery.Description,
			ImageUrl:    gallery.ImageUrl,
			CreatedAt:   gallery.CreatedAt.Format("2006-01-02"),
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages == 0 && total > 0 { // Untuk memastikan 1 halaman jika ada data tapi pageSize sangat besar
		totalPages = 1
	}

	response := &dto.GalleryPaginationResponse{
		Data:       galleryResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
	logger.Log.Info("Fetched gallery with pagination", zap.String("tenantId", tenantID), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.Int64("total", total))
	return response, nil
}

func (s *GalleryService) GetGalleryByID(id string, tenant string) (*models.Gallery, error) {
	return s.GalleryRepo.FindByID(id, tenant)
}

func (s *GalleryService) GetGalleryByIDAndTenant(id string, tenant string) (*models.Gallery, error) {
	return s.GalleryRepo.FindByIDAndTenant(id, tenant)
}
