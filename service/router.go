package main

import (
	"errors"
	"log"
	"strings"

	"github.com/afoninsky/makeomatic/common"
)

type router struct {
	// map service name to its instance
	services map[string]common.ServiceProvider
	// map command to service name
	commands map[string]common.CommandInfo
	// information about available commands
	logger   *log.Logger
}

func (r *router) RegisterService(name string, ctx *common.AppContext, service common.ServiceProvider) error {
	// init target service
	if err := service.Init(ctx); err != nil {
		return err
	}
	r.services[name] = service

	// load command information
	for _, cmd := range service.ListCommands() {
		cmd.Service = name
		r.commands[cmd.Name] = cmd
	}
	r.logger.Printf("service \"%s\" is enabled", name)
	return nil
}

func (r *router) ListCommands() map[string]common.CommandInfo {
	return r.commands
}

func (r *router) EmitEvent(event common.Event) error {
	for name, service := range r.services {
		if err := service.OnEvent(event); err != nil {
			r.logger.Printf("\"%s returned error: %s\"", name, err.Error())
		}
	}
	return nil
}

func (r *router) ExecuteCommandString(message, messageID, sender string) error {
	parts := strings.Split(message, " ")
	for i := len(parts); i > 0; i-- {
		testCmd := strings.Join(parts[:i], " ")
		info, found := r.commands[testCmd]
		if !found {
			continue
		}
		command := common.Command{
			ID: messageID,
			Name: testCmd,
			Args: parts[i:],
			Sender: sender,
		}
		serviceName := info.Service
		return r.services[serviceName].DoCommand(command)
	}
	return errors.New("I don't know that command")
}

func InitServiceRouter() common.ServiceRouter {
	return &router{
		logger:   common.CreateLogger("router"),
		services: make(map[string]common.ServiceProvider),
		commands: make(map[string]common.CommandInfo),
	}
}
