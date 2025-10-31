package models

import "time"

type Link struct {
	ID        string    `json:"id" gorm:"column:id"`
	Name      string    `json:"name" gorm:"column:name"`
	Value     string    `json:"value" gorm:"column:value"`
	Type      string    `json:"type" gorm:"column:type"`
	Icon      *string   `json:"icon,omitempty" gorm:"column:icon"`
	Order     *int      `json:"order,omitempty" gorm:"column:order"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:createdAt"`
	TenantId  *string   `json:"tenantId,omitempty" gorm:"column:tenantId"`
}

func (Link) TableName() string {
	return "Link"
}
