package main

//go:generate protoc -I../../proto --go_out=plugins=grpc,paths=source_relative:../../proto/ chat.proto

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"mse/pkg"
	"mse/proto"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	host := flag.String("chat-service-host", "chat-service", "chat service host")
	port := flag.Uint("chat-service-port", 0, "chat service port")
	//msg := flag.String("chater-msg", "", "chater msg")
	flag.Parse()

	// delay 10 sec
	time.Sleep(10 * time.Second)

	addr := fmt.Sprintf("%s:%d", *host, *port)
	client := pkg.NewChatClient(addr)

	r := rand.Int()
	msg := fmt.Sprintf("%d", r)
	for {
		log.Println("Send Msg:", msg)
		err := client.Say(&proto.SayReq{Msg: msg})
		if err != nil {
			log.Fatalln(err)
		}
		time.Sleep(time.Second)
	}
}
