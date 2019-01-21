package main
import (
	"log"
	"github.com/afoninsky/makeomatic/common"
)

type router struct {
	services map[string]common.ServiceProvider
	logger *log.Logger
}

func (r *router) RegisterService(name string, ctx *common.AppContext, service common.ServiceProvider) error {
	if err := service.Init(ctx); err != nil {
		return err
	}
	defer service.Close()
	r.services[name] = service
	r.logger.Printf("service \"%s\" is enabled", name)
	return nil
}

func (r *router) EmitEvent(event common.Event) error {
	return nil
}

func (r *router) ExecuteCommand(receiver string, command common.Command) error {
	return nil
}

func InitServiceRouter() common.ServiceRouter {
	return &router{
		logger: common.CreateLogger("router"),
		services:  make(map[string]common.ServiceProvider),
	}
}