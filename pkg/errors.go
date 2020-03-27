package pkg

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	LoginErr     = status.Error(codes.PermissionDenied, "username or password invalid")
	JWTErr       = status.Error(codes.Internal, "create jwt failed")
	MissingToken = status.Error(codes.Unauthenticated, "missing token")
	InvalidToken = status.Error(codes.Unauthenticated, "invalid token")
)
