package models

import "time"

type Booked struct {
	ID         string    `json:"id" gorm:"column:id"`
	UserId     string    `json:"userId" gorm:"column:userId"`
	ScheduleId string    `json:"scheduleId" gorm:"column:scheduleId"`
	BookedAt   time.Time `json:"bookedAt" gorm:"column:bookedAt"`
	TenantId   *string   `json:"tenantId,omitempty" gorm:"column:tenantId"`
}

func (Booked) TableName() string {
	return "Booked"
}
