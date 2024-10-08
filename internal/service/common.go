package service

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo"
	"avitoTech/internal/repo/repoerrs"
	"context"
	"errors"
)

func GetUserById(userRepo repo.User, id string) (entity.User, error) {
	user, err := userRepo.GetById(context.Background(), id)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.User{}, ErrUserNotExists
		}
		return entity.User{}, err
	}
	return user, nil
}

func GetUserByName(userRepo repo.User, name string) (entity.User, error) {
	user, err := userRepo.GetByName(context.Background(), name)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.User{}, ErrUserNotExists
		}
		return entity.User{}, err
	}
	return user, nil
}

func IsUserResponsibleByTenderId(responsibleRepo repo.Responsible, userId, tenderId string) error {
	isResponsible, err := responsibleRepo.IsUserResponsibleForOrganizationByTenderId(context.Background(), userId, tenderId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return ErrUserIsNotResposible
		}
		return err
	}
	if !isResponsible {
		return ErrUserIsNotResposible
	}
	return nil
}

func IsUserResponsibleByBidId(responsibleRepo repo.Responsible, userId, bidId string) error {
	isResponsible, err := responsibleRepo.IsUserResponsibleForOrganizationByBidId(context.Background(), userId, bidId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return ErrUserIsNotResposible
		}
		return err
	}
	if !isResponsible {
		return ErrUserIsNotResposible
	}
	return nil
}

func IsUserResponsibleByOrganizationId(responsibleRepo repo.Responsible, userId, organizationId string) error {
	isResponsible, err := responsibleRepo.IsUserResponsibleForOrganizationByOrganizationId(context.Background(), userId, organizationId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return ErrUserIsNotResposible
		}
		return err
	}
	if !isResponsible {
		return ErrUserIsNotResposible
	}
	return nil
}

func IsUserMadeBidForTender(bidRepo repo.Bid, userId, tenderId string) error {
	isUserMadeBid, err := bidRepo.IsUserMadeBid(context.Background(), userId, tenderId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return ErrUserDoseNotMadeBidForTender
		}
		return err
	}
	if !isUserMadeBid {
		return ErrUserDoseNotMadeBidForTender
	}

	return nil
}

func IsTenderExists(tenderRepo repo.Tender, tenderId string) error {
	exists, err := tenderRepo.IsTenderExists(context.Background(), tenderId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return ErrTenderNotFound
		}
		return err
	}

	if !exists {
		return ErrTenderNotFound
	}
	return nil
}

func IsBidExists(bidRepo repo.Bid, bidId string) error {
	exists, err := bidRepo.IsBidExists(context.Background(), bidId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return ErrBidNotFound
		}
		return err
	}

	if !exists {
		return ErrBidNotFound
	}
	return nil
}
