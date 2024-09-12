package v1

import (
	"avitoTech/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type bidRoutes struct {
	bidService service.Bid
}

func newBidRoutes(r chi.Router, bidService service.Bid) {
	routes := &bidRoutes{
		bidService: bidService,
	}

	r.Post("/new", routes.createBid)

	r.Get("/my", routes.getUserBids)
	r.Get("/{tenderId}/list}", routes.getBidsForTender)
	r.Get("/{bidId}/status", routes.getBidStatus)

	r.Put("/{bidId}/status", routes.updateBidStatus)

	r.Patch("/{bidId}/edit", routes.editBid)

	r.Put("/{bidId}/submit_decision", routes.submitBidDecision)
	r.Put("/{bidId}/feedback", routes.submitBidFeedback)
	r.Put("/{bidId}/rollback/{version}", routes.rollbackBid)

	r.Get("/{tenderId}/reviews", routes.getBidReviews)
}

func (br *bidRoutes) getUserBids(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (br *bidRoutes) createBid(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (br *bidRoutes) editBid(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (br *bidRoutes) submitBidFeedback(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (br *bidRoutes) rollbackBid(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (br *bidRoutes) getBidStatus(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (br *bidRoutes) updateBidStatus(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (br *bidRoutes) submitBidDecision(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (br *bidRoutes) getBidsForTender(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}

func (br *bidRoutes) getBidReviews(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, "not implemented", http.StatusBadRequest)
}
