package utils

import "context"

func GetOrganizationIdFromContext(c context.Context) (*uint64, error) {
	organizationId := c.Value("profileId")
	if organizationId == nil {
		return nil, nil
	}
	return organizationId.(*uint64), nil
}
