package utils

type WebResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data"`
}
