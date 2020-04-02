package main

import (
	"log"
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

	rdsPS := pkg.NewRedisPubSub(pkg.RedisAddr.Addr(), "notify")

	cs := NewChatService(rdsPS)
	defer cs.Close()

	done := make(chan bool)
	go func() {
		defer close(done)
		starter.StartServer(pkg.ChatAddr.Addr(), pkg.CertsPath.Pem(), pkg.CertsPath.Key(), func(gs *grpc.Server) {
			proto.RegisterChatServer(gs, cs)
		})
	}()
	cs.Run(done)
}
