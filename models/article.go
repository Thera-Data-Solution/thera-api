package models

import "time"

type Article struct {
	ID          string    `json:"id" gorm:"column:id"`
	Title       string    `json:"title" gorm:"column:title"`
	Slug        string    `json:"slug" gorm:"column:slug"`
	CoverImage  *string   `json:"coverImage,omitempty" gorm:"column:coverImage"`
	Body        *string   `json:"body,omitempty" gorm:"column:body;type:json"`
	Excerpt     *string   `json:"excerpt,omitempty" gorm:"column:excerpt"`
	IsPublished bool      `json:"isPublished" gorm:"column:isPublished"`
	ArticleType *string   `json:"articleType,omitempty" gorm:"column:articleType"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updatedAt"`
	TenantId    *string   `json:"tenantId,omitempty" gorm:"column:tenantId"`
}

func (Article) TableName() string {
	return "Article"
}
