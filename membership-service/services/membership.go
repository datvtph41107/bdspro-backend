package services

import (
	_db "common/db"
	_jwt "common/jwt"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"
	"time"

	"membership/data"
	"membership/models"
	"membership/repo"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type MembershipService struct {
	repo      *repo.MembershipRepo
	planRepo  *repo.PlanRepo
	payClient *PaymentClient
	Key       string
}

func NewMembershipService(
	repo *repo.MembershipRepo,
	planRepo *repo.PlanRepo,
	payClient *PaymentClient,
) *MembershipService {
	return &MembershipService{
		repo:      repo,
		planRepo:  planRepo,
		payClient: payClient,
		Key:       viper.GetString("payment.secretKey"),
	}
}

func (m *MembershipService) ListPlan(c *gin.Context) ([]models.Plan, error) {
	return m.planRepo.ListForUser()
}

func (m *MembershipService) SetupMembership(c *gin.Context) (*models.Membership, error) {
	return nil, nil
}

// ra hạn thêm
func (m *MembershipService) CreateOrder(c *gin.Context) (*models.Membership, error) {
	var memDTO models.Membership
	if err := c.ShouldBindJSON(&memDTO); err != nil {
		return nil, err
	}

	if memDTO.PlanId == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Vui lòng chọn gói",
		}
	}

	profileId := _jwt.GetProfileId(c)
	plan, _ := m.planRepo.GetByID(c, memDTO.PlanId)
	if plan == nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Gói không đúng",
		}
	}

	membership, _ := m.repo.GetByProfileID(profileId)

	if membership == nil {
		membership = &models.Membership{
			ProfileId: profileId,
			PlanId:    memDTO.PlanId,
			// OrderId:   mem.OrderId,
		}
		_db.SaveWithAudit(c, membership)
	} else if membership.PlanId == memDTO.PlanId &&
		membership.FinishedDate.After(time.Now()) {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Gói vẫn còn hạn. Nếu muốn bạn có thể nâng cấp gói",
		}
	}

	payload := data.OrderRequest{
		Amount:      int(plan.Price),
		Description: "UP_PLAN",
		TargetId:    &membership.ID,
		Type:        2,
	}
	m.payClient.RequestOrder(c, payload)
	return membership, nil
}

func (m *MembershipService) VerifyOrder(c *gin.Context) (*models.Membership, error) {
	var mem data.NotiComplete
	if err := c.ShouldBindJSON(&mem); err != nil {
		return nil, err
	}

	plantText := strconv.FormatUint(mem.TargetId, 10) + ":" + strconv.FormatInt(int64(mem.Type), 10)
	signature := _utils.GenerateHmacSHA256(plantText, m.Key)

	if signature != mem.Signature {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Giao dịch không chính xác",
		}
	}

	membership, err := m.repo.GetByID(c, mem.TargetId)
	if err != nil {
		return nil, err
	}

	plan, _ := m.planRepo.GetByID(c, membership.PlanId)

	membership.PaidAt = mem.PaidAt
	expiredAt := mem.PaidAt.Add(time.Duration(plan.Duration * int64(time.Hour)))
	membership.FinishedDate = &expiredAt
	m.repo.Update(c, mem.TargetId, membership)
	return membership, nil
}
