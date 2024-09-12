package service

import (
	"avitoTech/internal/repo"
)

type BidService struct {
	bidRepo repo.Bid
}

func NewBidService(bidRepo repo.Bid) *BidService {
	return &BidService{
		bidRepo: bidRepo,
	}
}
