package algorithms

import (
	"github.com/sirupsen/logrus"
	"testing"
	"tfs-trading-bot/internal/domain"

	"github.com/stretchr/testify/assert"
)

type Test struct {
	price  domain.Price
	symbol domain.TickerSymbol
}

func newTestTicker(test Test) domain.Ticker {
	return domain.Ticker{
		ProductId: test.symbol,
		Bid:       test.price,
		Ask:       test.price,
		MarkPrice: test.price,
	}
}

func newTestOrder(symbol domain.TickerSymbol, price domain.Price, side string) domain.Order {
	return domain.Order{
		OrderType:  "ioc",
		Symbol:     symbol,
		Side:       side,
		Size:       1,
		LimitPrice: price,
	}
}

func TestProcessTickers(t *testing.T) {
	tickers := []Test{
		{91.13, "xbtusd"},
		{91.19, "xbtusd"},
		{91.15, "xbtusd"},
		{91.24, "xbtusd"},
		{91.16, "xbtusd"},
		{91.01, "xbtusd"},
		{91.06, "xbtusd"},
		{91.02, "xbtusd"},
		{90.96, "xbtusd"},
		{90.98, "xbtusd"},
		{90.97, "xbtusd"},
		{91.08, "xbtusd"},
		{91.13, "xbtusd"},
		{98.14, "xbtusd"},
		{93.02, "xbtusd"},
		{91.13, "xbtusd"},
		{91.03, "xbtusd"},
	}

	expect := []domain.Order{
		newTestOrder("xbtusd", 91, "buy"),
		newTestOrder("xbtusd", 93, "sell"),
	}

	log := logrus.New()
	log.SetLevel(logrus.TraceLevel)
	ema := NewEMAAlgo(5, log)

	in := make(chan domain.Ticker)

	go func() {
		defer close(in)
		for _, t := range tickers {
			in <- newTestTicker(t)
		}
	}()

	var actual []domain.Order
	for o := range ema.ProcessTickers(in) {
		actual = append(actual, o)
	}

	assert.EqualValues(t, expect, actual)
}
