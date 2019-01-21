package keel

import (
	"github.com/afoninsky/makeomatic/common"
)

// Service ...
type Service struct {
}

// Init ...
func (c *Service) Init(ctx *common.AppContext) (map[string]string, error) {
	help := map[string]string{}
	return help, nil
}

// OnEvent ...
func (c *Service) OnEvent(event common.Event) error {
	return nil
}

// OnCommand ...
func (c *Service) OnCommand(command common.Command) error {
	return nil
}
