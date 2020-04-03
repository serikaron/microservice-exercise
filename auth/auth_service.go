package main

import (
	"context"
	"log"
	"mse/pkg"
	"mse/pkg/jwt_token"
	"mse/proto"
	"net"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

//go:generate protoc -I../../proto --go_out=plugins=grpc,paths=source_relative:../../proto/ auth.proto

type AuthService struct {
}

func (as *AuthService) Run(addr string, pemFile string, keyFile string) {
	defer log.Printf("AuthService Quit")
	log.Printf("AuthService listening to %v\n", addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	creds, err := credentials.NewServerTLSFromFile(pemFile, keyFile)
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}
	s := grpc.NewServer(grpc.Creds(creds))
	proto.RegisterAuthServer(s, as)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (as *AuthService) Login(_ context.Context, in *proto.LoginReq) (*proto.LoginRsp, error) {
	log.Println("Login", in)

	if in.Username != in.Password {
		log.Println(pkg.LoginErr)
		return nil, pkg.LoginErr
	}

	key := jwt_token.NewHS256Key("1", pkg.SignKey)
	tokenString, err := jwt_token.Gen(pkg.Identity{Name: in.Username}, 86400, key)
	log.Println("AuthService.Login name:%s token:%s", in.Username, tokenString)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &proto.LoginRsp{Jwt: tokenString}, nil
}
