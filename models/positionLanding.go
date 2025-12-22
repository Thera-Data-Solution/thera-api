package models

type PositionLanding struct {
	ID       string  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name     string  `json:"name" gorm:"not null"`
	Order    int     `json:"order"`
	IsActive bool    `json:"isActive" gorm:"default:true"`
	TenantId *string `json:"tenantId,omitempty" gorm:"index"`
}
