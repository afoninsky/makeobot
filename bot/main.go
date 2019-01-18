package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	tpl "text/template"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type Sevice struct {
	config           *viper.Viper
	router           *mux.Router
	bot              *tgbotapi.BotAPI
	tplNewDeployment *tpl.Template
}

func (s *Sevice) init() {

	// init http
	s.router.HandleFunc("/", s.httpHealthHandler).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/keel/deployment", s.httpKeelNewDeploymentHandler).Methods("POST", "OPTIONS")

	// init telegram
	log.Printf("Authorized on account %s", s.bot.Self.UserName)
	s.bot.Debug = true

	// add templates
	s.tplNewDeployment = tpl.Must(tpl.New("keel-update").Parse("*keel: {{ .Name }}*\n{{ .Message }}"))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := s.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	// wait for updates and clear backlog of old messages
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	go func() {
		for update := range updates {
			// ignore any non-Message Updates
			if update.Message == nil {
				continue
			}

			// TODO: ignore messages not from specified channel

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			parts := strings.Split(update.Message.Text, " ")

			// ignore not mine message
			if parts[0] != "keel" {
				continue
			}

			if len(parts) < 3 {
				msg := createErrorMessage(update.Message, errors.New("format"))
				s.bot.Send(msg)
				continue
			}

			_, command, payload := parts[0], parts[1], parts[2]

			// ignore as message not to me

			switch command {

			case "deploy":
				// keel deploy <image:tag>
				if err := s.updateKeelDeployment(payload); err != nil {
					msg := createErrorMessage(update.Message, err)
					s.bot.Send(msg)
					continue
				}
			default:
				msg := createErrorMessage(update.Message, errors.New("command does not supported"))
				s.bot.Send(msg)
			}
		}
	}()
}

func createErrorMessage(message *tgbotapi.Message, err error) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("invalid: %s", err.Error()))
	msg.ReplyToMessageID = message.MessageID
	return msg
}

func main() {

	// init configuration
	config := viper.New()
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config.SetDefault("http.listen", "localhost:8000")
	config.SetDefault("telegram.api", "")
	config.SetDefault("telegram.receiver", "498146361")
	config.SetDefault("keel.address", "http://keel.default.svc.cluster.local:9300")

	tgBot, err := tgbotapi.NewBotAPI(config.GetString("telegram.api"))
	if err != nil {
		log.Panic(err)
	}

	service := &Sevice{
		config: config,
		router: mux.NewRouter(),
		bot:    tgBot,
	}

	log.Println("111")
	service.init()
	log.Println(config.GetString("http.listen"))
	log.Fatal(http.ListenAndServe(config.GetString("http.listen"), service.router))
}
