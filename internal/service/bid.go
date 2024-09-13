package service

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo"
)

type BidService struct {
	bidRepo         repo.Bid
	userRepo        repo.User
	responsibleRepo repo.Responsible
}

func (bs BidService) CreateBid(input CreateBidInput) (entity.Bid, error) {
	//isResponsible, err := bs.responsibleRepo.IsUserResponsibleForOrganizationByTenderId(context.Background(), input.AuthorId, input.TenderId)
	//
	//if err != nil {
	//	if err == repoerrs.ErrNotFound {
	//		return entity.Bid{}, ErrUserNotExists
	//	}
	//	return entity.Bid{}, err
	//}
	//
	//isResponsible, err := s.responsibleRepo.IsUserResponsibleForOrganization(context.Background(), u.Id, ct.OrganizationId)
	//if err != nil {
	//	if err == repoerrs.ErrNotExists {
	//		return entity.Tender{}, ErrUserIsNotResposible
	//	}
	//	return entity.Tender{}, ErrCannotCreateTender
	//}
	//
	//if isResponsible == false {
	//	return entity.Tender{}, ErrUserIsNotResposible
	//}
	//
	//t, err := s.tenderRepo.CreateTender(context.Background(), ct.Name, ct.Description, ct.ServiceType, ct.Status, ct.OrganizationId)

	return entity.Bid{}, nil
}

func NewBidService(bidRepo repo.Bid, userRepo repo.User, responsibleRepo repo.Responsible) *BidService {
	return &BidService{
		bidRepo:         bidRepo,
		userRepo:        userRepo,
		responsibleRepo: responsibleRepo,
	}
}
