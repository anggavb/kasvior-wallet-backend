package dto

type TopupRequest struct {
	Amount          uint   `json:"amount" binding:"required,gt=0"`
	TypeTransaction string `json:"type" binding:"required,oneof=topup transfer receiver"`
	PaymentMethodId int    `json:"payment_method_id" binding:"required,gt=0"`
	Discount        *int   `json:"discount" binding:"required,gte=0"`
	Tax             *int   `json:"tax" binding:"required,gte=0"`
	SubTotal        *int   `json:"sub_total" binding:"required,gte=0"`
}

type TopupResponse struct {
	Amount        uint   `json:"amount"`
	PaymentMethod string `json:"payment_method"`
	Discount      int    `json:"discount"`
	Tax           int    `json:"tax"`
	SubTotal      int    `json:"sub_total"`
}
