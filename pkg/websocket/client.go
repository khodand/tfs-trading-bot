package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type Client struct {
	connLock sync.RWMutex
	conn     *websocket.Conn
	url      string
	duration time.Duration
}

func NewWebSocketClient(url string, readDuration time.Duration) *Client {
	return &Client{
		url:      url,
		duration: readDuration,
	}
}

func (client *Client) Connect() {
	client.connLock.Lock()
	defer client.connLock.Unlock()

	if client.conn == nil {
		ws, _, err := websocket.DefaultDialer.Dial(client.url, nil)
		if err != nil {
			log.Fatal(err)
			return
		}
		client.conn = ws
	}
}

func (client *Client) readMessage() (p []byte, err error) {
	client.connLock.RLock()
	defer client.connLock.RUnlock()

	_, messages, err := client.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	safeMessages := append(make([]byte, 0, len(messages)), messages...)
	return safeMessages, nil
}

func (client *Client) Listen() <-chan []byte {
	client.Connect()

	out := make(chan []byte)
	go func() {
		defer close(out)

		ticker := time.NewTicker(client.duration)
		var message []byte
		for range ticker.C {
			var err error
			message, err = client.readMessage()
			log.Println(string(message))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			out <- message
		}
	}()

	return out
}

func (client *Client) debugListen() {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		client.Connect()
		for {
			byteMsg, err := client.readMessage()
			if err != nil {
				client.Close()
				break
			}
			fmt.Println(string(byteMsg))
		}
	}
}

func (client *Client) Close() {
	client.connLock.Lock()
	defer client.connLock.Unlock()
	if client.conn == nil {
		return
	}
	err := client.conn.Close()
	if err != nil {
		return
	}
}

func (client *Client) WriteJSON(json interface{}) error {
	client.Connect()

	client.connLock.RLock()
	defer client.connLock.RUnlock()
	return client.conn.WriteJSON(json)
}

type Messageasdfas struct {
	Event      string   `json:"event"`
	Feed       string   `json:"feed"`
	ProductIds []string `json:"product_ids"`
}

func main() {
	wc := NewWebSocketClient("wss://demo-futures.kraken.com/ws/v1?chart", time.Second)
	wc.Connect()
	err := wc.WriteJSON(Messageasdfas{
		Event: "subscribe",
		Feed:  "candles_trade_1m",
		//ProductIds: []string{"PI_XBTUSD"},
		ProductIds: []string{"IN_ETHUSD"},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	wc.debugListen()

	for true {
	}
}
