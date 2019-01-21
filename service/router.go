package main

import (
	"errors"
	"log"
	"fmt"

	"github.com/afoninsky/makeomatic/common"
)

const helpMessage = "" +
	"[ChatOps bot](https://github.com/afoninsky/makeobot) welcomes you. Available commands are:\n\n"

type router struct {
	// map service name to its instance
	services map[string]common.ServiceProvider
	// map command to service name
	commands map[string]string
	// information about available commands
	help map[string]string
	logger   *log.Logger
}

func (r *router) RegisterService(name string, ctx *common.AppContext, service common.ServiceProvider) error {
	help, err := service.Init(ctx)
	if err != nil {
		return err
	}
	r.services[name] = service
	for command, description := range help {
		r.commands[command] = name
		r.help[command] = description
	}
	r.logger.Printf("service \"%s\" is enabled", name)
	return nil
}

func (r *router) EmitEvent(event common.Event) error {
	for name, service := range r.services {
		if err := service.OnEvent(event); err != nil {
			r.logger.Printf("\"%s returned error: %s\"", name, err.Error())
		}
	}
	return nil
}

func (r *router) ExecuteCommand(command common.Command)  error {
	switch command.Name {
	// display help
	case "help":
		message := helpMessage
		for name, desc := range r.help {
			message += fmt.Sprintf("/%s - %s", name, desc)
		}
		event := common.Event{
			Message: message,
			RootID: command.ID,
		}
		if err := r.EmitEvent(event); err != nil {
			return err
		}
	// emit "pong" event
	case "ping":
		event := common.Event{
			Message: "pong",
			RootID: command.ID,
		}
		if err := r.EmitEvent(event); err != nil {
			return err
		}
	default:
		return errors.New("I don't know that command")
	}
	return nil
}

func InitServiceRouter() common.ServiceRouter {
	help := map[string]string{
		"ping": "check liveness",
	}
	return &router{
		logger:   common.CreateLogger("router"),
		services: make(map[string]common.ServiceProvider),
		commands: make(map[string]string),
		help: help,
	}
}
