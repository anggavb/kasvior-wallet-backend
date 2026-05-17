package dto

import "time"

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Id          int        `json:"id"`
	Fullname    string     `json:"fullname"`
	Email       string     `json:"email"`
	PhoneNumber string     `json:"phone_number"`
	Photo       string     `json:"photo"`
	IsVerified  bool       `json:"is_verified"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Token       string     `json:"token,omitempty"`
}
