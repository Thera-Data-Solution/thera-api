package models

import "time"

type ResetPasswordRequest struct {
	ID        string    `json:"id" gorm:"column:id"`
	UserId    string    `json:"userId" gorm:"column:userId"`
	Code      string    `json:"code" gorm:"column:code"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"column:expiresAt"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:createdAt"`
}

func (ResetPasswordRequest) TableName() string {
	return "ResetPasswordRequest"
}
