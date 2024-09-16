package pgrepo

import (
	"avitoTech/internal/entity"
	"avitoTech/internal/repo/repoerrs"
	"avitoTech/internal/storage/postgres"
	"context"
)

type BidRepo struct {
	*postgres.Postgres
}

func NewBidRepo(pg *postgres.Postgres) *BidRepo {
	return &BidRepo{pg}
}

func (r *BidRepo) CreateBid(ctx context.Context, name string, description string, tenderId string, authorType string, authorId string) (entity.Bid, error) {
	sql := `
	INSERT INTO bid (name, description, tender_id, author_type, author_id)
	VALUES ($1, $2, $3, UPPER($4)::authore_type, $5) 
	RETURNING id, name, description, tender_id, INITCAP(author_type::text), author_id, version, created_at
	`

	var bid entity.Bid
	err := r.Pool.QueryRow(ctx, sql, name, description, tenderId, authorType, authorId).Scan(
		&bid.Id,
		&bid.Name,
		&bid.Description,
		&bid.TenderId,
		&bid.AuthorType,
		&bid.AuthorId,
		&bid.Version,
		&bid.CreatedAt,
	)
	if err != nil {
		return entity.Bid{}, repoerrs.ErrUnableToInsert
	}

	return bid, nil
}
