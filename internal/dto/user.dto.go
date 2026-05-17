package dto

type UserProfileResponse struct {
	Fullname *string `json:"fullname"`
	Email    string  `json:"email"`
	Photo    *string `json:"photo"`
}

type UserUpdateProfileRequest struct {
	Fullname    *string `json:"fullname"`
	PhoneNumber *string `json:"phone_number"`
	Photo       *string `json:"photo"`
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
	Pin string `json:"pin" binding:"required,len=6"`
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
	Pin string `json:"pin" binding:"required,len=6"`
}

type UserCheckPinResponse struct {
	IsValid bool `json:"is_valid"`
}
