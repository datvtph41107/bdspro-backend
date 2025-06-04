package shared_usecase

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	user_repo "bdspro/internal/repo/user"
	_routes "common/routes"
	_utils "common/utils"

	"github.com/gin-gonic/gin"
)

type ProductChildUsecase struct {
	productChildRepo  repo.UProductChildRepo
	productRepo       user_repo.UProductRepo
	postRepo          repo.UPostRepo
	assetRepo         user_repo.UAssetRepo
	mediaRepo         repo.ProductMediaRepo
	productAccessRepo repo.ProductAccessRepo
	houseInfoRepo     repo.HouseInfoRepo
	priceRepo         repo.ProductPriceRepo
	productPrivate    repo.ProductPrivateRepo
	region            repo.RegionRepo
	docTypeRepo       repo.DocTypeRepo
	propertyType      repo.PropertyTypeRepo
	amenity           repo.AmenityRepo
	ProductUsecase    *ProductUsecase
}

func MakeProductChildUsecase(
	productChildRepo repo.UProductChildRepo,
	productRepo user_repo.UProductRepo,
	postRepo repo.UPostRepo,
	assetRepo user_repo.UAssetRepo,
	mediaRepo repo.ProductMediaRepo,
	productAccessRepo repo.ProductAccessRepo,
	houseInfoRepo repo.HouseInfoRepo,
	priceRepo repo.ProductPriceRepo,
	productPrivate repo.ProductPrivateRepo,
	region repo.RegionRepo,
	docTypeRepo repo.DocTypeRepo,
	propertyType repo.PropertyTypeRepo,
	amenity repo.AmenityRepo,
	// UserProductService *UProductUsecase,
) *ProductChildUsecase {
	return &ProductChildUsecase{
		productChildRepo:  productChildRepo,
		productRepo:       productRepo,
		postRepo:          postRepo,
		assetRepo:         assetRepo,
		productAccessRepo: productAccessRepo,
		mediaRepo:         mediaRepo,
		houseInfoRepo:     houseInfoRepo,
		priceRepo:         priceRepo,
		productPrivate:    productPrivate,
		region:            region,
		docTypeRepo:       docTypeRepo,
		propertyType:      propertyType,
		amenity:           amenity,
		// UserProductService: UserProductService,
	}
}

func (s *ProductChildUsecase) RequiredBeforeSave(c *gin.Context, parentID uint64) (*domain.Product, error) {
	parent, err := s.productChildRepo.GetByID(c, parentID)
	if err != nil {
		return nil, err
	}

	// if parent.AssetId == nil {
	// 	return nil, &_routes.Except{
	// 		Code:    400,
	// 		Message: "Sản phẩm không đủ điều kiến để tạo sản phẩm con",
	// 	}
	// }

	profileId := _utils.GetProfileIdWithContext(c)

	if *parent.OwnerUserID != profileId {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Sản phẩm không thuộc quyền sở hữu của bạn",
		}
	}

	return parent, nil
}

func (s *ProductChildUsecase) InfoArea(c *gin.Context, id *uint64) (*dto.InfoAreaResponse, error) {
	_, err := s.RequiredBeforeSave(c, *id)
	if err != nil {
		return nil, err
	}

	profileId := _utils.GetProfileIdWithContext(c)
	result, err := s.productChildRepo.InfoArea(c, &profileId, id)
	return result, err
}

func (s *ProductChildUsecase) Create(c *gin.Context, body *dto.SaveProductChildRequest) (*dto.CreateResponse, error) {
	parent, err := s.RequiredBeforeSave(c, *body.ParentID)
	if err != nil {
		return nil, err
	}
	totalAreaDevide := s.productChildRepo.TotalAreaDevide(body.ParentID)
	if totalAreaDevide+body.Area > parent.Area {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Diện tích không hợp lệ",
		}
	}

	body.ProductSaveRequest.ParentID = body.ParentID
	result, err := s.ProductUsecase.CreateProduct(c, &body.ProductSaveRequest)
	s.productChildRepo.CreateHistory(c, enums.EProductHistory_NewChild,
		body.ParentID,
		*body.ParentID,
		[]uint64{result.Product.ID},
		[]float64{result.Product.Area},
	)
	return result, err
}

