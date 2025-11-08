package models

import "time"

type Booked struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserId     string    `json:"userId" gorm:"not null;index"`
	ScheduleId string    `json:"scheduleId" gorm:"not null;index"`
	BookedAt   time.Time `json:"bookedAt" gorm:"autoCreateTime"`
	TenantId   *string   `json:"tenantId,omitempty" gorm:"index"`

	User     User      `json:"user" gorm:"foreignKey:UserId;references:ID"`
	Schedule Schedules `json:"schedule" gorm:"foreignKey:ScheduleId;references:ID"`
}
