package chat

import (
	"context"
	"github.com/dgrijalva/jwt-go/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"mse/pkg"
	"strings"
)

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	//log.Printf("context:%v req:%v info:%v handler:%v", ctx, req, info, handler)
	//return nil, nil
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, pkg.MissingToken
	}
	if !valid(md["authorization"]) {
		return nil, pkg.InvalidToken
	}
	md.Set("name", "marry")
	ctx = metadata.NewIncomingContext(context.Background(), md)
	m, err := handler(ctx, req)
	if err != nil {
		log.Printf("RPC failed with error %v", err)
	}
	return m, err
}

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	tokenString := strings.TrimPrefix(authorization[0], "Bearer ")

	type UserClaims struct {
		Username string `json:"username"`
		jwt.StandardClaims
	}

	token, _ := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret-key"), nil
	})

	//if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
	//	log.Printf("%v %v", claims.Username, claims.StandardClaims.ExpiresAt)
	//} else {
	//	log.Println(err)
	//return false
	//}

	claims, _ := token.Claims.(*UserClaims)
	log.Printf("%v %v", claims.Username, claims.StandardClaims.ExpiresAt)

	return true
}
