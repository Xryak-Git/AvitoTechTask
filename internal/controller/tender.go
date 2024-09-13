package controller

import (
	"avitoTech/internal/service"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	log "log/slog"
	"net/http"
)

type TenderController struct {
	tenderService service.Tender
}

func NewTenderController(tenderService service.Tender) TenderController {
	return TenderController{
		tenderService: tenderService,
	}
}

func (tr *TenderController) CreateTender(w http.ResponseWriter, r *http.Request) {

	t, err := ParseJSONBody[service.CreateTenderInput](r, w)

	if err != nil {
		ErrorResponse(w, "invalid request body", http.StatusBadRequest)
		log.Debug("err: " + err.Error())
		return
	}

	tender, err := tr.tenderService.CreateTender(*t)

	if err != nil {
		if err == service.ErrUserIsNotResposible || err == service.ErrUserNotExists {
			ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Debug("err: ", err.Error())
		ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, tender)
}

func (tr *TenderController) GetTenders(w http.ResponseWriter, r *http.Request) {

	gtp, err := DecodeFormParams[service.GetTendersParams](r)

	if err != nil {
		ErrorResponse(w, "invalid params", http.StatusBadRequest)
		log.Debug("err: " + err.Error())
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

	SendJSONResponse(w, tenders)
}

func (tr *TenderController) GetUserTenders(w http.ResponseWriter, r *http.Request) {

	gutp, err := DecodeFormParams[service.GetUserTendersParams](r)
	if err != nil {
		ErrorResponse(w, "invalid params", http.StatusBadRequest)
		log.Debug("err: " + err.Error())
		return
	}

	tenders, err := tr.tenderService.GetUserTenders(*gutp)

	if err != nil {
		if err == service.ErrUserNotExists {
			ErrorResponse(w, "user not exists", http.StatusUnauthorized)
			return
		}
		log.Debug("err: %v", err.Error())
		ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
		return
	}

	if len(tenders) == 0 {
		ErrorResponse(w, "tenders not found", http.StatusBadRequest)
		return
	}

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

func (tr *TenderController) GetTenderStatus(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid params", http.StatusBadRequest)
		return
	}

	gtsp := new(service.GetTenderStatusParams)
	if err := schema.NewDecoder().Decode(gtsp, r.Form); err != nil {
		http.Error(w, "invalid params", http.StatusInternalServerError)
		return
	}

	tenderId := chi.URLParam(r, "tenderId")
	fmt.Println(tenderId)
	status, err := tr.tenderService.GetTenderStatus(*gtsp, tenderId)

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
	json.NewEncoder(w).Encode(status)
}

func (tr *TenderController) EditTender(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (tr *TenderController) RollbackTender(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (tr *TenderController) UpdateTenderStatus(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}
