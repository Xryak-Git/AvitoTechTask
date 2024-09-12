package service

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo"
	"strings"
)

type CreateTenderInput struct {
	Name            string `json:"name" validate:"required,min=3"`
	Description     string `json:"description" validate:"required"`
	ServiceType     string `json:"serviceType" validate:"required"`
	Status          string `json:"status" validate:"required"`
	OrganizationId  string `json:"organizationId" validate:"required"`
	CreatorUsername string `json:"creatorUsername" validate:"required"`
}

func (t *CreateTenderInput) toUpper() {
	t.ServiceType = strings.ToUpper(t.ServiceType)
	t.Status = strings.ToUpper(t.Status)
}

type Tender interface {
	CreateTender(input CreateTenderInput) (entity.Tender, error)
}

type Services struct {
	Tender Tender
}

func NewServices(r *repo.Repositories) *Services {
	return &Services{
		Tender: NewTenderService(r.Tender, r.User, r.Responsible),
	}
}
