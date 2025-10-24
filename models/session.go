package models

import "time"

type Session struct {
	ID           string    `json:"id" gorm:"column:id"`
	Token        string    `json:"token" gorm:"column:token"`
	UserId       *string   `json:"userId,omitempty" gorm:"column:userId"`
	TenantUserId *string   `json:"tenantUserId,omitempty" gorm:"column:tenantUserId"`
	Device       *string   `json:"device,omitempty" gorm:"column:device"`
	IP           *string   `json:"ip,omitempty" gorm:"column:ip"`
	ExpiresAt    time.Time `json:"expiresAt" gorm:"column:expiresAt"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:createdAt"`
	TenantId     string    `json:"tenantId" gorm:"column:tenantId"`
}
