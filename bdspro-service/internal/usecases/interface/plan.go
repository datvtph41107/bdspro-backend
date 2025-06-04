package iusecase

import (
	"bdspro/internal/domain"
	"context"
)

type IPlanUC interface {
	CurrentPlan(c context.Context) *domain.MemberPlan
}
