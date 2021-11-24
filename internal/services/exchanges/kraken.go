package exchanges

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"tfs-trading-bot/internal/domain"
	pkghttp "tfs-trading-bot/pkg/http"
	"tfs-trading-bot/pkg/websocket"
)

type KrakenFuturesExchange struct {
	socket       *websocket.Client
	client       pkghttp.Client
	restAddr     string
	apiPublicKey string
	apiSecretKey string
	log          *logrus.Logger
	cancelFunc   context.CancelFunc
}

func NewKrakenExchange(websocketAddr, restAddr, apiPublicKey, apiSecretKey string, logger *logrus.Logger) *KrakenFuturesExchange {
	ctx, cancel := context.WithCancel(context.Background())
	e := KrakenFuturesExchange{
		socket: websocket.NewWebSocketClient(websocketAddr, time.Second*10, ctx),
		client: pkghttp.Client{
			Client: http.Client{Timeout: time.Second},
		},
		restAddr:     restAddr,
		apiPublicKey: apiPublicKey,
		apiSecretKey: apiSecretKey,
		log:          logger,
		cancelFunc:   cancel,
	}

	e.socket.Connect()
	return &e
}

type SendOrderResponse struct {
	Result     string `json:"result"`
	SendStatus struct {
		OrderId      string    `json:"order_id"`
		Status       string    `json:"status"`
		ReceivedTime time.Time `json:"receivedTime"`
	} `json:"sendStatus"`
	ServerTime time.Time `json:"serverTime"`
}

func (exc *KrakenFuturesExchange) SendOrder(order domain.Order) error {
	v := url.Values{}
	v.Add("orderType", order.OrderType)
	v.Add("symbol", string(order.Symbol))
	v.Add("side", order.Side)
	v.Add("size", strconv.Itoa(order.Size))
	v.Add("limitPrice", fmt.Sprintf("%f", order.LimitPrice))
	queryString := v.Encode()
	req, err := http.NewRequest(http.MethodPost, exc.restAddr+"/sendorder"+"?"+queryString, nil)
	if err != nil {
		panic(err)
	}

	signature := encodeAuth(queryString, "/api/v3/sendorder", exc.apiSecretKey)

	req.Header.Add("APIKey", exc.apiPublicKey)
	req.Header.Add("Authent", signature)

	resp := exc.client.PostRequest(req)
	res, _ := io.ReadAll(resp.Body)
	var jsonResponse SendOrderResponse
	err = json.Unmarshal(res, &jsonResponse)
	if err != nil {
		return err
	}
	exc.log.Debug("Kraken server response:", string(res))
	if jsonResponse.SendStatus.Status != "placed" {
		return errors.New("THE ORDER WAS NOT PLACED")
	}
	return nil
}

type SubscribeMessage struct {
	Event      string   `json:"event"`
	Feed       string   `json:"feed"`
	ProductIds []string `json:"product_ids"`
}

func (exc *KrakenFuturesExchange) Subscribe(symbol domain.TickerSymbol) {
	exc.log.Debug("Subscribes to ", symbol)
	err := exc.socket.WriteJSON(SubscribeMessage{
		Event:      "subscribe",
		Feed:       "ticker",
		ProductIds: []string{string(symbol)},
	})
	if err != nil {
		exc.log.Error(err)
	}
}

func (exc *KrakenFuturesExchange) GetTickersChan() <-chan domain.Ticker {
	out := make(chan domain.Ticker)
	go func() {
		defer close(out)
		exc.log.Info("Kraken waits for tickers")
		for msg := range exc.socket.Listen() {
			exc.log.Trace(string(msg))
			var ticker domain.Ticker
			err := json.Unmarshal(msg, &ticker)
			if err != nil {
				continue
			}
			if ticker.Bid == 0 {
				continue
			}
			out <- ticker
		}
	}()

	return out
}

func (exc *KrakenFuturesExchange) Stop() {
	exc.cancelFunc()
}

func encodeAuth(postData string, endpointPath string, apiSecretKey string) string {
	data := []byte(postData + endpointPath)
	sha := sha256.New()
	sha.Write(data)

	apiDecode, _ := base64.StdEncoding.DecodeString(apiSecretKey)

	h := hmac.New(sha512.New, apiDecode)
	h.Write(sha.Sum(nil))

	out := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return out
}
