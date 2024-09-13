package controller

import (
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, message string, code int) {
	// Возвращаем ошибку 400 Bad Request
	http.Error(w, "{ \"message\": \""+message+"\" }", code)
}
