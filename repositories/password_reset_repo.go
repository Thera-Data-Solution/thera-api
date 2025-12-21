package repositories

import (
	"thera-api/models"
	"time"

	"gorm.io/gorm"
)

type PasswordResetRepository struct {
	DB *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) *PasswordResetRepository {
	return &PasswordResetRepository{DB: db}
}

func (r *PasswordResetRepository) Create(token *models.PasswordResetToken) error {
	return r.DB.Create(token).Error
}

func (r *PasswordResetRepository) FindByToken(token string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	err := r.DB.Where("token = ? AND used = false AND expires_at > ?", token, time.Now()).First(&resetToken).Error
	return &resetToken, err
}

func (r *PasswordResetRepository) MarkAsUsed(token string) error {
	return r.DB.Model(&models.PasswordResetToken{}).
		Where("token = ?", token).
		Update("used", true).Error
}

// CountRequestsByIP counts reset requests from an IP in the last 24 hours
func (r *PasswordResetRepository) CountRequestsByIP(ipAddress string) (int64, error) {
	var count int64
	oneDayAgo := time.Now().Add(-24 * time.Hour)
	err := r.DB.Model(&models.PasswordResetToken{}).
		Where("ip_address = ? AND created_at > ?", ipAddress, oneDayAgo).
		Count(&count).Error
	return count, err
}

// GetLastRequestTime gets the time of the last request from an IP
func (r *PasswordResetRepository) GetLastRequestTime(ipAddress string) (*time.Time, error) {
	var resetToken models.PasswordResetToken
	err := r.DB.Where("ip_address = ?", ipAddress).
		Order("created_at DESC").
		First(&resetToken).Error
	
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &resetToken.CreatedAt, nil
}

// CleanupExpiredTokens removes expired tokens (optional, for maintenance)
func (r *PasswordResetRepository) CleanupExpiredTokens() error {
	return r.DB.Where("expires_at < ? OR used = true", time.Now()).
		Delete(&models.PasswordResetToken{}).Error
}



