package services

import (
	"math"
	"thera-api/dto"
	"thera-api/logger"
	"thera-api/models"
	"thera-api/repositories"

	"go.uber.org/zap"
)

type PartnerService struct {
	PartnerRepo *repositories.PartnerRepository
}

func (s *PartnerService) CreatePartner(logo *string, tenantId *string) (*models.Partner, error) {

	partner := &models.Partner{
		Logo:     *logo,
		TenantId: *tenantId,
	}
	if err := s.PartnerRepo.Create(partner); err != nil {
		return nil, err
	}
	return partner, nil
}

func (s *PartnerService) UpdatePartner(id string, logo *string, tenantId *string) (*models.Partner, error) {

	partner, err := s.PartnerRepo.FindByID(id, *tenantId)

	if err != nil {
		return nil, err
	}

	if logo != nil {
		partner.Logo = *logo
	}

	if err := s.PartnerRepo.Update(partner); err != nil {
		return nil, err
	}
	return partner, nil
}

func (s *PartnerService) GetAllPartners(tenantID string, page, pageSize int) (*dto.PartnerPaginationResponse, error) {
	partner, total, err := s.PartnerRepo.FindAllWithPagination(tenantID, page, pageSize)
	if err != nil {
		logger.Log.Error("Failed to fetch all partner with pagination", zap.String("tenantId", tenantID), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.Error(err))
		return nil, err
	}

	partnerResponses := make([]dto.PartnerResponse, len(partner))
	for i, partner := range partner {
		partnerResponses[i] = dto.PartnerResponse{
			ID:       partner.ID,
			Logo:     &partner.Logo,
			TenantId: &partner.TenantId,
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages == 0 && total > 0 { // Untuk memastikan 1 halaman jika ada data tapi pageSize sangat besar
		totalPages = 1
	}

	response := &dto.PartnerPaginationResponse{
		Data:       partnerResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
	logger.Log.Info("Fetched partner with pagination", zap.String("tenantId", tenantID), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.Int64("total", total))
	return response, nil
}

func (s *PartnerService) Delete(id string, tenantId string) error {
	return s.PartnerRepo.Delete(id, tenantId)
}
