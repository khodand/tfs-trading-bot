package pkg

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Dsn             string `json:"dsn"`
	Telegram        string `json:"telegram"`
	AlgoSellPeriod  int `json:"algoPeriod"`
	KrakenWebsocket string `json:"krakenWebsocket"`
	KrakenREST      string `json:"krakenREST"`
	KrakenPublicKey string `json:"krakenPublicKey"`
	KrakenSecretKey string `json:"krakenSecretKey"`
}

func ReadConfig(filename string) (Config, error) {
	fmt.Println(os.Getwd())
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
