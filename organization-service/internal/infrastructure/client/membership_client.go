package client

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type MembershipClient struct {
	Plans map[uint64]*entity.PlanEntity
}

func NewMembershipClient() *MembershipClient {
	return &MembershipClient{
		Plans: make(map[uint64]*entity.PlanEntity),
	}
}

func (c *MembershipClient) GetPlan(ctx context.Context, id uint64) (*entity.PlanEntity, error) {

	// return c.Plans[id], nil
	// todo: cái này e sẽ tự implement giờ cứ fix cứng trong code

	return &entity.PlanEntity{
		ID:                id,
		Name:              "Free",
		MemberBranchLimit: 10,
	}, nil
}
