package models

import "time"

type User struct {
	ID        string     `json:"id"`
	Avatar    *string    `json:"avatar,omitempty"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	FullName  string     `json:"fullName"`
	Phone     *string    `json:"phone,omitempty"`
	Address   *string    `json:"address,omitempty"`
	IG        *string    `json:"ig,omitempty"`
	FB        *string    `json:"fb,omitempty"`
	Disable   bool       `json:"disable"`
	TenantId  *string    `json:"tenantId,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type TenantUser struct {
	ID       string  `json:"id"`
	Avatar   *string `json:"avatar,omitempty"`
	Email    string  `json:"email"`
	Password string  `json:"-"`
	FullName string  `json:"fullName"`
	Role     *string `json:"role,omitempty"`
	TenantId *string `json:"tenantId,omitempty"`
}

type Session struct {
	ID           string  `json:"id"`
	UserId       *string `json:"userId,omitempty"`
	TenantUserId *string `json:"tenantUserId,omitempty"`
	Token        string  `json:"token"`
	Device       *string `json:"device,omitempty"`
	IP           *string `json:"ip,omitempty"`
	ExpiresAt    string  `json:"expiresAt"`
	TenantId     *string `json:"tenantId,omitempty"`
}
