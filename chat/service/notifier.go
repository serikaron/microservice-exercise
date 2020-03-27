package main

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	pb "mse/proto"
)

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
