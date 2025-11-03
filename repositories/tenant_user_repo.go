package repositories

import (
	"thera-api/models"

	"gorm.io/gorm"
)

type TenantUserRepository struct {
	DB *gorm.DB
}

func (r *TenantUserRepository) FindByEmailAndTenant(email, tenantId string) (*models.TenantUser, error) {
	var u models.TenantUser
	err := r.DB.Where(`email = ? AND tenant_id = ?`, email, tenantId).First(&u).Error
	return &u, err
}

func (r *TenantUserRepository) Create(u *models.TenantUser) error {
	return r.DB.Create(u).Error
}

func (r *TenantUserRepository) FindByID(idToken string) (*models.TenantUser, error) {
	var user models.TenantUser
	err := r.DB.Where("id = ?", idToken).First(&user).Error
	return &user, err
}
