package services

import (
	"context"
	"log"
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
	SendOrder(order domain.Order)
}

type TradingAlgorithm interface {
	ProcessTickers(tickers <-chan domain.Ticker) <-chan domain.Order
}

type Trader struct {
	exchange TradingExchange
	algo     TradingAlgorithm
	database repository.TradingDatabase
}

func (t *Trader) ChangeAlgo(algo TradingAlgorithm) {
	// TODO: implement this
	panic("implement me")
}

func NewTrader(exch TradingExchange, alg TradingAlgorithm, db repository.TradingDatabase) *Trader {
	return &Trader{
		exchange: exch,
		algo:     alg,
		database: db,
	}
}

func (t *Trader) ProcessOrders() <-chan domain.Order {
	out := make(chan domain.Order)
	go func() {
		defer close(out)
		for order := range t.algo.ProcessTickers(t.exchange.GetTickersChan()) {
			log.Println(order)
			t.exchange.SendOrder(order)
			err := t.database.InsertOrder(context.Background(), order)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Inserted to database")
			out <- order
		}
	}()

	return out
}

func (t *Trader) TradeTicker(symbol domain.TickerSymbol) {
	t.exchange.Subscribe(symbol)
}
