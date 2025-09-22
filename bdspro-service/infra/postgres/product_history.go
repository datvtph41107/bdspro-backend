package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	_db "common/db"
	"context"
	"encoding/json"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.ProductHistoryRepo
type PostgreHistoryRepo struct {
	DB *gorm.DB
}

func NewPostgreHistoryRepo(DB *gorm.DB) *PostgreHistoryRepo {
	return &PostgreHistoryRepo{
		DB: DB,
	}
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

func (s *PostgreHistoryRepo) CreateHistory(
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

func (s *PostgreHistoryRepo) History(
	c context.Context,
	profileId uint64,
	productId uint64,
	dto dto.ProductHistorySearch,
) (*[]domain.ProductHistory, int64, error) {
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

	if err := s.DB.Debug().Raw(query, productId, dto.GetLimit(), dto.GetOffset()).Scan(&raws).Error; err != nil {
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
