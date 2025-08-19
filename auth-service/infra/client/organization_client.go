package client

import (
	"auth/internal/dto"
	"context"
	"net/http"
)

// @bind: auth/internal/interface/provider.OrganizationProvider
type OrganizationClient struct {
	HttpClient *http.Client
}

func NewOrganizationClient(httpClient *http.Client) *OrganizationClient {
	return &OrganizationClient{HttpClient: httpClient}
}

func (c *OrganizationClient) SwitchOrganization(ctx context.Context, organizationID uint64) (*dto.SwitchOrganizationResponse, error) {
	return nil, nil
}
