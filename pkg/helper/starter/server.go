package starter

import (
	"context"
	"log"
	"mse/pkg"
	"mse/pkg/jwt_token"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type registerFunc func(gs *grpc.Server)

func StartServer(addr string, pemFile string, keyFile string, rf registerFunc) {
	lis, err := net.Listen("tcp", addr)
	log.Println("listening at addr -", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile(pemFile, keyFile)
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor), grpc.Creds(creds))
	//s := grpc.NewServer()

	rf(s)

	log.Println("start grpc server")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Println("grpc server stop")
}

func unaryInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	//log.Printf("context:%v req:%v info:%v handler:%v", ctx, req, info, handler)
	//return nil, nil
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, pkg.MissingToken
	}
	id, err := valid(md["authorization"])
	if err != nil {
		return nil, pkg.InvalidToken
	}
	md.Set("name", id.Name)
	ctx = metadata.NewIncomingContext(context.Background(), md)
	m, err := handler(ctx, req)
	if err != nil {
		log.Printf("RPC failed with error %v", err)
	}
	return m, err
}

func valid(authorization []string) (*pkg.Identity, error) {
	if len(authorization) < 1 {
		return nil, pkg.MissingToken
	}
	tokenString := strings.TrimPrefix(authorization[0], "Bearer ")

	id, err := jwt_token.Parse(tokenString, func(kid string) *jwt_token.Key {
		return jwt_token.NewHS256Key("1", pkg.SignKey)
	})

	if err != nil {
		return nil, err
	}

	return &id, nil
}
