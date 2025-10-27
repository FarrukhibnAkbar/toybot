package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Postgres struct {
	Pool   *pgxpool.Pool
	Logger *zap.Logger
}

func NewPostgres(ctx context.Context, dsn string, logger *zap.Logger) (*Postgres, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	logger.Info("✅ Connected to PostgreSQL")
	return &Postgres{Pool: pool, Logger: logger}, nil
}

func (p *Postgres) Close() {
	p.Pool.Close()
}
