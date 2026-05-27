package dto

import "mime/multipart"

type UserProfileResponse struct {
	Fullname *string `json:"fullname"`
	Email    string  `json:"email"`
	Photo    *string `json:"photo"`
}

type UserUpdateProfileRequest struct {
	Fullname    *string               `form:"fullname" binding:"omitnil,min=3"`
	PhoneNumber *string               `form:"phone_number" binding:"omitnil,min=10,max=15"`
	Photo       *multipart.FileHeader `form:"photo" binding:"omitnil,image_max_size=2097152,image_type"` // 2MB
}

type UserUpdateProfileResponse struct {
	Fullname    *string `json:"fullname"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	Photo       *string `json:"photo"`
}

type UserUpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

type UserUpdatePinRequest struct {
	Pin string `json:"pin" binding:"required,len=6,numeric"`
}

type UserDashboardInformationResponse struct {
	Balance float64 `json:"balance"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

type UserTransactionReportResponse struct {
	Day     string  `json:"day"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

type UserCheckPinRequest struct {
	Pin           string `json:"pin" binding:"required,len=6,numeric"`
	TransactionId *int   `json:"transaction_id" binding:"omitnil,gt=0"`
}

type UserCheckPinResponse struct {
	IsValid bool `json:"is_valid"`
}

type TransactionReportQueryRequest struct {
	Duration string `form:"duration" binding:"oneof=7d"`
	Type     string `form:"type" binding:"oneof=all income expense"`
}
