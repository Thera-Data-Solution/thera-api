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

func (r *BookedRepository) GetAll(tenantId string) ([]models.Booked, error) {
	var booked []models.Booked
	err := r.DB.Where(`"tenantId" = ?`, tenantId).Find(&booked).Error
	return booked, err
}

func (r *BookedRepository) GetByUser(tenantId string, userId string) ([]models.Booked, error) {
	var booked []models.Booked
	err := r.DB.Where(`"tenantId" = ? AND "userId" = ?`, tenantId, userId).Find(&booked).Error
	return booked, err
}

func (r *BookedRepository) GetById(id string, tenantId string) (*models.Booked, error) {
	var booked models.Booked
	err := r.DB.Where(`id = ? AND "tenantId" = ?`, id, tenantId).First(&booked).Error
	return &booked, err
}

func (r *BookedRepository) Update(booked *models.Booked) error {
	return r.DB.
		Where(`id = ? AND "tenantId" = ?`, booked.ID, booked.TenantId).
		Save(booked).
		Error
}

func (r *BookedRepository) Delete(id string, tenantId string) error {
	return r.DB.
		Where(`id = ? AND "tenantId" = ?`, id, tenantId).
		Delete(&models.Booked{}).
		Error
}
