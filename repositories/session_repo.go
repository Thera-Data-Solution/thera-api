package repositories

import (
	"thera-api/models"

	"gorm.io/gorm"
)

type SessionRepository struct {
	DB *gorm.DB
}

func (r *SessionRepository) CreateSession(s *models.Session) error {
	return r.DB.Create(s).Error
}

func (r *SessionRepository) DeleteByTenantUserId(s *models.Session) error {
	return r.DB.Where(`"tenantUserId" = ?`, s.TenantUserId).Delete(&models.Session{}).Error
}

func (r *SessionRepository) FindByToken(token string) (*models.Session, error) {
	var session models.Session
	err := r.DB.Where("token = ?", token).First(&session).Error

	return &session, err
}
