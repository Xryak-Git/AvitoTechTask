package handlers

import (
	"avitoTech/internal/repo"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type createTenderInput struct {
	Name            string `json:"name" validate:"required,min=3"`
	Description     string `json:"description" validate:"required"`
	ServiceType     string `json:"serviceType" validate:"required"`
	Status          string `json:"status" validate:"required"`
	OrganizationId  int    `json:"organizationId" validate:"required"`
	CreatorUsername string `json:"creatorUsername" validate:"required"`
}

func CreateTender(repo *repo.Repositories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		id, err := repo.Tender.New(r.Context(), t.Name, t.Description, t.ServiceType, t.Status, t.OrganizationId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(id)

	}
}
