package api

import (
	"fmt"
	"net/http"
)

type MessageResponse struct {
	code    int
	message string
}

func (err *MessageResponse) Error() string {
	return fmt.Sprintf("API error: code=%d, message=%s", err.code, err.message)
}

func (err MessageResponse) Write(w *http.ResponseWriter) error {
	(*w).WriteHeader(err.code)
	_, returnErr := (*w).Write([]byte(err.message))
	return returnErr
}

func SendCode(code int, message string) *MessageResponse {
	return &MessageResponse{code, message}
}

func SendOk(message string) *MessageResponse {
	return SendCode(http.StatusOK, message)
}

func SendInternalServerError() *MessageResponse {
	return SendCode(http.StatusInternalServerError, "Internal Server Error")
}

func SendBadRequest(message string) *MessageResponse {
	return SendCode(http.StatusBadRequest, message)
}

func SendUnauthorized(message string) *MessageResponse {
	return SendCode(http.StatusUnauthorized, message)
}

func SendForbidden(message string) *MessageResponse {
	return SendCode(http.StatusForbidden, message)
}

func SendNotFound(message string) *MessageResponse {
	return SendCode(http.StatusNotFound, message)
}
