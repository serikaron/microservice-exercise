package main

import (
	"log"
	"mse/chat/internal"
	"mse/pkg"
	"mse/pkg/helper/starter"
	"mse/proto"

	"google.golang.org/grpc"
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
		pkg.CertsPath,
		pkg.IntegrationKey,
		pkg.IntegrationEnable,
	})

	cs := internal.NewChatService(pkg.RedisAddr.Addr())
	defer cs.Close()

	starter.StartServer(pkg.ChatAddr.Addr(), pkg.CertsPath.Pem(), pkg.CertsPath.Key(), func(gs *grpc.Server) {
		proto.RegisterChatServer(gs, cs)
	})
}
