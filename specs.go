package eventkit

type EventKit interface {
	Subscribe(payload any) error
	SubscribeFunc(listener string, callback any) error
	Trigger(name string, data ...any) error
}

type ErrEventCallbacks interface {
	Error() string
	Errors() []error
}
