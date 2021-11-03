package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"tfs-trading-bot/internal/domain"
	"tfs-trading-bot/pkg"
	pkgpostgres "tfs-trading-bot/pkg/postgres"
)

func TestOrderQueries(t *testing.T) {
	config, err := pkg.ReadConfig("../../config.json")
	if err != nil {
		t.Skip("Skipping test. No config file.")
		return
	}
	pool, err := pkgpostgres.NewPool(config.Dsn, logrus.New())
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
