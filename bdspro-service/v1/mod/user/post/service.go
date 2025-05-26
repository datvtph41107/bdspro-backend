package u_post

import (
	"bdspro/internal/enums"
	user_dto "bdspro/v1/mod/user/dto/request"
	u_post_dto "bdspro/v1/mod/user/post/dto"
	u_product "bdspro/v1/mod/user/product"
	iu_product "bdspro/v1/mod/user/product/interface"
	"bdspro/v1/models"
	"bdspro/v1/providers"
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

type UserPostService struct {
	crud2.BaseService[models.Post, *PostRepo]
	planService *providers.MemberPlanProvider
	// UserProductService *u_product.UserProductService
	transmit *providers.TransmitProvider
	// productRepo        *u_product.UserProductRepo
	UserProductService iu_product.IUserProductService
	ProductRepo        *u_product.UserProductRepo
}

func NewUserPostService(repo *PostRepo,
	planService *providers.MemberPlanProvider,
	UserProductService iu_product.IUserProductService,
	transmit *providers.TransmitProvider,
	productRepo *u_product.UserProductRepo,
	// serviceManager *manager.ServiceManager,
) *UserPostService {
	return &UserPostService{
		BaseService: crud2.BaseService[models.Post, *PostRepo]{
			Repo: repo,
		},
		planService:        planService,
		UserProductService: UserProductService,
		transmit:           transmit,
		ProductRepo:        productRepo,
		// serviceManager:     serviceManager,
	}
}

func (s *UserPostService) CreatePost(c *gin.Context, dto *u_post_dto.PostSaveRequest) (*models.Post, error) {
	plan := s.planService.CurrentPlan(c)
	if plan == nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: "Vui lòng thử lại sau",
		}
	}

	// _, err := s.UserProductService.RequiredOwner(c, dto.ProductID)
	_, err := s.UserProductService.RequiredOwnerProduct(c, dto.ProductID)
	if err != nil {
		return nil, err
	}

	principal := _jwt.GetPrincipal(c)
	// profileId := _jwt.GetProfileId(c)
	if plan.LimitPost <= s.Repo.CountPostFromDate(c, time.Now().Add(-30)) {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Số lượng tạo tin đã đạt giới hạn. Vui lòng nâng cấp gói",
			// Bạn đã hết hạn mức tin
		}
	}
	// plan := s.planService.Plans
	if existed := s.Repo.ExistedPostWorking(c, dto.ProductID, principal.ProfileId); existed {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Tin đăng vẫn còn hạn",
		}
	}

	numDate := dto.NumDate
	if numDate <= 0 {
		numDate = 15
	}
	if numDate > 30 {
		numDate = 30
	}
	entity := &models.Post{
		ProductId:       dto.ProductID,
		Title:           dto.Title,
		Content:         dto.Content,
		Price:           dto.Price,
		TransactionType: dto.TransactionType,
		PackageVisible:  dto.PackageVisible,
		Visibility:      dto.Visibility,
		NumDate:         numDate,
		ExpiredAt:       dto.ExpiredAt,
		Status:          enums.PostPublic,
		Hidden:          dto.Hidden,
	}
	if entity.ExpiredAt == nil {
		expired := time.Now().Add(time.Duration(numDate) * 24 * time.Hour)
		entity.ExpiredAt = &expired
	}

	err = _db.SaveWithAudit(c, entity)
	visibleAt := entity.ExpiredAt.Add(-3 * 24 * time.Hour)
	s.transmit.NotiClient.Send(c, &pb_noti.NotiNewRequest{
		Title:     "Hết hạn tin đăng",
		Message:   "Tin đăng của bạn sắp hết hạn. Vui lòng kiểm tra lại",
		Type:      3,
		UserID:    principal.ProfileId,
		VisibleAt: timestamppb.New(visibleAt),
	})

	return entity, err
}

func (s *UserPostService) UpdateHidden(c *gin.Context, dto u_post_dto.UpdateHiddenRequest) (*models.Post, error) {
	entity, err := s.RequiredOwner(c, dto.PostID)
	if err != nil {
		return nil, err
	}

	entity.Hidden = dto.Hidden

	err = _db.SaveWithAudit(c, entity)
	return entity, err
}

func (s *UserPostService) Action(c *gin.Context, dto u_post_dto.ActionRequest) (*models.Post, error) {
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

func (s *UserPostService) UpdatePost(c *gin.Context, id uint64, dto u_post_dto.PostUpdateRequest) (*models.Post, error) {
	entity, err := s.RequiredOwner(c, id)
	if err != nil {
		return nil, err
	}

	copier.Copy(&entity, &dto)

	err = _db.SaveWithAudit(c, entity)
	return entity, err
}

func (s *UserPostService) RequiredOwner(c *gin.Context, id uint64) (*models.Post, error) {
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

func (s *UserPostService) MakeNewExpired(c *gin.Context, dto u_post_dto.PostExpiredRequest) (*models.Post, error) {
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

func (s *UserPostService) Detail(c *gin.Context, id uint64) (*u_post_dto.PostDetail, error) {
	post, err := s.Repo.GetByID(c, id)
	if err != nil {
		return nil, err
	}
	product, _ := s.ProductRepo.GetByID(c, post.ProductId)
	if err != nil {
		return nil, err
	}

	return &u_post_dto.PostDetail{
		Post:        post,
		ProductInfo: product,
	}, nil
}

func (s *UserPostService) PersonalPost(c *gin.Context, dto u_post_dto.SearchRequest) (_routes.ResponseDTO, error) {
	profileId := _jwt.GetProfileId(c)
	entities, total, err := s.Repo.Search(c, profileId, dto)
	if err != nil {
		return _routes.ResponseDTO{}, &_routes.Except{
			Code:    500,
			Message: err.Error(),
		}
	}

	return _routes.ResponseDTO{
		Code:          0,
		Data:          entities,
		TotalElements: &total,
	}, nil
}

func (s *UserPostService) GlobalPost(c *gin.Context, dto u_post_dto.SearchRequest) (*_routes.ResponseDTO, error) {
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

func (s *UserPostService) ProductPostLink(c *gin.Context, id uint64, dto *user_dto.PostLinkSearch) (*[]models.PostItem, error) {
	entities, err := s.Repo.SearchLinkToProduct(c, id, dto)
	return &entities, err
}
