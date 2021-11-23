package pkg

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Dsn             string `json:"dsn"`
	Telegram        string `json:"telegram"`
	KrakenWebsocket string `json:"krakenWebsocket"`
	KrakenREST      string `json:"krakenREST"`
	KrakenPublicKey string `json:"krakenPublicKey"`
	KrakenSecretKey string `json:"krakenSecretKey"`
}

func ReadConfig(filename string) Config {
	file, err := os.Open(filename)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
