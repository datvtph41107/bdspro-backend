package u_product_provider

// import (
// 	"bdspro/internal/enums"
// 	"bdspro/v1/mod/global/member_plan"
// 	"bdspro/v1/mod/global/transmit"
// 	u_post "bdspro/v1/mod/user/post"
// 	u_post_dto "bdspro/v1/mod/user/post/dto"
// 	"bdspro/v1/models"
// 	_db "common/db"
// 	_jwt "common/jwt"
// 	_routes "common/routes"
// 	pb_noti "pb/types/notification"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"google.golang.org/protobuf/types/known/timestamppb"
// )

// type UserProductProvider struct {
// 	Repo        *u_post.PostRepo
// 	planService *member_plan.MemberPlanService
// 	transmit    *transmit.TransmitService
// }

// func (s *UserProductProvider) CreatePost(c *gin.Context, dto *u_post_dto.PostSaveRequest) (*models.Post, error) {
// 	plan := s.planService.CurrentPlan(c)
// 	if plan == nil {
// 		return nil, &_routes.Except{
// 			Code:    500,
// 			Message: "Vui lòng thử lại sau",
// 		}
// 	}

// 	// _, err := s.productService.RequiredOwner(c, dto.ProductID)
// 	_, err := s.productProvider.RequiredOwnerProduct(c, dto.ProductID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	principal := _jwt.GetPrincipal(c)
// 	// profileId := _jwt.GetProfileId(c)
// 	if plan.LimitPost <= s.Repo.CountPostFromDate(c, time.Now().Add(-30)) {
// 		return nil, &_routes.Except{
// 			Code:    400,
// 			Message: "Số lượng tạo tin đã đạt giới hạn. Vui lòng nâng cấp gói",
// 			// Bạn đã hết hạn mức tin
// 		}
// 	}
// 	// plan := s.planService.Plans
// 	if existed := s.Repo.ExistedPostWorking(c, dto.ProductID, principal.ProfileId); existed {
// 		return nil, &_routes.Except{
// 			Code:    400,
// 			Message: "Tin đăng vẫn còn hạn",
// 		}
// 	}

// 	numDate := dto.NumDate
// 	if numDate <= 0 {
// 		numDate = 15
// 	}
// 	if numDate > 30 {
// 		numDate = 30
// 	}
// 	entity := &models.Post{
// 		ProductId:       dto.ProductID,
// 		Title:           dto.Title,
// 		Content:         dto.Content,
// 		Price:           dto.Price,
// 		TransactionType: dto.TransactionType,
// 		PackageVisible:  dto.PackageVisible,
// 		Visibility:      dto.Visibility,
// 		NumDate:         numDate,
// 		ExpiredAt:       dto.ExpiredAt,
// 		Status:          enums.PostPublic,
// 		Hidden:          dto.Hidden,
// 	}
// 	if entity.ExpiredAt == nil {
// 		expired := time.Now().Add(time.Duration(numDate) * 24 * time.Hour)
// 		entity.ExpiredAt = &expired
// 	}

// 	err = _db.SaveWithAudit(c, entity)
// 	visibleAt := entity.ExpiredAt.Add(-3 * 24 * time.Hour)
// 	s.transmit.NotiClient.Send(c, &pb_noti.NotiNewRequest{
// 		Title:     "Hết hạn tin đăng",
// 		Message:   "Tin đăng của bạn sắp hết hạn. Vui lòng kiểm tra lại",
// 		Type:      3,
// 		UserID:    principal.ProfileId,
// 		VisibleAt: timestamppb.New(visibleAt),
// 	})

// 	return entity, err
// }
