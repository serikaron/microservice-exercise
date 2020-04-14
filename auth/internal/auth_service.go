package internal

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

	id, err := IdentifyWithPassword(in.Username, in.Password)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	key := jwt_token.NewHS256Key("1", pkg.SignKey)
	tokenString, err := jwt_token.Gen(*id, 86400, key)
	log.Printf("AuthService.Login name:%s token:%s", in.Username, tokenString)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &proto.LoginRsp{Jwt: tokenString}, nil
}

func (as *AuthService) OAuthLogin(_ context.Context, in *proto.OAuthLoginReq) (*proto.OAuthLoginRsp, error) {
	log.Println("OauthLogin", in)

	id, err := IdentifyWithOAuth(in.Code)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	tokenString, err := jwt_token.Gen(*id, 86400, jwt_token.NewHS256Key("1", pkg.SignKey))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &proto.OAuthLoginRsp{Jwt: tokenString}, nil
}
