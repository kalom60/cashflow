package response

type ErrorResponseFormat struct {
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Message    string       `json:"message,omitempty"`
	FieldError []FieldError `json:"field_error,omitempty"`
}

type FieldError struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SuccessResponse struct {
	Data any `json:"data"`
}
