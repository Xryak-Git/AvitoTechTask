package controller

import (
	"avitoTech/internal/service"
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

func (bc *BidController) GetUserBids(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (bc *BidController) CreateBid(w http.ResponseWriter, r *http.Request) {

	b, err := ParseJSONBody[service.CreateBidInput](r, w)
	if err != nil {
		HandleRequestError(w, err)
		return
	}

	bid, err := bc.BidService.CreateBid(*b)

	if err != nil {
		if err == service.ErrUserIsNotResposible || err == service.ErrUserNotExists {
			ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Debug("err: ", err.Error())
		ErrorResponse(w, "interanl server error", http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, bid)
}

func (bc *BidController) EditBid(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (bc *BidController) SubmitBidFeedback(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (bc *BidController) RollbackBid(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (bc *BidController) GetBidStatus(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (bc *BidController) UpdateBidStatus(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (bc *BidController) SubmitBidDecision(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (bc *BidController) GetBidsForTender(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (bc *BidController) GetBidReviews(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}
