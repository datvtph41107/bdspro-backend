package _provider

import (
	_dto "common/domain/dto"
	"context"
	sharepb "pb/types/shared"
)

type CrmProvider interface {
	GetRelationShip(ctx context.Context, currentId uint64, targetId uint64) (*_dto.FriendItemDTO, error)
	GetRelationShipV3(ctx context.Context, currentId uint64, targetId uint64) (*_dto.FriendV3DTO, error)
	GetUserRateStats(ctx context.Context, userId uint64) (*_dto.RateStats, error)
	GetUserRateStatsV3(ctx context.Context, userId uint64) (*_dto.RateStateV3DTO, error)
	GetOriginID(ctx context.Context, profileID uint64) (uint64, error)
	GetFriendInfo(ctx context.Context, profileID uint64) (*_dto.ContactInfoV3DTO, error)
	GetContactByOriginId(ctx context.Context, originId uint64) (*sharepb.ContactDTO, error)
	GetOriginProfileByOriginIds(ctx context.Context, originIds []uint64) (map[uint64]*sharepb.OriginProfile, error)
	HasContactRelation(ctx context.Context, ownerID uint64, contactOriginProfileID uint64) (bool, error)
}
