package service

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo"
	"context"
)

type BidService struct {
	bidRepo         repo.Bid
	userRepo        repo.User
	responsibleRepo repo.Responsible
}

func NewBidService(bidRepo repo.Bid, userRepo repo.User, responsibleRepo repo.Responsible) *BidService {
	return &BidService{
		bidRepo:         bidRepo,
		userRepo:        userRepo,
		responsibleRepo: responsibleRepo,
	}
}

func (bs BidService) CreateBid(i CreateBidInput) (entity.Bid, error) {
	return bs.bidRepo.CreateBid(context.Background(), i.Name, i.Description, i.TenderId, i.AuthorType, i.AuthorId)
}

func (bs BidService) GetUserBids(params GetUserBidParams) ([]entity.Bid, error) {
	//TODO implement me
	panic("implement me")
}
