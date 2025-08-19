package provider

import (
	"auth/internal/dto"
	"context"
)

type OrganizationProvider interface {
	SwitchOrganization(ctx context.Context, organizationID uint64) (*dto.SwitchOrganizationResponse, error)
}
