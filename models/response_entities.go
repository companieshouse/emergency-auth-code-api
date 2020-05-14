package models

// ResponseResource is the object returned in an error case
type ResponseResource struct {
	Message string `json:"message"`
}

// NewMessageResponse - convenience function for creating a response resource
func NewMessageResponse(message string) *ResponseResource {
	return &ResponseResource{Message: message}
}
