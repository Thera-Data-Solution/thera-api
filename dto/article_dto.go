// dto/article_dto.go
package dto

import "time"

type CreateArticleDTO struct {
	Title       string  `json:"title" binding:"required"`
	Slug        string  `json:"slug" binding:"required"`
	CoverImage  *string `json:"coverImage"`
	Body        *string `json:"body"`
	Excerpt     *string `json:"excerpt"`
	IsPublished bool    `json:"isPublished"`
	ArticleType *string `json:"articleType" binding:"required"` // PAGE atau BLOG
}

type UpdateArticleDTO struct {
	Title       *string `json:"title"`
	Slug        *string `json:"slug"`
	CoverImage  *string `json:"coverImage"`
	Body        *string `json:"body"`
	Excerpt     *string `json:"excerpt"`
	IsPublished *bool   `json:"isPublished"`
	ArticleType *string `json:"articleType"`
}

// Response structure for paginated results
type ArticlePaginationResponse struct {
	Data       []ArticleResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"pageSize"`
	TotalPages int               `json:"totalPages"`
}

// ArticleResponse for returning article data (e.g. without tenantId)
type ArticleResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	CoverImage  *string   `json:"coverImage,omitempty"`
	Body        *string   `json:"body,omitempty"`
	Excerpt     *string   `json:"excerpt,omitempty"`
	IsPublished bool      `json:"isPublished"`
	ArticleType *string   `json:"articleType,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
