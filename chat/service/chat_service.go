package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"math/rand"
	"mse/chat/pkg"
	pb "mse/chat/proto"
	"time"
)

//go:generate protoc -I../proto --go_out=plugins=grpc,paths=source_relative:../proto/ chat.proto

type ChatNotifier struct{}

func (cn *ChatNotifier) Notify(stream interface{}, message interface{}) error {
	s, ok := stream.(pb.Chat_ListenServer)
	if !ok {
		log.Printf("ChatNotifier.Notify Not a chat stream, got:%v\n", stream)
		return status.Error(codes.Internal, "Not a chat stream")
	}
	m, ok := message.(*pb.ListenRsp)
	if !ok {
		log.Printf("ChatNotifier.Notify Not a ListenRsp, got:%v\n", message)
		return status.Error(codes.Internal, "Not a stream message")
	}
	return s.Send(m)
}

type ChatService struct {
	name   string
	hub    *ListenerHub
	pubsub pkg.PubSub
}

func NewChatService(notifier Notifier, pubsub pkg.PubSub) *ChatService {
	r := rand.Int()
	name := fmt.Sprintf("%d", r)
	hub := NewListenerHub(notifier)

	return &ChatService{
		name:   name,
		hub:    hub,
		pubsub: pubsub,
	}
}

func (cs *ChatService) Close() {
	cs.pubsub.Close()
}

func (cs *ChatService) Run(done chan bool) {
	d := make(chan bool)
	go func() {
		defer close(d)
		cs.hub.Run(done)
	}()

	chn := cs.pubsub.Subscribe()
	for {
		select {
		case <-d:
			return
		case msg := <-chn:
			cs.notify(msg)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (cs *ChatService) Listen(_ *pb.ListenReq, stream pb.Chat_ListenServer) error {
	log.Printf("ChatService[%s].Listen", cs.name)
	done := make(chan error)
	cs.hub.AddListener(&Listener{
		done:   done,
		stream: stream,
		name:   fmt.Sprintf("%d", rand.Int()),
	})
	err := <-done
	if err != nil {
		log.Printf("ChatServer.Listen err:%v\n", err)
	}
	return err
}

func (cs *ChatService) Say(_ context.Context, in *pb.SayReq) (rsp *pb.SayRsp, err error) {
	log.Printf("ChatService[%s].Say in.Msg:%s", cs.name, in.Msg)
	defer func() {
		if err := recover(); err != nil {
			log.Printf("ChatService[%s] %v", cs.name, err)
			rsp = nil
			err = status.Error(codes.Internal, "Say failed")
		}
	}()
	inf := &pb.ListenRsp{Msg: in.Msg}
	data, err := proto.Marshal(inf)
	if err != nil {
		panic(fmt.Errorf("pb marshal failed, err:%v", err))
	}
	log.Printf("ChatService[%s].Say data.len:%d", cs.name, len(data))
	err = cs.pubsub.Publish(data)
	if err != nil {
		panic(fmt.Errorf("publish failed, err:%v", err))
	}
	return &pb.SayRsp{Msg: in.Msg}, nil
}

func (cs *ChatService) notify(data []byte) {
	log.Printf("ChatService[%s].notify data.len:%d", cs.name, len(data))
	rsp := &pb.ListenRsp{}
	if err := proto.Unmarshal(data, rsp); err != nil {
		log.Printf("ChatService[%s].notify unmarshal pb failed, err:%v", cs.name, err)
		return
	}
	cs.hub.Notify(rsp)
}
