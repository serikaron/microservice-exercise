package pkg

//go:generate protoc -I../proto --go_out=plugins=grpc,paths=source_relative:../proto/ auth.proto

import (
	"auth/proto"
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

type AuthClient struct {
	conn   *grpc.ClientConn
	client proto.AuthClient
	ctx    context.Context
}

func NewAuthClient(addr string) *AuthClient {
	log.Println("try connect to auth-service:", addr)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		cancel()
		log.Fatalln("grpc.Dial failed:", err)
	}

	return &AuthClient{
		conn:   conn,
		client: proto.NewAuthClient(conn),
		ctx:    ctx,
	}
}

func (ac *AuthClient) Close() {
	_ = ac.conn.Close()
}

func (ac *AuthClient) Login(req *proto.LoginReq) (*proto.LoginRsp, error) {
	rsp, err := ac.client.Login(ac.ctx, req)
	if err != nil {
		log.Printf("AuthClient.Login failed: %v", err)
		return nil, err
	}
	return rsp, nil
}
