package algorithms

import (
	"math"

	"github.com/sirupsen/logrus"

	"tfs-trading-bot/internal/domain"
)

type EMAAlgo struct {
	sellPeriod     int
	tickersAverage map[domain.TickerSymbol]tickerMA
	logger         *logrus.Logger
}

func NewEMAAlgo(sellPeriod int, logger *logrus.Logger) EMAAlgo {
	return EMAAlgo{
		sellPeriod:     sellPeriod,
		tickersAverage: make(map[domain.TickerSymbol]tickerMA),
		logger:         logger,
	}
}

type tickerMA struct {
	ma     domain.Price
	last   domain.Price
	period int
}

func getTickerMa(tickersAverage map[domain.TickerSymbol]tickerMA, id domain.TickerSymbol, price domain.Price) tickerMA {
	var average tickerMA
	if tmp, ok := tickersAverage[id]; ok {
		average = tmp
	} else {
		average = tickerMA{
			ma:     price,
			last:   price,
			period: 0,
		}
	}
	return average
}

func (m EMAAlgo) ProcessTickers(tickers <-chan domain.Ticker) <-chan domain.Order {
	tickersAverage := make(map[domain.TickerSymbol]tickerMA)
	out := make(chan domain.Order)

	go func() {
		defer close(out)
		for ticker := range tickers {
			price := ticker.MarkPrice
			average := getTickerMa(tickersAverage, ticker.ProductId, price)
			tickersAverage[ticker.ProductId] = average
			if average.last == price {
				// We want to skip tickers with the same market price to avoid a "sideways" chart
				continue
			}
			average.period++

			m.logger.Debugf("EMA: %s last:%f ma:%f price:%f", ticker.ProductId, average.last, average.ma, price)
			order := generateOrder(m.sellPeriod, average.period, average.last, average.ma, ticker)
			if (order != domain.Order{}) {
				out <- order
			}

			average.last = price
			average.ma = emaFormula(price, average.ma, m.sellPeriod)
			tickersAverage[ticker.ProductId] = average
		}
	}()

	return out
}

func (m EMAAlgo) ProcessTicker(ticker domain.Ticker) (domain.Order, bool) {
	price := ticker.MarkPrice
	average := getTickerMa(m.tickersAverage, ticker.ProductId, price)
	m.tickersAverage[ticker.ProductId] = average
	if average.last == price {
		// We want to skip tickers with the same market price to avoid a "sideways" chart
		return domain.Order{}, true
	}
	average.period++

	m.logger.Debugf("EMA: %s last:%f ma:%f price:%f", ticker.ProductId, average.last, average.ma, price)
	order := generateOrder(m.sellPeriod, average.period, average.last, average.ma, ticker)

	average.last = price
	average.ma = emaFormula(price, average.ma, m.sellPeriod)
	m.tickersAverage[ticker.ProductId] = average

	return order, order == domain.Order{}
}

func generateOrder(sellPeriod int, period int, last domain.Price, ma domain.Price, ticker domain.Ticker) domain.Order {
	if sellPeriod < period {
		if ticker.MarkPrice > ma && last < ma { // cross high
			return buy(ticker.ProductId, 1, more(ticker.Ask))
		}
		if ticker.MarkPrice < ma && last > ma { // cross low
			return sell(ticker.ProductId, 1, less(ticker.Bid))
		}
	}
	return domain.Order{}
}

func smoothingConstant(period int) domain.Price {
	return domain.Price(2.0 / (1 + float64(period)))
}

func emaFormula(price domain.Price, ma domain.Price, sellPeriod int) domain.Price {
	k := smoothingConstant(sellPeriod)
	return (price * k) + (ma * (1 - k))
}

func sell(symbol domain.TickerSymbol, size int, price domain.Price) domain.Order {
	return domain.Order{
		OrderType:     "ioc",
		Symbol:        symbol,
		Side:          "sell",
		Size:          size,
		LimitPrice:    price,
		StopPrice:     0,
		TriggerSignal: "",
	}
}

func buy(symbol domain.TickerSymbol, size int, price domain.Price) domain.Order {
	return domain.Order{
		OrderType:     "ioc",
		Symbol:        symbol,
		Side:          "buy",
		Size:          size,
		LimitPrice:    price,
		StopPrice:     0,
		TriggerSignal: "",
	}
}

func more(price domain.Price) domain.Price {
	return domain.Price(math.Round(float64(price + (price / 1000))))
}

func less(price domain.Price) domain.Price {
	return domain.Price(math.Round(float64(price - (price / 1000))))
}
