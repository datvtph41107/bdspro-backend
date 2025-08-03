package shared_usecase

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	delivery_dto "bdspro/internal/dto/delivery"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	"bdspro/internal/usecases"
	iusecase "bdspro/internal/usecases/interface"
	_db "common/db"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"time"
)

// type IPostUsecase interface {
// 	Save(c context.Context, entity *domain.Post) error
// }

type PostUsecase struct {
	// This           IPostUsecase
	ClientUC       iusecase.INotificationClient
	ProductUsecase *ProductUsecase // forward reference
	ProductRepo    repo.ProductRepo
	Repo           repo.PostRepo
	PostMediaRepo  repo.PostMediaRepo
	AssetRepo      repo.AssetRepo
	PlanUsecase    usecases.PlanUsecase
	RequestUsecase *RequestUsecase
}

func NewSharedPostUsecase(repo repo.PostRepo,
	// planService *providers.MemberPlanProvider,
	productRepo repo.ProductRepo,
	postMediaRepo repo.PostMediaRepo,
	assetRepo repo.AssetRepo,
	planUsecase usecases.PlanUsecase,
	clientUsecase iusecase.INotificationClient,
	// serviceManager *manager.ServiceManager,
) *PostUsecase {
	return &PostUsecase{
		// planService: planService,
		ProductRepo:   productRepo,
		Repo:          repo,
		PostMediaRepo: postMediaRepo,
		AssetRepo:     assetRepo,
		PlanUsecase:   planUsecase,
		ClientUC:      clientUsecase,
		// serviceManager:     serviceManager,
	}
}

func (s *PostUsecase) CreatePost(c context.Context, dto *dto.PostSaveRequest) (*domain.Post, error) {
	plan := s.PlanUsecase.CurrentPlan(c)
	// todo: nhớ check plan
	// if plan == nil {
	// 	return nil, &_routes.Except{
	// 		Code:    500,
	// 		Message: "Vui lòng thử lại sau",
	// 	}
	// }

	_, err := s.ProductUsecase.RequiredOwner(c, dto.ProductID)
	if err != nil {
		return nil, err
	}

	profileId := _utils.GetProfileIdWithContext(c)
	// profileId := _jwt.GetProfileId(c)
	if plan != nil && plan.LimitPost <= s.Repo.CountPostFromDate(c, time.Now().Add(-30)) {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Số lượng tạo tin đã đạt giới hạn. Vui lòng nâng cấp gói",
			// Bạn đã hết hạn mức tin
		}
	}
	// plan := s.planService.Plans
	if existed := s.Repo.ExistedPostWorking(c, dto.ProductID, profileId, dto.TransactionType); existed {
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

	ownerId, ownerType, err := s.RequestUsecase.GetOwnerIDAndType(c)
	if err != nil {
		return nil, err
	}

	// status := enums.EPostSelling
	// if dto.TransactionType == 2 {
	// 	status = enums.EPostRenting
	// }
	entity := &domain.Post{
		ProductId:       dto.ProductID,
		Title:           dto.Title,
		Content:         dto.Content,
		Price:           dto.Price,
		TransactionType: dto.TransactionType,
		PriceType:       dto.PriceType,
		PackageVisible:  dto.PackageVisible,
		Visibility:      dto.Visibility,
		NumDate:         numDate,
		ExpiredAt:       dto.ExpiredAt,
		// Status:          enums.EPostPublic,
		Hidden: dto.Hidden,
		// todo:
		OwnerType: ownerType,
	}
	if ownerId != nil {
		entity.OwnerID = *ownerId
	}
	if entity.ExpiredAt == nil {
		expired := time.Now().Add(time.Duration(numDate) * 24 * time.Hour)
		entity.ExpiredAt = &expired
	}

	err = s.Repo.CreatePost(c, entity)

	if err := s.PostMediaRepo.UpdateMediaItems(c, entity.ID, dto.Medatadatas); err != nil {
		return nil, err
	}

	visibleAt := entity.ExpiredAt.Add(-3 * 24 * time.Hour)
	s.ClientUC.SendNoti(c, &delivery_dto.NotiDTO{
		Title:     "Hết hạn tin đăng",
		Message:   "Tin đăng của bạn sắp hết hạn. Vui lòng kiểm tra lại",
		Type:      3,
		UserID:    profileId,
		VisibleAt: &visibleAt,
	})

	return entity, err
}

func (s *PostUsecase) UpdateHidden(c context.Context, dto dto.UpdateHiddenRequest) (*domain.Post, error) {
	entity, err := s.RequiredOwner(c, dto.PostID)
	if err != nil {
		return nil, err
	}

	entity.Hidden = dto.Hidden

	err = _db.SaveWithContext(c, entity)
	return &entity.Post, err
}

