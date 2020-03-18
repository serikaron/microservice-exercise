package main

import (
	"context"
	"echo/proto"
	"google.golang.org/grpc"
	"log"
	"os"
)

//go:generate protoc -I./ --go_out=plugins=grpc,paths=source_relative:./proto/ echo.proto

//const addr = "localhost:55555"

func main() {
	addr := os.Args[1]

	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal("grpc.Dial:", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	client := proto.NewEchoClient(conn)
	req := &proto.EchoReq{Msg: "an echo message"}
	rsp, err := client.Echo(context.Background(), req)
	if err != nil {
		log.Fatal("echo failed: ", err)
	}

	log.Println("rsp: ", rsp)
	if rsp.Msg != "an echo message" {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
