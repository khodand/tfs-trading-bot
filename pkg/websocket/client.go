package websocket

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	connLock sync.RWMutex
	conn     *websocket.Conn
	url      string
	duration time.Duration
	ctx      context.Context
}

func NewWebSocketClient(url string, readDuration time.Duration, ctx context.Context) *Client {
	return &Client{
		url:      url,
		duration: readDuration,
		ctx:      ctx,
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
		for {
			select {
			case <-client.ctx.Done():
				return
			case <-ticker.C:
				var err error
				message, err = client.readMessage()
				if err != nil {
					client.Connect()
				}
				out <- message
			}
		}
	}()

	return out
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
