package models

import (
	"time"

	"gorm.io/datatypes"
)

type Events struct {
	ID            string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name          string         `json:"name" gorm:"not null"`
	NameEn        string         `json:"nameEn,omitempty"`
	Description   *string        `json:"description,omitempty"`
	DescriptionEn *string        `json:"descriptionEn,omitempty"`
	Location      string         `json:"location,omitempty"`
	Image         *string        `json:"image,omitempty"`
	Slug          string         `json:"slug" gorm:"uniqueIndex;not null"`
	Price         float64        `json:"price" gorm:"default:0"`
	StartAt       time.Time      `json:"startAt" gorm:"not null"`
	EndAt         time.Time      `json:"endAt" gorm:"not null"`
	Capacity      int            `json:"capacity" gorm:"default:0"`
	Status        string         `json:"status" gorm:"default:'available'"`
	TenantId      *string        `json:"tenantId,omitempty" gorm:"index"`
	CustomFields  datatypes.JSON `json:"customFields,omitempty" gorm:"type:jsonb"`
}
