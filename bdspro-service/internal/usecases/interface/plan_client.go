package iusecase

import (
	"bdspro/internal/domain"
	"context"
)

type IPlanClient interface {
	CurrentPlan(c context.Context) *domain.MemberPlan
}
