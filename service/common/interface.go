package common

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
	Channel  string
	Name     string
	Message  string
	Link     string
	Severity SeverityLevel
}

// Notifier is the interface which notification providers must implement
type Notifier interface {

	// sends notification about event
	Notify(event Event) error

	// disconnects from notification channel
	Close()
}

// Provider is the interface for communication with other services
type Provider interface {

	// perform some action with provider's resource
	Action() error

	// available commands
	Help()

	// disconnects from provider
	Close()
}
