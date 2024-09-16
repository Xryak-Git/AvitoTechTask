package service

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo"
	"avitoTech/internal/repo/repoerrs"
	"context"
)

type BidService struct {
	bidRepo         repo.Bid
	userRepo        repo.User
	responsibleRepo repo.Responsible
	tenderRepo      repo.Tender
}

func NewBidService(bidRepo repo.Bid, userRepo repo.User, responsibleRepo repo.Responsible, tenderRepo repo.Tender) *BidService {
	return &BidService{
		bidRepo:         bidRepo,
		userRepo:        userRepo,
		responsibleRepo: responsibleRepo,
		tenderRepo:      tenderRepo,
	}
}

func (bs *BidService) CreateBid(bi CreateBidInput) (entity.Bid, error) {

	user, err := bs.userRepo.GetById(context.Background(), bi.AuthorId)
	if err != nil {
		return entity.Bid{}, err
	}

	isResponsibe, err := bs.responsibleRepo.IsUserResponsibleForOrganizationByTenderId(context.Background(), user.Id, bi.TenderId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrUserIsNotResposible
		}
		return entity.Bid{}, err
	}
	if !isResponsibe {
		return entity.Bid{}, ErrUserIsNotResposible
	}

	exists, err := bs.tenderRepo.IsTenderExists(context.Background(), bi.TenderId)
	if err != nil || !exists {
		return entity.Bid{}, ErrTenderOrVersionNotFound
	}

	bid, err := bs.bidRepo.CreateBid(context.Background(), bi.Name, bi.Description, bi.TenderId, bi.AuthorType, bi.AuthorId)
	if err != nil {
		return entity.Bid{}, err
	}

	return bid, nil
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
