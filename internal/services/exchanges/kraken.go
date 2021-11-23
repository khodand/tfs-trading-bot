package exchanges

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
	"tfs-trading-bot/pkg"
	"time"

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
}

func NewKrakenExchange(websocketAddr, restAddr, apiPublicKey, apiSecretKey string) *KrakenFuturesExchange {
	e := KrakenFuturesExchange{
		socket: websocket.NewWebSocketClient(websocketAddr, time.Second),
		client: pkghttp.Client{
			Client: http.Client{Timeout: time.Second * 5},
		},
		restAddr:     restAddr,
		apiPublicKey: apiPublicKey,
		apiSecretKey: apiSecretKey,
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
	log.Println("SERVER RESPONSE", string(res))
	log.Println("SERVER RESPONSE", jsonResponse)
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
	log.Println("Subscribes to ", symbol)
	err := exc.socket.WriteJSON(SubscribeMessage{
		Event:      "subscribe",
		Feed:       "ticker",
		ProductIds: []string{string(symbol)},
	})
	fmt.Println(err)
}

func (exc *KrakenFuturesExchange) GetTickersChan() <-chan domain.Ticker {
	out := make(chan domain.Ticker)
	go func() {
		defer close(out)
		for msg := range exc.socket.Listen() {
			log.Println(string(msg))
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

func (exc *KrakenFuturesExchange) GetAccounts() {
	req, err := http.NewRequest(http.MethodGet, exc.restAddr+"/accounts", nil)
	if err != nil {
		panic(err)
	}

	signature := encodeAuth("", "/api/v3/accounts", exc.apiSecretKey)

	req.Header.Add("APIKey", exc.apiPublicKey)
	req.Header.Add("Authent", signature)

	_, err = httputil.DumpRequestOut(req, true)
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(b))

	resp, err := exc.client.Do(req)
	if err != nil {
		panic(err)
	}

	_, err = httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	log.Println("SERVER RESPONSE", resp.Body)
}

func encodeAuth(postData string, endpointPath string, apiSecretKey string) string {
	data := []byte(postData + endpointPath)
	sha := sha256.New()
	sha.Write(data)

	apiDecode, _ := base64.StdEncoding.DecodeString(apiSecretKey)

	h := hmac.New(sha512.New, apiDecode)
	h.Write(sha.Sum(nil))

	out := base64.StdEncoding.EncodeToString(h.Sum(nil))
	//fmt.Println(out)
	return out
}

func main() {
	//encodeAuth("orderType=lmt&symbol=pi_xbtusd&side=buy&size=10000&limitPrice=9400", "", "/api/v3/sendorder")
	//encodeAuth("", "", "/api/v3/cancelallorders")
	config := pkg.ReadConfig("config.json")
	e := NewKrakenExchange(config.KrakenWebsocket, config.KrakenREST, config.KrakenPublicKey, config.KrakenSecretKey)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range e.socket.Listen() {
			var ticker domain.Ticker
			err := json.Unmarshal(msg, &ticker)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println(ticker)
		}
	}()

	e.Subscribe("pi_ethusd")
	//err := e.SendOrder(domain.Order{
	//	OrderType:  "ioc",
	//	Symbol:     "pi_ethusd",
	//	Side:       "buy",
	//	Size:       1,
	//	LimitPrice: 4400,
	//})
	//log.Println("SEDN" ,err)
	//e.GetAccounts()

	wg.Wait()
	fmt.Println("EnD")
}
