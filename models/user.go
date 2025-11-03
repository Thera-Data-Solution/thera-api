package models

type User struct {
	ID       string  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Avatar   *string `json:"avatar,omitempty"`
	Email    string  `json:"email" gorm:"not null;index:idx_tenant_email,unique"`
	Password string  `json:"password,omitempty" gorm:"not null"`
	FullName string  `json:"fullName" gorm:"not null"`
	Phone    string  `json:"phone" gorm:"not null"`
	Address  *string `json:"address,omitempty"`
	Ig       *string `json:"ig,omitempty"`
	Fb       *string `json:"fb,omitempty"`
	Disable  bool    `json:"disable" gorm:"default:false"`
	TenantId string  `json:"tenantId" gorm:"not null;index:idx_tenant_email,unique"`
}
