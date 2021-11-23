package main

import (
	"fmt"
	"log"
	"net/http"
	"tfs-trading-bot/internal/controllers/rest"
	"tfs-trading-bot/internal/controllers/telegram"
	"tfs-trading-bot/internal/repository"
	"tfs-trading-bot/internal/services"
	"tfs-trading-bot/internal/services/algorithms"
	"tfs-trading-bot/internal/services/exchanges"
	"tfs-trading-bot/pkg"
	pkgpostgres "tfs-trading-bot/pkg/postgres"
)

func main() {
	config := pkg.ReadConfig("config.json")

	pool, err := pkgpostgres.NewPool(config.Dsn)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pool.Close()

	repo := repository.NewRepository(pool)
	algo := algorithms.NewEMAAlgo(10)
	ex := exchanges.NewKrakenExchange(config.KrakenWebsocket, config.KrakenREST, config.KrakenPublicKey,
		config.KrakenSecretKey)

	trader := services.NewTrader(ex, algo, repo)

	bot := telegram.NewTelegramBot(config.Telegram, trader)
	bot.Start()

	handler := rest.NewServer(trader)
	log.Fatal(http.ListenAndServe(":5000", handler.Router()))
}
