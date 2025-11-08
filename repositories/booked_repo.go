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

	query := r.DB.Model(&models.Booked{}).
		Where("tenant_id = ?", tenantId).
		Preload("User").
		Preload("Schedule").
		Preload("Schedule.Categories").
		Order("booked_at DESC")

	// hitung total untuk pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ambil data sesuai pagination
	if err := query.Limit(limit).Offset(offset).Find(&booked).Error; err != nil {
		return nil, 0, err
	}

	return booked, total, nil
}

func (r *BookedRepository) GetByUser(tenantId string, userId string) ([]models.Booked, error) {
	var booked []models.Booked
	err := r.DB.Where(`tenant_id = ? AND "userId" = ?`, tenantId, userId).Find(&booked).Error
	return booked, err
}

func (r *BookedRepository) GetById(id string, tenantId string) (*models.Booked, error) {
	var booked models.Booked
	err := r.DB.Where(`id = ? AND tenant_id = ?`, id, tenantId).First(&booked).Error
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
