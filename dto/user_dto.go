package dto

type UserNih struct {
	ID         string  `json:"id"`
	FullName   string  `json:"fullName"`
	Email      string  `json:"email"`
	Phone      string  `json:"phone"`
	Address    *string `json:"address"`
	Ig         *string `json:"ig"`
	Fb         *string `json:"fb"`
	Avatar     *string `json:"avatar"`
	Disable    bool    `json:"disable"`
	TenantName string  `json:"tenantName"`
}

type UserResponse struct {
	Data       []UserNih `json:"data"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	PageSize   int       `json:"pageSize"`
	TotalPages int       `json:"totalPages"`
}
