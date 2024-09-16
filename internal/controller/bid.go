package controller

import (
	"avitoTech/internal/service"
	"github.com/go-chi/chi/v5"
	log "log/slog"
	"net/http"
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
	log.Debug("CreateBid err: ", err)

	//if err != nil {
	//	if err == service.ErrUserIsNotResposible || err == service.ErrUserNotExists {
	//		ErrorResponse(w, err.Error(), http.StatusBadRequest)
	//		return
	//	}
	//	log.Debug("err: ", err.Error())
	//	ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
	//	return
	//}

	SendJSONResponse(w, bid)
}

func (bc *BidController) GetUserBids(w http.ResponseWriter, r *http.Request) {
	ubp, err := DecodeFormParams[service.GetUserBidParams](r)
	if err != nil {
		HandleRequestError(w, err)
		return
	}

	bids, err := bc.BidService.GetUserBids(*ubp)
	log.Debug("GetUserBids err: ", err)

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
	log.Debug("GetBidsForTender err: ", err)

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
	log.Debug("GetBidStatus err: ", err)

	SendJSONResponse(w, status)
}

func (bc *BidController) UpdateBidStatus(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (bc *BidController) EditBid(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (bc *BidController) SubmitBidDecision(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (bc *BidController) SubmitBidFeedback(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (bc *BidController) RollbackBid(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (bc *BidController) GetBidReviews(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}
