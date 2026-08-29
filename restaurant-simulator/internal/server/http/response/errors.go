package response

import "errors"

type ErrorDetails struct {
	Message string `json:"message"`
	RootErr string `json:"err"`
}

type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

func NewErrorResponse(err error, message string) ErrorResponse {
	rootErr := err
	for {
		unwrapped := errors.Unwrap(rootErr)
		if unwrapped == nil {
			break
		}
		rootErr = unwrapped
	}

	return ErrorResponse{
		Error: ErrorDetails{
			Message: message,
			RootErr: rootErr.Error(),
		},
	}
}
