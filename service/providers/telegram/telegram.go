package telegram

import (
	"log"
	"time"
	"errors"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/afoninsky/makeomatic/common"
)

type Service struct {
	// bot *tgbotapi.BotAPI
	logger *log.Logger
	ctx *common.AppContext
	bot *tgbotapi.BotAPI
}

func (c *Service) Init(ctx *common.AppContext) error {
	c.logger = common.CreateLogger("telegram")
	c.ctx = ctx

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


	return nil
}
func (c *Service) OnEvent(event common.Event) error {
	return nil
}

func (c *Service) OnCommand(command common.Command) error {
	return nil
}

func (c *Service) Close() error {
	return nil
}



// func New(config *viper.Viper, router *mux.Router) (*context, error) {
// 	logger := common.CreateLogger("telegram")
// 	apiKey := config.GetString("telegram.api")
// 	if apiKey == "" {
// 		return nil, errors.New("telegram.api is not set")
// 	}

// 	bot, err := tgbotapi.NewBotAPI(apiKey)
// 	if err != nil {
// 		return nil, err
// 	}

// 	provider := &context{
// 		logger: logger,
// 		bot: bot,
// 	}

// 	logger.Printf("authorized on account %s", bot.Self.UserName)

// 	u := tgbotapi.NewUpdate(0)
// 	u.Timeout = 60
// 	updates, err := bot.GetUpdatesChan(u)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// wait for updates and clear backlog of old messages
// 	time.Sleep(time.Millisecond * 500)
// 	updates.Clear()



// 	return provider, nil
// }
