package v1

import (
	"avitoTech/internal/service"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"io"
	log "log/slog"
	"net/http"
)

type tenderRoutes struct {
	tenderService service.Tender
}

func newTenderRoutes(r chi.Router, tenderService service.Tender) {
	routes := &tenderRoutes{
		tenderService: tenderService,
	}

	r.Get("/", routes.getTenders)

	r.Post("/new", routes.createTender)

	r.Get("/my", routes.getUserTenders)
	r.Get("/{tenderId}/status", routes.getTenderStatus)

	r.Put("/{tendersId}/status", routes.updateTenderStatus)

	r.Patch("/{tenderId}/edit", routes.editTender)

	r.Put("/{tenderId}/rollback/{version}", routes.rollbackTender)

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
		log.Debug("err: ", err.Error())
		ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tender)
}

func (tr *tenderRoutes) getTenders(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid params", http.StatusBadRequest)
		return
	}

	gtp := new(service.GetTendersParams)
	if err := schema.NewDecoder().Decode(gtp, r.Form); err != nil {
		http.Error(w, "invalid params", http.StatusInternalServerError)
		return
	}

	tenders, err := tr.tenderService.GetTenders(*gtp)

	if err != nil {
		if err == service.ErrTendersNotFound {
			ErrorResponse(w, "tenders not found", http.StatusBadRequest)
			return
		}
		log.Debug("err: %v", err.Error())
		ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenders)
}

func (tr *tenderRoutes) getUserTenders(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (tr *tenderRoutes) editTender(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (tr *tenderRoutes) rollbackTender(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (tr *tenderRoutes) getTenderStatus(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (tr *tenderRoutes) updateTenderStatus(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}
