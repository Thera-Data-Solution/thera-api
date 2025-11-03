package models

import "time"

type Gallery struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title       *string   `json:"title,omitempty"`
	Description *string   `json:"description,omitempty"`
	ImageUrl    string    `json:"imageUrl" gorm:"not null"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	TenantId    *string   `json:"tenantId,omitempty" gorm:"index"`
}