func (s *PostUsecase) Action(c context.Context, dto dto.ActionRequest) (*domain.Post, error) {
	entity, err := s.RequiredOwner(c, dto.PostID)
	if err != nil {
		return nil, err
	}

	if dto.Action == 1 {
		entity.NumView = entity.NumView + 1
		err = s.Repo.UpdatePost(c, &entity.Post)
	}

	return nil, err
}

func (s *PostUsecase) UpdatePost(c context.Context, id uint64, dto *dto.PostUpdateRequest) (*domain.Post, error) {
	entity, err := s.RequiredOwner(c, id)
	if err != nil {
		return nil, err
	}

	// entity.TransactionType = dto.TransactionType
	entity.Title = dto.Title
	entity.PriceType = dto.PriceType
	entity.Price = dto.Price
	entity.Content = dto.Content
	entity.NumDate = dto.NumDate
	entity.Visibility = dto.Visibility
	entity.PackageVisible = dto.PackageVisible

	err = s.Repo.UpdatePost(c, &entity.Post)

	if err := s.PostMediaRepo.UpdateMediaItems(c, entity.ID, dto.MediaList); err != nil {
		return nil, err
	}
	return &entity.Post, err
}

func (s *PostUsecase) RequiredOwner(c context.Context, id uint64) (*domain.PostItem, error) {
	post, err := s.Repo.GetPostItemByID(c, id)
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

func (s *PostUsecase) MakeNewExpired(c context.Context, dto dto.PostExpiredRequest) (*domain.Post, error) {
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
	s.ClientUC.SendNoti(c, &delivery_dto.NotiDTO{
		Title:     "Hết hạn tin đăng",
		Message:   "Tin đăng của bạn sắp hết hạn. Vui lòng kiểm tra lại",
		Type:      3,
		UserID:    profileId,
		VisibleAt: &visibleAt,
	})

	return &entity.Post, err
}

func (s *PostUsecase) Detail(c context.Context, id uint64) (*dto.PostDetail, error) {
	post, err := s.Repo.GetPostItemByID(c, id)
	if err != nil {
		return nil, err
	}
	product, _ := s.ProductUsecase.Detail(c, post.ProductId)
	if err != nil {
		return nil, err
	}

	var asset *domain.Asset
	if post.ProductId != 0 {
		asset, _ = s.AssetRepo.GetAssetByProductID(c, &post.ProductId)
	}

	s.MapPostStatus(post)

	return &dto.PostDetail{
		Post:        post,
		ProductInfo: product,
		AssetInfo:   asset,
	}, nil
}

func (s *PostUsecase) MapPostStatus(post *domain.PostItem) enums.EPostStatus {
	// product := post.Product
	transactionType := post.TransactionType
	if transactionType == enums.TransactionTypeSale &&
		post.TransactionStatus == enums.TransactionStatusCompleted {
		post.Status = enums.EPostSold
	} else if transactionType == enums.TransactionTypeSale {
		post.Status = enums.EPostSelling
	} else if transactionType == enums.TransactionTypeRent &&
		post.TransactionStatus == enums.TransactionStatusCompleted {
		post.Status = enums.EPostRented
	} else {
		post.Status = enums.EPostRenting
	}
	return post.Status
}

func (s *PostUsecase) PersonalPost(c context.Context, dto *dto.PostSearchRequest) ([]domain.PostItem, int64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	entities, total, err := s.Repo.Search(c, profileId, dto)
	if err != nil {
		return nil, 0, err
	}

	for i, entity := range entities {
		entities[i].Status = s.MapPostStatus(&entity)
	}

	return entities, total, nil
}

func (s *PostUsecase) GlobalPost(c context.Context, dto dto.PostSearchRequest) ([]domain.Post, int64, error) {
	entities, total, err := s.Repo.Global(c, dto)
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (s *PostUsecase) ProductPostLink(c context.Context, id uint64, dto *dto.PostLinkSearch) (*[]domain.PostItem, error) {
	entities, err := s.Repo.SearchLinkToProduct(c, id, dto)
	if err != nil {
		return nil, err
	}

	for i, entity := range entities {
		entities[i].Status = s.MapPostStatus(&entity)
	}

	return &entities, nil
}

func (s *PostUsecase) Delete(c context.Context, id uint64) error {
	requiredOwner, err := s.RequiredOwner(c, id)
	if err != nil {
		return err
	}

	if requiredOwner.Status == enums.EPostSelling || requiredOwner.Status == enums.EPostRenting {
		return &_routes.Except{
			Code:    400,
			Message: "Tin đăng đã được đăng",
		}
	}

	profileId := _utils.GetProfileIdWithContext(c)
	err = s.Repo.Delete(c, profileId, id)
	return err
}

func (s *PostUsecase) GetPublishByProfileID(c context.Context, profileId uint64, dto dto.PostPublishSearch) ([]domain.Post, int64, error) {
	entities, total, err := s.Repo.GetPublishByProfileID(c, profileId, dto)
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// func (s *PostUsecase) Save(c context.Context, entity *domain.Post) error {
// 	return _db.SaveWithContext(c, entity)
// }
