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

func (bs BidService) GetBidStatus(u UserParam, bidId string) (string, error) {
	_ = u

	status, err := bs.bidRepo.GetBidStatus(context.Background(), bidId)

	return status, err
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

func (bs BidService) GetUserBids(ubp GetUserBidParams) ([]entity.Bid, error) {
	return bs.bidRepo.GetUserBids(context.Background(), ubp.Username, ubp.Limit, ubp.Offset)
}

func (bs BidService) GetBidsForTender(bftp GetBidsForTenderParams, tenderId string) ([]entity.Bid, error) {
	return bs.bidRepo.GetBidsForTender(context.Background(), tenderId, bftp.Limit, bftp.Offset)
}
