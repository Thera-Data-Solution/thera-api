package dto

import "time"

type ReviewUpsertRequest struct {
	BookingId   string `json:"bookingId" binding:"required"`
	Content     string `json:"content" binding:"required"`
	IsAnonymous bool   `json:"isAnonymous"`
}

type PublicReviewResponse struct {
	UserName string    `json:"userName"`
	Content  string    `json:"content"`
	Type     string    `json:"type"`
	Target   string    `json:"targetName"`
	Date     time.Time `json:"date"`
}

type UserHistoryResponse struct {
	BookingId    string      `json:"bookingId"`
	Status       string      `json:"status"`
	BookedAt     time.Time   `json:"bookedAt"`
	ItemName     string      `json:"itemName"`
	ItemImage    string      `json:"itemImage"`
	ScheduleDate time.Time   `json:"scheduleDate"`
	Review       *ReviewInfo `json:"review"`
}

type ReviewInfo struct {
	Content     string `json:"content"`
	IsAnonymous bool   `json:"isAnonymous"`
}

type AdminReviewResponse struct {
	ID          string    `json:"id"`
	BookingId   string    `json:"bookingId"`
	UserName    string    `json:"userName"`
	TargetName  string    `json:"targetName"`
	TargetType  string    `json:"targetType"`
	Content     string    `json:"content"`
	IsApproved  bool      `json:"isApproved"`
	IsAnonymous bool      `json:"isAnonymous"`
	CreatedAt   time.Time `json:"createdAt"`
}

type AdminUpdateReviewRequest struct {
	Content    *string `json:"content"`
	IsApproved *bool   `json:"isApproved"`
}
