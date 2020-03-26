package pkg

import (
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
	"mse/proto"
	"time"
)

type ChatClient struct {
	conn   *grpc.ClientConn
	client proto.ChatClient
	ctx    context.Context
}

func NewChatClient(addr string) *ChatClient {
	log.Println("try connect to chat-service:", addr)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		cancel()
		log.Fatalln("grpc.Dial failed:", err)
	}

	return &ChatClient{
		conn:   conn,
		client: proto.NewChatClient(conn),
		ctx:    ctx,
	}
}

func (cc *ChatClient) Listen(done chan bool) (chan string, error) {
	stream, err := cc.client.Listen(context.Background(), &proto.ListenReq{})
	if err != nil {
		log.Println("ChatClient.Listen failed:", err)
		return nil, err
	}

	c := make(chan string)
	go func() {
		defer close(c)

		for {
			select {
			case <-done:
				return
			default:
				inf, err := stream.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					log.Printf("ChatClient.Listen, stream.Recv failed %v", err)
					return
				}
				c <- inf.Msg
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
	return c, nil
}

func (cc *ChatClient) Say(req *proto.SayReq) error {
	rsp, err := cc.client.Say(context.Background(), req)
	if err != nil {
		log.Println("ChatClient.Say failled:", err)
	} else {
		log.Println("ChatClient.Say rsp.Msg:", rsp.Msg)
	}
	return err
}
