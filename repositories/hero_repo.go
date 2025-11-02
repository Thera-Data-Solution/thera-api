package repositories

import (
	"thera-api/models"

	"gorm.io/gorm"
)

type HeroRepository struct {
	DB *gorm.DB
}

func (r *HeroRepository) Create(hero *models.Hero) error {
	return r.DB.Create(hero).Error
}

func (r *HeroRepository) FindAll(tenant string) ([]models.Hero, error) {
	var heroes []models.Hero
	err := r.DB.Where(`"tenantId" = ?`, tenant).Find(&heroes).Error
	return heroes, err
}

func (r *HeroRepository) FindByID(id string, tenant string) (*models.Hero, error) {
	var hero models.Hero
	err := r.DB.First(&hero, `id = ? AND "tenantId" = ?`, id, tenant).Error
	if err != nil {
		return nil, err
	}
	return &hero, nil
}

func (r *HeroRepository) Update(hero *models.Hero) error {
	return r.DB.
		Where(`id = ? AND "tenantId" = ?`, hero.ID, hero.TenantId).
		Save(hero).
		Error
}

func (r *HeroRepository) Delete(id string, tenantId string) error {
	return r.DB.Delete(&models.Hero{}, `id = ? AND "tenantId" = ?`, id, tenantId).Error
}
