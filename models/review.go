package models

import "time"

type Review struct {
	ID          string    `json:"id" gorm:"column:id"`
	UserId      string    `json:"userId" gorm:"column:userId"`
	BookedId    string    `json:"bookedId" gorm:"column:bookedId"`
	IsAnonymous bool      `json:"isAnonymous" gorm:"column:isAnonymous"`
	Comment     *string   `json:"comment,omitempty" gorm:"column:comment"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:createdAt"`
	TenantId    *string   `json:"tenantId,omitempty" gorm:"column:tenantId"`
}

func (Review) TableName() string {
	return "Review"
}
