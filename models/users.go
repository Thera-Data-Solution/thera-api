package models

import (
	"github.com/google/uuid"
)

type RoleList string

const (
	USER  RoleList = "USER"
	ADMIN RoleList = "ADMIN"
	Slug  RoleList = "SU"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Avatar   *string   `json:"avatar"`
	Email    string    `gorm:"unique" json:"email"`
	Password string    `json:"password"`
	FullName string    `gorm:"column:\"fullName\"" json:"fullName"`
	Phone    string    `gorm:"unique" json:"phone"`
	Address  string    `json:"address"`
	Role     RoleList  `gorm:"type:text;default:USER" json:"role"`
	Ig       string    `json:"ig"`
	Fb       string    `json:"fb"`
	Disable  bool      `gorm:"default:false" json:"disable"`
}

func (User) TableName() string {
	return "User"
}
