package eventkit

import (
	log "github.com/sirupsen/logrus"
	"reflect"
)

// callback : represents a callback of an event
type callback struct {
	disabled  bool
	name      string
	from      string
	numIn     int
	reflectFn reflect.Value
}

// From : get the event name
func (c *callback) From() string {
	return c.from
}

// Name : get the callback name
func (c *callback) Name() string {
	return c.name
}

// NumIn : get the number of arguments
func (c *callback) NumIn() int {
	return c.numIn
}

// Disabled : check if the callback is disabled
func (c *callback) Disabled() bool {
	return c.disabled
}

// Disable : disable the callback
func (c *callback) Disable() {
	c.disabled = true
}

// Fn : get the reflection value of the callback
func (c *callback) Fn() reflect.Value {
	return c.reflectFn
}

// Call : call the callback
func (c *callback) Call(data []any) (err error) {
	if len(data[0].([]any)) != c.NumIn() {
		log.WithFields(wrapLogFields(log.Fields{
			"from": c.From(),
			"name": c.Name(),
			"want": c.NumIn(),
			"got":  len(data[0].([]any)),
			"kind": c.Fn().String(),
		})).Errorf("event callback argument mismatch")

		c.Disable()

		return
	}

	argv := make([]reflect.Value, c.NumIn())
	for i := 0; i < c.NumIn(); i++ {
		if data[0].([]any)[i] != nil {
			argv[i] = reflect.ValueOf(data[0].([]any)[i])
			continue
		}

		argv[i] = reflect.New(c.Fn().Type().In(i)).Elem()
	}

	err = try(func() {
		c.Fn().Call(argv)
	})

	if err != nil {
		return err
	}

	return
}

// newCallBack : create a new callback instance
func newCallBack(name string, from string, reflectFn reflect.Value) *callback {
	return &callback{
		name:      name,
		from:      from,
		numIn:     reflectFn.Type().NumIn(),
		reflectFn: reflectFn,
	}
}
