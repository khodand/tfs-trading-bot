package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type WebSocketClient struct {
	connLock sync.RWMutex
	conn     *websocket.Conn
	url      string
	duration time.Duration
}

func NewWebSocketClient(url string,readDuration time.Duration) *WebSocketClient {
	return &WebSocketClient{
		url:      url,
		duration: readDuration,
	}
}

func (client *WebSocketClient) Connect() {
	client.connLock.Lock()
	defer client.connLock.Unlock()

	if client.conn == nil {
		ws, _, err := websocket.DefaultDialer.Dial(client.url, nil)
		if err != nil {
			return
		}
		client.conn = ws
	}
}

func (client *WebSocketClient) readMessage() (p []byte, err error) {
	client.connLock.RLock()
	defer client.connLock.RUnlock()

	_, messages, err := client.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	safeMessages := append(make([]byte, 0, len(messages)), messages...)
	return safeMessages, nil
}

func (client *WebSocketClient) Listen() <-chan []byte {
	out := make(chan []byte)
	go func() {
		defer close(out)

		ticker := time.NewTicker(client.duration)
		var message []byte
		for range ticker.C {
			var err error
			message, err = client.readMessage()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
		out <- message
	}()

	return out
}

func (client *WebSocketClient) debugListen() {
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

func (client *WebSocketClient) Close() {
	client.connLock.Lock()
	defer client.connLock.Unlock()
	err := client.conn.Close()
	if err != nil {
		return
	}
}

func (client *WebSocketClient) WriteJSON(json interface{}) error {
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
	wc := NewWebSocketClient("wss://futures.kraken.com/ws/v1", time.Second)
	wc.Connect()
	err := wc.WriteJSON(Messageasdfas{
		Event:      "subscribe",
		Feed:       "ticker_lite",
		ProductIds: []string{"PI_XBTUSD"},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	wc.debugListen()

	for true {}
}
