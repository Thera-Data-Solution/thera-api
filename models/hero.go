package models

type Hero struct {
	ID          string  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title       string  `json:"title" gorm:"not null"`
	Subtitle    *string `json:"subtitle,omitempty"`
	Description *string `json:"description,omitempty"`
	Image       *string `json:"image,omitempty"`
	ButtonText  *string `json:"buttonText,omitempty"`
	ButtonLink  *string `json:"buttonLink,omitempty"`
	ThemeType   *string `json:"themeType,omitempty"`
	IsActive    bool    `json:"isActive" gorm:"default:true"`
	TenantId    *string `json:"tenantId,omitempty" gorm:"index"`
}
