package main

import (
	"chat/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":12345"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterChatServer(s, NewChatService())
	log.Println("server start, listening at port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Println("server stop")
}
