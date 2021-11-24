package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"tfs-trading-bot/internal/domain"
	"tfs-trading-bot/internal/repository/queries"
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
	OrdersBySide(ctx context.Context, side string) ([]domain.Order, error)
	DeleteOrdersBySide(ctx context.Context, side string) error
}
