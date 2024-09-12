package pgrepo

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo/repoerrs"
	"avitoTech/internal/storage/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	log "log/slog"
)

type TenderRepo struct {
	*postgres.Postgres
}

func NewTenderRepo(pg *postgres.Postgres) *TenderRepo {
	return &TenderRepo{pg}
}

func (r *TenderRepo) CreateTender(ctx context.Context, name, description, serviceType, status, organizationId string) (entity.Tender, error) {
	const fn = "repo.pgrepo.tender.CreateTender"

	sql := `
	INSERT INTO tender (name, description, service_type, status, organization_id)
	VALUES ($1, $2, UPPER($3)::service_type, UPPER($4)::tender_status, $5) 
	RETURNING id, name, description, INITCAP(service_type::text) AS service_type, INITCAP(status::text) AS status, organization_id, version, created_at
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
				return t, repoerrs.ErrAlreadyExists
			}
		}
		return t, fmt.Errorf("%s: %v", fn, err)
	}

	log.Debug("CreateTender tender: ", "tender", t)

	return t, nil

}

func (r *TenderRepo) GetTenders(ctx context.Context, limit, offset int, serviceType []string) ([]entity.Tender, error) {
	const fn = "repo.pgrepo.tender.GetTenders"

	sql := `
	SELECT id, name, description, INITCAP(service_type::text) AS service_type, INITCAP(status::text) AS status, organization_id, version, created_at
	FROM tender
	WHERE service_type::text = ANY($1)
	LIMIT $2
	OFFSET $3
	`

	rows, err := r.Pool.Query(ctx, sql, serviceType, limit, offset)

	if err != nil {
		log.Debug("err: ", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return []entity.Tender{}, repoerrs.ErrNotFound
		}
		return []entity.Tender{}, fmt.Errorf("%s: %v", fn, err)
	}

	defer rows.Close()

	var tenders []entity.Tender
	for rows.Next() {
		var t entity.Tender
		err := rows.Scan(
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
			return []entity.Tender{}, fmt.Errorf("%s: %v", err)
		}
		tenders = append(tenders, t)
	}

	log.Debug("GetTenders: ", tenders)

	return tenders, nil
}

func (r *TenderRepo) GetUserTenders(ctx context.Context, username string, limit int, offset int) ([]entity.Tender, error) {
	const fn = "repo.pgrepo.tender.GetUserTenders"

	sql := `
		SELECT id, name, description, INITCAP(service_type::text) AS service_type, INITCAP(status::text) AS status, organization_id, version, created_at
		FROM tender
		WHERE organization_id in (
			SELECT o.id
			FROM organization_responsible ores
					 JOIN organization o ON ores.organization_id = o.id
					 JOIN employee e ON ores.user_id = e.id
			WHERE e.username = $1)
		LIMIT $2
		OFFSET $3
		`

	rows, err := r.Pool.Query(ctx, sql, username, limit, offset)

	if err != nil {
		log.Debug("err: ", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return []entity.Tender{}, repoerrs.ErrNotFound
		}
		return []entity.Tender{}, fmt.Errorf("%s: %v", fn, err)
	}

	defer rows.Close()

	var tenders []entity.Tender
	for rows.Next() {
		var t entity.Tender
		err := rows.Scan(
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
			return []entity.Tender{}, fmt.Errorf("%s: %v", fn, err)
		}
		tenders = append(tenders, t)
	}

	log.Debug("GetUserTenders: ", tenders)

	return tenders, nil
}
