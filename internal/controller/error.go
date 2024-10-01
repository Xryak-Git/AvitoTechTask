package controller

import (
	"avitoTech/internal/service"
	log "log/slog"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, message string, code int) {
	http.Error(w, "{ \"message\": \""+message+"\" }", code)
	return
}

func HandleRequestError(w http.ResponseWriter, err error) {
	ErrorResponse(w, "invalid params", http.StatusBadRequest)
	log.Debug("err: " + err.Error())
	return
}

var ExpectedErrors = map[error]int{
	service.ErrUserNotExists:               http.StatusUnauthorized,
	service.ErrUserIsNotResposible:         http.StatusForbidden,
	service.ErrUserDoseNotMadeBidForTender: http.StatusForbidden,
	service.ErrBidNotFound:                 http.StatusNotFound,
	service.ErrBidReviewsNotFound:          http.StatusNotFound,
	service.ErrBidOrVersionNotFound:        http.StatusNotFound,
	service.ErrTenderNotFound:              http.StatusNotFound,
	service.ErrTendersNotFound:             http.StatusNotFound,
	service.ErrTenderOrVersionNotFound:     http.StatusNotFound,
}
