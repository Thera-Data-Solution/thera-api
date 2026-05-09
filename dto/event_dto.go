package dto

import "time"

type EventResponse struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	NameEn        string        `json:"nameEn"`
	Description   *string       `json:"description,omitempty"`
	DescriptionEn *string       `json:"descriptionEn,omitempty"`
	Image         *string       `json:"image,omitempty"`
	Slug          string        `json:"slug"`
	Price         float64       `json:"price"`
	StartAt       time.Time     `json:"startAt"`
	EndAt         time.Time     `json:"endAt"`
	Capacity      int           `json:"capacity"`
	Status        string        `json:"status"`
	Location      string        `json:"location,omitempty"`
	CustomFields  []CustomField `json:"customFields"`
}
