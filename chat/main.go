package main

import (
	"google.golang.org/grpc"
	"log"
	chat "mse/chat/service"
	"mse/pkg"
	"mse/proto"
	"net"
)

//const (
//	port = ":12345"
//)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	defer log.Println("service exit")

	pkg.ParseItem([]pkg.FlagItem{
		pkg.ChatAddr,
		pkg.RedisAddr,
	})

	lis, err := net.Listen("tcp", pkg.ChatAddr.Addr())
	log.Println("listening at addr -", pkg.ChatAddr.Addr())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	rdsPS := pkg.NewRedisPubSub(pkg.RedisAddr.Addr(), "notify")
	//creds, err := credentials.NewServerTLSFromFile("res/certs/server.pem", "res/certs/server.key")
	//if err != nil {
	//	log.Fatalf("failed to create credentials: %v", err)
	//}
	//s := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor), grpc.Creds(creds))
	s := grpc.NewServer()
	cs := chat.NewChatService(rdsPS)
	defer cs.Close()
	proto.RegisterChatServer(s, cs)

	done := make(chan bool)

	go func() {
		defer close(done)

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		log.Println("service stop")
	}()

	cs.Run(done)
}
