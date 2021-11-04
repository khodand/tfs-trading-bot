package main

import "tfs-trading-bot/internal/domain"

type TradingExchange interface {
	Subscribe(symbol domain.TickerSymbol)
	GetTickersChan() <-chan domain.Ticker
	SendOrder(order domain.Order)
}

type TradingAlgorithm interface {
	ProcessTickers(tickers <-chan domain.Ticker) <-chan domain.Order
}

type TradingDatabase interface {
	putOrder()

}

type Trader struct {
	exchange TradingExchange
	algo TradingAlgorithm
	//database TradingDatabase
}

func NewTrader() {
	var trader Trader
	trader.algo.ProcessTickers(trader.exchange.GetTickersChan())
}

func (t *Trader) processOrders(orders <-chan domain.Order) {
	go func() {
		for order := range orders {
			t.exchange.SendOrder(order)
		}
	}()
}

func (t *Trader) TradeTicker(symbol domain.TickerSymbol) {
	t.exchange.Subscribe(symbol)
}

func main() {

}
