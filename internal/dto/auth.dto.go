package dto

import "time"

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"anggavb8@gmail.com"`
	Password string `json:"password" binding:"required,min=8" example:"secreto123"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"anggavb8@gmail.com"`
	Password string `json:"password" binding:"required" example:"secreto123"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"anggavb8@gmail.com"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8" example:"secreto123"`
}

type AuthResponse struct {
	Id          int        `json:"id,omitempty"`
	Fullname    string     `json:"fullname,omitempty"`
	Email       string     `json:"email,omitempty"`
	PhoneNumber string     `json:"phone_number,omitempty"`
	Photo       string     `json:"photo,omitempty"`
	HasPin      *bool      `json:"has_pin,omitempty"`
	IsVerified  bool       `json:"is_verified,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	Token       string     `json:"token,omitempty"`
}
