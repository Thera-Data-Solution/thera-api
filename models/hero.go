package models

type Hero struct {
	ID          string  `json:"id" gorm:"column:id"`
	Title       string  `json:"title" gorm:"column:title"`
	Subtitle    *string `json:"subtitle,omitempty" gorm:"column:subtitle"`
	Description *string `json:"description,omitempty" gorm:"column:description"`
	Image       *string `json:"image,omitempty" gorm:"column:image"`
	ButtonText  *string `json:"buttonText,omitempty" gorm:"column:buttonText"`
	ButtonLink  *string `json:"buttonLink,omitempty" gorm:"column:buttonLink"`
	ThemeType   *string `json:"themeType,omitempty" gorm:"column:themeType"`
	IsActive    bool    `json:"isActive" gorm:"column:isActive"`
	TenantId    *string `json:"tenantId,omitempty" gorm:"column:tenantId"`
}

func (Hero) TableName() string {
	return "Hero"
}
