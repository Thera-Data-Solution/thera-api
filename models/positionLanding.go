package models

type PositionLanding struct {
	ID       string  `json:"id" gorm:"column:id"`
	Name     string  `json:"name" gorm:"column:name"`
	Order    int     `json:"order" gorm:"column:order"`
	IsActive bool    `json:"isActive" gorm:"column:isActive"`
	TenantId *string `json:"tenantId,omitempty" gorm:"column:tenantId"`
}

func (PositionLanding) TableName() string {
	return "PositionLanding"
}
