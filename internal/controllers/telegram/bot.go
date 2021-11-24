package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"

	"tfs-trading-bot/internal/domain"
	"tfs-trading-bot/internal/services"
)

type Bot struct {
	service   services.TradingService
	bot       *tgbotapi.BotAPI
	isStarted bool
	log       *logrus.Logger
}

func NewTelegramBot(token string, trader services.TradingService, logger *logrus.Logger) *Bot {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot := Bot{
		service:   trader,
		bot:       botAPI,
		isStarted: false,
		log:       logger,
	}

	return &bot
}

func (t *Bot) sendOrders(chatID int64, orders <-chan domain.Order) {
	go func() {
		t.log.Info("Telegram bot waits for orders")
		for order := range orders {
			msg := tgbotapi.NewMessage(chatID, order.String())
			_, _ = t.bot.Send(msg)
		}
	}()
}

func (t *Bot) Start() {
	t.bot.Debug = t.log.IsLevelEnabled(logrus.TraceLevel)

	t.log.Infof("Authorized on account %s", t.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := t.bot.GetUpdatesChan(u)
	if err != nil {
		t.log.Panic(err)
	}

	go func() {
		for update := range updates {
			if update.Message == nil { // ignore any non-Message Updates
				continue
			}

			t.log.Debugf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if !t.isStarted && update.Message.Text == "/start" {
				t.isStarted = true
				t.sendOrders(update.Message.Chat.ID, t.service.ProcessOrders())
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				_, _ = t.bot.Send(msg)
			}

		}
	}()
}
