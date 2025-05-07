package u_post

import (
	"bdspro/v1/mod/global/member_plan"
	"bdspro/v1/mod/global/transmit"
	u_product "bdspro/v1/mod/user/product"
	"bdspro/v1/models"
	"common/case/crud2"
	_db "common/db"
	_jwt "common/jwt"
	_routes "common/routes"
	pb_noti "pb/types/notification"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostService struct {
	crud2.BaseService[models.Post, *PostRepo]
	planService    *member_plan.MemberPlanService
	productService *u_product.ProductService
	transmit       *transmit.TransmitService
}

func NewPostService(repo *PostRepo,
	planService *member_plan.MemberPlanService,
	productService *u_product.ProductService,
	transmit *transmit.TransmitService,
) *PostService {
	return &PostService{
		BaseService: crud2.BaseService[models.Post, *PostRepo]{
			Repo: repo,
		},
		planService:    planService,
		productService: productService,
		transmit:       transmit,
	}
}

func (s *PostService) CreatePost(c *gin.Context, dto PostSaveRequest) (*models.Post, error) {
	plan := s.planService.CurrentPlan(c)
	if plan == nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: "Vui lòng thử lại sau",
		}
	}

	_, err := s.productService.RequiredOwner(c, dto.ProductID)
	if err != nil {
		return nil, err
	}

	profileId := _jwt.GetProfileId(c)
	if plan.LimitPost <= s.Repo.CountPost(c, profileId) {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Số lượng tạo tin đã đạt giới hạn. Vui lòng nâng cấp gói",
		}
	}
	// plan := s.planService.Plans
	if existed := s.Repo.ExistedPostWorking(c, dto.ProductID, profileId); existed {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Tin đăng vẫn còn hạn",
		}
	}
	entity := &models.Post{
		ProductId:       dto.ProductID,
		Content:         dto.Content,
		TransactionType: dto.TransactionType,
		Title:           dto.Title,
		ExpiredAt:       dto.ExpiredAt,
	}
	if entity.ExpiredAt == nil {
		expired := time.Now().Add(15 * 24 * time.Hour)
		entity.ExpiredAt = &expired
	}

	err = _db.SaveWithAudit(c, entity)
	visibleAt := entity.ExpiredAt.Add(-3 * 24 * time.Hour)
	s.transmit.NotiClient.Send(c, &pb_noti.NotiNewRequest{
		Title:     "Hết hạn tin đăng",
		Message:   "Tin đăng của bạn sắp hết hạn. Vui lòng kiểm tra lại",
		Type:      3,
		UserID:    profileId,
		VisibleAt: timestamppb.New(visibleAt),
	})

	return entity, err
}

func (s *PostService) UpdateHidden(c *gin.Context, dto UpdateHiddenRequest) (*models.Post, error) {
	entity, err := s.RequiredOwner(c, dto.PostID)
	if err != nil {
		return nil, err
	}

	entity.Hidden = dto.Hidden

	err = _db.SaveWithAudit(c, entity)
	return entity, err
}

func (s *PostService) Action(c *gin.Context, dto ActionRequest) (*models.Post, error) {
	entity, err := s.RequiredOwner(c, dto.PostID)
	if err != nil {
		return nil, err
	}

	if dto.Action == 1 {
		entity.NumView = entity.NumView + 1
		err = _db.SaveWithAudit(c, entity)
	}

	return nil, err
}

func (s *PostService) UpdatePost(c *gin.Context, id uint64, dto PostUpdateRequest) (*models.Post, error) {
	entity, err := s.RequiredOwner(c, id)
	if err != nil {
		return nil, err
	}

	copier.Copy(&entity, &dto)

	err = _db.SaveWithAudit(c, entity)
	return entity, err
}

func (s *PostService) RequiredOwner(c *gin.Context, id uint64) (*models.Post, error) {
	post, err := s.Repo.GetByID(c, id)
	if err != nil {
		return nil, err
	}

	profileId := _jwt.GetProfileId(c)

	if *post.CreatedBy != profileId {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Bạn không có quyền truy cập",
		}
	}
	return post, nil
}

func (s *PostService) MakeNewExpired(c *gin.Context, dto PostExpiredRequest) (*models.Post, error) {
	// plan := s.planService.CurrentPlan(c)
	// if plan == nil {
	// 	return nil, &_routes.Except{
	// 		Code:    500,
	// 		Message: "Vui lòng thử lại sau",
	// 	}
	// }

	entity, err := s.RequiredOwner(c, dto.PostID)
	if err != nil {
		return nil, err
	}

	profileId := _jwt.GetProfileId(c)
	// // plan := s.planService.Plans
	// if existed := s.Repo.ExistedPostWorking(c, dto.ProductID, profileId); existed {
	// 	return nil, &_routes.Except{
	// 		Code:    400,
	// 		Message: "Tin đăng vẫn còn hạn",
	// 	}
	// }
	// entity := &models.Post{
	// 	ProductId:       dto.ProductID,
	// 	Content:         dto.Content,
	// 	TransactionType: dto.TransactionType,
	// 	Title:           dto.Title,
	// 	ExpiredAt:       dto.ExpiredAt,
	// }
	newExpired := entity.ExpiredAt.Add(time.Duration(time.Duration(dto.NumDay) * 24 * time.Hour))
	entity.ExpiredAt = &newExpired

	err = _db.SaveWithAudit(c, entity)
	visibleAt := entity.ExpiredAt.Add(-3 * 24 * time.Hour)
	s.transmit.NotiClient.Send(c, &pb_noti.NotiNewRequest{
		Title:     "Hết hạn tin đăng",
		Message:   "Tin đăng của bạn sắp hết hạn. Vui lòng kiểm tra lại",
		Type:      3,
		UserID:    profileId,
		VisibleAt: timestamppb.New(visibleAt),
	})

	return entity, err
}

func (s *PostService) Detail(c *gin.Context, id uint64) (*PostDetail, error) {
	post, err := s.Repo.GetByID(c, id)
	if err != nil {
		return nil, err
	}
	product, _ := s.productService.ProductRepo.GetByID(c, post.ProductId)
	if err != nil {
		return nil, err
	}

	return &PostDetail{
		Post:        post,
		ProductInfo: product,
	}, nil
}

func (s *PostService) PersonalPost(c *gin.Context, dto SearchRequest) (*_routes.ResponseDTO, error) {
	profileId := _jwt.GetProfileId(c)
	entities, total, err := s.Repo.Search(c, profileId, dto)
	if err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}

	return &_routes.ResponseDTO{
		Code:          0,
		Data:          entities,
		TotalElements: &total,
	}, nil
}

func (s *PostService) GlobalPost(c *gin.Context, dto SearchRequest) (*_routes.ResponseDTO, error) {
	entities, total, err := s.Repo.Global(c, dto)
	if err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}

	return &_routes.ResponseDTO{
		Code:          0,
		Data:          entities,
		TotalElements: &total,
	}, nil
}
