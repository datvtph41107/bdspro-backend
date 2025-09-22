package handler

import (
	"bdspro/internal/domain"
	_utils "common/utils"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// @bind: bdspro/internal/usecases.PlanUsecase
type MemberPlanUsecase struct {
	Plans map[uint64]*domain.MemberPlan
}

func NewMemberPlanUsecase() *MemberPlanUsecase {
	planService := &MemberPlanUsecase{}
	planService.loadPlans()
	return planService
}

func (s *MemberPlanUsecase) CurrentPlan(c context.Context) *domain.MemberPlan {
	planId := _utils.GetPlanIdFromContext(c)
	return s.Plans[planId]
}

func (s *MemberPlanUsecase) loadPlans() {
	// var entities []MemberPlan
	url := "http://14.225.210.29:8000/v1/membership/plan/list"

	// Gọi GET request
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	result := []domain.MemberPlan{}
	// v := resp.Body
	body, err := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return
	}

	// DB.Find(&entities)
	cache := make(map[uint64]*domain.MemberPlan)
	for _, e := range result {
		cache[e.ID] = &e
	}
	s.Plans = cache
}
