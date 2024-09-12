package pgrepo

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo/repoerrs"
	"avitoTech/internal/storage/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type ResponsibleRepo struct {
	*postgres.Postgres
}

func NewResponsibleRepo(pg *postgres.Postgres) *ResponsibleRepo {
	return &ResponsibleRepo{pg}
}

func (r *ResponsibleRepo) GetAllResponsiblesByUserId(ctx context.Context, userId string) ([]entity.Responsible, error) {
	const fn = "repo.pgrepo.responsible.GetAllResponsiblesByUserId"

	sql := `
	SELECT *
	FROM organization_responsible
	WHERE user_id=$1::uuid
	`

	rows, err := r.Pool.Query(ctx, sql, userId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []entity.Responsible{}, repoerrs.ErrNotFound
		}
		return []entity.Responsible{}, fmt.Errorf("%s: %v", fn, err)
	}

	defer rows.Close()

	var responsibles []entity.Responsible
	for rows.Next() {
		var responsible entity.Responsible
		err := rows.Scan(
			&responsible.Id,
			&responsible.OrganizationId,
			&responsible.UserId,
		)
		if err != nil {
			return []entity.Responsible{}, fmt.Errorf("%s: %v", err)
		}
		responsibles = append(responsibles, responsible)
	}

	return responsibles, nil
}

func (r *ResponsibleRepo) IsUserResponsibleForOrganization(ctx context.Context, userId, organizationId string) (bool, error) {
	const fn = "repo.pgrepo.responsible.IsUserResponsibleForOrganization"

	responsibles, err := r.GetAllResponsiblesByUserId(ctx, userId)

	if err != nil {
		if err == repoerrs.ErrNotFound {
			return false, nil
		}
		return false, fmt.Errorf("%s: %v", fn, err)
	}

	for _, responsible := range responsibles {
		if responsible.OrganizationId == organizationId {
			return true, nil
		}
	}

	return false, nil
}
