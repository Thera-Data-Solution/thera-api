package models

import "time"

type Tenant struct {
	ID        string    `json:"id" gorm:"column:id"`
	Name      string    `json:"name" gorm:"column:name"`
	Logo      *string   `json:"logo,omitempty" gorm:"column:logo"`
	IsActive  bool      `json:"isActive" gorm:"column:isActive"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updatedAt"`
}

func (Tenant) TableName() string {
	return "Tenant"
}
