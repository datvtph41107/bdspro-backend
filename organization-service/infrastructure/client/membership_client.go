package client

import (
	"context"

	"organization/internal/domain/entity"
)

// @bind: organization/internal/usecase.IMembershipClient
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
