package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"tfs-trading-bot/internal/domain"
	"tfs-trading-bot/internal/services"
)

type Bot struct {
	service   services.TradingService
	bot       *tgbotapi.BotAPI
	isStarted bool
}

func NewTelegramBot(token string, trader services.TradingService) *Bot {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	tbot := Bot{
		service:   trader,
		bot:       botAPI,
		isStarted: false,
	}

	return &tbot
}

func (t *Bot) sendOrders(chatID int64, orders <-chan domain.Order) {
	go func() {
		for order := range orders {
			msg := tgbotapi.NewMessage(chatID, order.String())
			_, _ = t.bot.Send(msg)
		}
	}()
}

func (t *Bot) Start() {
	t.bot.Debug = true

	log.Printf("Authorized on account %s", t.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := t.bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if !t.isStarted && update.Message.Text == "/start" {
			t.isStarted = true
			t.sendOrders(update.Message.Chat.ID, t.service.ProcessOrders())
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			_, _ = t.bot.Send(msg)
		}

	}
}

//func main() {
//	bot, err := tgbotapi.NewBotAPI(pkg.ReadConfig("config.json").Telegram)
//	if err != nil {
//		log.Panic(err)
//	}
//
//	bot.Debug = true
//
//	log.Printf("Authorized on account %s", bot.Self.UserName)
//
//	u := tgbotapi.NewUpdate(0)
//	u.Timeout = 60
//
//	updates, err := bot.GetUpdatesChan(u)
//
//	for update := range updates {
//		if update.Message == nil { // ignore any non-Message Updates
//			continue
//		}
//
//		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
//
//		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
//		msg.ReplyToMessageID = update.Message.MessageID
//
//		bot.Send(msg)
//	}
//}
