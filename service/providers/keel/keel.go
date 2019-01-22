package keel

import (
	"log"
	"errors"
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
	
	"github.com/afoninsky/makeomatic/common"
)

// Service ..
type Service struct {
	logger *log.Logger
	ctx *common.AppContext
}

// Init ..
func (c *Service) Init(ctx *common.AppContext) (error) {
	c.logger = common.CreateLogger("keel")
	c.ctx = ctx
	return nil
}

// OnEvent ...
func (c *Service) OnEvent(event common.Event) error {
	//
	return nil
}

// ListCommands ...
func (c *Service) ListCommands() []common.CommandInfo {
	return []common.CommandInfo{
		common.CommandInfo{
			Name: "keel update",
			Description: "update services controlled by keel.sh",
			Example: "keel update {image} {tag}",
		},
		common.CommandInfo{
			Name: "keel approvals",
			Description: "{not implemented yet}",
		},
		common.CommandInfo{
			Name: "keel approve",
			Description: "{not implemented yet}",
			Example: "keel approve {id}",
		},
		common.CommandInfo{
			Name: "keel reject",
			Description: "{not implemented yet}",
			Example: "keel reject {id}",
		},
	}
}

// DoCommand ...
func (c *Service) DoCommand(cmd common.Command) error {
	switch cmd.Name {
	case "keel update":
		if len(cmd.Args) != 2 {
			return errors.New("expect image and tag")
		}
		return c.updateKeelDeployment(cmd.Args[0], cmd.Args[1])
	}

	return errors.New("I don't know that command")
}

func (c *Service) updateKeelDeployment(name, tag string) error {
	url := fmt.Sprintf("%s/v1/webhooks/native", c.ctx.Config.GetString("keel.host"))
	values := map[string]string{"name": name, "tag": tag}
	jsonValue, _ := json.Marshal(values)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	return err
}
