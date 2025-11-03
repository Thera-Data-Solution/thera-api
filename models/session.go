package models

import "time"

type Session struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Token        string    `json:"token" gorm:"uniqueIndex;not null"`
	UserId       *string   `json:"userId,omitempty"`
	TenantUserId *string   `json:"tenantUserId,omitempty"`
	Device       *string   `json:"device,omitempty"`
	IP           *string   `json:"ip,omitempty"`
	ExpiresAt    time.Time `json:"expiresAt" gorm:"not null"`
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime"`
	TenantId     string    `json:"tenantId" gorm:"index;not null"`
}
