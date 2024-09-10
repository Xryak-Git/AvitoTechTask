package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const (
	defaultMaxPoolSize  = 1
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

type PgxPool interface {
	Close()
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Ping(ctx context.Context) error
}

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Pool PgxPool
}

func New(url string, opts ...Option) (*Postgres, error) {
	const fn = "storage.postgres.New"

	pg := &Postgres{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(pg)
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", err)
	}

	err = pg.Ping()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", err)
	}

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

func (p *Postgres) Ping() error {
	err := p.Pool.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (p Postgres) TendersNew(ctx context.Context, name, description, serviceType, status string, organizationId int, creatorUsername string) (int, error) {
	const fn = "storage.postgres.TendersNew"

	sql := `
	INSERT INTO tender (name, description, service_type, status, organization_id)
	VALUES ($1, $2, $3::service_type, $4::tender_status, 1)
	`

	var id int
	err := p.Pool.QueryRow(ctx, sql).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("ProductRepo.CreateProduct - r.Pool.QueryRow: %v", err)
	}

	return id, nil

}
