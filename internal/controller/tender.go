package controller

import (
	"avitoTech/internal/service"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	log "log/slog"
	"net/http"
	"strings"
)

type TenderController struct {
	tenderService service.Tender
}

func NewTenderController(tenderService service.Tender) TenderController {
	return TenderController{
		tenderService: tenderService,
	}
}

func (tc *TenderController) CreateTender(w http.ResponseWriter, r *http.Request) {

	t, err := ParseJSONBody[service.CreateTenderInput](r, w)

	if err != nil {
		HandleRequestError(w, err)
		return
	}

	tender, err := tc.tenderService.CreateTender(*t)

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

func (tc *TenderController) GetTenders(w http.ResponseWriter, r *http.Request) {

	gtp, err := DecodeFormParams[service.GetTendersParams](r)

	if err != nil {
		HandleRequestError(w, err)
		return
	}

	tenders, err := tc.tenderService.GetTenders(*gtp)

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

func (tc *TenderController) GetUserTenders(w http.ResponseWriter, r *http.Request) {

	gutp, err := DecodeFormParams[service.GetUserTendersParams](r)
	if err != nil {
		HandleRequestError(w, err)
		return
	}

	fmt.Println(gutp)

	tenders, err := tc.tenderService.GetUserTenders(*gutp)

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

func (tc *TenderController) GetTenderStatus(w http.ResponseWriter, r *http.Request) {

	u, err := DecodeFormParams[service.UserParam](r)
	if err != nil {
		HandleRequestError(w, err)
		return
	}

	tenderId := chi.URLParam(r, "tenderId")

	status, err := tc.tenderService.GetTenderStatus(*u, tenderId)

	if err != nil {
		if err == service.ErrTenderNotFound || err == service.ErrUserNotExists || err == service.ErrUserIsNotResposible {
			ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Debug("err: %v", err.Error())
		ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, status)
}

func (tc *TenderController) EditTender(w http.ResponseWriter, r *http.Request) {
	//TODO: implement me fully
	u, err := DecodeFormParams[service.UserParam](r)
	if err != nil {
		HandleRequestError(w, err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		HandleRequestError(w, err)
		return
	}

	params := make(map[string]interface{})
	err = json.Unmarshal(body, &params)
	if err != nil {
		HandleRequestError(w, err)
		return
	}

	if val, ok := params["serviceType"]; ok {
		params["service_type"] = strings.ToUpper(val.(string)) // Добавляем новое значение с новым ключом
		delete(params, "serviceType")                          // Удаляем старый ключ
	}

	if val, ok := params["organizationId"]; ok {
		params["organization_id"] = val.(string) // Добавляем новое значение с новым ключом
		delete(params, "serviceType")            // Удаляем старый ключ
	}

	fmt.Println(params)

	tenderId := chi.URLParam(r, "tenderId")

	tender, err := tc.tenderService.EditTender(*u, tenderId, params)

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

func (tc *TenderController) RollbackTender(w http.ResponseWriter, r *http.Request) {
	//TODO: implement me
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (tc *TenderController) UpdateTenderStatus(w http.ResponseWriter, r *http.Request) {
	utsp, err := DecodeFormParams[service.UpdateTenderStatusParams](r)
	if err != nil {
		HandleRequestError(w, err)
		return
	}

	tenderId := chi.URLParam(r, "tenderId")

	tender, err := tc.tenderService.UpdateTenderStatus(*utsp, tenderId)

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
