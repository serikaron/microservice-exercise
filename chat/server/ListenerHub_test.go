package main

import (
	"chat/proto"
	"errors"
	"fmt"
	"sync"
	"testing"
)

func (l *Listener) equal(another *Listener) bool {
	return l.name == another.name
}

func (l *Listener) String() string {
	return l.name
}

func MapEqual(lhs map[string]*Listener, rhs map[string]*Listener) error {
	if len(lhs) != len(rhs) {
		return errors.New(fmt.Sprintf("map size not equal, lhs:%d, rhs:%d", len(lhs), len(rhs)))
	}
	for key, lvalue := range lhs {
		rvalue, ok := rhs[key]
		if !ok {
			return errors.New(fmt.Sprintf("key:[%s] found in lhs but not in rhs", key))
		}
		if !lvalue.equal(rvalue) {
			return errors.New(fmt.Sprintf("listener not match, key:%s, lhs:%v, rhs:%v", key, lvalue, rvalue))
		}
	}
	return nil
}

func TestListenerHub_AddListener(t *testing.T) {
	tests := []struct {
		name          string
		listenerMap   map[string]*Listener
		listenerToAdd *Listener
		wantMap       map[string]*Listener
	}{
		{
			name:          "expect success",
			listenerMap:   map[string]*Listener{},
			listenerToAdd: &Listener{stream: nil, name: "listener", done: nil},
			wantMap: map[string]*Listener{
				"listener": {stream: nil, name: "listener", done: nil},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := sync.WaitGroup{}
			t.Helper()
			//lh := ListenerHub{listenerMap: tt.listenerMap, addChan: make(chan *Listener), notifier: &FakeNotifier{}, notifyChan: make(chan interface{})}
			lh := NewListenerHub(&FakeNotifier{})
			done := make(chan bool)
			go func() {
				defer wg.Done()
				wg.Add(1)
				lh.Run(done)
			}()
			lh.AddListener(tt.listenerToAdd)
			close(done)
			wg.Wait()
			err := MapEqual(lh.listenerMap, tt.wantMap)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

type FakeNotifier struct {
	gotMsgList []string
}

func (fn *FakeNotifier) Notify(stream interface{}, message interface{}) error {
	ss, ok := stream.(string)
	if !ok {
		return errors.New("stream error")
	}
	var ms string
	switch message.(type) {
	case string:
		ms = message.(string)
	case *proto.ListenRsp:
		ms = message.(*proto.ListenRsp).Msg
	default:
		return errors.New("message error")
	}
	fn.gotMsgList = append(fn.gotMsgList, fmt.Sprintf("%s-%s", ss, ms))
	return nil
}

func (fn *FakeNotifier) Check(wantMsgList []string) error {
	if len(fn.gotMsgList) != len(wantMsgList) {
		return errors.New(fmt.Sprintf("message list size not equal, got:%d, want:%d", len(fn.gotMsgList), len(wantMsgList)))
	}
	for i, g := range fn.gotMsgList {
		w := wantMsgList[i]
		if w != g {
			return errors.New(fmt.Sprintf("message list not equal, at:%d got:%s want:%s, full list got:%v want:%v", i, g, w, fn.gotMsgList, wantMsgList))
		}
	}
	return nil
}

func TestListenerHub_Notify(t *testing.T) {
	errChan := make(chan error, 1)
	tests := []struct {
		name            string
		listenerMap     map[string]*Listener
		notifier        *FakeNotifier
		notifyMessage   string
		wantMsgList     []string
		wantListenerMap map[string]*Listener
		wantErrMsg      string
	}{
		{
			"expect success",
			map[string]*Listener{
				"1st-listener": {"1st-listener", "1st-stream", nil},
				"2nd-listener": {"2nd-listener", "2nd-stream", nil},
			},
			&FakeNotifier{make([]string, 0, 2)},
			"notify-message",
			[]string{
				"1st-stream-notify-message",
				"2nd-stream-notify-message",
			},
			map[string]*Listener{
				"1st-listener": {"1st-listener", "1st-stream", nil},
				"2nd-listener": {"2nd-listener", "2nd-stream", nil},
			},
			"",
		},
		{
			"send notify error",
			map[string]*Listener{
				"1st-listener": {"1st-listener", "1st-stream", nil},
				"2nd-listener": {"2nd-listener", nil, errChan},
			},
			&FakeNotifier{make([]string, 0, 2)},
			"notify-message",
			[]string{
				"1st-stream-notify-message",
			},
			map[string]*Listener{
				"1st-listener": {"1st-listener", "1st-stream", nil},
			},
			"error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			wg := sync.WaitGroup{}

			hub := &ListenerHub{
				listenerMap: tt.listenerMap,
				addChan:     nil,
				notifier:    tt.notifier,
				notifyChan:  make(chan interface{}),
			}
			done := make(chan bool)
			go func() {
				defer wg.Done()
				wg.Add(1)
				hub.Run(done)
			}()

			hub.Notify(tt.notifyMessage)
			close(done)
			wg.Wait()
			//time.Sleep(50 * time.Millisecond)
			if err := tt.notifier.Check(tt.wantMsgList); err != nil {
				t.Error(err)
			}
			if err := MapEqual(hub.listenerMap, tt.wantListenerMap); err != nil {
				t.Error(err)
			}
			if tt.wantErrMsg != "" {
				err := <-errChan
				if err.Error() != tt.wantErrMsg {
					t.Error(fmt.Sprintf("check errChan failed, got:%s want:%s", err.Error(), tt.wantErrMsg))
				}
			}
		})
	}
}