func (s *ProductChildUsecase) MergeProductChild(c *gin.Context, body *dto.MergeProductChildRequest) (*domain.Product, error) {
	_, err := s.RequiredBeforeSave(c, *body.ParentID)
	if err != nil {
		return nil, err
	}

	validate := s.productChildRepo.ValidatorMerge(*body.ParentID, *body.ChildIDs)
	if !validate {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Tồn tại sản phẩm con không đạt điều kiện để gộp",
		}
	}

	s.productChildRepo.CreateHistory(c,
		enums.EProductHistory_Merge,
		nil,
		*body.ParentID,
		*body.ChildIDs,
		[]float64{},
	)
	// body.ProductSaveRequest.ParentID = body.ParentID
	return nil, s.productChildRepo.MergeChilds(*body.ChildIDs)
}

func (s *ProductChildUsecase) MergeAllProductChild(c *gin.Context, parentID *uint64) (*domain.Product, error) {
	_, err := s.RequiredBeforeSave(c, *parentID)
	if err != nil {
		return nil, err
	}

	validate := s.productChildRepo.ValidListChild(*parentID)
	if !validate {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Tồn tại sản phẩm con không đạt điều kiện để gộp",
		}
	}

	ids := s.productChildRepo.AllChildId(*parentID)
	err = s.productChildRepo.MergeAll(*parentID)
	if err != nil {
		return nil, err
	}
	// body.ProductSaveRequest.ParentID = body.ParentID

	s.productChildRepo.CreateHistory(c,
		enums.EProductHistory_MergeAll,
		nil,
		*parentID, ids,
		[]float64{},
	)
	return nil, nil
}

func (s *ProductChildUsecase) DevideProductChild(c *gin.Context, devideData *dto.DevideChildRequest) (*[]domain.Product, error) {
	child, err := s.RequiredBeforeSave(c, *devideData.ParentID)
	if err != nil {
		return nil, err
	}

	if child.ParentId == nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Chia sản phẩm con chỉ áp dụng đối với sản phẩm con",
		}
	}

	totalDevide := 0.0
	if devideData.Childs != nil {
		for i := 0; i < len(devideData.Childs); i++ {
			totalDevide += devideData.Childs[i].Area
			devideData.Childs[i].ParentID = child.ParentId
		}
	}

	if child.Area <= totalDevide {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Diện tích không hợp lệ",
		}
	}
	validate := s.productChildRepo.ValidListChild(*devideData.ParentID)
	if !validate {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Sản phẩm con không đạt điều kiện để chia",
		}
	}

	// devideu_product_dto.ProductSaveRequest.ParentID = child.ParentId
	s.productRepo.UpdateArea(c, child.ID, child.Area-totalDevide)
	var results []domain.Product
	for i := 0; i < len(devideData.Childs); i++ {
		result, err := s.ProductUsecase.CreateProduct(c, &devideData.Childs[i])
		if err != nil {
			return nil, err
		}
		results = append(results, *result.Product)
	}

	var ids []uint64
	var areas []float64

	for _, r := range results {
		ids = append(ids, r.ID)
		areas = append(areas, r.Area)
	}
	s.productChildRepo.CreateHistory(c, enums.EProductHistory_Devide,
		&child.ID,
		*child.ParentId,
		ids,
		areas,
	)
	return &results, nil
}

func (s *ProductChildUsecase) History(c *gin.Context, id uint64, dto *dto.ProductHistorySearch) (_routes.ResponseDTO, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	result, total, err := s.productChildRepo.History(c, profileId, id, dto)
	if err != nil {
		return _routes.ResponseDTO{}, err
	}
	return _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, nil
}
