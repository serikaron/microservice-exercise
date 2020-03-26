package main

import (
	"chat/proto"
	"context"
	"sync"
	"testing"
	"time"
)

type FakePubSub struct {
	c chan []byte
}

func (fps *FakePubSub) Publish(data []byte) error {
	fps.c <- data
	return nil
}

func (fps *FakePubSub) Subscribe() chan []byte {
	return fps.c
}

func (fps *FakePubSub) Close() {}

var (
	fakePS = &FakePubSub{c: make(chan []byte)}
)

func TestChatService_Say(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantMsg string
	}{
		{
			"expect ok",
			"message from test",
			"message from test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			cs := NewChatService(fakePS)
			fakeNotifier := &FakeNotifier{gotMsgList: make([]string, 0)}
			cs.hub.notifier = fakeNotifier
			ld := make(chan error)
			cs.hub.listenerMap = map[string]*Listener{
				"listener": {"listener", "stream", ld},
			}
			wg := sync.WaitGroup{}
			done := make(chan bool)
			go func() {
				defer wg.Done()
				wg.Add(1)
				cs.Run(done)
			}()
			go func() { <-ld }()
			got, err := cs.Say(context.Background(), &proto.SayReq{Msg: tt.msg})
			time.AfterFunc(time.Second, func() {
				close(done)
			})
			wg.Wait()
			if err != nil {
				t.Error(err)
			}
			if got.Msg != tt.wantMsg {
				t.Error("error")
			}
			err = fakeNotifier.Check([]string{"stream-message from test"})
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestChatService_Listen(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"simple check",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			cs := NewChatService(fakePS)
			wg := sync.WaitGroup{}
			done := make(chan bool)
			go func() {
				defer wg.Done()
				wg.Add(1)
				cs.Run(done)
			}()
			time.AfterFunc(time.Second, func() {
				close(done)
			})
			err := cs.Listen(nil, nil)
			wg.Wait()
			if len(cs.hub.listenerMap) != 1 {
				t.Errorf("len of listenerMap invalid, want:%d got:%d", 1, len(cs.hub.listenerMap))
			}
			if err != nil {
				t.Error(err)
			}
		})
	}
}
