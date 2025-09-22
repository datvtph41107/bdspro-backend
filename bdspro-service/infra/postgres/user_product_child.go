package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	_db "common/db"
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.UProductChildRepo
type PostgreProductChild struct {
	DB *gorm.DB
	// crud2.BaseRepo[domain.Product]
	// PriceRepo *product_price.ProductPriceRepo
	// ProductAccess  *product_accecss.ProductAccessRepo
	// ProductPrivate *product_private.ProductPrivateRepo
}

func NewPostgreProductChildRepo(db *gorm.DB,

// PriceRepo *product_price.ProductPriceRepo,
// ProductAccess *product_accecss.ProductAccessRepo,
// ProductPrivate *product_private.ProductPrivateRepo,
) *PostgreProductChild {
	return &PostgreProductChild{
		DB: db,
		// BaseRepo: crud2.BaseRepo[domain.Product]{DB: db},
		// PriceRepo: PriceRepo,
		// ProductAccess:  ProductAccess,
		// ProductPrivate: ProductPrivate,
	}
}

func (r *PostgreProductChild) GetByID(ctx context.Context, id uint64) (*domain.Product, error) {
	var entity domain.Product
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *PostgreProductChild) InfoArea(ctx context.Context, profileId, id *uint64) (*dto.InfoAreaResponse, error) {
	entity := dto.InfoAreaResponse{}
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
		Where("products.id = ? AND products.owner_id = ? AND products.deleted_at IS NULL", id, profileId).
		Scan(&entity).Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (s *PostgreProductChild) TotalAreaDevide(parentID *uint64) float64 {
	var totalArea sql.NullFloat64
	// Tạo câu SQL động
	s.DB.Model(&domain.Product{}).
		Where("parent_id = ?", parentID).
		Select("SUM(area)").
		Scan(&totalArea)

	return totalArea.Float64
}

func (s *PostgreProductChild) ValidatorMerge(parentID uint64, childIds []uint64) bool {
	var count int64
	s.DB.Model(&domain.Product{}).
		Where("parent_id = ? AND id IN (?) AND sale_status = 10 AND rent_status = 40", parentID, childIds).
		Count(&count)
	return int(count) == len(childIds)
}

func (s *PostgreProductChild) MergeChilds(childIds []uint64) error {
	return s.DB.Model(&domain.Product{}).
		Where("id in (?)", childIds).
		Update("deleted_at", time.Now()).
		Error
}

func (s *PostgreProductChild) ValidListChild(parentID uint64) bool {
	var countAll int64
	var countValid int64
	s.DB.Model(&domain.Product{}).
		Where("parent_id = ?", parentID).
		Count(&countAll)

	s.DB.Model(&domain.Product{}).
		Where("parent_id = ? AND sale_status = 10 AND rent_status = 40", parentID).
		Count(&countValid)

	return countAll == countValid
}

func (s *PostgreProductChild) AllChildId(parentId uint64) []uint64 {
	var ids []uint64
	s.DB.Debug().Model(&domain.Product{}).
		Where("parent_id = ? and deleted_at is null", parentId).
		Pluck("id", &ids)
	return ids
}

func (s *PostgreProductChild) MergeAll(parentId uint64) error {
	return s.DB.Model(&domain.Product{}).
		Where("parent_id = ?", parentId).
		Update("deleted_at", time.Now()).
		Error
}

func (s *PostgreProductChild) ValidatorDevide(parentID uint64) bool {
	var countAll int64
	s.DB.Model(&domain.Product{}).
		Where("id = ? AND sale_status = 10 AND rent_status = 40", parentID).
		Count(&countAll)

	return countAll > 0
}

// func ConvertToPQArray(childIDs []uint64) pq.Int64Array {
// 	result := make(pq.Int64Array, len(childIDs))
// 	for i, id := range childIDs {
// 		result[i] = int64(id)
// 	}
// 	return result
// }
// func ConvertToPQArrayFloat(areas []float64) pq.Float64Array {
// 	result := make(pq.Float64Array, len(areas))
// 	for i, area := range areas {
// 		result[i] = float64(area)
// 	}
// 	return result
// }

func (s *PostgreProductChild) CreateHistory(
	c context.Context,
	action enums.EProductHistory,
	productID *uint64,
	parentID uint64,
	childIDs []uint64,
	areas []float64,
) error {
	return _db.SaveWithContext(c, &domain.ProductHistory{
		Action:    action,
		ProductID: productID,
		ParentID:  &parentID,
		ChildIds:  ConvertToPQArray(childIDs),
		Areas:     ConvertToPQArrayFloat(areas),
	})
}

func (s *PostgreProductChild) History(
	c context.Context,
	profileId uint64,
	id uint64,
	dto *dto.ProductHistorySearch) (*[]domain.ProductHistory, int64, error) {
	// var results []domain.ProductChildHistory

	type rawResult struct {
		domain.ProductHistory
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

	// size := dto.Size
	// page := dto.Page // Nếu page là 0, sẽ gán mặc định là 1
	// if size == 0 {
	// 	size = 20
	// }

	// offset := page * size

	if err := s.DB.Debug().Raw(query, id, dto.GetLimit(), dto.GetOffset()).Scan(&raws).Error; err != nil {
		return nil, 0, err
	}

	var result []domain.ProductHistory
	for _, r := range raws {
		// var childs []domain.ProductInfo
		if r.ChildsRaw != nil {
			if err := json.Unmarshal(r.ChildsRaw, &r.ProductHistory.Childs); err != nil {
				return nil, 0, err
			}
		}

		if r.Parent != nil {
			if err := json.Unmarshal(r.Parent, &r.ProductHistory.Parent); err != nil {
				return nil, 0, err
			}
		}

		if r.Target != nil {
			if err := json.Unmarshal(r.Target, &r.ProductHistory.Target); err != nil {
				return nil, 0, err
			}
		}

		if r.CreatedUser != nil {
			if err := json.Unmarshal(r.CreatedUser, &r.ProductHistory.CreatedUser); err != nil {
				return nil, 0, err
			}
		}

		// r.ProductChildHistory.Childs = childs
		result = append(result, r.ProductHistory)
	}

	var totalElements int64
	if len(raws) > 0 {
		totalElements = raws[0].TotalElement
	}
	// err := s.DB.
	// 	// Model(&domain.ProductChildHistory{}).
	// 	Preload("ProductData").
	// 	Where("created_by = ? AND deleted_at is null", profileId).
	// 	Find(&results).Error
	return &result, totalElements, nil
}
