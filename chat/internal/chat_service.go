package internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"mse/pkg"
	pb "mse/proto"

	"github.com/go-redis/redis/v7"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const rdsChannel = "rdsChannel"

type ChatService struct {
	name string
	hub  *SafeHub
	rds  *redis.Client
	ps   *redis.PubSub
}

func NewChatService(rdsAddr string) *ChatService {
	log.Printf("NewChatService rdsAddr:%s", rdsAddr)
	r := rand.Int()
	name := fmt.Sprintf("%d", r)

	rds := redis.NewClient(&redis.Options{Addr: rdsAddr})
	ps := rds.Subscribe(rdsChannel)

	cs := &ChatService{
		name: name,
		hub:  NewSafeHub(),
		rds:  rds,
		ps:   ps,
	}

	go cs.run()

	return cs
}

func (cs *ChatService) Close() {
	_ = cs.ps.Close()
	_ = cs.rds.Close()
	cs.hub.Close()
}

func (cs *ChatService) run() {
	c := cs.ps.Channel()
	for msg := range c {
		cs.notify(msg.Payload)
	}
}

func nameFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata not found in context")
	}
	names := md["name"]
	if len(names) < 1 {
		return "", errors.New("name not found in metadata")
	}
	return names[0], nil
}

func (cs *ChatService) Listen(_ *pb.ListenReq, stream pb.Chat_ListenServer) error {
	log.Printf("ChatService[%s].Listen", cs.name)
	name, err := nameFromContext(stream.Context())
	if err != nil {
		log.Println(err)
		return pkg.IdentityNotFound
	}
	l := NewListener(name)
	cs.hub.Add(l)
	err = l.Listen(func(msg string) error {
		rsp := &pb.ListenRsp{Msg: msg}
		return stream.Send(rsp)
	})
	if err != nil {
		cs.hub.Remove(name)
		log.Println(err)
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}

func (cs *ChatService) Say(ctx context.Context, in *pb.SayReq) (rsp *pb.SayRsp, err error) {
	log.Printf("ChatService[%s].Say in.Msg:%s", cs.name, in.Msg)
	defer func() {
		if err := recover(); err != nil {
			log.Printf("ChatService[%s] %v", cs.name, err)
			rsp = nil
			err = status.Error(codes.Internal, "Say failed")
		}
	}()

	name, err := nameFromContext(ctx)
	if err != nil {
		panic(err)
	}
	msg := Chat(name, in.Msg)
	cs.rds.Publish(rdsChannel, msg)
	return &pb.SayRsp{Msg: msg}, nil
}

func (cs *ChatService) notify(msg string) {
	log.Printf("ChatService[%s].notify msg:%s", cs.name, msg)
	cs.hub.Notify(msg)
}
