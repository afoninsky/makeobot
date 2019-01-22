package common

import (
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// SeverityLevel specify severity of event
type SeverityLevel string

const (
	// Info event means that no remedial action is required
	Info SeverityLevel = "info"
	// Warning event means that investigate to decide whether any action is required
	Warning = "warn"
	// Minor event means that action is required, but the situation is not serious at this time
	Minor = "minor"
	// Major event means that action is required immediately
	Major = "major"
	// Critical event means that action is required immediately because the scope of the problem has increased, investigate critical alerts or events immediately
	Critical = "critical"
	// Fatal event means that an error has occurred but it is too late to take any remedial action to address it
	Fatal = "fatal"
)

// Event describes emittable events
type Event struct {
	Service  string        `json:"service"` // event service source
	Name     string        `json:"name"`    // event name / type
	Message  string        `json:"message"` // event content
	Link     string        `json:"link"`    // link to external system
	Severity SeverityLevel `json:"level"`
	RootID   string        `json:"root_id"` // id of the root event or command which has created this event
}

// Command describe available commands to the services
type Command struct {
	ID     string // optional command id, can be used as root id for the event
	Name   string // command name: "ping", "deploy", "helm upgrade" etc
	Args   []string
	Sender string // username if command come from telegram
}

type CommandInfo struct {
	Service     string
	Name        string
	Example     string
	Description string
}

// Router instance routes incoming events and commads between services
type ServiceRouter interface {
	// registers service for routing
	RegisterService(name string, ctx *AppContext, service ServiceProvider) error
	// sends an event to all registered services
	EmitEvent(event Event) error
	// finds receiver and execute command returning result
	ExecuteCommandString(message, messageID, sender string) error
	ListCommands() map[string]CommandInfo
}

// ServiceProvider instance provides interface for executing commands
type ServiceProvider interface {
	// init service on start
	Init(ctx *AppContext) error
	// return available commands
	ListCommands() []CommandInfo
	// handle incoming event
	OnEvent(event Event) error
	// execute command returned in .Help
	DoCommand(command Command) error
}

type AppContext struct {
	Config *viper.Viper
	HTTP   *mux.Router
	Router ServiceRouter
}
