package services

import (
	"context"

	"github.com/sirupsen/logrus"

	"tfs-trading-bot/internal/domain"
	"tfs-trading-bot/internal/repository"
)

type TradingService interface {
	TradeTicker(symbol domain.TickerSymbol)
	ProcessOrders() <-chan domain.Order
	ChangeAlgo(algo TradingAlgorithm)
}

type TradingExchange interface {
	Subscribe(symbol domain.TickerSymbol)
	GetTickersChan() <-chan domain.Ticker
	SendOrder(order domain.Order) error
}

type TradingAlgorithm interface {
	ProcessTickers(tickers <-chan domain.Ticker) <-chan domain.Order
}

type Trader struct {
	exchange TradingExchange
	algo     TradingAlgorithm
	database repository.TradingDatabase
	log      *logrus.Logger
}

func (t *Trader) ChangeAlgo(algo TradingAlgorithm) {
	// TODO: implement this
	panic("implement me")
}

func NewTrader(exc TradingExchange, alg TradingAlgorithm, db repository.TradingDatabase, logger *logrus.Logger) *Trader {
	return &Trader{
		exchange: exc,
		algo:     alg,
		database: db,
		log:      logger,
	}
}

func (t *Trader) ProcessOrders() <-chan domain.Order {
	out := make(chan domain.Order)
	go func() {
		defer close(out)
		t.log.Info("Trader waits for tickers")
		for order := range t.algo.ProcessTickers(t.exchange.GetTickersChan()) {
			t.log.Debug("Sending order ot exchange:", order)
			if err := t.exchange.SendOrder(order); err == nil {
				t.log.Debug("Inserting order to database:", order)
				if err := t.database.InsertOrder(context.Background(), order); err != nil {
					t.log.Fatal(err)
				}
				out <- order
			}
		}
	}()

	return out
}

func (t *Trader) TradeTicker(symbol domain.TickerSymbol) {
	t.exchange.Subscribe(symbol)
}
