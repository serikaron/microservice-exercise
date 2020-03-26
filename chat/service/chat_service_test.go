package main

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"mse/chat/proto"
	"reflect"
	"testing"
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

func (fps *FakePubSub) Close() {
	close(fps.c)
}

var (
	fakePS    = &FakePubSub{c: make(chan []byte)}
	TypeError = status.Error(codes.Internal, "type error")
	monitor   = &Monitor{c: make(chan string, 1)}
)

type Monitor struct {
	c chan string
}

func (m *Monitor) Notify(_ interface{}, message interface{}) error {
	//_, ok := stream.(proto.Chat_ListenServer)
	//if !ok {
	//	log.Printf("stream type error, want:%s got:%v", "proto.Chat_listenServer", reflect.TypeOf(stream))
	//	return TypeError
	//}
	rsp, ok := message.(*proto.ListenRsp)
	if !ok {
		log.Printf("message type error, want:%s got:%v", "*proto.ListenRsp", reflect.TypeOf(message))
		return TypeError
	}

	m.c <- rsp.Msg
	return nil
}

func TestChatService(t *testing.T) {
	cs := NewChatService(monitor, fakePS)
	defer cs.Close()
	cs.hub.notifier = monitor

	done := make(chan bool)
	go cs.Run(done)

	go func() {
		_ = cs.Listen(nil, nil)
	}()

	wantMsg := "grettings"

	rsp, err := cs.Say(nil, &proto.SayReq{Msg: "grettings"})
	if err != nil {
		t.Fatal(err)
	}

	if rsp.Msg != wantMsg {
		t.Fatalf("say response invalid, want:%s got:%s", wantMsg, rsp.Msg)
	}

	monitorGot := <-monitor.c
	if monitorGot != wantMsg {
		t.Fatalf("monitor message error, want:%s got:%s", wantMsg, monitorGot)
	}
}
