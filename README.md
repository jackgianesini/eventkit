[![Go](https://github.com/kitstack/eventkit/actions/workflows/coverage.yml/badge.svg)](https://github.com/kitstack/eventkit/actions/workflows/coverage.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kitstack/eventkit)
[![Go Report Card](https://goreportcard.com/badge/github.com/kitstack/eventkit)](https://goreportcard.com/report/github.com/kitstack/eventkit)
[![codecov](https://codecov.io/gh/kitstack/eventkit/branch/main/graph/badge.svg?token=3JRL5ZLSIH)](https://codecov.io/gh/kitstack/eventkit)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/kitstack/eventkit/blob/main/LICENSE)
[![Github tag](https://badgen.net/github/release/kitstack/eventkit)](https://github.com/kitstack/eventkit/releases)

# Overview

The `eventkit` package provides a simple event-driven programming mechanism in Go. It allows you to create and manage events, subscribe to events using struct methods or functions, and trigger events with optional data.

## Usage

### Creating a new event instance

To create a new event instance, use the `New()` function:

```go
eventInstance := eventkit.New()
```

### Subscribing to an event

#### Using struct methods

To subscribe to an event using struct methods, use the `Subscribe()` function. It takes a struct instance as the payload and automatically registers all methods starting with a specified prefix (default: "On") as event listeners.

```go
type MyEventPayload struct {
	// struct fields
}

func (m *MyEventPayload) OnEventName() {
	// event handler logic
}

// Subscribe to the event
eventInstance.Subscribe(&MyEventPayload{})
```

#### Using functions

To subscribe to an event using a function, use the `SubscribeFunc()` function. Provide a unique listener identifier and the callback function.

```go
func myEventListener(data ...interface{}) {
	// event handler logic
}

// Subscribe to the event
eventInstance.SubscribeFunc("myListener", myEventListener)
```

### Triggering an event

To trigger an event, use the `Trigger()` function. Provide the name of the event and optional data as arguments.

```go
err := eventInstance.Trigger("eventName", eventData)
if err != nil {
	// handle error
}
```

Note: The package utilizes various external dependencies, such as `logrus`, `golang.org/x/text/cases`, and `sync`, among others.

For more details on the package, please refer to the code documentation and comments within the source files.

