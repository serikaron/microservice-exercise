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
	s := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor), grpc.StreamInterceptor(streamInterceptor), grpc.Creds(creds))
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
	m, err := handler(context.WithValue(ctx, "name", id.Name), req)
	if err != nil {
		log.Printf("RPC failed with error %v", err)
	}
	return m, err
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// Call 'handler' to invoke the stream handler before this function returns
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return pkg.MissingToken
	}
	id, err := valid(md["authorization"])
	if err != nil {
		return pkg.InvalidToken
	}
	err = handler(srv, &wrappedStream{stream, context.WithValue(stream.Context(), "name", id.Name)})
	if err != nil {
		log.Printf("RPC failed with error %v", err)
	}
	return err
}

func valid(authorization []string) (*pkg.Identity, error) {
	if len(authorization) < 1 {
		return nil, pkg.MissingToken
	}
	tokenString := strings.TrimPrefix(authorization[0], "Bearer ")

	if pkg.IntegrationEnable.Val && pkg.IntegrationKey.Val == tokenString {
		return &pkg.Identity{Name: "integration-name"}, nil
	}

	id, err := jwt_token.Parse(tokenString, func(kid string) *jwt_token.Key {
		return jwt_token.NewHS256Key("1", pkg.SignKey)
	})

	if err != nil {
		return nil, err
	}

	return &id, nil
}
