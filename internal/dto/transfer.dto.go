package dto

type ReceiverResponse struct {
	Id          int     `json:"id"`
	Photo       *string `json:"photo"`
	Receiver    *string `json:"receiver"`
	PhoneNumber *string `json:"phone_number"`
}

type PaginationMetaResponse struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type ReceiverListResponse struct {
	Items []ReceiverResponse     `json:"items"`
	Meta  PaginationMetaResponse `json:"meta"`
}

type FindReceiversQueryRequest struct {
	Search string `form:"search" binding:"max=100"`
	Page   *int   `form:"page" binding:"gte=1"`
	Limit  *int   `form:"limit" binding:"gte=1,lte=100"`
}
