package services

import (
	_jwt "common/jwt"
	"membership/models"
	"membership/repo"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MembershipService struct {
	repo     *repo.MembershipRepo
	planRepo *repo.PlanRepo
}

func NewMembershipService(repo *repo.MembershipRepo, planRepo *repo.PlanRepo) *MembershipService {
	return &MembershipService{
		repo:     repo,
		planRepo: planRepo,
	}
}

func (m *MembershipService) ListPlan(c *gin.Context) ([]models.Plan, error) {
	return m.planRepo.ListForUser()
}

func (m *MembershipService) CreateOrder(c *gin.Context) (*models.Membership, error) {
	profileId := _jwt.GetProfileId(c)
	membershipId, _ := strconv.ParseUint(c.Param("membershipId"), 10, 32)

	membership := &models.Membership{
		ProfileId:    profileId,
		MembershipId: membershipId,
	}
	m.repo.Create(c, membership)
	return membership, nil
}
