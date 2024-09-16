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

func (bs *BidService) CreateBid(i CreateBidInput) (entity.Bid, error) {
	return bs.bidRepo.CreateBid(context.Background(), i.Name, i.Description, i.TenderId, i.AuthorType, i.AuthorId)
}

func (bs *BidService) GetUserBids(ubp GetUserBidParams) ([]entity.Bid, error) {
	return bs.bidRepo.GetUserBids(context.Background(), ubp.Username, ubp.Limit, ubp.Offset)
}

func (bs *BidService) GetBidsForTender(bftp GetBidsForTenderParams, tenderId string) ([]entity.Bid, error) {
	return bs.bidRepo.GetBidsForTender(context.Background(), tenderId, bftp.Limit, bftp.Offset)
}

func (bs *BidService) UpdateBidStatus(params UpdateBidStatusParams, bidId string) (entity.Bid, error) {

	return bs.bidRepo.UpdateBidStatus(context.Background(), params.Status, bidId)
}

func (bs *BidService) GetBidStatus(u UserParam, bidId string) (string, error) {
	_ = u

	status, err := bs.bidRepo.GetBidStatus(context.Background(), bidId)

	return status, err
}

func (bs *BidService) EditBid(param UserParam, bidId string, params map[string]interface{}) (entity.Bid, error) {
	return bs.bidRepo.EditBid(context.Background(), bidId, params)
}

func (bs *BidService) SubmitBidFeedback(params SubmitBidFeedbackParams, bidId string) (entity.Bid, error) {
	_ = bs.bidRepo.CreateBidFeedback(context.Background(), params.BidFeedback, bidId)

	bid, err := bs.bidRepo.GetBid(context.Background(), bidId)

	return bid, err

}

func (bs *BidService) RollbackBid(param UserParam, bidId string, version int) (entity.Bid, error) {
	return bs.bidRepo.RollbackBidVersion(context.Background(), bidId, version)
}
