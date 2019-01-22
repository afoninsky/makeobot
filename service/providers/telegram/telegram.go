package telegram

import (
	"log"
	"time"
	"errors"
	"strconv"
	"bytes"
	"fmt"

	tpl "text/template"

	"github.com/patrickmn/go-cache"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/afoninsky/makeomatic/common"
)

const helpMessage = "" +
	"[ChatOps bot](https://github.com/afoninsky/makeobot) welcomes you. Available commands are:\n\n"

// Service ..
type Service struct {
	logger *log.Logger
	ctx *common.AppContext
	bot *tgbotapi.BotAPI
	tplMessage *tpl.Template
	cache *cache.Cache
}

// Init ..
func (c *Service) Init(ctx *common.AppContext) (error) {
	c.logger = common.CreateLogger("telegram")
	c.ctx = ctx
	c.cache = cache.New(1*time.Minute, 10*time.Minute)
	
	c.tplMessage = tpl.Must(tpl.New("message").Parse(ctx.Config.GetString("telegram.template")))

	apiKey := ctx.Config.GetString("telegram.api")
	if apiKey == "" {
		return errors.New("telegram.api is not set")
	}

	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return err
	}
	c.bot = bot


	c.logger.Printf("authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}
	// wait for updates and clear backlog of old messages
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	go c.handleIncomingMessages(updates)

	return nil
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

// ListCommands ...
func (c *Service) ListCommands() []common.CommandInfo {
	return []common.CommandInfo{
		common.CommandInfo{
			Name: "help",
			Description: "display this help",
		},
		common.CommandInfo{
			Name: "ping",
			Description: "ensure chatbot is available",
			// Example: "qwe asd",
		},
	}
}

// DoCommand ...
func (c *Service) DoCommand(cmd common.Command) error {
	switch cmd.Name {
	case "help":
		message := helpMessage
		for _, info := range c.ctx.Router.ListCommands() {
			message += fmt.Sprintf("/%s - %s\n", defaults(info.Example, info.Name), info.Description)
		}
		event := common.Event{
			Message: message,
			RootID: cmd.ID,
		}
		return c.ctx.Router.EmitEvent(event)
	case "ping":
		event := common.Event{
			Message: "I'm alive",
			RootID: cmd.ID,
		}
		return c.ctx.Router.EmitEvent(event)
	}

	return errors.New("I don't know that command")
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
		c.cache.Set(messageID, chatID, cache.DefaultExpiration,)

		// send incoming command to the router
		message := update.Message.Command() + " " + update.Message.CommandArguments()
		sender := update.Message.From.UserName
		if err := c.ctx.Router.ExecuteCommandString(message, messageID, sender); err != nil {
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

func defaults(value, def string) string {
	if value == "" {
		return def
	}
	return value
}