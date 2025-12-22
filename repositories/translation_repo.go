package repositories

import (
	"thera-api/models"

	"gorm.io/gorm"
)

type TranslationRepository struct {
	DB *gorm.DB
}

func (r *TranslationRepository) Create(translation *models.Translation) error {
	return r.DB.Create(translation).Error
}

func (r *TranslationRepository) FindAll(tenantId string) ([]models.Translation, error) {
	var translations []models.Translation
	err := r.DB.Where("tenant_id = ?", tenantId).Find(&translations).Error
	return translations, err
}

func (r *TranslationRepository) FindByID(id string, tenantId string) (*models.Translation, error) {
	var translation models.Translation
	err := r.DB.Where("tenant_id = ? AND id = ?", tenantId, id).First(&translation).Error
	if err != nil {
		return nil, err
	}
	return &translation, nil
}

func (r *TranslationRepository) Update(translation *models.Translation) error {
	return r.DB.Save(translation).Error
}

func (r *TranslationRepository) Delete(id string, tenantId string) error {
	return r.DB.Where("tenant_id = ? AND id = ?", tenantId, id).Delete(&models.Translation{}).Error
}
