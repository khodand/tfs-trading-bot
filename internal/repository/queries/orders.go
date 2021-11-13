package queries

import (
	"context"
	"tfs-trading-bot/internal/domain"
)

const insertOrdersQuery = `INSERT INTO orders (ordertype, symbol, side, size, limitprice, stopprice, triggersignal)
VALUES ($1, $2, $3, $4, $5, $6, $7)`

func (q *Queries) InsertOrder(ctx context.Context, order domain.Order) error {
	_, err := q.pool.Query(ctx, insertOrdersQuery, order.OrderType, order.Symbol, order.Side, order.Size,
		order.LimitPrice, order.StopPrice, order.TriggerSignal)
	return err
}
