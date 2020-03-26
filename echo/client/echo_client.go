package main

import (
	"context"
	"echo/proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

//go:generate protoc -I./ --go_out=plugins=grpc,paths=source_relative:./proto/ echo.proto

const addr = "echo-service:55555"

func main() {
	//addr := os.Args[1]
	log.Println("connecting to", addr)

	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal("grpc.Dial:", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	client := proto.NewEchoClient(conn)
	for {
		req := &proto.EchoReq{Msg: "an echo message"}
		rsp, err := client.Echo(context.Background(), req)
		if err != nil {
			log.Fatal("echo failed: ", err)
		}

		log.Println(rsp.Msg)
		time.Sleep(time.Second)
	}
}
