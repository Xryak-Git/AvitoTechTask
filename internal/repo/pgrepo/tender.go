package pgrepo

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo"
	"avitoTech/internal/storage/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	log "log/slog"
)

type TenderRepo struct {
	*postgres.Postgres
}

func NewTenderRepo(pg *postgres.Postgres) *TenderRepo {
	return &TenderRepo{pg}
}

func (r *TenderRepo) New(ctx context.Context, name, description, serviceType, status, organizationId string) (string, error) {
	const fn = "repo.pgrepo.tender.New"

	sql := `
	INSERT INTO tender (name, description, service_type, status, organization_id)
	VALUES ($1, $2, $3::service_type, $4::tender_status, $5) 
	RETURNING id, name, description, service_type, status, organization_id, version, created_at
	`

	var t entity.Tender
	err := r.Pool.QueryRow(ctx, sql, name, description, serviceType, status, organizationId).Scan(
		&t.Id,
		&t.Name,
		&t.Description,
		&t.ServiceType,
		&t.Status,
		&t.OrganizationId,
		&t.Version,
		&t.CreatedAt,
	)

	if err != nil {
		log.Info("err: ", err)
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return "", repo.ErrAlreadyExists
			}
		}
		return "", fmt.Errorf("%s: %v", fn, err)
	}

	log.Info("New tender: ", "tender", t)

	return t.Id, nil

}
