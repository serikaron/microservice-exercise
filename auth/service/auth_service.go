package main

import (
	"context"
	"github.com/dgrijalva/jwt-go/v4"
	"google.golang.org/grpc"
	"log"
	"mse/pkg"
	"mse/proto"
	"net"
)

//go:generate protoc -I../../proto --go_out=plugins=grpc,paths=source_relative:../../proto/ auth.proto

type AuthService struct {
}

func (as *AuthService) Run(addr string) {
	defer log.Printf("AuthService Quit")
	log.Printf("AuthService listening to %v\n", addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
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

	type UserClaims struct {
		Username string `json:"username"`
		jwt.StandardClaims
	}

	customClaim := UserClaims{
		Username: in.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.NewTime(86400),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaim)

	tokenString, err := token.SignedString([]byte(signKey))
	if err != nil {
		log.Printf("%v.Login create token failed, %v\n", as, err)
		return nil, pkg.JWTErr
	}

	log.Println("tokenString:", tokenString)

	return &proto.LoginRsp{Jwt: tokenString}, nil
}
