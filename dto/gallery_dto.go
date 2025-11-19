package dto

type GalleryResponse struct {
	ID          string  `json:"id"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	ImageUrl    string  `json:"imageUrl"`
	CreatedAt   string  `json:"createdAt"`
}

type GalleryPaginationResponse struct {
	Data       []GalleryResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"pageSize"`
	TotalPages int               `json:"totalPages"`
}
