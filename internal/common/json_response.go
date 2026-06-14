package common

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/MarcelArt/passwordless/internal/enums"
	"gorm.io/gorm"
)

type JSONResponse struct {
	Items     any    `json:"items"`
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}

func NewJSONResponse(items any, message string) *JSONResponse {
	err, ok := items.(error)
	if ok {
		if message == "" {
			message = err.Error()
		} else {
			message = fmt.Sprintf("%s: %s", message, err.Error())
		}
		return &JSONResponse{
			Items:     nil,
			IsSuccess: false,
			Message:   message,
		}
	}

	return &JSONResponse{
		Items:     items,
		IsSuccess: true,
		Message:   message,
	}
}

func StatusCodeFromError(err error) int {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound
	}

	if errors.Is(err, enums.ErrAlreadyRegsitered) {
		return http.StatusUnauthorized
	}

	return http.StatusInternalServerError
}
