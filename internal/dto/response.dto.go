package dto

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"is_success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}
