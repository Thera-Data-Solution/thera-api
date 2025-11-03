package models

import "time"

type ResetPasswordRequest struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserId    string    `json:"userId" gorm:"not null;index"`
	Code      string    `json:"code" gorm:"not null"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
}
