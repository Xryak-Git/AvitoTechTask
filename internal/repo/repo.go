package repo

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo/pgrepo"
	"avitoTech/internal/storage/postgres"
	"context"
)

type Tender interface {
	CreateTender(ctx context.Context, name, description, serviceType, status, organizationId string) (entity.Tender, error)
	GetTenders(ctx context.Context, limit, offset int, serviceType []string) ([]entity.Tender, error)
}

type Bid interface {
}

type User interface {
	GetByName(ctx context.Context, username string) (entity.User, error)
}

type Responsible interface {
	GetAllResponsiblesByUserId(ctx context.Context, userId string) ([]entity.Responsible, error)
	IsUserResponsibleForOrganization(ctx context.Context, userId, organizationId string) (bool, error)
}

type Repositories struct {
	Tender
	Bid
	User
	Responsible
}

func NewRepos(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Tender:      pgrepo.NewTenderRepo(pg),
		User:        pgrepo.NewUserRepo(pg),
		Responsible: pgrepo.NewResponsibleRepo(pg),
	}
}
