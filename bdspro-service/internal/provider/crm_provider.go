package provider

import (
	"context"
	sharepb "pb/types/shared"
)

type CrmProvider interface {
	GetContactById(ctx context.Context, id uint64) (*sharepb.ContactDTO, error)
	GetOriginProfileByOriginIds(ctx context.Context, originIDs []uint64) ([]*sharepb.OriginProfile, error)
}
