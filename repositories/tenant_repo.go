package repositories

import (
	"thera-api/models"

	"gorm.io/gorm"
)

type TenantRepository struct {
	DB *gorm.DB
}

func (r *TenantRepository) Create(tenant *models.Tenant) error {
	return r.DB.Create(tenant).Error
}

func (r *TenantRepository) FindAll() ([]models.Tenant, error) {
	var tenants []models.Tenant
	err := r.DB.Find(&tenants).Error
	return tenants, err
}

func (r *TenantRepository) FindByID(id string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.DB.First(&tenant, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepository) Update(tenant *models.Tenant) error {
	return r.DB.Save(tenant).Error
}

func (r *TenantRepository) Delete(id string) error {
	return r.DB.Delete(&models.Tenant{}, "id = ?", id).Error
}
