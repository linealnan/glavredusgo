package telegrambot

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	blevesearch "github.com/linealnan/glavredusgo/internal/blevesearch"
	conf "github.com/linealnan/glavredusgo/internal/config"
)

type TelegramBot struct {
	conf          *conf.AppConfig
	bot           *tgbotapi.BotAPI
	searchService *blevesearch.BleaveSearch
}

func NewTelegramBot(c *conf.AppConfig, s *blevesearch.BleaveSearch) *TelegramBot {
	token := conf.InitWithDotEnv().TelegramBotApiToken
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &TelegramBot{conf: c, bot: bot, searchService: s}
}

func (tb *TelegramBot) SubsribeUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tb.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			groups := strings.Split(update.Message.Text, "\n")
			log.Println(groups)
			res := tb.searchService.TgSearch(update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, res)
			msg.ReplyToMessageID = update.Message.MessageID

			tb.bot.Send(msg)
		}
	}
}
