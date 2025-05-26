package history

import (
	"bdspro/internal/enums"
	"bdspro/v1/models"
	_db "common/db"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AssetHistoryRepo struct {
	DB *gorm.DB
}

func NewAssetHistoryRepo(DB *gorm.DB) *AssetHistoryRepo {
	return &AssetHistoryRepo{
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

// deprecated
func (s *AssetHistoryRepo) CreateHistory(
	c *gin.Context,
	action enums.AssetHistory,
	assetID *uint64,
	// parentID uint64,
	BeforeData *string,
	AfterData *string,
) error {
	return _db.SaveWithAudit(c, &models.AssetHistory{
		AssetID:    assetID,
		ActionType: action,
		BeforeData: BeforeData,
		AfterData:  AfterData,
		// ParentID: &parentID,
		// ChildAssetIDs: ConvertToPQArray(childIDs),
	})
}

// func (s *AssetHistoryRepo) CreateHistoryLog(
// 	c *gin.Context,
// 	assetID *uint64,
// 	action enums.AssetHistory,
// 	beforeData *string,
// 	afterData *string,
// 	// parentID uint64,
// ) error {
// 	return _db.SaveWithAudit(c, &models.AssetHistory{
// 		AssetID:    assetID,
// 		ActionType: action,
// 		BeforeData: beforeData,
// 		AfterData:  afterData,
// 		// ParentID: &parentID,
// 		// ChildAssetIDs: ConvertToPQArray(childIDs),
// 	})
// }
