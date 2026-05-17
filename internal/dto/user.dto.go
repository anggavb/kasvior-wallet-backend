package dto

type UserProfileResponse struct {
	Fullname *string `json:"fullname"`
	Email    string  `json:"email"`
	Photo    *string `json:"photo"`
}
