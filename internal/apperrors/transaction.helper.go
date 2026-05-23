package apperrors

import "errors"

var (
	InvalidSubtotal          = errors.New("invalid calculate subtotal")
	InvalidPaymentMethodType = errors.New("invalid payment method type")
)
