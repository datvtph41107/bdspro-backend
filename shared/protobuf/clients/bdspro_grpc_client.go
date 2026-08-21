package clients

import (
	_dto "common/domain/dto"
	"context"
	"errors"
	"fmt"
	"log"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"

	"google.golang.org/grpc"
)

type BdsproGrpcClient struct {
	Client                bdspropb.ProductServiceClient
	InternalClient        bdspropb.BdsproInternalServiceClient
	PostClient            bdspropb.PostServiceClient
	AdminProjectClient    bdspropb.AdminProjectServiceClient
	AdminRegionClient     bdspropb.AdminRegionServiceClient
	AdminAreaRegionClient bdspropb.AdminAreaRegionServiceClient
}

func BindBdsproGrpcClient(conn grpc.ClientConnInterface) *BdsproGrpcClient {
	return &BdsproGrpcClient{
		Client:                bdspropb.NewProductServiceClient(conn),
		InternalClient:        bdspropb.NewBdsproInternalServiceClient(conn),
		PostClient:            bdspropb.NewPostServiceClient(conn),
		AdminProjectClient:    bdspropb.NewAdminProjectServiceClient(conn),
		AdminRegionClient:     bdspropb.NewAdminRegionServiceClient(conn),
		AdminAreaRegionClient: bdspropb.NewAdminAreaRegionServiceClient(conn),
	}
}

func (c *BdsproGrpcClient) GetProductByIds(ctx context.Context, ids []uint64) (*bdspropb.ProductInternalResponse, error) {
	if c.Client == nil {
		return nil, errors.New("bdspro client not avaiable")
	}

	response, err := c.Client.GetDetailByIds(ctx, &bdspropb.IdRequest{Ids: ids})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *BdsproGrpcClient) GetProductAttachmentByIds(ctx context.Context, ids []uint64) (*bdspropb.ProductAttachmentResponse, error) {
	if c.Client == nil {
		return nil, errors.New("bdspro client not avaiable")
	}

	response, err := c.Client.GetProductAttachmentByIds(ctx, &bdspropb.IdRequest{Ids: ids})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *BdsproGrpcClient) GetCountByOwner(ctx context.Context, req *bdspropb.GetCountByOwnerRequest) (*bdspropb.GetCountByOwnerResponse, error) {
	if c.InternalClient == nil {
		return nil, errors.New("bdspro internal client not available")
	}

	response, err := c.InternalClient.GetCountByOwner(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetDashboardByProfileId lấy dashboard stats theo profileId
func (c *BdsproGrpcClient) GetDashboardByProfileId(ctx context.Context, profileId uint64) (*sharepb.BDSProDashboardProto, error) {
	if c.InternalClient == nil {
		return nil, errors.New("bdspro internal client not available")
	}

	resp, err := c.InternalClient.GetDashboardByProfileId(ctx, &sharepb.IdRequest{Id: profileId})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *BdsproGrpcClient) GetSuggest(ctx context.Context, content string) (*_dto.ProductV3DTO, error) {
	if c.InternalClient == nil {
		return nil, errors.New("bdspro internal client not available")
	}

	fmt.Println("content: " + content)

	response, err := c.InternalClient.GetSuggest(ctx, &sharepb.RequestV3Proto{Text: content})
	if err != nil {
		return nil, err
	}

	fmt.Println("response: ....")

	result := &_dto.ProductV3DTO{
		ID:          response.Id,
		Name:        response.Name,
		Description: response.Description,
	}

	// if response.Address != nil {
	// 	result.Address = &_dto.AddressV3DTO{
	// 		ProvinceID:   response.Address.ProvinceId,
	// 		ProvinceName: response.Address.ProvinceName,
	// 		DistrictID:   response.Address.DistrictId,
	// 		DistrictName: response.Address.DistrictName,
	// 		WardID:       response.Address.WardId,
	// 		WardName:     response.Address.WardName,
	// 	}
	// }

	if response.PropertyType != nil {
		result.PropertyType = &_dto.ItemDTO{
			ID:   response.PropertyType.Id,
			Name: response.PropertyType.Name,
		}
	}

	if response.DocType != nil {
		result.DocType = &_dto.ItemDTO{
			ID:   response.DocType.Id,
			Name: response.DocType.Name,
		}
	}

	if response.Amenities != nil {
		result.Amenities = make([]_dto.ItemDTO, len(response.Amenities))
		for i, amenity := range response.Amenities {
			result.Amenities[i] = _dto.ItemDTO{
				ID:   amenity.Id,
				Name: amenity.Name,
			}
		}
	}

	return result, nil
}

func (c *BdsproGrpcClient) GetDealMember(ctx context.Context, dealId uint64, userId uint64) (*sharepb.DealMember, error) {
	if c.InternalClient == nil {
		return nil, errors.New("bdspro internal client not available")
	}

	resp, err := c.InternalClient.GetDealMembers(ctx, &bdspropb.DealMemberRequest{
		DealId:  dealId,
		UserIds: []uint64{userId},
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, nil
	}

	return resp.Data[0], nil
}

func (c *BdsproGrpcClient) GetResourceIds(ctx context.Context, resource string, ids uint64) ([]uint64, error) {
	resp, err := c.InternalClient.GetResourceIds(ctx, &bdspropb.ResourceRequest{
		Resource: resource,
	})
	if err != nil {
		log.Printf("err: %s", err)
		return nil, err
	}
	return resp.Ids, err
}

func (c *BdsproGrpcClient) GetUpdatedAtOfId(ctx context.Context, resource string, id uint64) (int64, error) {
	resp, err := c.InternalClient.GetUpdatedAtOfId(ctx, &bdspropb.ResourceRequest{
		Resource: resource,
		Id:       id,
	})
	if err != nil {
		log.Printf("err: %s", err)
		return 0, err
	}
	return resp.Timestamp, err
}
