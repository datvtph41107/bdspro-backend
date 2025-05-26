package u_product_child

import (
	"bdspro/internal/enums"
	a_doc_type "bdspro/v1/mod/admin/doc_type"
	a_property_type "bdspro/v1/mod/admin/property_type"
	"bdspro/v1/mod/global/amenity"
	"bdspro/v1/mod/global/house_info"
	product_accecss "bdspro/v1/mod/global/product_access"
	"bdspro/v1/mod/global/product_media"
	"bdspro/v1/mod/global/product_price"
	"bdspro/v1/mod/global/product_private"
	g_region "bdspro/v1/mod/global/region"
	u_asset "bdspro/v1/mod/user/asset"
	u_post "bdspro/v1/mod/user/post"
	u_product "bdspro/v1/mod/user/product"
	u_product_dto "bdspro/v1/mod/user/product/dto"
	"bdspro/v1/models"
	_jwt "common/jwt"
	_routes "common/routes"

	"github.com/gin-gonic/gin"
)

type ProductChildService struct {
	productChildRepo   *UserProductChildRepo
	productRepo        *u_product.UserProductRepo
	postRepo           *u_post.PostRepo
	assetRepo          *u_asset.AssetRepo
	mediaRepo          *product_media.ProductMediaRepo
	productAccessRepo  *product_accecss.ProductAccessRepo
	houseInfoRepo      *house_info.HouseInfoRepo
	priceRepo          *product_price.ProductPriceRepo
	productPrivate     *product_private.ProductPrivateRepo
	region             *g_region.RegionRepo
	docTypeRepo        *a_doc_type.DocTypeRepo
	propertyType       *a_property_type.PropertyTypeRepo
	amenity            *amenity.AmenityRepo
	UserProductService *u_product.UserProductService
}

func NewProductChildService(productChildRepo *UserProductChildRepo,
	productRepo *u_product.UserProductRepo,
	postRepo *u_post.PostRepo,
	assetRepo *u_asset.AssetRepo,
	mediaRepo *product_media.ProductMediaRepo,
	productAccessRepo *product_accecss.ProductAccessRepo,
	houseInfoRepo *house_info.HouseInfoRepo,
	priceRepo *product_price.ProductPriceRepo,
	productPrivate *product_private.ProductPrivateRepo,
	region *g_region.RegionRepo,
	docTypeRepo *a_doc_type.DocTypeRepo,
	propertyType *a_property_type.PropertyTypeRepo,
	amenity *amenity.AmenityRepo,
	UserProductService *u_product.UserProductService,
) *ProductChildService {
	return &ProductChildService{
		productChildRepo:   productChildRepo,
		productRepo:        productRepo,
		postRepo:           postRepo,
		assetRepo:          assetRepo,
		productAccessRepo:  productAccessRepo,
		mediaRepo:          mediaRepo,
		houseInfoRepo:      houseInfoRepo,
		priceRepo:          priceRepo,
		productPrivate:     productPrivate,
		region:             region,
		docTypeRepo:        docTypeRepo,
		propertyType:       propertyType,
		amenity:            amenity,
		UserProductService: UserProductService,
	}
}

func (s *ProductChildService) RequiredBeforeSave(c *gin.Context, parentID uint64) (*models.Product, error) {
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

	profileId := _jwt.GetProfileId(c)

	if *parent.OwnerUserID != profileId {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Sản phẩm không thuộc quyền sở hữu của bạn",
		}
	}

	return parent, nil
}

func (s *ProductChildService) InfoArea(c *gin.Context, id *uint64) (*InfoAreaResponse, error) {
	_, err := s.RequiredBeforeSave(c, *id)
	if err != nil {
		return nil, err
	}

	profileId := _jwt.GetProfileId(c)
	result, err := s.productChildRepo.InfoArea(c, &profileId, id)
	return result, err
}

func (s *ProductChildService) Create(c *gin.Context, body *SaveProductChildRequest) (*u_product_dto.CreateResponse, error) {
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
	result, err := s.UserProductService.CreateProduct(c, &body.ProductSaveRequest)
	s.productChildRepo.CreateHistory(c, enums.ProductHistory_NewChild,
		body.ParentID,
		*body.ParentID,
		[]uint64{result.Product.ID},
		[]float64{result.Product.Area},
	)
	return result, err
}

func (s *ProductChildService) MergeProductChild(c *gin.Context, body *MergeProductChildRequest) (*models.Product, error) {
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
		enums.ProductHistory_Merge,
		nil,
		*body.ParentID,
		*body.ChildIDs,
		[]float64{},
	)
	// body.ProductSaveRequest.ParentID = body.ParentID
	return nil, s.productChildRepo.MergeChilds(*body.ChildIDs)
}

func (s *ProductChildService) MergeAllProductChild(c *gin.Context, parentID *uint64) (*models.Product, error) {
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
		enums.ProductHistory_MergeAll,
		nil,
		*parentID, ids,
		[]float64{},
	)
	return nil, nil
}

func (s *ProductChildService) DevideProductChild(c *gin.Context, devideData *DevideChildRequest) (*[]models.Product, error) {
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
	var results []models.Product
	for i := 0; i < len(devideData.Childs); i++ {
		result, err := s.UserProductService.CreateProduct(c, &devideData.Childs[i])
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
	s.productChildRepo.CreateHistory(c, enums.ProductHistory_Devide,
		&child.ID,
		*child.ParentId,
		ids,
		areas,
	)
	return &results, nil
}

func (s *ProductChildService) History(c *gin.Context, id uint64, dto *ProductHistorySearch) (_routes.ResponseDTO, error) {
	profileId := _jwt.GetProfileId(c)
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
