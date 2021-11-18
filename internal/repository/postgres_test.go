package repository

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"tfs-trading-bot/internal/domain"

	pkgpostgres "tfs-trading-bot/pkg/postgres"
)

const dsn = "postgres://postgres:1234@localhost:5432/postgres" +
	"?sslmode=disable"

func TestOrderQueries(t *testing.T) {
	pool, err := pkgpostgres.NewPool(dsn)
	if err != nil {
		t.Skip("Skipping test. No connection to database.")
		return
	}

	repo := NewRepository(pool)
	err = repo.InsertOrder(context.Background(), domain.Order{
		OrderType:     "ioc",
		Symbol:        "pi_test",
		Side:          "test",
		Size:          10,
		LimitPrice:    domain.Price(100),
		StopPrice:     domain.Price(100),
		TriggerSignal: "test",
	})

	assert.NoError(t, err)
	orders, err := repo.OrdersBySide(context.Background(), "test")
	assert.NoError(t, err)
	fmt.Println(orders)
	assert.Equal(t, 1, len(orders))
	err = repo.DeleteOrdersBySide(context.Background(), "test")
	assert.NoError(t, err)
}
