package controller

import (
	"avitoTech/internal/service"
	"fmt"
	"github.com/go-chi/chi/v5"
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
		HandleRequestError(w, err)
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
		HandleRequestError(w, err)
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
		HandleRequestError(w, err)
		return
	}

	fmt.Println(gutp)

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

	SendJSONResponse(w, tenders)
}

func (tr *TenderController) GetTenderStatus(w http.ResponseWriter, r *http.Request) {

	u, err := DecodeFormParams[service.UserParam](r)
	if err != nil {
		HandleRequestError(w, err)
		return
	}

	tenderId := chi.URLParam(r, "tenderId")

	status, err := tr.tenderService.GetTenderStatus(*u, tenderId)

	if err != nil {
		if err == service.ErrTenderNotFound || err == service.ErrUserNotExists {
			ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Debug("err: %v", err.Error())
		ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, status)
}

func (tr *TenderController) EditTender(w http.ResponseWriter, r *http.Request) {
	//TODO: implement me fully
	u, err := DecodeFormParams[service.UserParam](r)
	if err != nil {
		HandleRequestError(w, err)
		return
	}

	pt, err := ParseJSONBody[service.PatchTenderInput](r, w)
	if err != nil {
		HandleRequestError(w, err)
		return
	}

	tenderId := chi.URLParam(r, "tenderId")

	tender, err := tr.tenderService.PathTender(*u, tenderId, *pt)

	if err != nil {
		if err == service.ErrUserNotExists || err == service.ErrTenderNotFound {
			ErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}
		log.Debug("err: %v", err.Error())
		ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, tender)

}

func (tr *TenderController) RollbackTender(w http.ResponseWriter, r *http.Request) {
	//TODO: implement me
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (tr *TenderController) UpdateTenderStatus(w http.ResponseWriter, r *http.Request) {
	utsp, err := DecodeFormParams[service.UpdateTenderStatusParams](r)
	if err != nil {
		HandleRequestError(w, err)
		return
	}

	tenderId := chi.URLParam(r, "tenderId")

	tender, err := tr.tenderService.UpdateTenderStatus(*utsp, tenderId)

	if err != nil {
		if err == service.ErrUserNotExists || err == service.ErrTenderNotFound {
			ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Debug("err: ", err.Error())
		ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, tender)
}
