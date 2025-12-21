package models

import "time"

type PasswordResetToken struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null"`
	Email     string    `json:"email" gorm:"not null;index"`
	IPAddress string    `json:"ipAddress" gorm:"not null;index"`
	UserType  string    `json:"userType" gorm:"not null"` // "user" or "admin"
	TenantId  string    `json:"tenantId" gorm:"not null;index"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"not null"`
	Used      bool      `json:"used" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

