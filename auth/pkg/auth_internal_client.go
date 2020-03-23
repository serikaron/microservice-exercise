package pkg

import (
	"auth/proto"
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

type AuthInternalClient struct {
	conn   *grpc.ClientConn
	client proto.AuthInternalClient
	ctx    context.Context
}

func NewAuthInternalClient(addr string) *AuthInternalClient {
	log.Println("try connect to auth-service:", addr)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		cancel()
		log.Fatalln("grpc.Dial failed:", err)
	}

	return &AuthInternalClient{
		conn:   conn,
		client: proto.NewAuthInternalClient(conn),
		ctx:    ctx,
	}
}

func (ac *AuthInternalClient) Close() {
	_ = ac.conn.Close()
}

func (ac *AuthInternalClient) GetSignKey(req *proto.GetSignKeyReq) (*proto.GetSignKeyRsp, error) {
	rsp, err := ac.client.GetSignKey(ac.ctx, req)
	if err != nil {
		log.Printf("AuthInternalClient.GetSignKey failed: %v", err)
		return nil, err
	}
	return rsp, nil
}
