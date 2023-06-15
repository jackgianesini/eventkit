package eventkit

import "fmt"

// errEventCallbacks : Error event callbacks
type errEventCallbacks struct {
	event  string
	errors []error
}

// Error : Error message
func (e *errEventCallbacks) Error() string {
	return fmt.Sprintf("event `%s` executed with %d errors", e.event, len(e.errors))
}

// Errors : List of errors
func (e *errEventCallbacks) Errors() []error {
	return e.errors
}

// NewErrEventCallbacks : Create new error event callbacks
func NewErrEventCallbacks(event string, errors []error) ErrEventCallbacks {
	return &errEventCallbacks{
		event:  event,
		errors: errors,
	}
}
