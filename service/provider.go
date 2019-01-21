package main

type Provider interface {

	// calls once during initialisation
	Init() error

	// perform some action with provider's resource
	Action(event Event) error

	// disconnects from provider
	Close()
}
