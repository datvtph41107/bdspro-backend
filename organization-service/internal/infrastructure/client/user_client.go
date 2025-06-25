package client

import (
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	"organization/env"
	organizationpb "pb/types/organization"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"time"

	"google.golang.org/grpc"
)

type IUserClient interface {
	GetProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error)
	GetMapProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (map[uint64]*sharepb.ProfileItem, error)
	MapMemberToDealPb(ctx context.Context, deals []*organizationpb.GroupDeal)
}

// @bind: organization/internal/infrastructure/client.IUserClient
type UserClient struct {
	client userpb.ProfileServiceClient
}

func NewUserClient() *UserClient {
	userClient := &UserClient{client: nil}
	go func() {
		url := "localhost"
		name := "user_grpc"
		if env.ENV_RUNTIME == "release" {
			url = name
		}

		conn, err := grpc.Dial(fmt.Sprintf("%s:8201", url),
			grpc.WithInsecure(),
			grpc.WithBackoffMaxDelay(5*time.Second),
			grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
		)
		if err != nil {
			log.Fatalf("failed to dial: %v", err)
		}

		ctx, _ := context.WithCancel(context.Background())

		// Khởi chạy goroutine để theo dõi trạng thái kết nối
		go _middleware.MonitorConnection(ctx, name, conn)

		client := userpb.NewProfileServiceClient(conn)
		userClient.client = client
	}()
	return userClient
	// conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", config.USER_SERVICE_HOST, config.USER_SERVICE_PORT), grpc.WithInsecure())
	// if err != nil {
	// 	panic(fmt.Sprintf("Failed to connect to chat service: %v", err))
	// }
	// return &userClient{
	// 	client: userpb.NewProfileServiceClient(conn),
	// }
}

func (c *UserClient) GetProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error) {
	return c.client.GetProfileByIds(ctx, in)
}

func (c *UserClient) GetMapProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (map[uint64]*sharepb.ProfileItem, error) {
	profiles, err := c.client.GetProfileByIds(ctx, in)
	if err != nil {
		return nil, err
	}
	profileMap := make(map[uint64]*sharepb.ProfileItem)
	for _, profile := range profiles.Profiles {
		profileMap[profile.Id] = profile
	}
	return profileMap, nil
}

func (c *UserClient) MapMemberToDealPb(ctx context.Context, deals []*organizationpb.GroupDeal) {
	profileIds := make([]uint64, 0)
	for _, deal := range deals {
		if deal.Members != nil {
			for _, member := range deal.Members {
				profileIds = append(profileIds, uint64(member.Id))
			}
		}
		if deal.Customers != nil {
			for _, customer := range deal.Customers {
				profileIds = append(profileIds, uint64(customer.Id))
			}
		}
		if deal.Partners != nil {
			for _, partner := range deal.Partners {
				profileIds = append(profileIds, uint64(partner.Id))
			}
		}
	}

	profiles, err := c.client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: profileIds})
	if err != nil {
		return
	}
	profileMap := make(map[uint64]*sharepb.ProfileItem)
	for _, profile := range profiles.Profiles {
		profileMap[profile.Id] = profile
	}

	for _, deal := range deals {
		if deal.Members != nil {
			members := make([]*sharepb.ProfileItem, len(deal.Members))
			for i, member := range deal.Members {
				members[i] = profileMap[uint64(member.Id)]
			}
			deal.Members = members
		}
		if deal.Customers != nil {
			customers := make([]*sharepb.ProfileItem, len(deal.Customers))
			for i, customer := range deal.Customers {
				customers[i] = profileMap[uint64(customer.Id)]
			}
			deal.Customers = customers
		}
		if deal.Partners != nil {
			partners := make([]*sharepb.ProfileItem, len(deal.Partners))
			for i, partner := range deal.Partners {
				partners[i] = profileMap[uint64(partner.Id)]
			}
			deal.Partners = partners
		}
	}

}
