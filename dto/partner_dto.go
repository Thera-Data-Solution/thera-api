package dto

type PartnerResponse struct {
	ID       string  `json:"id"`
	Logo     *string `json:"logo"`
	TenantId *string `json:"tenantId"`
}

type PartnerPaginationResponse struct {
	Data       []PartnerResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"pageSize"`
	TotalPages int               `json:"totalPages"`
}
