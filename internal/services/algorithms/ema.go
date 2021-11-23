package algorithms

import (
	"log"
	"math"
	"tfs-trading-bot/internal/domain"
)

type EMAAlgo struct {
	sellPeriod int
}

func NewEMAAlgo(sellPeriod int) EMAAlgo {
	return EMAAlgo{
		sellPeriod: sellPeriod,
	}
}

type tickerMA struct {
	ma     domain.Price
	last   domain.Price
	period int
}

func (m EMAAlgo) ProcessTickers(tickers <-chan domain.Ticker) <-chan domain.Order {
	tickersAverage := make(map[domain.TickerSymbol]tickerMA)
	out := make(chan domain.Order)

	go func() {
		defer close(out)
		for ticker := range tickers {
			log.Println("ema", ticker)
			var average tickerMA
			price := ticker.MarkPrice
			if tmp, ok := tickersAverage[ticker.ProductId]; ok {
				average = tmp
				average.period++
				k := smoothingConstant(m.sellPeriod)
				average.ma = (price * k) + (average.ma * (1 - k))
			} else {
				average = tickerMA{
					ma:     price,
					last:   price,
					period: 1,
				}
			}
			log.Println("EMA:", average.last, average.ma, price)

			if m.sellPeriod < average.period {
				if price > average.ma && average.last < average.ma { // cross high
					log.Println("EMA BUY")
					out <- buy(ticker.ProductId, 1, more(ticker.Ask))
				}
				if price < average.ma && average.last > average.ma { // cross low
					log.Println("EMA SELL")
					out <- sell(ticker.ProductId, 1, less(ticker.Bid))
				}
			}
			average.last = price
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
	return domain.Price(math.Round(float64((price+(price/1000))*1000)) / 1000)
}

func less(price domain.Price) domain.Price {
	return domain.Price(math.Round(float64((price-(price/1000))*1000)) / 1000)
}
