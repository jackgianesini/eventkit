package eventkit

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"unicode"
)

var (
	delimiter = "."
	prefix    = "On"
)

type event struct {
	data map[string][]*callback
	sync.RWMutex
}

// New : Create a new event instance
func New() EventKit {
	return &event{
		data: make(map[string][]*callback),
	}
}

// AddNewEventCallback : Add a new event callback
func (e *event) AddNewEventCallback(identifier string, callback *callback) {
	e.Lock()
	defer e.Unlock()

	e.data[identifier] = append(e.data[identifier], callback)
}

// Subscribe : Subscribe to an event with struct methods
func (e *event) Subscribe(payload any) (err error) {
	value := reflect.ValueOf(payload)

	if !value.IsValid() ||
		value.IsValid() && value.Kind() != reflect.Ptr ||
		value.IsValid() && value.Kind() == reflect.Ptr && value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("payload must be a struct")
	}

	for i := 0; i < value.NumMethod(); i++ {
		method := value.Method(i)
		methodName := value.Type().Method(i).Name
		if methodName[:len(prefix)] != prefix {
			continue
		}
		_ = e.GenericSubscribe(methodName[len(prefix):], method.Interface(), 2)
	}

	return
}

// SubscribeFunc : Subscribe to an event with a function
func (e *event) SubscribeFunc(listener string, callback any) error {
	return e.GenericSubscribe(listener, callback, 1)
}

// Trigger : Trigger an event
func (e *event) Trigger(name string, data ...any) error {
	return e.GenericTrigger(name, data)
}

// Resolve : Resolve the event name
func (e *event) Resolve(name string) string {
	s := strings.Split(name, delimiter)
	var f []string
	for _, v := range s {
		f = append(f, cases.Title(language.Und, cases.NoLower).String(v))
	}
	return strings.Join(f, "")
}

// ReverseResolve : Reverse resolve the event name
func (e *event) ReverseResolve(s string) string {
	var resolved []string
	for _, word := range s {
		if unicode.IsUpper(word) && len(resolved) > 0 {
			resolved = append(resolved, delimiter)
		}
		resolved = append(resolved, strings.ToLower(string(word)))
	}

	return strings.Join(resolved, "")
}

// GenericTrigger : Generic trigger
func (e *event) GenericTrigger(name string, data ...any) (err error) {
	log.WithFields(log.Fields{
		"trigger": name,
	}).Infof("event triggered")

	e.RLock()
	call, ok := e.data[e.Resolve(name)]
	e.RUnlock()

	if !ok {
		return
	}

	errAccumulator := make([]error, 0)
	for _, v := range call {

		if v.Disabled() {
			return
		}

		err := v.Call(data)
		if err != nil {
			errAccumulator = append(errAccumulator, err)
		}
	}

	if len(errAccumulator) > 0 {
		err = NewErrEventCallbacks(name, errAccumulator)
	}

	return
}

// GenericSubscribe : Generic subscribe
func (e *event) GenericSubscribe(eventName string, callback any, caller int) error {
	value := reflect.ValueOf(callback)

	if !value.IsValid() || value.IsValid() && value.Kind() != reflect.Func {
		return errors.New("callback must be a function")
	}

	from := ""
	info := reflect.TypeOf(callback)
	_, file, no, ok := runtime.Caller(caller)
	if ok {
		from = fmt.Sprintf("%s#%d->%s", file, no, info.String())
	}

	identifier := e.ReverseResolve(eventName)
	e.AddNewEventCallback(eventName, newCallBack(identifier, from, value))

	log.WithFields(wrapLogFields(log.Fields{
		"listener": identifier,
		"from":     from,
	})).Infof("event subscribed")

	return nil
}
