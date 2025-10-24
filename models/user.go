package models

type User struct {
	ID       string  `json:"id" gorm:"column:id"`
	Avatar   *string `json:"avatar,omitempty" gorm:"column:avatar"`
	Email    string  `json:"email" gorm:"column:email"`
	Password string  `json:"password,omitempty" gorm:"column:password"`
	FullName string  `json:"fullName" gorm:"column:fullName"`
	Phone    string  `json:"phone" gorm:"column:phone"`
	Address  *string `json:"address,omitempty" gorm:"column:address"`
	Ig       *string `json:"ig,omitempty" gorm:"column:ig"`
	Fb       *string `json:"fb,omitempty" gorm:"column:fb"`
	Disable  bool    `json:"disable" gorm:"column:disable"`
	TenantId string  `json:"tenantId" gorm:"column:tenantId"`
}
