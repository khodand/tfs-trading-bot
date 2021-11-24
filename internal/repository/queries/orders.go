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

const selectOrdersBySideQuery = `SELECT ordertype, symbol, side, size, limitprice, stopprice,
		triggersignal FROM orders WHERE side = $1`

func (q *Queries) OrdersBySide(ctx context.Context, side string) ([]domain.Order, error) {
	rows, err := q.pool.Query(ctx, selectOrdersBySideQuery, side)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		var limitPrice float64
		var stopPrice float64
		err = rows.Scan(&o.OrderType, &o.Symbol, &o.Side, &o.Size, &limitPrice, &stopPrice, &o.TriggerSignal)
		o.LimitPrice = domain.Price(limitPrice)
		o.StopPrice = domain.Price(stopPrice)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

const deleteOrdersBySideQuery = `DELETE FROM orders WHERE side = $1`

func (q *Queries) DeleteOrdersBySide(ctx context.Context, side string) error {
	_, err := q.pool.Query(ctx, deleteOrdersBySideQuery, side)
	return err
}
