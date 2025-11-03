package repositories

import (
	"thera-api/models"

	"gorm.io/gorm"
)

type SchedulesRepository struct {
	DB *gorm.DB
}

func (r *SchedulesRepository) Create(schedule *models.Schedules) error {
	return r.DB.Where(`tenant_id = ?`).Create(schedule).Error
}

func (r *SchedulesRepository) FindAll(tenantId string) ([]models.Schedules, error) {
	var schedules []models.Schedules
	err := r.DB.Where(`tenant_id = ?`, tenantId).Find(&schedules).Error
	return schedules, err
}

func (r *SchedulesRepository) FindByID(id string, tenantId string) (*models.Schedules, error) {
	var schedule models.Schedules
	err := r.DB.Where(`tenant_id = ? AND id = ?`, tenantId, id).First(&schedule).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *SchedulesRepository) Update(schedule *models.Schedules) error {
	return r.DB.Where(`tenant_id = ?`, schedule.TenantId).Save(schedule).Error
}

func (r *SchedulesRepository) Delete(id string, tenantId string) error {
	return r.DB.Where(`tenant_id = ? AND id = ?`, tenantId, id).Delete(&models.Schedules{}).Error
}
