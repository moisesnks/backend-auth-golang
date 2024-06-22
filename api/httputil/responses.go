package httputil

type ErrorResponse struct {
	Message string `json:"message"`
}

type StandardResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
