package controller

import (
	"avitoTech/internal/service"
	"encoding/json"
	"github.com/gorilla/schema"
	"io"
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

func (tr *TenderController) GetTenders(w http.ResponseWriter, r *http.Request) {
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

func (tr *TenderController) GetUserTenders(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid params", http.StatusBadRequest)
		return
	}

	gutp := new(service.GetUserTendersParams)
	if err := schema.NewDecoder().Decode(gutp, r.Form); err != nil {
		http.Error(w, "invalid params", http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenders)

}

func (tr *TenderController) EditTender(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (tr *TenderController) RollbackTender(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (tr *TenderController) GetTenderStatus(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (tr *TenderController) UpdateTenderStatus(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}
