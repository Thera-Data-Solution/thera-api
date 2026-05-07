package repositories

import (
	"thera-api/models"

	"gorm.io/gorm"
)

type BookedRepository struct {
	DB *gorm.DB
}

func NewBookedRepository(db *gorm.DB) *BookedRepository {
	return &BookedRepository{DB: db}
}

func (r *BookedRepository) Create(booked *models.Booked) error {
	return r.DB.Create(booked).Error
}

func (r *BookedRepository) GetAll(tenantId string, limit, offset int) ([]models.Booked, int64, error) {
	var booked []models.Booked
	var total int64

	query := r.DB.Model(&models.Booked{}).Where("tenant_id = ?", tenantId)
	query.Count(&total)

	err := query.
		Preload("User").
		Preload("Schedule.Categories").
		Preload("Event").
		Limit(limit).
		Offset(offset).
		Order("booked_at DESC").
		Find(&booked).Error

	return booked, total, err
}

func (r *BookedRepository) GetByIDWithDetails(id, tenantId string) (*models.Booked, error) {
	var booked models.Booked
	err := r.DB.
		Preload("User").
		Preload("Schedule.Categories").
		Preload("Event").
		Where("id = ? AND tenant_id = ?", id, tenantId).
		First(&booked).Error

	return &booked, err
}

func (r *BookedRepository) GetLatestByUserAndEvent(
	userId, eventId, tenantId string,
) (*models.Booked, error) {

	var booked models.Booked

	err := r.DB.
		Preload("User").
		Preload("Event").
		Where(`
			user_id = ? 
			AND event_id = ? 
			AND tenant_id = ?
			AND status IN ('PENDING', 'CONFIRMED')
			`, userId, eventId, tenantId).
		Order("booked_at DESC").
		First(&booked).Error

	if err != nil {
		return nil, err
	}

	return &booked, nil
}

func (r *BookedRepository) UpdateStatusBySchedule(scheduleId, oldStatus, newStatus string) error {
	return r.DB.Model(&models.Booked{}).
		Where("schedule_id = ? AND status = ?", scheduleId, oldStatus).
		Update("status", newStatus).Error
}

func (r *BookedRepository) GetLatestByUserAndEventWithoutStatus(
	userId, eventId, tenantId string,
) (*models.Booked, error) {

	var booked models.Booked

	err := r.DB.
		Preload("User").
		Preload("Event").
		Where(`
			user_id = ? 
			AND event_id = ? 
			AND tenant_id = ?
			`, userId, eventId, tenantId).
		Order("booked_at DESC").
		First(&booked).Error

	if err != nil {
		return nil, err
	}

	return &booked, nil
}

func (r *BookedRepository) GetLatestByUserAndSchedule(
	userId, scheduleId, tenantId string,
) (*models.Booked, error) {

	var booked models.Booked

	err := r.DB.
		Where(`
			user_id = ? 
			AND schedule_id = ? 
			AND tenant_id = ?
			AND status IN ('PENDING', 'CONFIRMED')
			`, userId, scheduleId, tenantId).
		Order("booked_at DESC").
		First(&booked).Error

	if err != nil {
		return nil, err
	}

	return &booked, nil
}

func (r *BookedRepository) GetLatestByUserAndScheduleWithoutStatus(
	userId, scheduleId, tenantId string,
) (*models.Booked, error) {

	var booked models.Booked

	err := r.DB.
		Where(`
			user_id = ? 
			AND schedule_id = ? 
			AND tenant_id = ?
			`, userId, scheduleId, tenantId).
		Order("booked_at DESC").
		First(&booked).Error

	if err != nil {
		return nil, err
	}

	return &booked, nil
}

func (r *BookedRepository) GetById(id string, tenantId string) (*models.Booked, error) {
	var booked models.Booked
	err := r.DB.
		Where(`id = ? AND tenant_id = ?`, id, tenantId).
		Preload("User").
		Preload("Schedule").
		Preload("Schedule.Categories").
		First(&booked).Error
	return &booked, err
}

func (r *BookedRepository) Update(booked *models.Booked) error {
	return r.DB.
		Where(`id = ? AND tenant_id = ?`, booked.ID, booked.TenantId).
		Save(booked).
		Error
}

func (r *BookedRepository) Delete(id string, tenantId string) error {
	return r.DB.
		Where(`id = ? AND tenant_id = ?`, id, tenantId).
		Delete(&models.Booked{}).
		Error
}
