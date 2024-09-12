package v1

import (
	"avitoTech/internal/service"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type tenderRoutes struct {
	tenderService service.Tender
}

func newTenderRoutes(r chi.Router, tenderService service.Tender) {
	routes := &tenderRoutes{
		tenderService: tenderService,
	}

	r.Post("/new", routes.createTender)
}

func (tr *tenderRoutes) createTender(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorResponse(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var t service.CreateTenderInput
	err = json.Unmarshal(body, &t)
	if err != nil {
		ErrorResponse(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tender, err := tr.tenderService.CreateTender(t)

	if err != nil {
		if err == service.ErrUserIsNotResposible || err == service.ErrUserNotExists {
			ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tender)

}
