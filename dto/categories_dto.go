package dto

type CategoriesResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	NameEn         string `json:"nameEn"`
	Description    string `json:"description"`
	DescriptionEn  string `json:"descriptionEn"`
	Slug           string `json:"slug"`
	Image          string `json:"image"`
	Start          int    `json:"start"`
	End            int    `json:"end"`
	Location       string `json:"location"`
	Price          int    `json:"price"`
	IsGroup        bool   `json:"isGroup"`
	IsFree         bool   `json:"isFree"`
	IsPayAsYouWish bool   `json:"isPayAsYouWish"`
	IsManual       bool   `json:"isManual"`
}

type AdminCategoriesResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	NameEn         string `json:"nameEn"`
	Description    string `json:"description"`
	DescriptionEn  string `json:"descriptionEn"`
	Slug           string `json:"slug"`
	Image          string `json:"image"`
	Start          int    `json:"start"`
	End            int    `json:"end"`
	Location       string `json:"location"`
	Price          int    `json:"price"`
	IsGroup        bool   `json:"isGroup"`
	IsFree         bool   `json:"isFree"`
	IsPayAsYouWish bool   `json:"isPayAsYouWish"`
	IsManual       bool   `json:"isManual"`
	Disable        bool   `json:"disable"`
}
