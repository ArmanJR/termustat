package dto

// ErrorResponse is used for any endpoint-level error payload
// swagger:model ErrorResponse
type ErrorResponse struct {
	Error string `json:"error"`
}
