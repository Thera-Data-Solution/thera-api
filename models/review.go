package models

import "time"

type Review struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	BookingId   string    `json:"bookingId" gorm:"uniqueIndex;not null"`
	UserId      string    `json:"userId" gorm:"not null"`
	TargetId    string    `json:"targetId" gorm:"index"`
	TargetType  string    `json:"targetType" gorm:"index"`
	Content     string    `json:"content"`
	IsApproved  bool      `json:"isApproved" gorm:"default:false"`
	IsAnonymous bool      `json:"isAnonymous" gorm:"default:false"`
	CreatedAt   time.Time `json:"createdAt"`
}
