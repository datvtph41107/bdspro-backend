package client

import (
	_middleware "common/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewGrpcClient tạo kết nối gRPC đến organization service
func NewGrpcClient() (*grpc.ClientConn, error) {
	// Kết nối đến organization service
	conn, err := grpc.Dial(
		"localhost:8207",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(_middleware.StreamClientInterceptor),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
