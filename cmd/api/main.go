package main

import (
	"net/http"

	"github.com/sirupsen/logrus"

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
	var log = logrus.New()
	log.SetLevel(logrus.TraceLevel)

	config, err := pkg.ReadConfig("config.json")
	if err != nil {
		log.Fatal("Impossible to read the config:", err)
	}

	pool, err := pkgpostgres.NewPool(config.Dsn, log)
	if err != nil {
		log.Fatal("Impossible to connect to the database:", err)
		return
	}
	defer pool.Close()

	repo := repository.NewRepository(pool)
	algo := algorithms.NewEMAAlgo(config.AlgoSellPeriod, log)
	ex := exchanges.NewKrakenExchange(config.KrakenWebsocket, config.KrakenREST, config.KrakenPublicKey,
		config.KrakenSecretKey, log)

	trader := services.NewTrader(ex, algo, repo, log)

	bot := telegram.NewTelegramBot(config.Telegram, trader, log)
	bot.Start()

	handler := rest.NewServer(trader, log)
	log.Info("Start listen")
	log.Fatal(http.ListenAndServe(":5000", handler.Router()))
}
