package mapper

import (
	shared_enum "pb/enums"
	pb_history "pb/types/notification"
	"time"

	"bdspro/internal/domain"
	"bdspro/internal/enums"
)

type HistoryMapper struct{}

func NewHistoryMapper() *HistoryMapper {
	return &HistoryMapper{}
}

func (m *HistoryMapper) HistoryToPb(history *domain.RecordHistory) *pb_history.HistoryDTO {
	if history == nil {
		return nil
	}

	// Convert domain enums to protobuf enums
	var targetType shared_enum.ETargetHistory
	switch history.RecordType {
	case enums.ERecordTypeProduct:
		targetType = shared_enum.TargetHistoryProduct
	case enums.ERecordTypeAsset:
		targetType = shared_enum.TargetHistoryAsset
	default:
		targetType = shared_enum.TargetHistoryProduct // default fallback
	}

	// Convert action type to protobuf enum
	var actionType shared_enum.EHistory
	switch history.ActionType {
	case enums.EActionProductCreate:
		actionType = shared_enum.HistoryProductCreate
	case enums.EActionProductUpdate:
		actionType = shared_enum.HistoryProductUpdate
	case enums.EActionProductDelete:
		actionType = shared_enum.HistoryProductDelete
	case enums.EActionAssetCreate:
		actionType = shared_enum.HistoryAssetCreate
	case enums.EActionAssetUpdate:
		actionType = shared_enum.HistoryAssetUpdate
	case enums.EActionAssetDelete:
		actionType = shared_enum.HistoryAssetDelete
	case enums.EActionAssetMerge:
		actionType = shared_enum.HistoryAssetMerge
	case enums.EActionAssetSplit:
		actionType = shared_enum.HistoryAssetSplit
	default:
		actionType = shared_enum.HistoryCreateStep // default fallback
	}

	// Convert content array to string slice
	var content []string
	if history.Content != nil {
		content = []string(history.Content)
	}

	// Format created time
	var createdAt string
	if history.CreatedAt != nil {
		createdAt = history.CreatedAt.Format(time.RFC3339)
	}

	return &pb_history.HistoryDTO{
		TargetId:   history.ID,
		TargetType: int32(targetType),
		ActionType: int32(actionType),
		Note:       content,
		OwnerId:    history.RecordID,
		CreatedAt:  createdAt,
		CreatedBy:  history.CreatedBy,
		Title:      getHistoryTitle(history.RecordType, history.ActionType),
		ActionUser: getActionUserName(history.CreatedBy), // This would need user service integration
		ActionName: getActionName(history.ActionType),
	}
}

// Helper functions to generate titles and names
func getHistoryTitle(recordType enums.ERecordType, actionType enums.EHistoryAction) string {
	recordTypeName := enums.RecordTypeMap[recordType]
	actionName := enums.HistoryActionMap[actionType]
	return actionName + " " + recordTypeName
}

func getActionName(actionType enums.EHistoryAction) string {
	return enums.HistoryActionMap[actionType]
}

func getActionUserName(userID *uint64) string {
	if userID == nil {
		return "Hệ thống"
	}
	// This would typically involve calling a user service to get the user name
	// For now, return a placeholder
	return "Người dùng"
}
