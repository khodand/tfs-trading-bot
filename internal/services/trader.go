package main

import "tfs-trading-bot/internal/domain"

type TradingExchange interface {
	Subscribe(symbol domain.TickerSymbol)
	GetTickersChan() <-chan domain.Ticker
}

type Trader struct {
	exchange TradingExchange
}



func main() {

}
