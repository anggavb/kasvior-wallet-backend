package dto

type UserProfileResponse struct {
	Fullname *string `json:"fullname"`
	Email    string  `json:"email"`
	Photo    *string `json:"photo"`
}

type UserDashboardInformationResponse struct {
	Balance float64 `json:"balance"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}
