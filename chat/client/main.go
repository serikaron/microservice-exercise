package main

import (
	"bufio"
	"chat/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
)

const addr = "localhost:12345"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal("grpc.Dial:", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	log.Println("connected to server: ", addr)

	client := proto.NewChatClient(conn)

	req := &proto.ListenReq{}
	stream, err := client.Listen(context.Background(), req)
	if err != nil {
		log.Fatal("echo failed: ", err)
	}
	go func() {
		log.Println("I'm listening")
		for {
			inf, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("%v.Listen(_) = _, %v", client, err)
			}
			fmt.Println(inf.Msg)
		}
	}()

	req1 := &proto.ListenReq{}
	stream1, err := client.AnotherListen(context.Background(), req1)
	if err != nil {
		log.Fatal("echo failed: ", err)
	}
	go func() {
		log.Println("I'm listening")
		for {
			inf, err := stream1.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("%v.Listen(_) = _, %v", client, err)
			}
			fmt.Println(inf.Msg)
		}
	}()

	for {
		fmt.Println("Say:")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]
		req := &proto.SayReq{Msg: text}
		_, err := client.Say(context.Background(), req)
		if err != nil {
			log.Fatalln(client, ".Say failed: ", err)
		}
		//fmt.Println(text)
	}
}
