// handles incoming messages via telegram
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
)

func initTgBot(config *viper.Viper) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(config.GetString("telegram.api"))
	if err != nil {
		return nil, err
	}

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

	go handleIncominMessages(updates, bot, config)

	return bot, nil
}

func handleIncominMessages(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI, config *viper.Viper) {
	keelAddress := config.GetString("keel.address")
	for update := range updates {
		// ignore any non-command messages
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}
		arguments := strings.Split(update.Message.CommandArguments(), " ")

		log.Printf("[%s in %d said] %s", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = "Markdown"
		msg.ReplyToMessageID = update.Message.MessageID

		switch update.Message.Command() {
		case "help":
			msg.Text = helpMessage
		case "ping":
			msg.Text = "pong"
		case "release":
			if len(arguments) != 2 {
				msg.Text = "valid command is: /release image tag"
			} else {
				if err := updateKeelDeployment(keelAddress, arguments[0], arguments[1]); err != nil {
					msg.Text = err.Error()
				}
			}
		default:
			msg.Text = "I don't know that command"
		}

		if msg.Text != "" {
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
			}
		}
	}
}

func updateKeelDeployment(keelHost string, name string, tag string) error {
	url := fmt.Sprintf("%s/v1/webhooks/native", keelHost)
	values := map[string]string{"name": name, "tag": tag}
	jsonValue, _ := json.Marshal(values)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	return err
}
