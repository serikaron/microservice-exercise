package chat

import (
	"testing"
)

type TestData struct {
	msgToSend string
	wantMsg   string
}

type Checker interface {
	Check() error
	Setup(data *TestData)
}

type Caller interface {
	Call(interface{}) (interface{}, error)
	Setup(data *TestData)
}

type Monitor interface {
	Checker
	Attach(caller Caller) error
}

type Sender interface {
	Checker
	Send(caller Caller) error
}

func TestChatService(t *testing.T, caller Caller, monitor Monitor, sender Sender) {
	data := &TestData{
		msgToSend: "greetings",
		wantMsg:   "greetings",
	}
	caller.Setup(data)
	monitor.Setup(data)
	sender.Setup(data)

	go func() {
		err := monitor.Attach(caller)
		if err != nil {
			t.Error(err)
		}
	}()

	if err := sender.Send(caller); err != nil {
		t.Error(err)
	}

	if err := monitor.Check(); err != nil {
		t.Error(err)
	}

	if err := sender.Check(); err != nil {
		t.Error(err)
	}
}
