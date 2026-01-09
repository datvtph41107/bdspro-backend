package client

import (
	"common"
	_middleware "common/middleware"
	"context"
	"crm/env"
	"fmt"
	"log"
	bdspropb "pb/types/bdspro"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
	"time"

	"google.golang.org/grpc"
)

// @bind: crm/internal/interface/provider.BdsproProvider
type BdsproClient struct {
	Client         bdspropb.ProductServiceClient
	InternalClient bdspropb.BdsproInternalServiceClient
}

func NewBdsproClient() *BdsproClient {
	bdsproClient := &BdsproClient{Client: nil, InternalClient: nil}
	go func() {
		name := "bdspro"
		url := common.GetPrefixProtobufUrl(env.ENV_RUNTIME, name)

		conn, err := grpc.Dial(fmt.Sprintf("%s:8202", url),
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

		client := bdspropb.NewProductServiceClient(conn)
		bdsproClient.Client = client
		internalClient := bdspropb.NewBdsproInternalServiceClient(conn)
		bdsproClient.InternalClient = internalClient
	}()
	return bdsproClient
}

func (c *BdsproClient) GetProductByIds(ctx context.Context, ids []uint64) ([]*bdspropb.ProductResponse, error) {
	request := &bdspropb.IdRequest{
		Ids: ids,
	}

	response, err := c.Client.GetDetailByIds(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

func (c *BdsproClient) MapProductPbByIds(ctx context.Context, lead *crmpb.LeadDetailDTO) {
	products, err := c.GetProductByIds(ctx, lead.ProductIds)
	if err != nil {
		return
	}

	lead.Products = products
}

func (c *BdsproClient) GetProductAttachmentByIds(ctx context.Context, ids []uint64) (*bdspropb.ProductAttachmentResponse, error) {
	request := &bdspropb.IdRequest{
		Ids: ids,
	}
	response, err := c.Client.GetProductAttachmentByIds(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *BdsproClient) GetDealById(ctx context.Context, dealId uint64) (*bdspropb.GroupDeal, error) {
	request := &sharepb.IdRequest{
		Id: dealId,
	}
	response, err := c.InternalClient.GetDealById(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *BdsproClient) GetDealsByIds(ctx context.Context, dealIds []uint64) ([]*bdspropb.GroupDeal, error) {
	request := &sharepb.IdRequest{
		Ids: dealIds,
	}
	response, err := c.InternalClient.GetDealsByIds(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

func (m *BdsproClient) MapDealMemberToContacts(ctx context.Context, dealId uint64, contacts []*sharepb.ContactDTO) error {
	profileIDSet := make(map[uint64]struct{})
	for _, contact := range contacts {
		if contact.ProfileId != nil {
			profileIDSet[*contact.ProfileId] = struct{}{}
		}
	}

	userIds := make([]uint64, 0, len(profileIDSet))
	for userId := range profileIDSet {
		userIds = append(userIds, userId)
	}

	response, err := m.InternalClient.GetDealMembers(ctx, &bdspropb.DealMemberRequest{
		DealId:  dealId,
		UserIds: userIds,
	})
	dealMembers := response.Data
	dealMemberMap := make(map[uint64]*sharepb.DealMember)
	for _, dealMember := range dealMembers {
		dealMemberMap[dealMember.MemberId] = dealMember
	}

	if err != nil {
		return err
	}

	for _, contact := range contacts {
		if contact.ProfileId == nil {
			continue
		}
		if dealMember, ok := dealMemberMap[*contact.ProfileId]; ok {
			contact.DealMember = dealMember
		} else {
			contact.DealMember = nil
		}
	}

	return nil
}

func (c *BdsproClient) GetDealContractsByIds(ctx context.Context, dealContractIds []uint64) ([]*bdspropb.DealContractItem, error) {
	request := &sharepb.IdRequest{
		Ids: dealContractIds,
	}
	response, err := c.InternalClient.GetDealContractsByIds(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

func (c *BdsproClient) UpdateScheduleCount(ctx context.Context, productID uint64, count uint32) error {
	_, err := c.InternalClient.UpdateScheduleCount(ctx, &bdspropb.UpdateScheduleCountRequest{
		ProductId: productID,
		Count:     count,
	})
	return err
}

func (c *BdsproClient) UpdateContactCount(ctx context.Context, productID uint64, count uint32) error {
	_, err := c.InternalClient.UpdateContactCount(ctx, &bdspropb.UpdateContactCountRequest{
		ProductId: productID,
		Count:     count,
	})
	return err
}
