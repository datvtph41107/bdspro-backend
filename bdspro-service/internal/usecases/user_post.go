package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	"bdspro/v1/providers"
	_db "common/db"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	pb_noti "pb/types/notification"
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UPostUsecase struct {
	// crud2.BaseService[domain.Post, *repo.UPostRepo]
	planService *providers.MemberPlanProvider
	// UserProductService *u_product.UserProductService
	transmit *providers.TransmitProvider
	// productRepo        *u_product.UserProductRepo
	UProductUsecase *UProductUsecase
	ProductRepo     repo.UProductRepo
	Repo            repo.UPostRepo
	PostMediaRepo   repo.PostMediaRepo
	AssetRepo       repo.AssetRepo
}

func NewUPostUsecase(repo repo.UPostRepo,
	planService *providers.MemberPlanProvider,
	// UserProductService *UProductUsecase,
	transmit *providers.TransmitProvider,
	productRepo repo.UProductRepo,
	postMediaRepo repo.PostMediaRepo,
	assetRepo repo.AssetRepo,
	// serviceManager *manager.ServiceManager,
) *UPostUsecase {
	return &UPostUsecase{
		// BaseService: crud2.BaseService[domain.Post, *repo.UPostRepo]{
		// 	Repo: repo,
		// },
		planService: planService,
		// UProductUsecase: UserProductService,
		transmit:      transmit,
		ProductRepo:   productRepo,
		Repo:          repo,
		PostMediaRepo: postMediaRepo,
		AssetRepo:     assetRepo,
		// serviceManager:     serviceManager,
	}
}

func (s *UPostUsecase) CreatePost(c context.Context, dto *dto.PostSaveRequest) (*domain.Post, error) {
	plan := s.planService.CurrentPlanFromContext(c)
	if plan == nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: "Vui lòng thử lại sau",
		}
	}

	// _, err := s.UserProductService.RequiredOwner(c, dto.ProductID)
	_, err := s.UProductUsecase.RequiredOwner(c, dto.ProductID)
	if err != nil {
		return nil, err
	}

	profileId := _utils.GetProfileIdWithContext(c)
	// profileId := _jwt.GetProfileId(c)
	if plan.LimitPost <= s.Repo.CountPostFromDate(c, time.Now().Add(-30)) {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Số lượng tạo tin đã đạt giới hạn. Vui lòng nâng cấp gói",
			// Bạn đã hết hạn mức tin
		}
	}
	// plan := s.planService.Plans
	if existed := s.Repo.ExistedPostWorking(c, dto.ProductID, profileId); existed {
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
	entity := &domain.Post{
		ProductId:       dto.ProductID,
		Title:           dto.Title,
		Content:         dto.Content,
		Price:           dto.Price,
		TransactionType: dto.TransactionType,
		PackageVisible:  dto.PackageVisible,
		Visibility:      dto.Visibility,
		NumDate:         numDate,
		ExpiredAt:       dto.ExpiredAt,
		Status:          enums.EPostPublic,
		Hidden:          dto.Hidden,
	}
	if entity.ExpiredAt == nil {
		expired := time.Now().Add(time.Duration(numDate) * 24 * time.Hour)
		entity.ExpiredAt = &expired
	}

	err = _db.SaveWithContext(c, entity)

	if err := s.PostMediaRepo.UpdateMediaItems(c, entity.ID, dto.MediaList); err != nil {
		return nil, err
	}

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

func (s *UPostUsecase) UpdateHidden(c context.Context, dto dto.UpdateHiddenRequest) (*domain.Post, error) {
	entity, err := s.RequiredOwner(c, dto.PostID)
	if err != nil {
		return nil, err
	}

	entity.Hidden = dto.Hidden

	err = _db.SaveWithContext(c, entity)
	return entity, err
}

func (s *UPostUsecase) Action(c context.Context, dto dto.ActionRequest) (*domain.Post, error) {
	entity, err := s.RequiredOwner(c, dto.PostID)
	if err != nil {
		return nil, err
	}

	if dto.Action == 1 {
		entity.NumView = entity.NumView + 1
		err = _db.SaveWithContext(c, entity)
	}

	return nil, err
}

func (s *UPostUsecase) UpdatePost(c context.Context, id uint64, dto dto.PostUpdateRequest) (*domain.Post, error) {
	entity, err := s.RequiredOwner(c, id)
	if err != nil {
		return nil, err
	}

	copier.Copy(&entity, &dto)

	err = _db.SaveWithContext(c, entity)

	if err := s.PostMediaRepo.UpdateMediaItems(c, entity.ID, dto.MediaList); err != nil {
		return nil, err
	}
	return entity, err
}

func (s *UPostUsecase) RequiredOwner(c context.Context, id uint64) (*domain.Post, error) {
	post, err := s.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	profileId := _utils.GetProfileIdWithContext(c)

	if *post.CreatedBy != profileId {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Bạn không có quyền truy cập",
		}
	}
	return post, nil
}

