package main

import (
	"errors"
	"log"

	"github.com/afoninsky/makeomatic/common"
)

type router struct {
	// map service name to its instance
	services map[string]common.ServiceProvider
	// map command to service name
	commands map[string]string
	logger   *log.Logger
}

func (r *router) RegisterService(name string, ctx *common.AppContext, service common.ServiceProvider) error {
	help, err := service.Init(ctx)
	if err != nil {
		return err
	}
	r.services[name] = service
	for command, _ := range help {
		r.commands[command] = name
	}
	r.logger.Printf("service \"%s\" is enabled", name)
	return nil
}

func (r *router) EmitEvent(event common.Event) error {
	return nil
}

func (r *router) ExecuteCommand(command common.Command)  error {
	switch command.Name {
	// display help
	case "help":
		event := common.Event{}
		r.EmitEvent(event)
	// emit "pong" event
	case "ping":
		event := common.Event{}
		r.EmitEvent(event)
	default:
		return errors.New("I don't know that command")
	}
	return nil
}

func InitServiceRouter() common.ServiceRouter {
	return &router{
		logger:   common.CreateLogger("router"),
		services: make(map[string]common.ServiceProvider),
	}
}
