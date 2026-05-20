package dto

type PaymentMethodResponse struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Logo   string `json:"logo"`
	Method string `json:"method"`
	Tax    int    `json:"tax"`
}
