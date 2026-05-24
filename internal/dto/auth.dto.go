package dto

import "time"

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
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
