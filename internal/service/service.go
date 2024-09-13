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

type PatchTenderInput struct {
	Name            string `json:"name" validate:"omitempty,min=3"`
	Description     string `json:"description" validate:"omitempty"`
	ServiceType     string `json:"serviceType" validate:"omitempty"`
	Status          string `json:"status" validate:"omitempty"`
	OrganizationId  string `json:"organizationId" validate:"omitempty"`
	CreatorUsername string `json:"creatorUsername" validate:"omitempty"`
}

type GetParams struct {
	Limit  int `schema:"limit" validate:"required,gt=0"`
	Offset int `schema:"offset" validate:"omitempty,gte=0"`
}

type GetTendersParams struct {
	GetParams
	ServiceType []string `schema:"service_type" validate:"required,dive,alpha"`
}

type GetUserTendersParams struct {
	GetParams
	Username string `schema:"username" validate:"required,alpha"`
}

type UserParam struct {
	Username string `schema:"username" validate:"required,alpha"`
}

type UpdateTenderStatusParams struct {
	UserParam
	Status string `json:"status" validate:"required"`
}

type Tender interface {
	CreateTender(input CreateTenderInput) (entity.Tender, error)
	GetTenders(gtp GetTendersParams) ([]entity.Tender, error)
	GetUserTenders(gutp GetUserTendersParams) ([]entity.Tender, error)
	GetTenderStatus(up UserParam, tenderId string) (string, error)
	PathTender(up UserParam, tenderId string, pti PatchTenderInput) (entity.Tender, error)
	UpdateTenderStatus(utsp UpdateTenderStatusParams, id string) (entity.Tender, error)
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
