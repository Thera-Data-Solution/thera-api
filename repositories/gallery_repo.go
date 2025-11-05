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

func (r *GalleryRepository) FindAll(tenant string) ([]models.Gallery, error) {
	var gallery []models.Gallery
	err := r.DB.Where(`tenant_id = ?`, tenant).Find(&gallery).Error
	return gallery, err
}
