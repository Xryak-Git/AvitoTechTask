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
	BidService service.Bid
}

func NewBidController(bidService service.Bid) BidController {
	return BidController{
		BidService: bidService,
	}
}

func (bc *BidController) CreateBid(w http.ResponseWriter, r *http.Request) {

	bi, err := ParseJSONBody[service.CreateBidInput](r, w)
	if err != nil {
		HandleRequestError(w, err)
		return
	}

	bid, err := bc.BidService.CreateBid(*bi)

	if err != nil {
		log.Debug("CreateBid err: ", err.Error())
		if err == service.ErrUserNotExists {
			ErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err == service.ErrUserIsNotResposible {
			ErrorResponse(w, err.Error(), http.StatusForbidden)
			return
		}
		if err == service.ErrTenderNotFound {
			ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
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

	bids, err := bc.BidService.GetUserBids(*ubp)
	log.Debug("GetUserBids err: ", err.Error())

	if err != nil {
		if err == service.ErrUserNotExists {
			ErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
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

	bids, err := bc.BidService.GetBidsForTender(*bftp, tenderId)
	if err != nil {
		log.Debug("GetBidsForTender err: ", err.Error())
		if err == service.ErrUserNotExists {
			ErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err == service.ErrUserIsNotResposible {
			ErrorResponse(w, err.Error(), http.StatusForbidden)
			return
		}
		if err == service.ErrTenderNotFound || err == service.ErrBidNotFound {
			ErrorResponse(w, "tender or bid was not found", http.StatusNotFound)
			return
		}
		ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
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

	status, err := bc.BidService.GetBidStatus(*u, bidId)
	if err != nil {
		log.Debug("GetBidStatus err: ", err.Error())
		if err == service.ErrUserNotExists {
			ErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err == service.ErrUserIsNotResposible {
			ErrorResponse(w, err.Error(), http.StatusForbidden)
			return
		}
		if err == service.ErrBidNotFound {
			ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
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

	bid, err := bc.BidService.UpdateBidStatus(*bs, bidId)
	if err != nil {
		log.Debug("UpdateBidStatus err: ", err.Error())
		if err == service.ErrUserNotExists {
			ErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err == service.ErrUserIsNotResposible {
			ErrorResponse(w, err.Error(), http.StatusForbidden)
			return
		}
		if err == service.ErrBidNotFound {
			ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
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

	bid, err := bc.BidService.EditBid(*u, bidId, params)
	if err != nil {
		log.Debug("EditBid err: ", err.Error())
		if err == service.ErrUserNotExists {
			ErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err == service.ErrUserIsNotResposible {
			ErrorResponse(w, err.Error(), http.StatusForbidden)
			return
		}
		if err == service.ErrBidNotFound {
			ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, bid)

}

func (bc *BidController) SubmitBidDecision(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (bc *BidController) SubmitBidFeedback(w http.ResponseWriter, r *http.Request) {
	bf, err := DecodeFormParams[service.SubmitBidFeedbackParams](r)
	if err != nil {
		HandleRequestError(w, err)
		return
	}
	bidId := chi.URLParam(r, "bidId")

	bid, err := bc.BidService.SubmitBidFeedback(*bf, bidId)
	if err != nil {
		log.Debug("SubmitBidDecision err: ", err.Error())
		if err == service.ErrUserNotExists {
			ErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err == service.ErrUserIsNotResposible {
			ErrorResponse(w, err.Error(), http.StatusForbidden)
			return
		}
		if err == service.ErrBidNotFound {
			ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
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

	bid, err := bc.BidService.RollbackBid(*u, bidId, versionInt)
	if err != nil {
		log.Debug("RollbackBid err: ", err.Error())
		if err == service.ErrUserNotExists {
			ErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err == service.ErrUserIsNotResposible || err == service.ErrUserDoseNotMadeBidForTender {
			ErrorResponse(w, err.Error(), http.StatusForbidden)
			return
		}
		if err == service.ErrBidOrVersionNotFound {
			ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
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

	reviews, err := bc.BidService.GetBidReviews(*params, tenderId)
	if err != nil {
		log.Debug("RollbackBid err: ", err.Error())
		if err == service.ErrUserNotExists {
			ErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err == service.ErrUserIsNotResposible {
			ErrorResponse(w, err.Error(), http.StatusForbidden)
			return
		}
		if err == service.ErrTenderNotFound || err == service.ErrBidReviewsNotFound {
			ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, reviews)
}
