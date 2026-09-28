package client

import (
	"context"
	"errors"
	"pb/clients"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
	"userv4/internal/dto"
)

// @bind: user/internal/interface/providers.BdsproProvider
type BdsproClient struct {
	*clients.BdsproGrpcClient
}

func NewBdsproClient(
	rpcClient *clients.BdsproGrpcClient,
) *BdsproClient {
	return &BdsproClient{
		BdsproGrpcClient: rpcClient,
	}
}

func (c *BdsproClient) GetDealMember(ctx context.Context, dealId uint64, userId uint64) (*dto.DealMember, error) {
	if c.BdsproGrpcClient == nil {
		return nil, errors.New("bdspro grpc client not available")
	}

	resp, err := c.BdsproGrpcClient.GetDealMember(ctx, dealId, userId)
	if err != nil {
		return nil, err
	}
	return &dto.DealMember{
		DealId:  resp.DealId,
		UserId:  resp.MemberId,
		RoleId:  resp.RoleId,
		RoleKey: resp.RoleKey,
	}, nil
}

// GetDashboardByProfileId lấy dashboard stats theo profileId (dùng bởi profile handler).
func (c *BdsproClient) GetDashboardByProfileId(ctx context.Context, profileID uint64) (*sharepb.BDSProDashboardProto, error) {
	if c.BdsproGrpcClient == nil || c.BdsproGrpcClient.InternalClient == nil {
		return nil, nil
	}

	resp, err := c.BdsproGrpcClient.GetDashboardByProfileId(ctx, profileID)
	if err != nil {
		return nil, nil
	}
	return resp, nil
}

// CreateProduct gọi bdspro ProductNew để tạo sản phẩm.
func (c *BdsproClient) CreateProduct(ctx context.Context, product *bdspropb.ProductSaveRequest) (*bdspropb.Response, error) {
	if c.BdsproGrpcClient == nil || c.Client == nil {
		return nil, errors.New("bdspro grpc client not available")
	}
	return c.Client.ProductNew(ctx, product)
}

// DeleteProduct gọi bdspro DeleteProduct để xoá sản phẩm.
func (c *BdsproClient) DeleteProduct(ctx context.Context, productID uint64) error {
	if c.BdsproGrpcClient == nil || c.Client == nil {
		return errors.New("bdspro grpc client not available")
	}
	_, err := c.Client.DeleteProduct(ctx, &bdspropb.IdRequest{Id: productID})
	return err
}

// ListMyProducts lấy danh sách sản phẩm của profile hiện tại (ProductMe).
func (c *BdsproClient) ListMyProducts(ctx context.Context, page, size uint32) (*bdspropb.SearchResponse, error) {
	if c.BdsproGrpcClient == nil || c.Client == nil {
		return nil, errors.New("bdspro grpc client not available")
	}
	return c.Client.ProductMe(ctx, &sharepb.SyncRequest{
		Page:      page,
		Size:      size,
		Timestamp: 0,
	})
}

// CreateGroupDeal tạo thương vụ cá nhân (ownerType=member).
func (c *BdsproClient) CreateGroupDeal(ctx context.Context, req *bdspropb.CreateGroupDealRequest) (uint64, error) {
	if c.BdsproGrpcClient == nil || c.DealClient == nil {
		return 0, errors.New("bdspro deal client not available")
	}
	resp, err := c.DealClient.CreateGroupDeal(ctx, req)
	if err != nil {
		return 0, err
	}
	if resp == nil {
		return 0, nil
	}
	return uint64(resp.Id), nil
}

// DeleteGroupDeal soft-delete thương vụ theo id.
func (c *BdsproClient) DeleteGroupDeal(ctx context.Context, dealID uint64) error {
	if c.BdsproGrpcClient == nil || c.DealClient == nil {
		return errors.New("bdspro deal client not available")
	}
	if dealID == 0 {
		return nil
	}
	_, err := c.DealClient.DeleteGroupDeal(ctx, &bdspropb.DeleteGroupDealRequest{Id: uint32(dealID)})
	return err
}
