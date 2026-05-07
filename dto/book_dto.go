package dto

import "time"

type GetAllBookingResponse struct {
	ID            string    `json:"id"`
	Avatar        *string   `json:"avatar,omitempty"`
	UserId        string    `json:"userId"`
	FullName      string    `json:"fullName"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	BookAt        time.Time `json:"bookAt"`
	ScheduleImage *string   `json:"scheduleImage,omitempty"`
	ScheduleId    *string   `json:"scheduleId,omitempty"`
	ScheduleName  *string   `json:"scheduleName,omitempty"`
	EventImage    *string   `json:"eventImage,omitempty"`
	EventId       *string   `json:"eventId,omitempty"`
	EventName     *string   `json:"eventName,omitempty"`
	Status        string    `json:"status"`
	// customAnswer []CustomAnswerResponse `json:"customAnswer,omitempty"`
}

// type CustomAnswerResponse struct {
// 	Question string `json:"question"`
// 	Answer   string `json:"answer"`
// }
