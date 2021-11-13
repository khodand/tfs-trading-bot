package domain

import (
	"fmt"
	"strconv"
)

type Price float64

type TickerSymbol string

// Order https://support.kraken.com/hc/en-us/articles/360022839691-Send-order
type Order struct {
	OrderType     string
	Symbol        TickerSymbol
	Side          string
	Size          int
	LimitPrice    Price
	StopPrice     Price  // not required
	TriggerSignal string // not required
}

func (o Order) String() string {
	return "Order: " + o.Side + " " + string(o.Symbol) + " " + strconv.Itoa(o.Size) + " " + o.OrderType + " " +
		fmt.Sprintf("%f", o.LimitPrice)
}

// Ticker https://support.kraken.com/hc/en-us/articles/360022839751-Ticker-Lite
type Ticker struct {
	Feed         string       `json:"feed"`
	ProductId    TickerSymbol `json:"product_id"`
	Bid          Price        `json:"bid"`
	Ask          Price        `json:"ask"`
	Change       float64      `json:"change"`
	Premium      float64      `json:"premium"`
	Volume       float64      `json:"volume"`
	Tag          string       `json:"tag"`
	Pair         string       `json:"pair"`
	Dtm          int          `json:"dtm"`
	MaturityTime int          `json:"maturityTime"`
}

func main() {

}
