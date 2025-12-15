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
	err := r.DB.
		Preload("Categories").
		Where(`tenant_id = ?`, tenantId).
		Find(&schedules).Error
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

func (r *SchedulesRepository) FindByCatID(id string, tenantId string, date string) ([]models.Schedules, error) {
	var schedules []models.Schedules

	query := r.DB.
		Preload("Categories").
		Where("tenant_id = ? AND category_id = ?", tenantId, id)

	if date != "" {
		query = query.Where(
			"DATE(date_time) = ?",
			date,
		)
	} else {
		query = query.Where("date_time > NOW()")
	}

	err := query.Find(&schedules).Error
	if err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *SchedulesRepository) Update(schedule *models.Schedules) error {
	return r.DB.Where(`tenant_id = ?`, schedule.TenantId).Save(schedule).Error
}

func (r *SchedulesRepository) Delete(id string, tenantId string) error {
	return r.DB.Where(`tenant_id = ? AND id = ?`, tenantId, id).Delete(&models.Schedules{}).Error
}
