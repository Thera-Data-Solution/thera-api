package models

import "time"

type Review struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserId      string    `json:"userId" gorm:"not null;index"`
	BookedId    string    `json:"bookedId" gorm:"not null;index"`
	IsAnonymous bool      `json:"isAnonymous" gorm:"default:false"`
	Comment     *string   `json:"comment,omitempty"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	TenantId    *string   `json:"tenantId,omitempty" gorm:"index"`
}
