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
	Stop()
}

type TradingExchange interface {
	Subscribe(symbol domain.TickerSymbol)
	GetTickersChan() <-chan domain.Ticker
	SendOrder(order domain.Order) error
	Stop()
}

type TradingAlgorithm interface {
	ProcessTicker(ticker domain.Ticker) (domain.Order, bool)
}

type Trader struct {
	exchange   TradingExchange
	algo       TradingAlgorithm
	database   repository.TradingDatabase
	changeAlgo chan TradingAlgorithm
	log        *logrus.Logger
}

func (t *Trader) ChangeAlgo(algo TradingAlgorithm) {
	t.changeAlgo <- algo
}

func (t *Trader) Stop() {
	t.exchange.Stop()
}

func NewTrader(exc TradingExchange, alg TradingAlgorithm, db repository.TradingDatabase, logger *logrus.Logger) *Trader {
	return &Trader{
		exchange:   exc,
		algo:       alg,
		database:   db,
		log:        logger,
		changeAlgo: make(chan TradingAlgorithm),
	}
}

func (t *Trader) ProcessOrders() <-chan domain.Order {
	out := make(chan domain.Order)
	go func() {
		defer close(out)
		t.log.Info("Trader waits for tickers")
		for order := range t.tickersToAlgo(t.exchange.GetTickersChan()) {
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

func (t *Trader) tickersToAlgo(tickers <-chan domain.Ticker) <-chan domain.Order {
	out := make(chan domain.Order)
	go func() {
		defer close(out)
		for {
			select {
			case algo := <-t.changeAlgo:
				t.algo = algo
			case ticker, ok := <-tickers:
				if !ok {
					return
				}
				if order, skip := t.algo.ProcessTicker(ticker); !skip {
					out <- order
				}
			}
		}
	}()

	return out
}

func (t *Trader) TradeTicker(symbol domain.TickerSymbol) {
	t.exchange.Subscribe(symbol)
}
