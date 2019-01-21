package telegram

import (
	"log"
	"time"
	"errors"
	"strings"
	"strconv"
	"bytes"

	tpl "text/template"

	"github.com/patrickmn/go-cache"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/afoninsky/makeomatic/common"
)


type Service struct {
	logger *log.Logger
	ctx *common.AppContext
	bot *tgbotapi.BotAPI
	tplMessage *tpl.Template
	cache *cache.Cache
}

func (c *Service) Init(ctx *common.AppContext) (common.Help, error) {
	c.logger = common.CreateLogger("telegram")
	c.ctx = ctx
	c.cache = cache.New(1*time.Minute, 10*time.Minute)
	
	help := common.Help{}

	c.tplMessage = tpl.Must(tpl.New("message").Parse(ctx.Config.GetString("telegram.template")))

	apiKey := ctx.Config.GetString("telegram.api")
	if apiKey == "" {
		return help, errors.New("telegram.api is not set")
	}

	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return help, err
	}
	c.bot = bot


	c.logger.Printf("authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return help, err
	}
	// wait for updates and clear backlog of old messages
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	go c.handleIncomingMessages(updates)

	return help, nil
}

// OnEvent displays rendered event in the chat
func (c *Service) OnEvent(event common.Event) error {
	var buffer bytes.Buffer
	if err := c.tplMessage.Execute(&buffer, event); err != nil {
		return err
	}
	chatID := c.guessChatID(event.RootID)
	msg := tgbotapi.NewMessageToChannel(chatID, buffer.String())
	msg.ParseMode = "Markdown"
	if _, err := c.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

// OnCommand does nothing as telegram service does not receive commands from the other services
func (c *Service) OnCommand(command common.Command) error {
	return nil
}

func (c *Service) guessChatID(rootID string) string {
	chatID, found := c.cache.Get(rootID)
	if found {
		return chatID.(string)
	}
	return c.ctx.Config.GetString("telegram.defaultChannel")
}

func (c *Service) handleIncomingMessages(updates tgbotapi.UpdatesChannel) {

	for update := range updates {
		// ignore any non-command messages
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		c.logger.Printf("[%s in %d said] %s", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text)

		// store command id so response can be forwarded to this user / chat
		messageID := strconv.Itoa(update.Message.MessageID)
		chatID := strconv.FormatInt(update.Message.Chat.ID, 10)

		c.cache.Set(messageID, chatID, cache.DefaultExpiration)

		// send incoming command to the router
		command := common.Command{
			ID: messageID,
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