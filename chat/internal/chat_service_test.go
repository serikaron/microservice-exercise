package internal_test

import (
	"context"
	"errors"
	"log"
	"mse/chat/internal"
	"mse/pkg"
	"mse/proto"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func init() {
	pkg.RedisAddr.Attach()
}

func Test_monitor_can_know_what_sender_say(t *testing.T) {
	log.Println("test start")
	cs := internal.NewChatService(pkg.RedisAddr.Addr())
	//time.Sleep(100 * time.Millisecond)
	marry := addSpyMonitor(cs, "Marry")
	cherry := addSpyMonitor(cs, "Cherry")

	time.Sleep(10 * time.Millisecond)
	say(cs, "John", "Hello")

	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		defer wg.Done()
		messageOf(marry.sc).withT(t).shouldBe("John: Hello")
	}()
	go func() {
		defer wg.Done()
		messageOf(cherry.sc).withT(t).shouldBe("John: Hello")
	}()
	go func() {
		defer wg.Done()
		errorOf(marry.ec).withT(t).shouldBe(nil)
	}()
	go func() {
		defer wg.Done()
		errorOf(cherry.ec).withT(t).shouldBe(nil)
	}()

	time.Sleep(10 * time.Millisecond)
	cs.Close()
	wg.Wait()
	log.Println("test end")
}

func say(cs *internal.ChatService, name string, msg string) {
	_, _ = cs.Say(contextWithName(name), &proto.SayReq{Msg: msg})
}

func contextWithName(name string) context.Context {
	return context.WithValue(context.Background(), "name", name)
}

type spyMonitor struct {
	name string
	sc   chan string
	ec   chan error
	grpc.ServerStream
}

func (m *spyMonitor) Send(r *proto.ListenRsp) error {
	m.sc <- r.Msg
	return nil
}

func (m *spyMonitor) Context() context.Context {
	return contextWithName(m.name)
}

func addSpyMonitor(cs *internal.ChatService, name string) *spyMonitor {
	m := &spyMonitor{
		name:         name,
		sc:           make(chan string, 1),
		ec:           make(chan error, 1),
		ServerStream: nil,
	}
	go func() {
		err := cs.Listen(nil, m)
		m.ec <- err
		close(m.sc)
	}()

	return m
}

type errMonitor struct {
	name string
	ec   chan error
	grpc.ServerStream
}

func (m *errMonitor) Send(r *proto.ListenRsp) error {
	return errors.New("send error")
}

func (m *errMonitor) Context() context.Context {
	return contextWithName(m.name)
}

func addErrMonitor(cs *internal.ChatService, name string) *spyMonitor {
	m := &spyMonitor{
		name:         name,
		sc:           make(chan string, 1),
		ec:           make(chan error, 1),
		ServerStream: nil,
	}
	go func() {
		err := cs.Listen(nil, m)
		m.ec <- err
		close(m.sc)
	}()

	return m
}
