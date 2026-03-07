package dto

import "time"

type BookingGetResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"`
	Review    *string   `json:"review,omitempty"`
	Anonymous bool      `json:"anonymous"`
}
