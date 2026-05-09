package models

type Hero struct {
	ID            string  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title         string  `json:"title" gorm:"not null"`
	TitleEn       *string `json:"titleEn,omitempty"`
	Subtitle      *string `json:"subtitle,omitempty"`
	SubtitleEn    *string `json:"subtitleEn,omitempty"`
	Description   *string `json:"description,omitempty"`
	DescriptionEn *string `json:"descriptionEn,omitempty"`
	Image         *string `json:"image,omitempty"`
	ButtonText    *string `json:"buttonText,omitempty"`
	ButtonTextEn  *string `json:"buttonTextEn,omitempty"`
	ButtonLink    *string `json:"buttonLink,omitempty"`
	ThemeType     *string `json:"themeType,omitempty"`
	IsActive      bool    `json:"isActive" gorm:"default:true"`
	TenantId      *string `json:"tenantId,omitempty" gorm:"index"`
}
