package service

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo"
)

type CreateTenderInput struct {
	Name            string `json:"name" validate:"required,min=3"`
	Description     string `json:"description" validate:"required"`
	ServiceType     string `json:"serviceType" validate:"required"`
	Status          string `json:"status" validate:"required"`
	OrganizationId  string `json:"organizationId" validate:"required"`
	CreatorUsername string `json:"creatorUsername" validate:"required"`
}

type GetParams struct {
	Limit  int `schema:"limit"`
	Offset int `schema:"offset"`
}

type GetTendersParams struct {
	GetParams
	ServiceType []string `schema:"service_type"`
}

type GetUserTendersParams struct {
	GetParams
	Username string `schema:"username"`
}

type Tender interface {
	CreateTender(input CreateTenderInput) (entity.Tender, error)
	GetTenders(gtp GetTendersParams) ([]entity.Tender, error)
	GetUserTenders(gtp GetUserTendersParams) ([]entity.Tender, error)
}

type Bid interface {
}

type Services struct {
	Tender Tender
	Bid    Bid
}

func NewServices(r *repo.Repositories) *Services {
	return &Services{
		Tender: NewTenderService(r.Tender, r.User, r.Responsible),
	}
}
