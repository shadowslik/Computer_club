package domain

// ErrorResponse is the standard error body returned by API endpoints.
type ErrorResponse struct {
	Error string `json:"error"`
}
