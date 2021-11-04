package main

import "tfs-trading-bot/internal/domain"

type ExponentialMovingAverageAlgo struct {
}

type tickerMA struct {
	ma domain.Price
	last domain.Price
	period int
}

func (m *ExponentialMovingAverageAlgo) ProcessTickers(tickers <-chan domain.Ticker) <-chan domain.Order {
	tickersAverage := make(map[domain.TickerSymbol]tickerMA)
	out := make(chan domain.Order)

	go func() {
		defer close(out)
		for ticker := range tickers {
			var average tickerMA
			price := ticker.Bid
			if tmp, ok := tickersAverage[ticker.ProductId]; ok {
				average = tmp
				average.period++
				k := smoothingConstant(average.period)
				average.ma = (price * k) + (average.ma * (1 - k))
			} else {
				average = tickerMA{
					ma:  price,
					last: price,
					period: 1,
				}
			}
			tickersAverage[ticker.ProductId] = average

			if price > average.ma && average.last < average.ma { // cross high
				out <- buy()
			} else {
				if price < average.ma && average.last > average.ma { // cross low
					out <- sell()
				}
			}
		}
	}()

	return out
}

func smoothingConstant(period int) domain.Price {
	return domain.Price(2.0 / (1 + float64(period)))
}

func sell() domain.Order {
	return domain.Order{
		OrderType:     "",
		Symbol:        "",
		Side:          "",
		Size:          0,
		LimitPrice:    0,
		StopPrice:     0,
		TriggerSignal: "",
	}
}

func buy() domain.Order {
	return domain.Order{
		OrderType:     "",
		Symbol:        "",
		Side:          "",
		Size:          0,
		LimitPrice:    0,
		StopPrice:     0,
		TriggerSignal: "",
	}
}