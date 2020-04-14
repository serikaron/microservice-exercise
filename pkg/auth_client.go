package pkg

//go:generate protoc -I../../proto --go_out=plugins=grpc,paths=source_relative:../../proto/ auth.proto

import (
	"context"
	"log"
	"mse/proto"
	"time"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

type AuthClient struct {
	conn   *grpc.ClientConn
	client proto.AuthClient
	ctx    context.Context
}

func NewAuthClient(addr string, pemFile string) *AuthClient {
	log.Println("try connect to auth-service:", addr)

	creds, err := credentials.NewClientTLSFromFile(pemFile, "serika-server")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(creds), grpc.WithBlock())
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
	rsp, err := ac.client.Login(context.Background(), req)
	if err != nil {
		log.Printf("AuthClient.Login failed: %v", err)
		return nil, err
	}
	return rsp, nil
}

func (ac *AuthClient) OAuthLogin(req *proto.OAuthLoginReq) (*proto.OAuthLoginRsp, error) {
	rsp, err := ac.client.OAuthLogin(context.Background(), req)
	if err != nil {
		log.Printf("AuthClient.OAuthLogin failed: %v", err)
		return nil, err
	}
	return rsp, nil
}
