package api

// ErrorResponse defines API response in case of error
type ErrorResponse struct {
	Message string        `json:"message"`
	Details []interface{} `json:"details,omitempty"`
}
