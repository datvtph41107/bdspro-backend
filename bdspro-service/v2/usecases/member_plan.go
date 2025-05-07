package usecases

import (
	"bdspro/v2/internal/domain"
	"context"
)

type PlanUsecase interface {
	CurrentPlan(c context.Context) *domain.MemberPlan
}
