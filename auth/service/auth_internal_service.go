package main

import (
	"context"
	"github.com/dgrijalva/jwt-go/v4"
	"google.golang.org/grpc"
	"log"
	"mse/proto"
	"net"
)

//go:generate protoc -I../../proto --go_out=plugins=grpc,paths=source_relative:../../proto/ auth.internal.proto

type AuthInternalService struct {
}

func (ais *AuthInternalService) Run(addr string) {
	defer log.Printf("AuthInternalService Quit")
	log.Printf("AuthInternalService listening to %v\n", addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterAuthInternalServer(s, ais)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (ais *AuthInternalService) GetSignKey(_ context.Context, _ *proto.GetSignKeyReq) (*proto.GetSignKeyRsp, error) {
	return &proto.GetSignKeyRsp{Key: signKey, Kid: 1, Alg: jwt.SigningMethodHS256.Name}, nil
}
