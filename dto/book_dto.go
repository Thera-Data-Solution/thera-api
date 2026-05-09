package dto

import (
	"time"
)

type GetAllBookingResponse struct {
	ID            string                 `json:"id"`
	Avatar        *string                `json:"avatar,omitempty"`
	UserId        string                 `json:"userId"`
	FullName      string                 `json:"fullName"`
	Email         string                 `json:"email"`
	Phone         string                 `json:"phone"`
	BookAt        time.Time              `json:"bookAt"`
	ScheduleImage *string                `json:"scheduleImage,omitempty"`
	ScheduleId    *string                `json:"scheduleId,omitempty"`
	ScheduleName  *string                `json:"scheduleName,omitempty"`
	ScheduleDate  *time.Time             `json:"scheduleDate,omitempty"`
	EventImage    *string                `json:"eventImage,omitempty"`
	EventId       *string                `json:"eventId,omitempty"`
	EventName     *string                `json:"eventName,omitempty"`
	EventDate     *time.Time             `json:"eventDate,omitempty"`
	Status        string                 `json:"status"`
	CustomAnswer  []CustomAnswerResponse `json:"customAnswer,omitempty"`
}

type CustomAnswerResponse struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
