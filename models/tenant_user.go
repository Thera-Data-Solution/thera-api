package models

type TenantUser struct {
	ID       string  `json:"id" gorm:"column:id"`
	Avatar   *string `json:"avatar,omitempty" gorm:"column:avatar"`
	Email    string  `json:"email" gorm:"column:email"`
	Password string  `json:"password,omitempty" gorm:"column:password"`
	FullName string  `json:"fullName" gorm:"column:fullName"`
	Role     string  `json:"role" gorm:"column:role"`
	TenantId string  `json:"tenantId" gorm:"column:tenantId"`
}
