package websocket

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

type Test struct {
	Event      string   `json:"event"`
	Feed       string   `json:"feed"`
	ProductIds []string `json:"product_ids"`
}

func TestConnect(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	url := "ws" + strings.TrimPrefix(s.URL, "http")
	ws := NewWebSocketClient(url, time.Second, context.Background())
	defer ws.Close()
	ws.Connect()

	err := ws.WriteJSON(Test{
		Event:      "subscribe",
		Feed:       "candles_trade_1m",
		ProductIds: []string{"PI_ETHUSD"},
	})
	assert.NoError(t, err)

	expected := "{\"event\":\"subscribe\",\"feed\":\"candles_trade_1m\",\"product_ids\":[\"PI_ETHUSD\"]}\n"
	assert.Equal(t, expected, string(<-ws.Listen()))
}