func (s *UPostUsecase) MakeNewExpired(c context.Context, dto dto.PostExpiredRequest) (*domain.Post, error) {
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

	profileId := _utils.GetProfileIdWithContext(c)
	// // plan := s.planService.Plans
	// if existed := s.Repo.ExistedPostWorking(c, dto.ProductID, profileId); existed {
	// 	return nil, &_routes.Except{
	// 		Code:    400,
	// 		Message: "Tin đăng vẫn còn hạn",
	// 	}
	// }
	// entity := &domain.Post{
	// 	ProductId:       dto.ProductID,
	// 	Content:         dto.Content,
	// 	TransactionType: dto.TransactionType,
	// 	Title:           dto.Title,
	// 	ExpiredAt:       dto.ExpiredAt,
	// }
	newExpired := entity.ExpiredAt.Add(time.Duration(time.Duration(dto.NumDay) * 24 * time.Hour))
	entity.ExpiredAt = &newExpired

	err = _db.SaveWithContext(c, entity)
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

func (s *UPostUsecase) Detail(c context.Context, id uint64) (*dto.PostDetail, error) {
	post, err := s.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	product, _ := s.ProductRepo.GetByIDContext(c, post.ProductId)
	if err != nil {
		return nil, err
	}

	var asset *domain.Asset
	if post.ProductId != 0 {
		asset, _ = s.AssetRepo.GetAssetByProductID(c, &post.ProductId)
	}

	return &dto.PostDetail{
		Post:        post,
		ProductInfo: product,
		AssetInfo:   asset,
	}, nil
}

func (s *UPostUsecase) PersonalPost(c context.Context, dto dto.PostSearchRequest) (_routes.ResponseDTO, error) {
	profileId := _utils.GetProfileIdWithContext(c)
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

func (s *UPostUsecase) GlobalPost(c context.Context, dto dto.PostSearchRequest) (*_routes.ResponseDTO, error) {
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

func (s *UPostUsecase) ProductPostLink(c context.Context, id uint64, dto *dto.PostLinkSearch) (*[]domain.PostItem, error) {
	entities, err := s.Repo.SearchLinkToProduct(c, id, dto)
	return &entities, err
}

func (s *UPostUsecase) Delete(c context.Context, id uint64) error {
	requiredOwner, err := s.RequiredOwner(c, id)
	if err != nil {
		return err
	}

	if requiredOwner.Status == enums.EPostPublic {
		return &_routes.Except{
			Code:    400,
			Message: "Tin đăng đã được đăng",
		}
	}

	profileId := _utils.GetProfileIdWithContext(c)
	err = s.Repo.Delete(c, profileId, id)
	return err
}

func (s *UPostUsecase) GetPublishByProfileID(c context.Context, profileId uint64, dto dto.PostPublishSearch) ([]domain.Post, int64, error) {
	entities, total, err := s.Repo.GetPublishByProfileID(c, profileId, dto)
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}
