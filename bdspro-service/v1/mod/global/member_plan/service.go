package member_plan

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MemberPlanService struct {
	Plans map[uint64]*MemberPlan
}

func NewMemberPlanService() *MemberPlanService {
	planService := &MemberPlanService{}
	planService.loadPlans()
	return planService
}

func (s *MemberPlanService) CurrentPlan(c *gin.Context) *MemberPlan {
	var id *uint64
	if val, ok := c.Get("planId"); ok && val != nil {
		id, _ = val.(*uint64)
	}
	return s.Plans[*id]
}

func (s *MemberPlanService) loadPlans() {
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
