package service

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo"
	"avitoTech/internal/repo/repoerrs"
	"context"
	"errors"
	"strings"
)

type TenderService struct {
	tenderRepo      repo.Tender
	userRepo        repo.User
	responsibleRepo repo.Responsible
}

func NewTenderService(tenderRepo repo.Tender, userRepo repo.User, responsibleRepo repo.Responsible) *TenderService {
	return &TenderService{
		tenderRepo:      tenderRepo,
		userRepo:        userRepo,
		responsibleRepo: responsibleRepo,
	}
}

func (ts *TenderService) CreateTender(params CreateTenderInput) (entity.Tender, error) {
	u, err := GetUserByName(ts.userRepo, params.CreatorUsername)
	if err != nil {
		return entity.Tender{}, err
	}

	err = IsUserResponsibleByOrganizationId(ts.responsibleRepo, u.Id, params.OrganizationId)
	if err != nil {
		return entity.Tender{}, err
	}

	return ts.tenderRepo.CreateTender(context.Background(), params.Name, params.Description, params.ServiceType, params.OrganizationId)
}

func (ts *TenderService) GetTenders(params GetTendersParams) ([]entity.Tender, error) {
	for i, st := range params.ServiceType {
		params.ServiceType[i] = strings.ToUpper(st)
	}

	tenders, err := ts.tenderRepo.GetTenders(context.Background(), params.Limit, params.Offset, params.ServiceType)

	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return []entity.Tender{}, ErrTendersNotFound
		}
		return []entity.Tender{}, err
	}

	return tenders, nil
}

func (ts *TenderService) GetUserTenders(params GetUserTendersParams) ([]entity.Tender, error) {
	_, err := GetUserByName(ts.userRepo, params.Username)
	if err != nil {
		return []entity.Tender{}, err
	}

	tenders, err := ts.tenderRepo.GetUserTenders(context.Background(), params.Username, params.Limit, params.Offset)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return []entity.Tender{}, ErrTendersNotFound
		}
		return []entity.Tender{}, err
	}

	return tenders, nil
}

func (ts *TenderService) GetTenderStatus(params UserParam, tenderId string) (string, error) {

	user, err := GetUserByName(ts.userRepo, params.Username)
	if err != nil {
		return "", err
	}

	err = IsTenderExists(ts.tenderRepo, tenderId)
	if err != nil {
		return "", err
	}

	status, err := ts.tenderRepo.GetTenderStatus(context.Background(), tenderId)
	if err != nil {
		return "", err
	}

	if status == "Published" {
		return status, nil
	}

	err = IsUserResponsibleByTenderId(ts.responsibleRepo, user.Id, tenderId)
	if err != nil {
		return "", err
	}

	return status, nil

}

func (ts *TenderService) EditTender(params UserParam, tenderId string, editFields map[string]interface{}) (entity.Tender, error) {
	user, err := GetUserByName(ts.userRepo, params.Username)
	if err != nil {
		return entity.Tender{}, err
	}

	err = IsTenderExists(ts.tenderRepo, tenderId)
	if err != nil {
		return entity.Tender{}, err
	}

	err = IsUserResponsibleByTenderId(ts.responsibleRepo, user.Id, tenderId)
	if err != nil {
		return entity.Tender{}, err
	}

	return ts.tenderRepo.UpdateTender(context.Background(), tenderId, editFields)
}

func (ts *TenderService) UpdateTenderStatus(params UpdateTenderStatusParams, tenderId string) (entity.Tender, error) {

	user, err := GetUserByName(ts.userRepo, params.Username)
	if err != nil {
		return entity.Tender{}, err
	}

	err = IsTenderExists(ts.tenderRepo, tenderId)
	if err != nil {
		return entity.Tender{}, err
	}

	err = IsUserResponsibleByTenderId(ts.responsibleRepo, user.Id, tenderId)
	if err != nil {
		return entity.Tender{}, err
	}

	return ts.tenderRepo.UpdateTenderStatus(context.Background(), params.Status, tenderId)

}

func (ts *TenderService) RollbackTender(params UserParam, tenderId string, version int) (entity.Tender, error) {
	user, err := GetUserByName(ts.userRepo, params.Username)
	if err != nil {
		return entity.Tender{}, err
	}

	err = IsUserResponsibleByTenderId(ts.responsibleRepo, user.Id, tenderId)
	if err != nil {
		return entity.Tender{}, err
	}

	err = IsTenderExists(ts.tenderRepo, tenderId)
	if err != nil {
		return entity.Tender{}, err
	}

	tender, err := ts.tenderRepo.RollbackTenderVersion(context.Background(), tenderId, version)

	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Tender{}, ErrTenderOrVersionNotFound
		}
		return entity.Tender{}, err
	}

	return tender, nil

}
