package dto

type TestimoniResponse struct {
	ID        string `json:"id"`
	Testimoni string `json:"testimoni"`
	User      string `json:"user"`
	Event     string `json:"event"`
}

type TestimoniAdminResponse struct {
	ID        string `json:"id"`
	User      string `json:"user"`
	Image     string `json:"image"`
	Testimoni string `json:"testimoni"`
	Event     string `json:"event"`
	ShowTesti bool   `json:"showTesti"`
	Anonymous bool   `json:"anonymous"`
}
