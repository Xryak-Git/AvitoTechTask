package service

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo"
	"avitoTech/internal/repo/repoerrs"
	"context"
	"errors"
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

func (bs *BidService) CreateBid(params CreateBidInput) (entity.Bid, error) {

	user, err := GetUserById(bs.userRepo, params.AuthorId)
	if err != nil {
		return entity.Bid{}, err
	}

	err = IsUserResponsibleByTenderId(bs.responsibleRepo, user.Id, params.TenderId)
	if err != nil {
		return entity.Bid{}, err
	}

	err = IsTenderExists(bs.tenderRepo, params.TenderId)
	if err != nil {
		return entity.Bid{}, err
	}

	bid, err := bs.bidRepo.CreateBid(context.Background(), params.Name, params.Description, params.TenderId, params.AuthorType, params.AuthorId)
	if err != nil {
		return entity.Bid{}, err
	}

	return bid, nil
}

func (bs *BidService) GetUserBids(params GetUserBidParams) ([]entity.Bid, error) {
	_, err := GetUserByName(bs.userRepo, params.Username)
	if err != nil {
		return []entity.Bid{}, err
	}

	return bs.bidRepo.GetUserBids(context.Background(), params.Username, params.Limit, params.Offset)
}

func (bs *BidService) GetBidsForTender(params GetBidsForTenderParams, tenderId string) ([]entity.Bid, error) {
	user, err := GetUserByName(bs.userRepo, params.Username)
	if err != nil {
		return []entity.Bid{}, err
	}

	err = IsUserResponsibleByTenderId(bs.responsibleRepo, user.Id, tenderId)
	if err != nil {
		return []entity.Bid{}, err
	}

	err = IsTenderExists(bs.tenderRepo, tenderId)
	if err != nil {
		return []entity.Bid{}, err
	}

	bids, err := bs.bidRepo.GetBidsForTender(context.Background(), tenderId, params.Limit, params.Offset)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return []entity.Bid{}, ErrBidNotFound
		}
		return []entity.Bid{}, err
	}

	return bids, err
}

func (bs *BidService) GetBidStatus(params UserParam, bidId string) (string, error) {
	user, err := bs.userRepo.GetByName(context.Background(), params.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return "", ErrUserNotExists
		}
		return "", err
	}

	isResponsibe, err := bs.responsibleRepo.IsUserResponsibleForOrganizationByBidId(context.Background(), user.Id, bidId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return "", ErrUserIsNotResposible
		}
		return "", err
	}
	if !isResponsibe {
		return "", ErrUserIsNotResposible
	}

	exists, err := bs.bidRepo.IsBidExists(context.Background(), bidId)
	if err != nil || !exists {
		return "", ErrTenderNotFound
	}

	return bs.bidRepo.GetBidStatus(context.Background(), bidId)
}

func (bs *BidService) UpdateBidStatus(params UpdateBidStatusParams, bidId string) (entity.Bid, error) {
	user, err := bs.userRepo.GetByName(context.Background(), params.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrUserNotExists
		}
		return entity.Bid{}, err
	}

	isResponsibe, err := bs.responsibleRepo.IsUserResponsibleForOrganizationByBidId(context.Background(), user.Id, bidId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrUserIsNotResposible
		}
		return entity.Bid{}, err
	}
	if !isResponsibe {
		return entity.Bid{}, ErrUserIsNotResposible
	}

	exists, err := bs.bidRepo.IsBidExists(context.Background(), bidId)
	if err != nil || !exists {
		return entity.Bid{}, ErrTenderNotFound
	}

	return bs.bidRepo.UpdateBidStatus(context.Background(), params.Status, bidId)
}

func (bs *BidService) EditBid(params UserParam, bidId string, editFields map[string]interface{}) (entity.Bid, error) {
	user, err := bs.userRepo.GetByName(context.Background(), params.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrUserNotExists
		}
		return entity.Bid{}, err
	}

	isResponsibe, err := bs.responsibleRepo.IsUserResponsibleForOrganizationByBidId(context.Background(), user.Id, bidId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrUserIsNotResposible
		}
		return entity.Bid{}, err
	}
	if !isResponsibe {
		return entity.Bid{}, ErrUserIsNotResposible
	}

	exists, err := bs.bidRepo.IsBidExists(context.Background(), bidId)
	if err != nil || !exists {
		return entity.Bid{}, ErrTenderNotFound
	}
	return bs.bidRepo.EditBid(context.Background(), bidId, editFields)
}

