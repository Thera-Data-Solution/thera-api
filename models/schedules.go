package models

import "time"

type Schedules struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	DateTime   time.Time `json:"dateTime" gorm:"not null;index"`
	CategoryId string    `json:"categoryId" gorm:"not null;index"`
	Status     string    `json:"status" gorm:"not null"`
	TenantId   *string   `json:"tenantId,omitempty" gorm:"index"`

	Categories Categories `json:"categories" gorm:"foreignKey:CategoryId;references:ID"`
}
