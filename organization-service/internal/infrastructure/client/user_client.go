package client

import (
	"context"
	"fmt"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"github.com/tuandq2112/organization-service/appconfig"
	"google.golang.org/grpc"
)

type UserClient interface {
	GetProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error)
}
type userClient struct {
	client userpb.ProfileServiceClient
}
func NewUserClient() UserClient {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", appconfig.USER_SERVICE_HOST, appconfig.USER_SERVICE_PORT), grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to chat service: %v", err))
	}
	return &userClient{
		client: userpb.NewProfileServiceClient(conn),
	}
}

func (c *userClient) GetProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error) {
	return c.client.GetProfileByIds(ctx, in)
}

