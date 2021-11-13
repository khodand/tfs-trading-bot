package repository

import (
	"context"
	"tfs-trading-bot/internal/domain"
	"tfs-trading-bot/internal/repository/queries"

	"github.com/jackc/pgx/v4/pgxpool"
)

type repo struct {
	*queries.Queries
	pool *pgxpool.Pool
}

func NewRepository(pgxPool *pgxpool.Pool) TradingDatabase {
	return &repo{
		Queries: queries.New(pgxPool),
		pool:    pgxPool,
	}
}

type TradingDatabase interface {
	InsertOrder(ctx context.Context, order domain.Order) error
}
