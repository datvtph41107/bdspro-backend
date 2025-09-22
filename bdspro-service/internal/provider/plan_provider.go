package provider

import (
	"bdspro/internal/domain"
	"context"
)

type PlanProvider interface {
	CurrentPlan(c context.Context) *domain.MemberPlan
}
