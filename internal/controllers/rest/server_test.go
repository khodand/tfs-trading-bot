package rest

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"tfs-trading-bot/internal/domain"
	"tfs-trading-bot/internal/services"

	"github.com/stretchr/testify/assert"
)

type StubTradingService struct {
}

func (s StubTradingService) TradeTicker(symbol domain.TickerSymbol) {}

func (s StubTradingService) ProcessOrders() <-chan domain.Order {
	return nil
}

func (s StubTradingService) ChangeAlgo(algo services.TradingAlgorithm) {}

func TestTradeTicker(t *testing.T) {
	server := NewServer(StubTradingService{})
	s := httptest.NewServer(server.Router())

	test := "PI_test"
	res, body := testRequest(t, s, http.MethodPost, "/trade/"+test, nil)

	assert.Equal(t, "Starting to trade the "+test, body)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

type ChangeAlgoTest struct {
	algo       string
	period     string
	statusCode int
}

func TestChangeAlgo(t *testing.T) {
	server := NewServer(StubTradingService{})
	s := httptest.NewServer(server.Router())
	tests := []ChangeAlgoTest{
		{"EMA", "100", http.StatusOK},
		{"qqqqq", "100", http.StatusBadRequest},
	}

	for _, test := range tests {
		res, _ := testRequest(t, s, http.MethodPost, "/algo/"+test.algo+"/"+test.period, nil)
		assert.Equalf(t, test.statusCode, res.StatusCode, test.algo)
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	defer resp.Body.Close()
	return resp, string(respBody)
}
