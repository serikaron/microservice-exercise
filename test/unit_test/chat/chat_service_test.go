package chat

import (
	"context"
	"google.golang.org/grpc"
	"mse/chat/service"
	"mse/proto"
	"testing"
)

func TestChatService(t *testing.T) {
	t.Run("message_from_say_will_be_notified", message_from_say_will_be_notified)
}

func message_from_say_will_be_notified(t *testing.T) {
	fakePS := &FakePubSub{c: make(chan []byte)}
	cs := chat.NewChatService(fakePS)
	defer cs.Close()

	done := make(chan bool)
	go cs.Run(done)

	s := &testStream{
		c: make(chan string),
	}
	go func() {
		_ = cs.Listen(nil, s)
	}()

	wantMsg := "grettings"

	_, err := cs.Say(context.Background(), &proto.SayReq{Msg: "grettings"})
	if err != nil {
		t.Fatal(err)
	}

	msgGot := <-s.c
	if msgGot != wantMsg {
		t.Fatalf("monitor message error, want:%s got:%s", wantMsg, msgGot)
	}
}

type testStream struct {
	grpc.ServerStream
	c chan string
}

func (x *testStream) Send(m *proto.ListenRsp) error {
	go func() { x.c <- m.Msg }()
	return nil
}

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
