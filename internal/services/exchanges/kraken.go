package exchanges

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
	"time"

	"tfs-trading-bot/internal/domain"
	pkghttp "tfs-trading-bot/pkg/http"
	"tfs-trading-bot/pkg/websocket"
)

const (
	websocketAddr = "wss://demo-futures.kraken.com/ws/v1"
	restAddr      = "https://demo-futures.kraken.com/derivatives/api/v3"
	apiPublicKey  = "52Hw5QQCv6O5X+zQFDeVH+9DpHV339mT+NW/EtjN0+krrwcWFFBiWke7"
	apiSecretKey  = "i465FHjhb0oKMaV4l+FIzZXfb9N3PYose1CP9qrBY8vdlVOC64Q9/M76ANyXml915TrBFzepJZ7Zdc/NelIOwDa7"
)

type KrakenFuturesExchange struct {
	socket *websocket.Client
	client pkghttp.Client
}

func NewKrakenExchange() *KrakenFuturesExchange {
	e := KrakenFuturesExchange{
		socket: websocket.NewWebSocketClient(websocketAddr, time.Second),
		client: pkghttp.Client{
			Client: http.Client{Timeout: time.Second * 5},
		},
	}
	e.socket.Connect()
	return &e
}

// https://futures.kraken.com/derivatives/api/v3/sendorder?orderType=lmt&symbol=pi_xbtusd&side=buy&size=10000&limitPrice=9400&reduceOnly=true

func (exc *KrakenFuturesExchange) SendOrder(order domain.Order) {
	v := url.Values{}
	v.Add("orderType", order.OrderType)
	v.Add("symbol", string(order.Symbol))
	v.Add("side", order.Side)
	v.Add("size", strconv.Itoa(order.Size))
	v.Add("limitPrice", fmt.Sprintf("%f", order.LimitPrice))
	queryString := v.Encode()
	req, err := http.NewRequest(http.MethodPost, restAddr+"/sendorder"+"?"+queryString, nil)
	if err != nil {
		panic(err)
	}

	signature := encodeAuth(queryString, "/api/v3/sendorder")

	req.Header.Add("APIKey", apiPublicKey)
	req.Header.Add("Authent", signature)

	resp := exc.client.PostRequest(req)
	log.Println("SERVER RESPONCE", resp.Body)
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
		Feed:       "ticker_lite",
		ProductIds: []string{string(symbol)},
	})
	fmt.Println(err)
}

func (exc *KrakenFuturesExchange) GetTickersChan() <-chan domain.Ticker {
	out := make(chan domain.Ticker)
	go func() {
		defer close(out)
		for msg := range exc.socket.Listen() {
			var ticker domain.Ticker
			err := json.Unmarshal(msg, &ticker)
			if err != nil {
				continue
			}
			if ticker.Bid == 0 {
				continue
			}
			log.Println(ticker)
			out <- ticker
		}
	}()

	return out
}

func (exc *KrakenFuturesExchange) GetAccounts() {
	req, err := http.NewRequest(http.MethodGet, restAddr+"/accounts", nil)
	if err != nil {
		panic(err)
	}

	signature := encodeAuth("", "/api/v3/accounts")

	req.Header.Add("APIKey", apiPublicKey)
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
	log.Println("SERVER RESPONCE", resp.Body)
}

func encodeAuth(postData string, endpointPath string) string {
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
	e := NewKrakenExchange()

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
	e.SendOrder(domain.Order{
		OrderType:  "ioc",
		Symbol:     "pi_ethusd",
		Side:       "buy",
		Size:       1,
		LimitPrice: 4342,
	})
	e.GetAccounts()

	wg.Wait()
	fmt.Println("EnD")
}
