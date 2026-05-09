package dto

import "time"

type ScheduleResponse struct {
	ID           string    `json:"id"`
	DateTime     time.Time `json:"dateTime"`
	CategoryId   string    `json:"categoryId"`
	CategoryName string    `json:"category"`
	Status       string    `json:"status"`
}

type ScheduleIdResponse struct {
	ID             string        `json:"id"`
	DateTime       time.Time     `json:"dateTime"`
	Name           string        `json:"name"`
	NameEn         string        `json:"nameEn"`
	Description    string        `json:"description"`
	DescriptionEn  string        `json:"descriptionEn"`
	Image          string        `json:"image"`
	Location       string        `json:"location"`
	Price          int           `json:"price"`
	Start          int           `json:"start"`
	End            int           `json:"end"`
	IsGroup        bool          `json:"isGroup"`
	IsFree         bool          `json:"isFree"`
	IsPayAsYouWish bool          `json:"isPayAsYouWish"`
	IsManual       bool          `json:"isManual"`
	Status         string        `json:"status"`
	CustomFields   []CustomField `json:"customFields"`
}
