package repositories

import (
	"thera-api/models"

	"gorm.io/gorm"
)

type SettingRepo interface {
	FindAll() ([]models.Setting, error)
	FindById(id string) (*models.Setting, error)
	FindByTenantId(tenantId string) (*models.Setting, error)
	Upsert(setting *models.Setting) (*models.Setting, error)
	Delete(id string) error
}

type settingRepo struct {
	DB *gorm.DB
}

func NewSettingRepo(db *gorm.DB) SettingRepo {
	return &settingRepo{DB: db}
}

func (r *settingRepo) FindAll() ([]models.Setting, error) {
	var settings []models.Setting

	if err := r.DB.Find(&settings).Error; err != nil {
		return nil, err
	}

	return settings, nil
}

func (r *settingRepo) FindById(id string) (*models.Setting, error) {
	var setting models.Setting

	if err := r.DB.Where("id = ?", id).First(&setting).Error; err != nil {
		return nil, err
	}

	return &setting, nil
}

func (r *settingRepo) FindByTenantId(tenantId string) (*models.Setting, error) {
	var setting models.Setting

	if err := r.DB.Where("tenant_id = ?", tenantId).First(&setting).Error; err != nil {
		return nil, err
	}

	return &setting, nil
}

func (r *settingRepo) Upsert(setting *models.Setting) (*models.Setting, error) {
    var existingSetting models.Setting
    if err := r.DB.Where("tenant_id = ?", setting.TenantId).First(&existingSetting).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            // Create new setting
            if err := r.DB.Create(setting).Error; err != nil {
                return nil, err
            }
            return setting, nil
        } else {
            return nil, err
        }
    } else {
        // Update existing setting
        if err := r.DB.Model(&existingSetting).Updates(setting).Error; err != nil {
            return nil, err
        }
        return setting, nil
    }
}

func (r *settingRepo) Delete(id string) error {
	return r.DB.Delete(&models.Setting{}, "id = ?", id).Error
}