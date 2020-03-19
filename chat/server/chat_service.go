package main

import (
	"chat/proto"
	"context"
	"log"
	"time"
)

//go:generate protoc -I../proto --go_out=plugins=grpc,paths=source_relative:../proto/ chat.proto

type ChatService struct {
	streamList []proto.Chat_ListenServer
}

func NewChatService() *ChatService {
	return &ChatService{streamList: make([]proto.Chat_ListenServer, 0)}
}

func (cs *ChatService) Listen(req *proto.ListenReq, stream proto.Chat_ListenServer) error {
	cs.streamList = append(cs.streamList, stream)
	for {
		time.Sleep(10 * time.Millisecond)
	}
	return nil
}

func (cs *ChatService) AnotherListen(req *proto.ListenReq, stream proto.Chat_AnotherListenServer) error {
	cs.streamList = append(cs.streamList, stream)
	for {
		time.Sleep(10 * time.Millisecond)
	}
	return nil
}

func (cs *ChatService) Say(_ context.Context, in *proto.SayReq) (*proto.SayRsp, error) {
	log.Println("req msg: ", in.Msg)
	inf := proto.ListenRsp{Msg: in.Msg}
	for _, stream := range cs.streamList {
		err := stream.Send(&inf)
		if err != nil {
			log.Println(stream, ".Send failed: ", err)
		}
	}
	return &proto.SayRsp{}, nil
}
