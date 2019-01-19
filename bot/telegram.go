// handles incoming messages via telegram
package main

import (
	"bytes"
	"encoding/json"
	"errors"
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
	for update := range updates {
		// ignore any non-message updates
		if update.Message == nil {
			continue
		}

		// TODO: ignore messages not from specified channel or not from unauthorized user
		// TODO: valid command workflow, ex.:
		// https://go-telegram-bot-api.github.io/examples/commands/
		// https://go-telegram-bot-api.github.io/examples/keyboard/
		log.Printf("[%s in %d] %s", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text)
		parts := strings.Split(update.Message.Text, " ")

		switch parts[0] {
		case "keel":
			// ensure command format is valud
			if len(parts) < 3 {
				msg := createErrorMessage(update.Message, errors.New("invalid command format"))
				bot.Send(msg)
				continue
			}
			_, command, payload := parts[0], parts[1], parts[2]

			switch command {

			case "deploy":
				// keel deploy <image:tag>
				if err := updateKeelDeployment(config.GetString("keel.address"), payload); err != nil {
					msg := createErrorMessage(update.Message, err)
					bot.Send(msg)
					continue
				}
			default:
				msg := createErrorMessage(update.Message, errors.New("don't support this command"))
				bot.Send(msg)
			}
		default:
			// do nothing
		}

	}
}

func createErrorMessage(message *tgbotapi.Message, err error) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("invalid: %s", err.Error()))
	msg.ReplyToMessageID = message.MessageID
	return msg
}

func updateKeelDeployment(keelHost string, image string) error {
	parts := strings.Split(image, ":")
	if len(parts) != 2 {
		return errors.New("expect {name:tag} as image")
	}
	name, tag := parts[0], parts[1]

	url := fmt.Sprintf("%s/v1/webhooks/native", keelHost)

	values := map[string]string{"name": name, "tag": tag}
	jsonValue, _ := json.Marshal(values)
	log.Println(values)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	return err
}
