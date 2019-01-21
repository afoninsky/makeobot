package telegram

import (
	"log"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/afoninsky/makeomatic/common"
)

type context struct {
	bot *tgbotapi.BotAPI
}

func (c *context) Notify(event common.Event) error {
	return nil
}

func (c *context) Close() {
	//
}

func New(config *viper.Viper, router *mux.Router) (*context, error) {
	provider := &context{}

	bot, err := tgbotapi.NewBotAPI(config.GetString("telegram.api"))
	if err != nil {
		return nil, err
	}
	provider.bot = bot

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}
	// wait for updates and clear backlog of old messages
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	

	return provider, nil
}
