package repositories

import (
	"thera-api/models"

	"gorm.io/gorm"
)

type PartnerRepository struct {
	DB *gorm.DB
}

func (r *PartnerRepository) Create(partner *models.Partner) error {
	return r.DB.Create(partner).Error
}

func (r *PartnerRepository) Update(partner *models.Partner) error {
	return r.DB.
		Where(`id = ? AND tenant_id = ?`, partner.ID, partner.TenantId).
		Save(partner).
		Error
}

func (r *PartnerRepository) Delete(id string, tenantId string) error {
	return r.DB.Delete(&models.Partner{}, `id = ? AND tenant_id = ?`, id, tenantId).Error
}

func (r *PartnerRepository) FindByID(id string, tenant string) (*models.Partner, error) {
	var partner models.Partner
	err := r.DB.First(&partner, `id = ? AND tenant_id = ?`, id, tenant).Error
	if err != nil {
		return nil, err
	}
	return &partner, nil
}

func (r *PartnerRepository) FindAllWithPagination(tenantID string, page, pageSize int) ([]models.Partner, int64, error) {
	var partners []models.Partner
	var total int64

	query := r.DB.Model(&models.Partner{}).Where("tenant_id = ?", tenantID)

	// Hitung total data
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Terapkan pagination
	offset := (page - 1) * pageSize
	err := query.
		Offset(offset).
		Limit(pageSize).
		Find(&partners).Error
	if err != nil {
		return nil, 0, err
	}

	return partners, total, nil
}
