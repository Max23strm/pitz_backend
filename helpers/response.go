package helpers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
	Code    string      `json:"code,omitempty"`
}

type PaginatedResponse struct {
	Response
	Meta PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func writeJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

func SuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	writeJSON(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func CreatedResponse(w http.ResponseWriter, message string, data interface{}) {
	writeJSON(w, http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string, err error, codes ...string) {
	code := ""

	if len(codes) > 0 {
		code = codes[0]
	}

	response := Response{
		Success: false,
		Message: message,
		Code:    code,
		Data:    nil,
	}

	if err != nil {
		response.Error = err.Error()
	}

	writeJSON(w, statusCode, response)
}

func BadRequestResponse(w http.ResponseWriter, message string, err error, code string) {
	ErrorResponse(w, http.StatusBadRequest, message, err, code)
}

func UnauthorizedResponse(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusUnauthorized, message, nil)
}

func ForbiddenResponse(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusForbidden, message, nil)
}

func NotFoundResponse(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusNotFound, message, nil)
}

func InternalServerErrorResponse(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusInternalServerError, message, nil)
}

func PaginatedSuccessResponse(
	w http.ResponseWriter,
	message string,
	data interface{},
	meta PaginationMeta,
) {
	writeJSON(w, http.StatusOK, PaginatedResponse{
		Response: Response{
			Success: true,
			Message: message,
			Data:    data,
		},
		Meta: meta,
	})
}
