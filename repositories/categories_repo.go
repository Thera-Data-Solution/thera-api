package repositories

import (
	"thera-api/models"

	"gorm.io/gorm"
)

type CategoriesRepository struct {
	DB *gorm.DB
}

func (r *CategoriesRepository) Create(category *models.Categories) error {
	return r.DB.Create(category).Error
}

func (r *CategoriesRepository) FindAll(tenant string) ([]models.Categories, error) {
	var categories []models.Categories
	err := r.DB.
		Where(`tenant_id = ? AND disable=false`, tenant).
		Order("name ASC").
		Find(&categories).
		Error
	return categories, err
}

func (r *CategoriesRepository) FindAllAsAdmin(tenant string) ([]models.Categories, error) {
	var categories []models.Categories
	err := r.DB.
		Where(`tenant_id = ?`, tenant).
		Order("name ASC").
		Order("disable ASC").
		Find(&categories).
		Error
	return categories, err
}

func (r *CategoriesRepository) FindByID(id string, tenant string) (*models.Categories, error) {
	var category models.Categories
	err := r.DB.First(&category, `id = ? AND tenant_id = ? AND disable = false`, id, tenant).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoriesRepository) FindByTenantSlug(slug string, tenant string) (*models.Categories, error) {
	var category models.Categories
	err := r.DB.First(&category, `slug = ? AND tenant_id = ? AND disable = false`, slug, tenant).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoriesRepository) FindBySlug(slug string) (*models.Categories, error) {
	var category models.Categories
	err := r.DB.First(&category, `slug = ?`, slug).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoriesRepository) FindByIDAndTenant(id string, tenant string) (*models.Categories, error) {
	var category models.Categories
	err := r.DB.First(&category, `id = ? AND tenant_id = ?`, id, tenant).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoriesRepository) Update(category *models.Categories) error {
	return r.DB.
		Where(`id = ? AND tenant_id = ?`, category.ID, category.TenantId).
		Save(category).
		Error
}

func (r *CategoriesRepository) Delete(id string, tenantId string) error {
	return r.DB.Delete(&models.Categories{}, `id = ? AND tenant_id = ?`, id, tenantId).Error
}
