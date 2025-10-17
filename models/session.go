package models

import "time"

type Session struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    string    `json:"userId"`
	Token     string    `gorm:"unique" json:"token"`
	Device    string    `json:"device"`
	IP        string    `json:"ip"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}
