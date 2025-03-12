package http

import (
	"errors"
	"ggltask/internal/task/domain/usecase"
	"net/http"
)

type ErrorResponse struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

// UseCaesErrorToErrorResp is a helper function that converts a usecase error to an error response.
// It returns the HTTP status code and the error response.
func UseCaesErrorToErrorResp(err error) (int, ErrorResponse) {
	var usecaseErr usecase.UseCaseError
	if !errors.As(err, &usecaseErr) {
		return http.StatusInternalServerError, ErrorResponse{
			ErrorCode:    "INTERNAL_SERVER_ERROR",
			ErrorMessage: "Internal Server Error",
		}
	}

	return usecaseErr.HTTPStatusCode(), ErrorResponse{
		ErrorCode:    usecaseErr.ErrorCode(),
		ErrorMessage: usecaseErr.ErrorMsg(),
	}
}

func InvalidRequestError() ErrorResponse {
	return ErrorResponse{
		ErrorCode:    "INVALID_REQUEST",
		ErrorMessage: "Invalid Request",
	}
}
