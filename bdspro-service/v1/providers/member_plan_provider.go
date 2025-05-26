package providers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MemberPlan struct {
	ID                 uint64   `json:"id"`
	PackageName        string   `json:"packageName"`
	Duration           int64    `json:"duration"` // day
	Price              int64    `json:"price"`
	CreateProduct      int      `json:"createProduct"`
	LimitPost          int64    `json:"limitPost"`
	PublicProductLimit int      `json:"publicProductLimit"`
	CreateProductAI    int      `json:"createProductAI"`
	FriendLimit        int      `json:"friendLimit"`
	CustomerLimit      int      `json:"customerLimit"`
	JoinOrg            int      `json:"joinOrg"`
	CreateOrganization int      `json:"createOrganization"`
	OrgMemberLimit     int      `json:"orgMemberLimit"`
	OrgProductLimit    int      `json:"orgProductLimit"`
	Color              string   `json:"color"`
	Features           []string `json:"features"`
	Visibility         bool     `json:"visibility"`
}

type MemberPlanProvider struct {
	Plans map[uint64]*MemberPlan
}

func NewMemberPlanService() *MemberPlanProvider {
	planService := &MemberPlanProvider{}
	planService.loadPlans()
	return planService
}

func (s *MemberPlanProvider) CurrentPlan(c *gin.Context) *MemberPlan {
	var id *uint64
	if val, ok := c.Get("planId"); ok && val != nil {
		id, _ = val.(*uint64)
	}
	return s.Plans[*id]
}

func (s *MemberPlanProvider) CurrentPlan2(c context.Context) *MemberPlan {
	id := c.Value("planId").(uint64)
	return s.Plans[id]
}

func (s *MemberPlanProvider) loadPlans() {
	// var entities []MemberPlan
	url := "http://14.225.210.29:8000/v1/membership/plan/list"

	// Gọi GET request
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	result := []MemberPlan{}
	// v := resp.Body
	body, err := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return
	}

	// DB.Find(&entities)
	cache := make(map[uint64]*MemberPlan)
	for _, e := range result {
		cache[e.ID] = &e
	}
	s.Plans = cache
}
