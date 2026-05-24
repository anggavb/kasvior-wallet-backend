package apperrors

import "errors"

var (
	ErrTokenAlreadyExpired       = errors.New("token already expired")
	ErrInvalidPasswordResetToken = errors.New("invalid password reset token")
)
