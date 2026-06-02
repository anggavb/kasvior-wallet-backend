package dto

import "time"

type TransactionHistoryQueryRequest struct {
	Q     string `form:"q" binding:"max=100"`
	Page  *int   `form:"page" binding:"gte=1"`
	Limit *int   `form:"limit" binding:"gte=1,lte=100"`
}

type TransactionHistoryItemResponse struct {
	Id                int       `json:"id"`
	Type              string    `json:"type"`
	Direction         string    `json:"direction"`
	Status            string    `json:"status"`
	Amount            float64   `json:"amount"`
	CounterpartyName  *string   `json:"counterparty_name"`
	CounterpartyPhone *string   `json:"counterparty_phone"`
	CounterpartyPhoto *string   `json:"counterparty_photo"`
	PaymentMethod     *string   `json:"payment_method"`
	Notes             *string   `json:"notes"`
	CreatedAt         time.Time `json:"created_at"`
}

type TransactionHistoryMetaResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type TransactionHistoryResponse struct {
	Items []TransactionHistoryItemResponse `json:"items"`
	Meta  TransactionHistoryMetaResponse   `json:"meta"`
}
