package apperrors

import "errors"

var (
	InvalidSubtotal          = errors.New("invalid calculate subtotal")
	InvalidPaymentMethodType = errors.New("invalid payment method type")
	ErrInsufficientBalance   = errors.New("insufficient balance")
	ErrInvalidRecipient      = errors.New("invalid recipient wallet")
	ErrSelfTransfer          = errors.New("cannot transfer to own wallet")
	ErrTransactionNotFound   = errors.New("transaction not found")
	ErrTransactionFinalized  = errors.New("transaction already finalized")
)