func (bs *BidService) SubmitBidFeedback(params SubmitBidFeedbackParams, bidId string) (entity.Bid, error) {

	user, err := bs.userRepo.GetByName(context.Background(), params.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrUserNotExists
		}
		return entity.Bid{}, err
	}

	isResponsibe, err := bs.responsibleRepo.IsUserResponsibleForOrganizationByBidId(context.Background(), user.Id, bidId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrUserIsNotResposible
		}
		return entity.Bid{}, err
	}
	if !isResponsibe {
		return entity.Bid{}, ErrUserIsNotResposible
	}

	exists, err := bs.bidRepo.IsBidExists(context.Background(), bidId)
	if err != nil || !exists {
		return entity.Bid{}, ErrTenderNotFound
	}

	err = bs.bidRepo.CreateBidFeedback(context.Background(), params.BidFeedback, bidId)
	if err != nil {
		return entity.Bid{}, err
	}

	return bs.bidRepo.GetBid(context.Background(), bidId)

}

func (bs *BidService) RollbackBid(params UserParam, bidId string, version int) (entity.Bid, error) {
	user, err := bs.userRepo.GetByName(context.Background(), params.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrUserNotExists
		}
		return entity.Bid{}, err
	}

	isResponsibe, err := bs.responsibleRepo.IsUserResponsibleForOrganizationByBidId(context.Background(), user.Id, bidId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrUserIsNotResposible
		}
		return entity.Bid{}, err
	}
	if !isResponsibe {
		return entity.Bid{}, ErrUserIsNotResposible
	}

	exists, err := bs.bidRepo.IsBidExists(context.Background(), bidId)
	if err != nil || !exists {
		return entity.Bid{}, ErrTenderNotFound
	}

	bid, err := bs.bidRepo.RollbackBidVersion(context.Background(), bidId, version)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Bid{}, ErrBidOrVersionNotFound
		}
		return entity.Bid{}, err
	}

	return bid, err
}

func (bs *BidService) GetBidReviews(params GetBidReviewsParams, tenderId string) ([]entity.BidReview, error) {
	responsible, err := bs.userRepo.GetByName(context.Background(), params.RequesterUsername)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return []entity.BidReview{}, ErrUserNotExists
		}
		return []entity.BidReview{}, err
	}

	isResponsible, err := bs.responsibleRepo.IsUserResponsibleForOrganizationByTenderId(context.Background(), responsible.Id, tenderId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return []entity.BidReview{}, ErrUserIsNotResposible
		}
		return []entity.BidReview{}, err
	}
	if !isResponsible {
		return []entity.BidReview{}, ErrUserIsNotResposible
	}

	exists, err := bs.tenderRepo.IsTenderExists(context.Background(), tenderId)
	if err != nil || !exists {
		return []entity.BidReview{}, ErrTenderNotFound
	}

	author, err := bs.userRepo.GetByName(context.Background(), params.AuthorUsername)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return []entity.BidReview{}, ErrUserNotExists
		}
		return []entity.BidReview{}, err
	}

	isUserMadeBid, err := bs.bidRepo.IsUserMadeBid(context.Background(), author.Id, tenderId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return []entity.BidReview{}, ErrUserDoseNotMadeBidForTender
		}
		return []entity.BidReview{}, err
	}
	if !isUserMadeBid {
		return []entity.BidReview{}, ErrUserDoseNotMadeBidForTender
	}

	reviews, err := bs.bidRepo.GetAuthorBidReviews(context.Background(), author.Id, params.Limit, params.Offset)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return []entity.BidReview{}, ErrBidReviewsNotFound
		}
		return []entity.BidReview{}, err
	}

	return reviews, nil
}

func (bs *BidService) SubmitBidDecision(params SubmitBidDecisionParams, bidId string) (entity.Bid, error) {
	user, _ := bs.userRepo.GetByName(context.Background(), params.Username)

	err := bs.bidRepo.SubmitBidDecision(context.Background(), user.Id, bidId, params.Decision)
	if err != nil {
		return entity.Bid{}, err
	}

	return bs.bidRepo.GetBid(context.Background(), bidId)
}
