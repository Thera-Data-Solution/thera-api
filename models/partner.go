package models

type Partner struct {
	ID       string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Logo     string `json:"logo" gorm:"not null"`
	TenantId string `json:"tenantId" gorm:"not null;index"`
}
