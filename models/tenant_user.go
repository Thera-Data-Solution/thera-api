package models

type TenantUser struct {
	ID       string  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Avatar   *string `json:"avatar,omitempty"`
	Email    string  `json:"email" gorm:"not null;index:idx_tenant_email,unique"`
	Password string  `json:"password,omitempty" gorm:"not null"`
	FullName string  `json:"fullName" gorm:"not null"`
	Role     string  `json:"role" gorm:"not null;default:'REGISTERED'"`
	TenantId string  `json:"tenantId" gorm:"not null;index:idx_tenant_email,unique"`
}
