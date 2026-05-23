package dto

import "time"

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Id          int        `json:"id,omitempty"`
	Fullname    string     `json:"fullname,omitempty"`
	Email       string     `json:"email"`
	PhoneNumber string     `json:"phone_number,omitempty"`
	Photo       string     `json:"photo,omitempty"`
	IsVerified  bool       `json:"is_verified,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Token       string     `json:"token,omitempty"`
}
