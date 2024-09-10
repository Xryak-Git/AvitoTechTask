package repo

import (
	"avitoTech/internal/repo/pgrepo"
	"avitoTech/internal/storage/postgres"
	"context"
)

type Tender interface {
	New(ctx context.Context, name, description, serviceType, status string, organizationId int) (int, error)
}

type Bid interface {
} // TODO: Add interfaces

type Repositories struct {
	Tender
	Bid
}

func NewRepos(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Tender: pgrepo.NewTenderRepo(pg),
	}
}
