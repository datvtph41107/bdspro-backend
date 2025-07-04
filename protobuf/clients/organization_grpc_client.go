package clients

import (
	_middleware "common/middleware"
	"context"
	"errors"
	"fmt"
	"log"
	organizationpb "pb/types/organization"
	"time"

	"google.golang.org/grpc"
)

type OrganizationClient struct {
	client organizationpb.InternalOrganizationServiceClient
}

func MakeOrganizationClient(env_runtime string) *OrganizationClient {
	organizationClient := &OrganizationClient{}
	go func() {
		url := "localhost"
		name := "organization"
		if env_runtime == "release" {
			url = name
		}

		conn, err := grpc.Dial(fmt.Sprintf("%s:8207", url),
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

		client := organizationpb.NewInternalOrganizationServiceClient(conn)
		organizationClient.client = client
	}()
	return organizationClient
}

func (c *OrganizationClient) GetOrganizationMapByIDs(ctx context.Context, organizationIDSet map[uint64]struct{}) (map[uint64]*organizationpb.OrganizationItem, error) {
	if c.client == nil {
		return nil, errors.New("organization client not avaiable")
	}

	organizationIDs := make([]uint64, 0, len(organizationIDSet))
	for id := range organizationIDSet {
		organizationIDs = append(organizationIDs, id)
	}

	organizations, err := c.client.GetOrganizationByIds(ctx, &organizationpb.GetOrganizationByIdsRequest{Ids: organizationIDs})
	if err != nil {
		return nil, err
	}

	organizationMap := make(map[uint64]*organizationpb.OrganizationItem, len(organizations.Organizations))
	for _, organization := range organizations.Organizations {
		organizationMap[organization.Id] = organization
	}

	return organizationMap, nil
}

func (c *OrganizationClient) GetGroupMapByIDs(ctx context.Context, groupIDSet map[uint64]struct{}) (map[uint64]*organizationpb.GroupItem, error) {
	if c.client == nil {
		return nil, errors.New("organization client not avaiable")
	}

	groupIDs := make([]uint64, 0, len(groupIDSet))
	for id := range groupIDSet {
		groupIDs = append(groupIDs, id)
	}

	groups, err := c.client.GetGroupByIds(ctx, &organizationpb.GetGroupByIdsRequest{Ids: groupIDs})
	if err != nil {
		return nil, err
	}

	groupMap := make(map[uint64]*organizationpb.GroupItem, len(groups.Groups))
	for _, group := range groups.Groups {
		groupMap[group.Id] = group
	}

	return groupMap, nil
}
