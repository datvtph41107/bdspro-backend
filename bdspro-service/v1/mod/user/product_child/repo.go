package u_product_child

import (
	"bdspro/internal/enums"
	product_accecss "bdspro/v1/mod/global/product_access"
	"bdspro/v1/mod/global/product_price"
	"bdspro/v1/mod/global/product_private"
	"bdspro/v1/models"
	"common/case/crud2"
	_db "common/db"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserProductChildRepo struct {
	crud2.BaseRepo[models.Product]
	PriceRepo      *product_price.ProductPriceRepo
	ProductAccess  *product_accecss.ProductAccessRepo
	ProductPrivate *product_private.ProductPrivateRepo
}

func NewUserProductChildRepo(db *gorm.DB,
	PriceRepo *product_price.ProductPriceRepo,
	ProductAccess *product_accecss.ProductAccessRepo,
	ProductPrivate *product_private.ProductPrivateRepo,
) *UserProductChildRepo {
	return &UserProductChildRepo{
		BaseRepo:       crud2.BaseRepo[models.Product]{DB: db},
		PriceRepo:      PriceRepo,
		ProductAccess:  ProductAccess,
		ProductPrivate: ProductPrivate,
	}
}

func (r *UserProductChildRepo) InfoArea(ctx *gin.Context, profileId, id *uint64) (*InfoAreaResponse, error) {
	entity := InfoAreaResponse{}
	// price := r.PriceRepo.FindByProductId(&id)

	err := _db.DB.
		Debug().
		Table("products").
		Select(`products.area as total_area, 
		products.area - COALESCE(child.total_area, 0) AS avaiable_area,
		COALESCE(child.total_area, 0) AS split_area
		`).
		Joins(`LEFT JOIN (SELECT parent_id, SUM(area) AS total_area 
						FROM products WHERE deleted_at is null GROUP BY parent_id) 
						AS child ON products.id = child.parent_id`).
		// Preload("PrivateData").
		Where("products.id = ? AND products.owner_user_id = ? AND products.deleted_at IS NULL", id, profileId).
		Scan(&entity).Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (s *UserProductChildRepo) TotalAreaDevide(parentID *uint64) float64 {
	var totalArea sql.NullFloat64
	// Tạo câu SQL động
	s.DB.Model(&models.Product{}).
		Where("parent_id = ?", parentID).
		Select("SUM(area)").
		Scan(&totalArea)

	return totalArea.Float64
}

func (s *UserProductChildRepo) ValidatorMerge(parentID uint64, childIds []uint64) bool {
	var count int64
	s.DB.Model(&models.Product{}).
		Where("parent_id = ? AND id IN (?) AND sale_status = 10 AND rent_status = 40", parentID, childIds).
		Count(&count)
	return int(count) == len(childIds)
}

func (s *UserProductChildRepo) MergeChilds(childIds []uint64) error {
	return s.DB.Model(&models.Product{}).
		Where("id in (?)", childIds).
		Update("deleted_at", time.Now()).
		Error
}

func (s *UserProductChildRepo) ValidListChild(parentID uint64) bool {
	var countAll int64
	var countValid int64
	s.DB.Model(&models.Product{}).
		Where("parent_id = ?", parentID).
		Count(&countAll)

	s.DB.Model(&models.Product{}).
		Where("parent_id = ? AND sale_status = 10 AND rent_status = 40", parentID).
		Count(&countValid)

	return countAll == countValid
}

func (s *UserProductChildRepo) AllChildId(parentId uint64) []uint64 {
	var ids []uint64
	s.DB.Debug().Model(&models.Product{}).
		Where("parent_id = ? and deleted_at is null", parentId).
		Pluck("id", &ids)
	return ids
}

func (s *UserProductChildRepo) MergeAll(parentId uint64) error {
	return s.DB.Model(&models.Product{}).
		Where("parent_id = ?", parentId).
		Update("deleted_at", time.Now()).
		Error
}

func (s *UserProductChildRepo) ValidatorDevide(parentID uint64) bool {
	var countAll int64
	s.DB.Model(&models.Product{}).
		Where("id = ? AND sale_status = 10 AND rent_status = 40", parentID).
		Count(&countAll)

	return countAll > 0
}

func ConvertToPQArray(childIDs []uint64) pq.Int64Array {
	result := make(pq.Int64Array, len(childIDs))
	for i, id := range childIDs {
		result[i] = int64(id)
	}
	return result
}
func ConvertToPQArrayFloat(areas []float64) pq.Float64Array {
	result := make(pq.Float64Array, len(areas))
	for i, area := range areas {
		result[i] = float64(area)
	}
	return result
}

func (s *UserProductChildRepo) CreateHistory(
	c *gin.Context,
	action enums.ProductHistory,
	productID *uint64,
	parentID uint64,
	childIDs []uint64,
	areas []float64,
) error {
	return _db.SaveWithAudit(c, &models.ProductChildHistory{
		Action:    action,
		ProductID: productID,
		ParentID:  &parentID,
		ChildIds:  ConvertToPQArray(childIDs),
		Areas:     ConvertToPQArrayFloat(areas),
	})
}

func (s *UserProductChildRepo) History(
	c *gin.Context,
	profileId uint64,
	id uint64,
	dto *ProductHistorySearch) (*[]models.ProductChildHistory, int64, error) {
	// var results []models.ProductChildHistory

	type rawResult struct {
		models.ProductChildHistory
		Parent       json.RawMessage `gorm:"column:parent" json:"parent"`
		Target       json.RawMessage `gorm:"column:target" json:"target"`
		ChildsRaw    json.RawMessage `gorm:"column:childs" json:"childs"`
		CreatedUser  json.RawMessage `gorm:"column:created_user" json:"createdUser"`
		TotalElement int64           `gorm:"column:total_elements" json:"totalElements"`
	}
	var raws []rawResult
	query := `SELECT 
    h.*,
    COUNT(*) OVER() AS total_elements,
    row_to_json(parent_product.*) as parent,
	row_to_json(target_product.*) as target,
	json_build_object(
		'profileId', profile.profile_id,
        'fullName', profile.full_name
    ) AS created_user,
    json_agg(p.*) AS childs
FROM 
    product_history  h
LEFT JOIN 
    products p ON p.id = ANY(h.child_ids)
LEFT JOIN 
    products target_product ON target_product.id = h.product_id
LEFT JOIN 
    products parent_product ON parent_product.id = h.parent_id
LEFT JOIN 
    profile_transfer profile ON profile.profile_id = h.created_by
WHERE 
    h.parent_id = ?
GROUP BY 
    h.id, parent_product.id,target_product.id,profile.profile_id
ORDER BY h.created_at DESC
	LIMIT ? OFFSET ? ;
`

	size := dto.Size
	page := dto.Page // Nếu page là 0, sẽ gán mặc định là 1
	if size == 0 {
		size = 20
	}

	offset := page * size

	if err := s.DB.Debug().Raw(query, id, size, offset).Scan(&raws).Error; err != nil {
		return nil, 0, err
	}

	var result []models.ProductChildHistory
	for _, r := range raws {
		// var childs []models.ProductInfo
		if r.ChildsRaw != nil {
			if err := json.Unmarshal(r.ChildsRaw, &r.ProductChildHistory.Childs); err != nil {
				return nil, 0, err
			}
		}

		if r.Parent != nil {
			if err := json.Unmarshal(r.Parent, &r.ProductChildHistory.Parent); err != nil {
				return nil, 0, err
			}
		}

		if r.Target != nil {
			if err := json.Unmarshal(r.Target, &r.ProductChildHistory.Target); err != nil {
				return nil, 0, err
			}
		}

		if r.CreatedUser != nil {
			if err := json.Unmarshal(r.CreatedUser, &r.ProductChildHistory.CreatedUser); err != nil {
				return nil, 0, err
			}
		}

		// r.ProductChildHistory.Childs = childs
		result = append(result, r.ProductChildHistory)
	}

	var totalElements int64
	if len(raws) > 0 {
		totalElements = raws[0].TotalElement
	}
	// err := s.DB.
	// 	// Model(&models.ProductChildHistory{}).
	// 	Preload("ProductData").
	// 	Where("created_by = ? AND deleted_at is null", profileId).
	// 	Find(&results).Error
	return &result, totalElements, nil
}
