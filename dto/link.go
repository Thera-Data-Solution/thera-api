package dto

type CreateLinkDTO struct {
	Name  string  `json:"name" binding:"required"`
	Value string  `json:"value" binding:"required"`
	Type  string  `json:"type" binding:"required"`
	Icon  *string `json:"icon,omitempty"`
}

type UpdateLinkDTO struct {
	Name  *string `json:"name,omitempty"`
	Value *string `json:"value,omitempty"`
	Type  *string `json:"type,omitempty"`
	Icon  *string `json:"icon,omitempty"`
	Order *int    `json:"order,omitempty"`
}

type UpdateLinkOrderDTO struct {
	Order int `json:"order"`
}
