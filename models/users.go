package models

import "time"

type RoleList string

const (
	USER  RoleList = "USER"
	ADMIN RoleList = "ADMIN"
	Slug  RoleList = "SU"
)

type User struct {
	ID       string   `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Avatar   string   `json:"avatar"`
	Email    string   `gorm:"unique" json:"email"`
	Password string   `json:"password"`
	FullName string   `json:"fullName"`
	Phone    string   `gorm:"unique" json:"phone"`
	Address  string   `json:"address"`
	Role     RoleList `gorm:"type:text;default:USER" json:"role"`
	Ig       string   `json:"ig"`
	Fb       string   `json:"fb"`
	Disable  bool     `gorm:"default:false" json:"disable"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
