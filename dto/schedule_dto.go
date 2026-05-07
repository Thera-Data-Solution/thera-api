package dto

import "time"

type ScheduleResponse struct {
	ID           string    `json:"id"`
	DateTime     time.Time `json:"dateTime"`
	CategoryId   string    `json:"categoryId"`
	CategoryName string    `json:"category"`
	Status       string    `json:"status"`
}
