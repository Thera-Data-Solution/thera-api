package repositories

import (
	"thera-api/models"

	"gorm.io/gorm"
)

type HeroRepository struct {
	DB *gorm.DB
}

func (r *HeroRepository) Create(tenantId *string, hero *models.Hero) error {
	var existingHero models.Hero
	if err := r.DB.Where("tenant_id = ? ", tenantId).First(&existingHero).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return r.DB.Create(hero).Error
		} else {
			return err
		}
	} else {
		return r.DB.Model(&existingHero).Updates(hero).Error
	}
}

func (r *HeroRepository) FindAll(tenant string) (models.Hero, error) {
	var heroes models.Hero
	err := r.DB.Where(`tenant_id = ?`, tenant).First(&heroes).Error
	return heroes, err
}

func (r *HeroRepository) FindByID(id string, tenant string) (*models.Hero, error) {
	var hero models.Hero
	err := r.DB.First(&hero, `id = ? AND tenant_id = ?`, id, tenant).Error
	if err != nil {
		return nil, err
	}
	return &hero, nil
}
