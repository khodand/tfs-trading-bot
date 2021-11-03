package domain

type Price float64

type TickerSymbol string

// SendOrder https://support.kraken.com/hc/en-us/articles/360022839691-Send-order
type SendOrder struct {
	orderType string
	symbol TickerSymbol
	side string
	size int
	limitPrice Price
	stopPrice Price // not required
	triggerSignal string // not required
}

// Ticker https://support.kraken.com/hc/en-us/articles/360022839751-Ticker-Lite
type Ticker struct {
	Feed         string  `json:"feed"`
	ProductId    string  `json:"product_id"`
	Bid          Price  `json:"bid"`
	Ask          Price 	`json:"ask"`
	Change       float64 `json:"change"`
	Premium      float64 `json:"premium"`
	Volume       float64     `json:"volume"`
	Tag          string  `json:"tag"`
	Pair         string  `json:"pair"`
	Dtm          int     `json:"dtm"`
	MaturityTime int     `json:"maturityTime"`
}

func main() {

}
