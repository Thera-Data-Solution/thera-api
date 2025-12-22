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

// FindAll simple fetch (tanpa filter & pagination), masih disimpan jika dibutuhkan di tempat lain
func (r *TranslationRepository) FindAll(tenantId string) ([]models.Translation, error) {
	var translations []models.Translation
	err := r.DB.Where("tenant_id = ?", tenantId).Find(&translations).Error
	return translations, err
}

// FindFiltered menyediakan advanced filter + pagination
func (r *TranslationRepository) FindFiltered(tenantId string, filter models.TranslationFilter, page, limit int) ([]models.Translation, int64, error) {
	if page < 1 {
		page = 1
	}

	db := r.DB.Model(&models.Translation{}).Where("tenant_id = ?", tenantId)

	if filter.Locale != nil && *filter.Locale != "" {
		db = db.Where("locale = ?", *filter.Locale)
	}
	if filter.Namespace != nil && *filter.Namespace != "" {
		db = db.Where("namespace = ?", *filter.Namespace)
	}
	if filter.Key != nil && *filter.Key != "" {
		db = db.Where("key = ?", *filter.Key)
	}
	if filter.Value != nil && *filter.Value != "" {
		db = db.Where("value = ?", *filter.Value)
	}
	if filter.Search != nil && *filter.Search != "" {
		s := "%" + *filter.Search + "%"
		db = db.Where(r.DB.Where("key ILIKE ?", s).Or("value ILIKE ?", s).Or("namespace ILIKE ?", s))
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var translations []models.Translation
	if err := db.
		Order("namespace ASC, key ASC").
		Offset(offset).
		Limit(limit).
		Find(&translations).Error; err != nil {
		return nil, 0, err
	}

	return translations, total, nil
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
