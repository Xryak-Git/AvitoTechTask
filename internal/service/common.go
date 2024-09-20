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
		if err == repoerrs.ErrNotFound {
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

func IsTenderExists(tenderRepo repo.Tender, tenderId string) error {
	exists, err := tenderRepo.IsTenderExists(context.Background(), tenderId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return ErrTendersNotFound
		}
		return err
	}

	if !exists {
		return ErrTendersNotFound
	}
	return nil
}

func GetUserByName(userRepo repo.User, name string) (entity.User, error) {
	user, err := userRepo.GetByName(context.Background(), name)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.User{}, ErrUserNotExists
		}
		return entity.User{}, err
	}
	return user, nil
}
