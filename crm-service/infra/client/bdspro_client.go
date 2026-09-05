package client

import (
	"context"
	"errors"
	"pb/clients"
	bdspropb "pb/types/bdspro"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
)

// @bind: crm/internal/interface/provider.BdsproProvider
type BdsproClient struct {
	Client                bdspropb.ProductServiceClient
	InternalClient        bdspropb.BdsproInternalServiceClient
	AdminProjectClient    bdspropb.AdminProjectServiceClient
	AdminRegionClient     bdspropb.AdminRegionServiceClient
	AdminAreaRegionClient bdspropb.AdminAreaRegionServiceClient
}

func NewBdsproClient(rpcClient *clients.BdsproGrpcClient) *BdsproClient {
	return &BdsproClient{
		Client:                rpcClient.Client,
		InternalClient:        rpcClient.InternalClient,
		AdminProjectClient:    rpcClient.AdminProjectClient,
		AdminRegionClient:     rpcClient.AdminRegionClient,
		AdminAreaRegionClient: rpcClient.AdminAreaRegionClient,
	}
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

func (c *BdsproClient) SearchProjects(ctx context.Context, text string, page uint32, size uint32) ([]*bdspropb.Project, int64, error) {
	if c.AdminProjectClient == nil {
		return nil, 0, errors.New("bdspro admin project client not available")
	}
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 20
	}
	resp, err := c.AdminProjectClient.GetList(ctx, &bdspropb.SearchQueryRequest{
		Page: page,
		Size: size,
		Text: text,
	})
	if err != nil {
		return nil, 0, err
	}
	return resp.GetData(), resp.GetTotalElements(), nil
}

func (c *BdsproClient) SearchRegions(ctx context.Context, text string, page uint32, size uint32) ([]*bdspropb.Region, int64, error) {
	if c.AdminRegionClient == nil {
		return nil, 0, errors.New("bdspro admin region client not available")
	}
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 20
	}
	resp, err := c.AdminRegionClient.GetRegionList(ctx, &bdspropb.RegionListRequest{
		Page: int32(page),
		Size: int32(size),
		Text: &text,
	})
	if err != nil {
		return nil, 0, err
	}
	return resp.GetData(), resp.GetTotal(), nil
}

func (c *BdsproClient) GetRegionByID(ctx context.Context, id uint64) (*bdspropb.Region, error) {
	if c.AdminRegionClient == nil {
		return nil, errors.New("bdspro admin region client not available")
	}
	return c.AdminRegionClient.GetRegionDetail(ctx, &bdspropb.RegionDetailRequest{Id: id})
}

func (c *BdsproClient) ListAreaRegions(ctx context.Context) ([]*bdspropb.AreaRegion, error) {
	if c.AdminAreaRegionClient == nil {
		return nil, errors.New("bdspro admin area region client not available")
	}
	resp, err := c.AdminAreaRegionClient.ListAreaRegions(ctx, &bdspropb.ListAreaRegionsRequest{})
	if err != nil {
		return nil, err
	}
	return resp.GetData(), nil
}

func (c *BdsproClient) GetAreaRegionByID(ctx context.Context, id uint64) (*bdspropb.AreaRegion, error) {
	if c.AdminAreaRegionClient == nil {
		return nil, errors.New("bdspro admin area region client not available")
	}
	resp, err := c.AdminAreaRegionClient.GetAreaRegionByID(ctx, &bdspropb.GetAreaRegionByIDRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return resp.GetRegion(), nil
}

func (c *BdsproClient) GetProductUsersByDistributeId(ctx context.Context, distributeId uint64) ([]*sharepb.ProductUser, error) {
	response, err := c.InternalClient.GetProductUsersByDistributeId(ctx, &sharepb.IdRequest{
		Id: distributeId,
	})
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

func (c *BdsproClient) GetProductUsersByProductId(ctx context.Context, productId uint64) ([]*sharepb.ProductUser, error) {
	response, err := c.InternalClient.GetProductUsersByProductId(ctx, &sharepb.IdRequest{
		Id: productId,
	})
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}
