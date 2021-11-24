package algorithms

import (
	"math"

	"github.com/sirupsen/logrus"

	"tfs-trading-bot/internal/domain"
)

type EMAAlgo struct {
	sellPeriod int
	logger     *logrus.Logger
}

func NewEMAAlgo(sellPeriod int, logger *logrus.Logger) EMAAlgo {
	return EMAAlgo{
		sellPeriod: sellPeriod,
		logger:     logger,
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
		average.period++
	} else {
		average = tickerMA{
			ma:     price,
			last:   price,
			period: 1,
		}
		tickersAverage[id] = average
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
			if average.last == price {
				// We want to skip tickers with the same market price to avoid a "sideways" chart
				continue
			}

			m.logger.Debugf("EMA: %s last:%f ma:%f price:%f", ticker.ProductId, average.last, average.ma, price)
			if m.sellPeriod < average.period {
				if price > average.ma && average.last < average.ma { // cross high
					out <- buy(ticker.ProductId, 1, more(ticker.Ask))
				}
				if price < average.ma && average.last > average.ma { // cross low
					out <- sell(ticker.ProductId, 1, less(ticker.Bid))
				}
			}

			average.last = price
			k := smoothingConstant(m.sellPeriod)
			average.ma = (price * k) + (average.ma * (1 - k))

			tickersAverage[ticker.ProductId] = average
		}
	}()

	return out
}

func smoothingConstant(period int) domain.Price {
	return domain.Price(2.0 / (1 + float64(period)))
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
