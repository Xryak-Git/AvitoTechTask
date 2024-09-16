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
	GetUserTenders(ctx context.Context, username string, limit int, offset int) ([]entity.Tender, error)
	GetTenderStatus(ctx context.Context, tenderId string) (string, error)
	UpdateTender(ctx context.Context, tenderId string, params map[string]interface{}) (entity.Tender, error)
	UpdateTenderStatus(ctx context.Context, tenderId, status string) (entity.Tender, error)
	RollbackTenderVersion(ctx context.Context, tenderId string, version int) error
	GetTenderById(ctx context.Context, tenderId string) (entity.Tender, error)
}

type Bid interface {
	CreateBid(ctx context.Context, name string, description string, tenderId string, authorType string, authorId string) (entity.Bid, error)
	GetUserBids(ctx context.Context, username string, limit int, offset int) ([]entity.Bid, error)
}

type User interface {
	GetByName(ctx context.Context, username string) (entity.User, error)
}

type Responsible interface {
	GetAllResponsiblesByUserId(ctx context.Context, userId string) ([]entity.Responsible, error)
	IsUserResponsibleForOrganizationByTenderId(ctx context.Context, userId, organizationId string) (bool, error)
	IsUserResponsibleForOrganizationByOrganizationId(ctx context.Context, userId, organizationId string) (bool, error)
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
		Bid:         pgrepo.NewBidRepo(pg),
	}
}
