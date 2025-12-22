package models

import "time"

type Article struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title       string    `json:"title" gorm:"not null"`
	Slug        string    `json:"slug" gorm:"uniqueIndex;not null"`
	CoverImage  *string   `json:"coverImage,omitempty"`
	Body        *string   `json:"body,omitempty" gorm:"type:text"`
	Excerpt     *string   `json:"excerpt,omitempty"`
	IsPublished bool      `json:"isPublished" gorm:"default:false"`
	ArticleType *string   `json:"articleType,omitempty"`
	TenantId    *string   `json:"tenantId,omitempty" gorm:"index"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
