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
	pkgpostgres "tfs-trading-bot/pkg/postgres"
)

const dsn = "postgres://postgres:1234@localhost:5432/postgres" +
	"?sslmode=disable"

const telegramToken = "2023059929:AAGEN8f835UIb8k4puoBn5n32nACjaRDSxE"

func main() {
	pool, err := pkgpostgres.NewPool(dsn)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pool.Close()

	repo := repository.NewRepository(pool)
	algo := algorithms.NewEMAAlgo(10)
	// TODO: api tokens as parameters
	ex := exchanges.NewKrakenExchange()

	trader := services.NewTrader(ex, algo, repo)

	bot := telegram.NewTelegramBot(telegramToken, trader)
	bot.Start()

	handler := rest.NewServer(trader)
	log.Fatal(http.ListenAndServe(":5000", handler.Router()))
}
