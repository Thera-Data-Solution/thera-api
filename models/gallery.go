package models

import "time"

type Gallery struct {
	ID          string    `json:"id" gorm:"column:id"`
	Title       *string   `json:"title,omitempty" gorm:"column:title"`
	Description *string   `json:"description,omitempty" gorm:"column:description"`
	ImageUrl    string    `json:"imageUrl" gorm:"column:imageUrl"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:createdAt"`
	TenantId    *string   `json:"tenantId,omitempty" gorm:"column:tenantId"`
}

func (Gallery) TableName() string {
	return "Gallery"
}
