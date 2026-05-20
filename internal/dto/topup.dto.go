package dto

type TopupRequest struct {
	Amount          uint `json:"amount" binding:"required,gt=0"`
	PaymentMethodId int  `json:"payment_method_id" binding:"required"`
	Discount        int  `json:"discount" binding:"required"`
	Tax             int  `json:"tax" binding:"required"`
	SubTotal        int  `json:"sub_total" binding:"required"`
}

type TopupResponse struct {
	Amount        uint   `json:"amount"`
	PaymentMethod string `json:"payment_method"`
	Discount      int    `json:"discount"`
	Tax           int    `json:"tax"`
	SubTotal      int    `json:"sub_total"`
}
