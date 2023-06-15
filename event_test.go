package eventkit

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"testing"
)

type eventTest struct {
}

func (test *eventTest) NoListener() {
}

func (test *eventTest) OnTest() {
	log.Print("eventTest.OnTest")
}

func (test *eventTest) OnTe() {
	log.Print("eventTest.OnTest")
}

type EventTestSuite struct {
	suite.Suite
	log      bytes.Buffer
	instance EventKit
}

func (test *EventTestSuite) SetupTest() {
	test.instance = New()
	test.log = bytes.Buffer{}
	log.SetOutput(&test.log)
}

func (test *EventTestSuite) TestSubscribeFuncErr() {
	test.EqualError(test.instance.SubscribeFunc("test", nil), "callback must be a function")
	test.EqualError(test.instance.SubscribeFunc("test", 1), "callback must be a function")
}

func (test *EventTestSuite) TestSubscribeFunc() {
	err := test.instance.SubscribeFunc("OnTest", func() {
		log.Print("calling")
	})
	test.NoError(err)
	test.Contains(test.log.String(), "eventkit/event.go")
	test.Contains(test.log.String(), "->func()\" listener=on.test package=eventkit")
}

func (test *EventTestSuite) TestSubscribeFuncWithArgs() {
	err := test.instance.SubscribeFunc("OnTestArgs", func(message string) {
		log.Print(message)
	})
	test.NoError(err)
	test.Contains(test.log.String(), "eventkit/event.go")
	test.Contains(test.log.String(), "->func(string)\" listener=on.test.args package=eventkit")
}

func (test *EventTestSuite) TestTriggerFunc() {
	test.TestSubscribeFunc()

	err := test.instance.Trigger("on.test")
	test.Contains(test.log.String(), "calling")
	test.NoError(err)
}

func (test *EventTestSuite) TestTriggerWithFuncPanic() {
	err := test.instance.SubscribeFunc("OnTestPanicWithText", func() {
		panic("panicWithText")
	})
	test.NoError(err)

	err = test.instance.Trigger("on.test.panic.with.text")
	test.Error(err)
	test.IsType(err, &errEventCallbacks{}, err)
	test.Contains(err.Error(), "event `on.test.panic.with.text` executed with 1 errors")
	test.Contains(err.(ErrEventCallbacks).Errors()[0].Error(), "panicWithText")

	err = test.instance.SubscribeFunc("OnTestPanicWithError", func() {
		panic(errors.New("panicWithError"))
	})
	test.NoError(err)

	err = test.instance.Trigger("on.test.panic.with.error")
	test.Error(err)
	fmt.Print(err.Error())
	test.IsType(err, &errEventCallbacks{}, err)
	test.Contains(err.Error(), "event `on.test.panic.with.error` executed with 1 errors")
	test.Contains(err.(ErrEventCallbacks).Errors()[0].Error(), "panicWithError")
}

func (test *EventTestSuite) TestTriggerFuncUnknown() {
	test.TestSubscribeFunc()

	err := test.instance.Trigger("on.unknown")
	test.NotContains(test.log.String(), "calling")
	test.NoError(err)
}

func (test *EventTestSuite) TestTriggerFuncWithArgs() {
	test.TestSubscribeFuncWithArgs()

	err := test.instance.Trigger("OnTestArgs", "with_args")
	test.Contains(test.log.String(), "with_args")
	test.NoError(err)
}

func (test *EventTestSuite) TestTriggerFuncWithNilArgs() {
	err := test.instance.SubscribeFunc("OnTestNilArgs", func(message *string) {
		log.Print(message)
	})
	test.NoError(err)

	err = test.instance.Trigger("OnTestNilArgs", nil)
	test.Contains(test.log.String(), "nil")
	test.NoError(err)
}

func (test *EventTestSuite) TestTriggerFuncWithBadArgs() {
	test.TestSubscribeFuncWithArgs()

	err := test.instance.Trigger("OnTestArgs")
	test.NotContains(test.log.String(), "with_args")
	test.Contains(test.log.String(), "event callback argument mismatch")
	fmt.Print(test.log.String())
	test.NoError(err)

	test.log.Reset()
	err = test.instance.Trigger("OnTestArgs")
	test.NoError(err)

	test.Contains(test.log.String(), "trigger=OnTestArgs")
	test.NotContains(test.log.String(), "event callback argument mismatch")
}

func (test *EventTestSuite) TestSubscribe() {
	err := test.instance.Subscribe(&eventTest{})
	test.NoError(err)

	test.Contains(test.log.String(), "event subscribed")
	test.Contains(test.log.String(), "eventkit/event_test.go")
	test.Contains(test.log.String(), "->func()\" listener=test package=eventkit")
}

func (test *EventTestSuite) TestSubscribeErr() {
	err := test.instance.Subscribe(nil) // invalid reflect
	test.EqualError(err, "payload must be a struct")

	err = test.instance.Subscribe(func() {}) // not a struct
	test.EqualError(err, "payload must be a struct")

	err = test.instance.Subscribe(eventTest{}) // not a struct
	test.EqualError(err, "payload must be a struct")
}

func TestEventTestSuite(t *testing.T) {
	suite.Run(t, new(EventTestSuite))
}
