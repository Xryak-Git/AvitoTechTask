package controller

import (
	"avitoTech/internal/service"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	log "log/slog"
	"net/http"
	"strconv"
	"strings"
)

type TenderController struct {
	tenderService  service.Tender
	expectedErrors map[error]int
}

func NewTenderController(tenderService service.Tender) TenderController {
	expectedErrors := map[error]int{
		service.ErrUserNotExists:           http.StatusUnauthorized,
		service.ErrUserIsNotResposible:     http.StatusForbidden,
		service.ErrTenderOrVersionNotFound: http.StatusNotFound,
		service.ErrTenderNotFound:          http.StatusNotFound,
	}

	return TenderController{
		tenderService:  tenderService,
		expectedErrors: expectedErrors,
	}
}

func (tc *TenderController) CreateTender(w http.ResponseWriter, r *http.Request) {

	t, err := ParseJSONBody[service.CreateTenderInput](r, w)

	if err != nil {
		log.Debug("err: " + err.Error())
		HandleRequestError(w, err)
		return
	}

	tender, err := tc.tenderService.CreateTender(*t)

	if err != nil {
		log.Debug("CreateTender err: ", err.Error())
		HandelServiceError(w, err)
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
		log.Debug("GetTenders err: ", err.Error())
		HandelServiceError(w, err)
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

	tenders, err := tc.tenderService.GetUserTenders(*gutp)

	if err != nil {
		log.Debug("GetUserTenders err: " + err.Error())
		HandelServiceError(w, err)
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
		log.Debug("GetTenderStatus err: ", err.Error())
		HandelServiceError(w, err)
		return
	}

	SendJSONResponse(w, status)
}

func (tc *TenderController) EditTender(w http.ResponseWriter, r *http.Request) {

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
		params["service_type"] = strings.ToUpper(val.(string))
		delete(params, "serviceType")
	}

	if val, ok := params["status"]; ok {
		params["status"] = strings.ToUpper(val.(string))
	}

	if val, ok := params["organizationId"]; ok {
		params["organization_id"] = val.(string)
		delete(params, "serviceType")
	}

	tenderId := chi.URLParam(r, "tenderId")

	tender, err := tc.tenderService.EditTender(*u, tenderId, params)

	if err != nil {
		log.Debug("EditTender err: ", err.Error())
		HandelServiceError(w, err)
		return
	}

	SendJSONResponse(w, tender)

}

func (tc *TenderController) RollbackTender(w http.ResponseWriter, r *http.Request) {

	u, err := DecodeFormParams[service.UserParam](r)
	if err != nil {
		HandleRequestError(w, err)
		return
	}
	tenderId := chi.URLParam(r, "tenderId")
	versionStr := chi.URLParam(r, "version")

	versionInt, err := strconv.Atoi(versionStr)
	if err != nil {
		HandleRequestError(w, err)
		return
	}

	tender, err := tc.tenderService.RollbackTender(*u, tenderId, versionInt)

	if err != nil {
		log.Debug("RollbackTender err: ", err.Error())
		HandelServiceError(w, err)
		return
	}

	SendJSONResponse(w, tender)

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
		log.Debug("UpdateTenderStatus err: ", err.Error())
		HandelServiceError(w, err)
		return
	}

	SendJSONResponse(w, tender)
}
