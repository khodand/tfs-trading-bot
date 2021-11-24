package websocket

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"runtime"
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
		var err = errors.New("")
		var ws *websocket.Conn
		for err != nil {
			ws, _, err = websocket.DefaultDialer.Dial(client.url, nil)
		}
		client.conn = ws
	}
}

func (client *Client) readMessage() (p []byte, err error) {
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
			if err != nil {
				client.Connect()
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
			log.Println(string(byteMsg))
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
	go func() {
		time.Sleep(time.Second * 10)
		err := wc.WriteJSON(Messageasdfas{
			Event: "subscribe",
			Feed:  "ticker",
			ProductIds: []string{"PI_ETHUSD"},
		})
		log.Println("Write")
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	go wc.debugListen()

	for true {
		runtime.Gosched()
	}
}
