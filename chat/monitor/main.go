package main

//go:generate protoc -I../proto --go_out=plugins=grpc,paths=source_relative:../proto/ chat.proto

import (
	"chat/pkg"
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	host := flag.String("chat-service-host", "chat-service", "chat service host")
	port := flag.Uint("chat-service-port", 0, "chat service port")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)
	client := pkg.NewChatClient(addr)

	done := make(chan bool)
	defer close(done)

	// delay 5 sec
	time.Sleep(5 * time.Second)

	c, err := client.Listen(done)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("monitoring msg")
	log.Println("--------------------------------")
	for msg := range c {
		log.Println(msg)
	}
}
