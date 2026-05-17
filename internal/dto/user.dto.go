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

type UserTransactionReportResponse struct {
	Day     string  `json:"day"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}
