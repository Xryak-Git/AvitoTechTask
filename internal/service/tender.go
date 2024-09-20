package service

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo"
	"avitoTech/internal/repo/repoerrs"
	"context"
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
	u, err := ts.userRepo.GetByName(context.Background(), params.CreatorUsername)

	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Tender{}, ErrUserNotExists
		}
		return entity.Tender{}, ErrCannotCreateTender
	}

	_, err = ts.responsibleRepo.IsUserResponsibleForOrganizationByOrganizationId(context.Background(), u.Id, params.OrganizationId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Tender{}, ErrUserIsNotResposible
		}
		return entity.Tender{}, ErrCannotCreateTender
	}

	t, err := ts.tenderRepo.CreateTender(context.Background(), params.Name, params.Description, params.ServiceType, params.Status, params.OrganizationId)

	return t, err
}

func (ts *TenderService) GetTenders(params GetTendersParams) ([]entity.Tender, error) {
	for i, st := range params.ServiceType {
		params.ServiceType[i] = strings.ToUpper(st)
	}

	tenders, err := ts.tenderRepo.GetTenders(context.Background(), params.Limit, params.Offset, params.ServiceType)

	if err != nil {
		if err == repoerrs.ErrNotFound {
			return []entity.Tender{}, nil
		}
		return []entity.Tender{}, err
	}

	return tenders, nil
}

func (ts *TenderService) GetUserTenders(params GetUserTendersParams) ([]entity.Tender, error) {
	_, err := ts.userRepo.GetByName(context.Background(), params.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return []entity.Tender{}, ErrUserNotExists
		}
		return []entity.Tender{}, err
	}

	tenders, err := ts.tenderRepo.GetUserTenders(context.Background(), params.Username, params.Limit, params.Offset)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return []entity.Tender{}, ErrTendersNotFound
		}
		return []entity.Tender{}, err
	}

	return tenders, nil
}

func (ts *TenderService) GetTenderStatus(params UserParam, tenderId string) (string, error) {

	user, err := ts.userRepo.GetByName(context.Background(), params.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return "", ErrUserNotExists
		}
		return "", err
	}

	exists, err := ts.tenderRepo.IsTenderExists(context.Background(), tenderId)
	if err != nil || !exists {
		return "", ErrTenderNotFound
	}

	status, err := ts.tenderRepo.GetTenderStatus(context.Background(), tenderId)
	if err != nil {
		return "", err
	}

	if status == "Published" {
		return status, nil
	}

	isResponsible, err := ts.responsibleRepo.IsUserResponsibleForOrganizationByTenderId(context.Background(), user.Id, tenderId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return "", ErrUserIsNotResposible
		}
		return "", ErrCannotGetTenderStatus
	}

	if isResponsible == false {
		return "", ErrUserIsNotResposible
	}

	return status, nil

}

func (ts *TenderService) EditTender(params UserParam, tenderId string, editFields map[string]interface{}) (entity.Tender, error) {
	user, err := ts.userRepo.GetByName(context.Background(), params.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Tender{}, ErrUserNotExists
		}
		return entity.Tender{}, err
	}

	exists, err := ts.tenderRepo.IsTenderExists(context.Background(), tenderId)
	if err != nil || !exists {
		return entity.Tender{}, ErrTenderNotFound
	}

	_, err = ts.responsibleRepo.IsUserResponsibleForOrganizationByTenderId(context.Background(), user.Id, tenderId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Tender{}, ErrUserIsNotResposible
		}
		return entity.Tender{}, err
	}

	tender, err := ts.tenderRepo.UpdateTender(context.Background(), tenderId, editFields)
	if err != nil {
		return entity.Tender{}, err
	}

	return tender, nil
}

func (ts *TenderService) UpdateTenderStatus(params UpdateTenderStatusParams, tenderId string) (entity.Tender, error) {

	user, err := ts.userRepo.GetByName(context.Background(), params.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Tender{}, ErrUserNotExists
		}
		return entity.Tender{}, err
	}

	exists, err := ts.tenderRepo.IsTenderExists(context.Background(), tenderId)
	if err != nil || !exists {
		return entity.Tender{}, ErrTenderNotFound
	}

	isResponsibe, err := ts.responsibleRepo.IsUserResponsibleForOrganizationByTenderId(context.Background(), user.Id, tenderId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Tender{}, ErrUserIsNotResposible
		}
		return entity.Tender{}, err
	}

	if !isResponsibe {
		return entity.Tender{}, ErrUserIsNotResposible
	}

	t, err := ts.tenderRepo.UpdateTenderStatus(context.Background(), params.Status, tenderId)

	if err != nil {
		return entity.Tender{}, err
	}

	return t, nil

}

func (ts *TenderService) RollbackTender(params UserParam, tenderId string, version int) (entity.Tender, error) {
	user, err := ts.userRepo.GetByName(context.Background(), params.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Tender{}, ErrUserNotExists
		}
		return entity.Tender{}, err
	}

	isResponsibe, err := ts.responsibleRepo.IsUserResponsibleForOrganizationByTenderId(context.Background(), user.Id, tenderId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Tender{}, ErrUserIsNotResposible
		}
		return entity.Tender{}, err
	}
	if !isResponsibe {
		return entity.Tender{}, ErrUserIsNotResposible
	}

	exists, err := ts.tenderRepo.IsTenderExists(context.Background(), tenderId)
	if err != nil || !exists {
		return entity.Tender{}, ErrTenderOrVersionNotFound
	}

	tender, err := ts.tenderRepo.RollbackTenderVersion(context.Background(), tenderId, version)

	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Tender{}, ErrTenderOrVersionNotFound
		}
		return entity.Tender{}, err
	}

	return tender, nil

}
