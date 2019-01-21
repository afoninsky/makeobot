package telegram

import (
	"log"
	"time"
	"errors"
	"strings"
	"strconv"

	tpl "text/template"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/afoninsky/makeomatic/common"
)

const helpMessage = "" +
	"[ChatOps bot](https://github.com/afoninsky/makeobot) welcomes you. Available commands are:\n\n" +
	"`/ping` - check my liveness\n" +
	"`/release image tag` - deploy new image using keel.sh" +
	""

const messageTemplate = "*{{ .Service }}: {{ .Name }}*\n{{ .Message }}"

type Service struct {
	logger *log.Logger
	ctx *common.AppContext
	bot *tgbotapi.BotAPI
	tplMessage *tpl.Template
}

func (c *Service) Init(ctx *common.AppContext) (common.Help, error) {
	c.logger = common.CreateLogger("telegram")
	c.ctx = ctx

	c.tplMessage = tpl.Must(tpl.New("message").Parse(messageTemplate))

	apiKey := ctx.Config.GetString("telegram.api")
	if apiKey == "" {
		return nil, errors.New("telegram.api is not set")
	}

	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return nil, err
	}
	c.bot = bot


	c.logger.Printf("authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}
	// wait for updates and clear backlog of old messages
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	go c.handleIncomingMessages(updates)

	help := common.Help{}
	return help, nil
}
func (c *Service) OnEvent(event common.Event) error {
	return nil
}

// do nothing as telegram service does not receive commands from the other services
func (c *Service) OnCommand(command common.Command) error {
	return nil
}

func (c *Service) handleIncomingMessages(updates tgbotapi.UpdatesChannel) {

	for update := range updates {
		// ignore any non-command messages
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		c.logger.Printf("[%s in %d said] %s", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text)

		// send incoming command to the router
		command := common.Command{
			Name: update.Message.Command(),
			Args: strings.Split(update.Message.CommandArguments(), " "),
			Sender: update.Message.From.UserName,
			Channel: strconv.FormatInt(update.Message.Chat.ID, 10),
		}
		err := c.ctx.Router.ExecuteCommand(command)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.ParseMode = "Markdown"
			msg.ReplyToMessageID = update.Message.MessageID	
			msg.Text = err.Error()
		if _, err := c.bot.Send(msg); err != nil {
			c.logger.Println(err)
		}
	}

	}
}