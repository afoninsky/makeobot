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
	Service  string // event service source
	Name     string // event name / type
	Message  string // event content
	Link     string // link to external system
	Severity SeverityLevel
	RootID   string // id of the root event or command which create this event
}

// Command describe available commands to the services
type Command struct {
	ID      string   // optional command id
	Name    string   // command name: "ping", "deploy" etc
	Args    []string // command additional arguments
	Sender  string   // username if command come from telegram
	Channel string   // channel id if command comes from telegram
}

// Router instance routes incoming events and commads between services
type ServiceRouter interface {
	// registers  service for routing
	RegisterService(name string, ctx *AppContext, service ServiceProvider) error
	// sends an event to all registered services
	EmitEvent(event Event) error
	// finds receiver and execute command returning result
	ExecuteCommand(command Command) error
}

// ServiceProvider instance provides interface for executing commands
type ServiceProvider interface {
	Init(ctx *AppContext) (map[string]string, error)
	OnEvent(event Event) error
	OnCommand(command Command) error
}

type AppContext struct {
	Config *viper.Viper
	HTTP   *mux.Router
	Router ServiceRouter
}
