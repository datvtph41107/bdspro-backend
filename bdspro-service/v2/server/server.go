package server

import (
	"context"
	pb_bdspro "pb/types/bdspro"
)

type ServerBdspro struct {
	pb_bdspro.UnimplementedBdsproServiceServer
}

func (s *ServerBdspro) ListProductInternal(context.Context, *pb_bdspro.ProductInternalRequest) (*pb_bdspro.ProductInternalResponse, error) {
	return nil, nil
}
