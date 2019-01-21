package common

import (
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// SeverityLevel specify severity of event
type SeverityLevel int

const (
	// Info event means that no remedial action is required
	Info SeverityLevel = 1 + iota
	// Warning event means that investigate to decide whether any action is required
	Warning
	// Minor event means that action is required, but the situation is not serious at this time
	Minor
	// Major event means that action is required immediately
	Major
	// Critical event means that action is required immediately because the scope of the problem has increased, investigate critical alerts or events immediately
	Critical
	// Fatal event means that an error has occurred but it is too late to take any remedial action to address it
	Fatal
)

// Event describes emittable events
type Event struct {
	Service  string
	Name     string
	Message  string
	Link     string
	Severity SeverityLevel
}

// Command describe available commands to the services
type Command struct {
	Command string
	Args    []string
	Sender  string
}

// Router instance routes incoming events and commads between services
type ServiceRouter interface {
	// registers  service for routing
	RegisterService(name string, ctx *AppContext, service ServiceProvider) error
	// sends an event to all registered services
	EmitEvent(event Event) error
	// finds receiver and execute command returning result
	ExecuteCommand(receiver string, command Command) error
}

// ServiceProvider instance provides interface for executing commands
type ServiceProvider interface {
	Init(ctx *AppContext) error
	OnEvent(event Event) error
	OnCommand(command Command) error
	Close() error
}

type AppContext struct {
	Config *viper.Viper
	HTTP   *mux.Router
	Router ServiceRouter
}
