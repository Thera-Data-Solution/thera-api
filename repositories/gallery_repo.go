package repositories

import (
	"thera-api/models"

	"gorm.io/gorm"
)

type GalleryRepository struct {
	DB *gorm.DB
}

func (r *GalleryRepository) Create(gallery *models.Gallery) error {
	return r.DB.Create(gallery).Error
}

func (r *GalleryRepository) Update(gallery *models.Gallery) error {
	return r.DB.
		Where(`id = ? AND tenant_id = ?`, gallery.ID, gallery.TenantId).
		Save(gallery).
		Error
}

func (r *GalleryRepository) Delete(id string, tenantId string) error {
	return r.DB.Delete(&models.Gallery{}, `id = ? AND tenant_id = ?`, id, tenantId).Error
}

func (r *GalleryRepository) FindByID(id string, tenant string) (*models.Gallery, error) {
	var gallery models.Gallery
	err := r.DB.First(&gallery, `id = ? AND tenant_id = ?`, id, tenant).Error
	if err != nil {
		return nil, err
	}
	return &gallery, nil
}

func (r *GalleryRepository) FindByIDAndTenant(id string, tenant string) (*models.Gallery, error) {
	var gallery models.Gallery
	err := r.DB.First(&gallery, `id = ? AND tenant_id = ?`, id, tenant).Error
	if err != nil {
		return nil, err
	}
	return &gallery, nil
}

func (r *GalleryRepository) FindAllWithPagination(tenantID string, page, pageSize int) ([]models.Gallery, int64, error) {
	var gallery []models.Gallery
	var total int64

	query := r.DB.Model(&models.Gallery{}).Where("tenant_id = ?", tenantID)

	// Hitung total data
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Terapkan pagination
	offset := (page - 1) * pageSize
	err := query.
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC"). // Urutkan berdasarkan tanggal terbaru
		Find(&gallery).Error

	if err != nil {
		return nil, 0, err
	}

	return gallery, total, nil
}
