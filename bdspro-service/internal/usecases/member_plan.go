package usecases

import (
	"bdspro/internal/domain"
	"context"
)

type PlanUsecase interface {
	CurrentPlan(c context.Context) *domain.MemberPlan
}
