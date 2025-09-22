package providers

import (
	"bdspro/internal/domain"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// type MemberPlan struct {
// 	ID                 uint64   `json:"id"`
// 	PackageName        string   `json:"packageName"`
// 	Duration           int64    `json:"duration"` // day
// 	Price              int64    `json:"price"`
// 	CreateProduct      int      `json:"createProduct"`
// 	LimitPost          int64    `json:"limitPost"`
// 	PublicProductLimit int      `json:"publicProductLimit"`
// 	CreateProductAI    int      `json:"createProductAI"`
// 	FriendLimit        int      `json:"friendLimit"`
// 	CustomerLimit      int      `json:"customerLimit"`
// 	JoinOrg            int      `json:"joinOrg"`
// 	CreateOrganization int      `json:"createOrganization"`
// 	OrgMemberLimit     int      `json:"orgMemberLimit"`
// 	OrgProductLimit    int      `json:"orgProductLimit"`
// 	Color              string   `json:"color"`
// 	Features           []string `json:"features"`
// 	Visibility         bool     `json:"visibility"`
// }

// @bind: bdspro/internal/provider.PlanProvider
type MemberPlanProvider struct {
	Plans map[uint64]*domain.MemberPlan
}

func NewMemberPlanService() *MemberPlanProvider {
	planService := &MemberPlanProvider{}
	planService.loadPlans()
	return planService
}

func (s *MemberPlanProvider) CurrentPlan(c context.Context) *domain.MemberPlan {
	// var id *uint64
	// if val, ok := c.Get("planId"); ok && val != nil {
	// 	id, _ = val.(*uint64)
	// }
	// return s.Plans[*id]
	planId := c.Value("planId")
	id := planId.(*uint64)
	if s.Plans == nil {
		s.loadPlans()
		return nil
	}
	if id == nil {
		return nil
	}
	return s.Plans[*id]
}

// func (s *MemberPlanProvider) CurrentPlanFromContext(c context.Context) *domain.MemberPlan {
// 	planId := c.Value("planId")
// 	id := planId.(*uint64)
// 	if s.Plans == nil {
// 		s.loadPlans()
// 		return nil
// 	}
// 	if id == nil {
// 		return nil
// 	}
// 	return s.Plans[*id]
// }

func (s *MemberPlanProvider) loadPlans() {
	// var entities []MemberPlan
	url := "http://14.225.205.222:8000/v1/membership/user/plans"

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
