package models

import "time"

type Link struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"not null"`
	Value     string    `json:"value" gorm:"not null"`
	Type      string    `json:"type" gorm:"not null"`
	Icon      *string   `json:"icon,omitempty"`
	Order     *int      `json:"order,omitempty"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	TenantId  *string   `json:"tenantId,omitempty" gorm:"index"`
}
