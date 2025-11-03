package repositories

import (
	"thera-api/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) FindByEmailAndTenant(email, tenantId string) (*models.User, error) {
	var user models.User
	err := r.DB.Where(`email = ? AND tenant_id = ?`, email, tenantId).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByPhoneAndTenant(phone, tenantId string) (*models.User, error) {
	var user models.User
	err := r.DB.Where(`phone = ? AND tenant_id = ?`, phone, tenantId).First(&user).Error
	return &user, err
}

func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) FindByID(idToken string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("id = ?", idToken).First(&user).Error
	return &user, err
}
