package provider

import (
	"context"
	"hub/infra/client"
	providers "hub/internal/provider"
)

type UpdateDataProviderImpl struct {
	BdsproClient *client.BDSProClient
}

func NewUpdateDataProviderImpl(BdsproClient *client.BDSProClient) providers.UpdateDataProvider {
	return &UpdateDataProviderImpl{
		BdsproClient: BdsproClient,
	}
}

func (p *UpdateDataProviderImpl) GetRefreshIds(c context.Context, resource string, ownerId uint64) ([]uint64, error) {
	switch resource {
	case "product", "property":
		return p.BdsproClient.GetResourceIds(c, resource, ownerId)
	default:
		break
	}
	return nil, nil
}

func (p *UpdateDataProviderImpl) GetUpdatedAtOfId(c context.Context, resource string, id uint64) (int64, error) {
	switch resource {
	case "product", "property":
		return p.BdsproClient.GetUpdatedAtOfId(c, resource, id)
	default:
		break
	}

	return 0, nil
}
