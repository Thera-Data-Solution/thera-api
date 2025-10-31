package models

import "time"

type Schedules struct {
	ID         string    `json:"id" gorm:"column:id"`
	DateTime   time.Time `json:"dateTime" gorm:"column:dateTime"`
	CategoryId string    `json:"categoryId" gorm:"column:categoryId"`
	Status     string    `json:"status" gorm:"column:status"`
	TenantId   *string   `json:"tenantId,omitempty" gorm:"column:tenantId"`
}

func (Schedules) TableName() string {
	return "Schedules"
}
