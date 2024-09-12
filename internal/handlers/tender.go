package handlers

import (
	"avitoTech/internal/repo"
	"encoding/json"
	"fmt"
	"github.com/gorilla/schema"
	"io"
	"net/http"
	"strings"
)

type createTenderInput struct {
	Name            string `json:"name" validate:"required,min=3"`
	Description     string `json:"description" validate:"required"`
	ServiceType     string `json:"serviceType" validate:"required"`
	Status          string `json:"status" validate:"required"`
	OrganizationId  string `json:"organizationId" validate:"required"`
	CreatorUsername string `json:"creatorUsername" validate:"required"`
}

func CreateTender(repo *repo.Repositories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(1)
		// Чтение тела запроса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
			return
		}

		// Создаем объект Tender
		var t createTenderInput

		// Парсим JSON в структуру Tender
		err = json.Unmarshal(body, &t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u, err := repo.User.GetByName(r.Context(), t.CreatorUsername)

		//TODO: Добавить ошибку не существования пользователя
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		isResponsible, err := repo.Responsible.IsUserResponsibleForOrganization(r.Context(), u.Id, t.OrganizationId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//TODO: Добавить ошибку отсутствие прав для пользователя за организацию
		if isResponsible == false {
			http.Error(w, ErrUserIsNotResposible, http.StatusBadRequest)
			return
		}

		id, err := repo.Tender.New(r.Context(),
			t.Name,
			t.Description,
			strings.ToUpper(t.ServiceType),
			strings.ToUpper(t.Status),
			t.OrganizationId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(id)

	}
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
