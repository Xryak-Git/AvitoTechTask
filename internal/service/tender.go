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

func (ts *TenderService) CreateTender(ct CreateTenderInput) (entity.Tender, error) {
	u, err := ts.userRepo.GetByName(context.Background(), ct.CreatorUsername)

	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Tender{}, ErrUserNotExists
		}
		return entity.Tender{}, ErrCannotCreateTender
	}

	_, err = ts.responsibleRepo.IsUserResponsibleForOrganizationByOrganizationId(context.Background(), u.Id, ct.OrganizationId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Tender{}, ErrUserIsNotResposible
		}
		return entity.Tender{}, ErrCannotCreateTender
	}

	t, err := ts.tenderRepo.CreateTender(context.Background(), ct.Name, ct.Description, ct.ServiceType, ct.Status, ct.OrganizationId)

	return t, err
}

func (ts *TenderService) GetTenders(gtp GetTendersParams) ([]entity.Tender, error) {
	for i, st := range gtp.ServiceType {
		gtp.ServiceType[i] = strings.ToUpper(st)
	}

	tenders, err := ts.tenderRepo.GetTenders(context.Background(), gtp.Limit, gtp.Offset, gtp.ServiceType)

	if err != nil {
		if err == repoerrs.ErrNotFound {
			return []entity.Tender{}, nil
		}
		return []entity.Tender{}, err
	}

	return tenders, nil
}

func (ts *TenderService) GetUserTenders(gutp GetUserTendersParams) ([]entity.Tender, error) {
	_, err := ts.userRepo.GetByName(context.Background(), gutp.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return []entity.Tender{}, ErrUserNotExists
		}
		return []entity.Tender{}, err
	}

	return ts.tenderRepo.GetUserTenders(context.Background(), gutp.Username, gutp.Limit, gutp.Offset)
}

func (ts *TenderService) GetTenderStatus(u UserParam, tenderId string) (string, error) {

	user, err := ts.userRepo.GetByName(context.Background(), u.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return "", ErrUserNotExists
		}
		return "", err
	}

	status, err := ts.tenderRepo.GetTenderStatus(context.Background(), tenderId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return "", ErrTenderNotFound
		}
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

func (ts *TenderService) EditTender(up UserParam, tenderId string, params map[string]interface{}) (entity.Tender, error) {
	user, err := ts.userRepo.GetByName(context.Background(), up.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Tender{}, ErrUserNotExists
		}
		return entity.Tender{}, err
	}

	exists, err := ts.tenderRepo.IsTenderExists(context.Background(), tenderId)
	if err != nil {
		return entity.Tender{}, err
	}

	if !exists {
		return entity.Tender{}, ErrTenderNotFound
	}

	_, err = ts.responsibleRepo.IsUserResponsibleForOrganizationByTenderId(context.Background(), user.Id, tenderId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Tender{}, ErrUserIsNotResposible
		}
		return entity.Tender{}, err
	}

	tender, err := ts.tenderRepo.UpdateTender(context.Background(), tenderId, params)
	if err != nil {
		return entity.Tender{}, err
	}

	return tender, nil
}

func (ts *TenderService) UpdateTenderStatus(utsp UpdateTenderStatusParams, tenderId string) (entity.Tender, error) {

	_, err := ts.userRepo.GetByName(context.Background(), utsp.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Tender{}, ErrUserNotExists
		}
		return entity.Tender{}, err
	}

	isResponsibe, err := ts.responsibleRepo.IsUserResponsibleForOrganizationByTenderId(context.Background(), utsp.Username, tenderId)

	if !isResponsibe {
		return entity.Tender{}, ErrUserIsNotResposible
	}

	t, err := ts.tenderRepo.UpdateTenderStatus(context.Background(), utsp.Status, tenderId)

	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Tender{}, ErrTenderNotFound
		}
		return entity.Tender{}, err
	}

	return t, nil

}

func (ts *TenderService) RollbackTender(u UserParam, tenderId string, version int) (entity.Tender, error) {
	_ = u

	err := ts.tenderRepo.RollbackTenderVersion(context.Background(), tenderId, version)

	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Tender{}, ErrTenderNotFound
		}
		return entity.Tender{}, err
	}

	t, err := ts.tenderRepo.GetTenderById(context.Background(), tenderId)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Tender{}, ErrTenderNotFound
		}
		return entity.Tender{}, err
	}

	return t, nil

}

//// Получить тендеры пользователя
//// (GET /tenders/my)
//GetUserTenders(w http.ResponseWriter, r *http.Request, params GetUserTendersParams)
//// Создание нового тендера
//// (POST /tenders/new)
//
//// Редактирование тендера
//// (PATCH /tenders/{tenderId}/edit)
//EditTender(w http.ResponseWriter, r *http.Request, tenderId TenderId, params EditTenderParams)
//// Откат версии тендера
//// (PUT /tenders/{tenderId}/rollback/{version})
//RollbackTender(w http.ResponseWriter, r *http.Request, tenderId TenderId, version int32, params RollbackTenderParams)
//// Получение текущего статуса тендера
//// (GET /tenders/{tenderId}/status)
//GetTenderStatus(w http.ResponseWriter, r *http.Request, tenderId TenderId, params UserParam)
//// Изменение статуса тендера
//// (PUT /tenders/{tenderId}/status)
//UpdateTenderStatus(w http.ResponseWriter, r *http.Request, tenderId TenderId, params UpdateTenderStatusParams)
