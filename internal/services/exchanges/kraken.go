package exchanges

import (
	"encoding/json"
	"time"

	"tfs-trading-bot/internal/domain"
	"tfs-trading-bot/internal/services/websocket"
)

const (
	websocketAddr = "wss://futures.kraken.com/ws/v1"
	restAddr = "https://futures.kraken.com/derivatives/api/v3"
)

type KrakenFuturesExchange struct {
	socket *websocket.WebSocketClient
	tickersOut chan domain.Ticker
}

func NewKrakenExchange() *KrakenFuturesExchange{
	return &KrakenFuturesExchange{
		socket:     websocket.NewWebSocketClient(websocketAddr, time.Second),
		tickersOut: nil,
	}
}

type Message struct {
	Event      string   `json:"event"`
	Feed       string   `json:"feed"`
	ProductIds []string `json:"product_ids"`
}

func (exc *KrakenFuturesExchange) Subscribe(symbol domain.TickerSymbol) {
	_ = exc.socket.WriteJSON(Message{
		Event:      "subscribe",
		Feed:       "ticker_lite",
		ProductIds: []string{string(symbol)},
	})
}

func (exc *KrakenFuturesExchange) GetTickersChan() <-chan domain.Ticker {
	return exc.tickersOut
}

func (exc *KrakenFuturesExchange) listenSocket() {
	go func() {
		for msg := range exc.socket.Listen() {
			var ticker domain.Ticker
			err := json.Unmarshal(msg, &ticker)
			if err != nil {
				continue
			}
			exc.tickersOut <- ticker
		}
	}()
}

func main() {

}
