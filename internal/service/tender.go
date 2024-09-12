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

func (s *TenderService) CreateTender(ct CreateTenderInput) (entity.Tender, error) {
	u, err := s.userRepo.GetByName(context.Background(), ct.CreatorUsername)

	if err != nil {
		if err == repoerrs.ErrNotExists {
			return entity.Tender{}, ErrUserNotExists
		}
		return entity.Tender{}, ErrCannotCreateTender
	}

	isResponsible, err := s.responsibleRepo.IsUserResponsibleForOrganization(context.Background(), u.Id, ct.OrganizationId)
	if err != nil {
		if err == repoerrs.ErrNotExists {
			return entity.Tender{}, ErrUserIsNotResposible
		}
		return entity.Tender{}, ErrCannotCreateTender
	}

	if isResponsible == false {
		return entity.Tender{}, ErrUserIsNotResposible
	}

	t, err := s.tenderRepo.CreateTender(context.Background(), ct.Name, ct.Description, ct.ServiceType, ct.Status, ct.OrganizationId)

	return t, err
}

func (s *TenderService) GetTenders(gtp GetTendersParams) ([]entity.Tender, error) {
	for i, st := range gtp.ServiceType {
		gtp.ServiceType[i] = strings.ToUpper(st)
	}

	tenders, err := s.tenderRepo.GetTenders(context.Background(), gtp.Limit, gtp.Offset, gtp.ServiceType)

	if err != nil {
		if err == repoerrs.ErrNotExists {
			return []entity.Tender{}, ErrTendersNotFound
		}
		return []entity.Tender{}, err
	}

	return tenders, nil
}

func (s *TenderService) GetUserTenders(gtp GetUserTendersParams) ([]entity.Tender, error) {
	_, err := s.userRepo.GetByName(context.Background(), gtp.Username)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return []entity.Tender{}, ErrUserNotExists
		}
		return []entity.Tender{}, err
	}

	return s.tenderRepo.GetUserTenders(context.Background(), gtp.Username, gtp.Limit, gtp.Offset)
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
//GetTenderStatus(w http.ResponseWriter, r *http.Request, tenderId TenderId, params GetTenderStatusParams)
//// Изменение статуса тендера
//// (PUT /tenders/{tenderId}/status)
//UpdateTenderStatus(w http.ResponseWriter, r *http.Request, tenderId TenderId, params UpdateTenderStatusParams)
