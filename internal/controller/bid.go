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

type BidController struct {
	bidService     service.Bid
	expectedErrors map[error]int
}

func NewBidController(bidService service.Bid) BidController {
	expectedErrors := map[error]int{
		service.ErrUserNotExists:               http.StatusUnauthorized,
		service.ErrUserIsNotResposible:         http.StatusForbidden,
		service.ErrUserDoseNotMadeBidForTender: http.StatusForbidden,
		service.ErrBidNotFound:                 http.StatusNotFound,
		service.ErrBidReviewsNotFound:          http.StatusNotFound,
		service.ErrBidOrVersionNotFound:        http.StatusNotFound,
		service.ErrTenderNotFound:              http.StatusNotFound,
	}

	return BidController{
		bidService:     bidService,
		expectedErrors: expectedErrors,
	}
}

func (bc *BidController) CreateBid(w http.ResponseWriter, r *http.Request) {

	bi, err := ParseJSONBody[service.CreateBidInput](r, w)
	if err != nil {
		HandleRequestError(w, err)
		return
	}

	bid, err := bc.bidService.CreateBid(*bi)

	if err != nil {
		log.Debug("CreateBid err: ", err.Error())
		HandelServiceError(w, bc.expectedErrors, err)
		return
	}

	SendJSONResponse(w, bid)
}

func (bc *BidController) GetUserBids(w http.ResponseWriter, r *http.Request) {
	ubp, err := DecodeFormParams[service.GetUserBidParams](r)
	if err != nil {
		HandleRequestError(w, err)
		return
	}

	bids, err := bc.bidService.GetUserBids(*ubp)

	if err != nil {
		log.Debug("GetUserBids err: ", err.Error())
		HandelServiceError(w, bc.expectedErrors, err)
		return
	}

	SendJSONResponse(w, bids)
}

func (bc *BidController) GetBidsForTender(w http.ResponseWriter, r *http.Request) {
	bftp, err := DecodeFormParams[service.GetBidsForTenderParams](r)
	if err != nil {
		HandleRequestError(w, err)
		return
	}
	tenderId := chi.URLParam(r, "tenderId")

	bids, err := bc.bidService.GetBidsForTender(*bftp, tenderId)
	if err != nil {
		log.Debug("GetBidsForTender err: ", err.Error())
		HandelServiceError(w, bc.expectedErrors, err)
		return
	}

	SendJSONResponse(w, bids)
}

func (bc *BidController) GetBidStatus(w http.ResponseWriter, r *http.Request) {
	u, err := DecodeFormParams[service.UserParam](r)
	if err != nil {
		HandleRequestError(w, err)
		return
	}
	bidId := chi.URLParam(r, "bidId")

	status, err := bc.bidService.GetBidStatus(*u, bidId)
	if err != nil {
		log.Debug("GetBidStatus err: ", err.Error())
		HandelServiceError(w, bc.expectedErrors, err)
		return
	}

	SendJSONResponse(w, status)
}

func (bc *BidController) UpdateBidStatus(w http.ResponseWriter, r *http.Request) {
	bs, err := DecodeFormParams[service.UpdateBidStatusParams](r)
	if err != nil {
		HandleRequestError(w, err)
		return
	}
	bidId := chi.URLParam(r, "bidId")

	bid, err := bc.bidService.UpdateBidStatus(*bs, bidId)
	if err != nil {
		log.Debug("UpdateBidStatus err: ", err.Error())
		HandelServiceError(w, bc.expectedErrors, err)
		return
	}

	SendJSONResponse(w, bid)

}

func (bc *BidController) EditBid(w http.ResponseWriter, r *http.Request) {
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

	if val, ok := params["status"]; ok {
		params["status"] = strings.ToUpper(val.(string))
	}

	if val, ok := params["authorType"]; ok {
		params["author_type"] = strings.ToUpper(val.(string))
		delete(params, "authorType")
	}

	if val, ok := params["authorId"]; ok {
		params["author_id"] = val.(string)
		delete(params, "authorId")
	}

	bidId := chi.URLParam(r, "bidId")

	bid, err := bc.bidService.EditBid(*u, bidId, params)
	if err != nil {
		log.Debug("EditBid err: ", err.Error())
		HandelServiceError(w, bc.expectedErrors, err)
		return
	}

	SendJSONResponse(w, bid)

}

func (bc *BidController) SubmitBidFeedback(w http.ResponseWriter, r *http.Request) {
	bf, err := DecodeFormParams[service.SubmitBidFeedbackParams](r)
	if err != nil {
		HandleRequestError(w, err)
		return
	}
	bidId := chi.URLParam(r, "bidId")

	bid, err := bc.bidService.SubmitBidFeedback(*bf, bidId)
	if err != nil {
		log.Debug("SubmitBidDecision err: ", err.Error())
		HandelServiceError(w, bc.expectedErrors, err)
		return
	}

	SendJSONResponse(w, bid)
}

func (bc *BidController) RollbackBid(w http.ResponseWriter, r *http.Request) {
	u, err := DecodeFormParams[service.UserParam](r)
	if err != nil {
		HandleRequestError(w, err)
		return
	}
	bidId := chi.URLParam(r, "bidId")
	versionStr := chi.URLParam(r, "version")

	versionInt, err := strconv.Atoi(versionStr)
	if err != nil {
		HandleRequestError(w, err)
		return
	}

	bid, err := bc.bidService.RollbackBid(*u, bidId, versionInt)
	if err != nil {
		log.Debug("RollbackBid err: ", err.Error())
		HandelServiceError(w, bc.expectedErrors, err)
		return
	}

	SendJSONResponse(w, bid)

}

func (bc *BidController) GetBidReviews(w http.ResponseWriter, r *http.Request) {
	params, err := DecodeFormParams[service.GetBidReviewsParams](r)
	if err != nil {
		HandleRequestError(w, err)
		return
	}
	tenderId := chi.URLParam(r, "tenderId")

	reviews, err := bc.bidService.GetBidReviews(*params, tenderId)
	if err != nil {
		log.Debug("RollbackBid err: ", err.Error())
		HandelServiceError(w, bc.expectedErrors, err)
		return
	}

	SendJSONResponse(w, reviews)
}

func (bc *BidController) SubmitBidDecision(w http.ResponseWriter, r *http.Request) {
	params, err := DecodeFormParams[service.SubmitBidDecisionParams](r)
	if err != nil {
		HandleRequestError(w, err)
		return
	}
	bidId := chi.URLParam(r, "bidId")

	bid, err := bc.bidService.SubmitBidDecision(*params, bidId)

	if err != nil {
		log.Debug("SubmitBidDecision err: ", err.Error())
		HandelServiceError(w, bc.expectedErrors, err)
		return
	}

	SendJSONResponse(w, bid)
}
