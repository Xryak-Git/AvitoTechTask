package service

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo"
	"avitoTech/internal/repo/repoerrs"
	"context"
	"fmt"
	"github.com/gorilla/schema"
	"net/http"
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
	ct.toUpper()

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
	t.Capitalize()

	return t, err
}

type GetTendersParams struct {
	Limit       int      `schema:"limit"`
	Offset      int      `schema:"offset"`
	ServiceType []string `schema:"service_type"`
}

func GetTenders(repo *repo.Repositories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Не удалост", http.StatusBadGateway)
		}

		// Создаем объект Tender
		params := new(GetTendersParams)
		if err := schema.NewDecoder().Decode(params, r.Form); err != nil {
			http.Error(w, "Не удалост", http.StatusBadGateway)
		}
		fmt.Println(params)
	}
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
