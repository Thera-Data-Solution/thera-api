package models

import (
	"time"

	"gorm.io/datatypes"
)

type Booked struct {
	ID     string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserId string `json:"userId" gorm:"not null;index"`

	ScheduleId *string `json:"scheduleId" gorm:"index"`
	EventId    *string `json:"eventId" gorm:"index"`

	Status   string    `json:"status" gorm:"default:'PENDING'"`
	BookedAt time.Time `json:"bookedAt" gorm:"autoCreateTime"`

	// Akan di delete setelah migrasi production, karena sudah tidak dipakai
	Testimoni *string `json:"testimoni,omitempty"`
	ShowTesti *bool   `json:"showTesti,omitempty"`
	Anonymous *bool   `json:"anonymous,omitempty"`

	TenantId *string `json:"tenantId,omitempty" gorm:"index"`

	CustomAnswer datatypes.JSON `json:"customAnswers,omitempty" gorm:"type:jsonb"`

	User     User       `json:"user" gorm:"foreignKey:UserId;references:ID"`
	Schedule *Schedules `json:"schedule" gorm:"foreignKey:ScheduleId;references:ID"`
	Event    *Events    `json:"event" gorm:"foreignKey:EventId;references:ID"`
	Review   *Review    `json:"review" gorm:"foreignKey:BookingId;references:ID"`
}
