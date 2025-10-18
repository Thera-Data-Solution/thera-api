package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"column:\"userId\"" json:"userId"` // << ini penting
	Token     string    `gorm:"unique" json:"token"`
	Device    string    `json:"device"`
	IP        string    `json:"ip"`
	ExpiresAt time.Time `gorm:"column:\"expiresAt\"" json:"expiresAt"`
	CreatedAt time.Time `gorm:"column:\"createdAt\";autoCreateTime" json:"createdAt"`
}

func (Session) TableName() string {
	return "Session"
}
