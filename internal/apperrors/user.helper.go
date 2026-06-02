package apperrors

import "errors"

var (
	ErrPinNotSet       = errors.New("pin not set")
	ErrInvalidPin      = errors.New("invalid pin")
	ErrInvalidPassword = errors.New("invalid password")
	ErrEmailNotFound   = errors.New("email not found")
)
