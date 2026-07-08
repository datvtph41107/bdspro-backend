package usecase

import (
	"context"
	"crm/internal/enums"
	sharepb "pb/types/shared"
	"sort"
)

type EnumUsecase struct{}

func NewEnumUsecase() *EnumUsecase {
	return &EnumUsecase{}
}

// GetEnumByName trả về danh sách ItemV3Proto dựa trên tên enum
func (u *EnumUsecase) GetEnumByName(ctx context.Context, name string) ([]*sharepb.ItemV3Proto, error) {
	var enumMap map[interface{}]string
	var keys []interface{}

	// Map các enum theo tên
	switch name {
	case "contact_tag":
		enumMap = convertMap(enums.TagContactMap)
		keys = getKeys(enums.TagContactMap)
	case "statusContact":
		enumMap = convertMap(enums.StatusContactMap)
		keys = getKeys(enums.StatusContactMap)
	case "visibility":
		enumMap = convertMap(enums.VisibilityMap)
		keys = getKeys(enums.VisibilityMap)
	case "priority":
		enumMap = convertMap(enums.PriorityMap)
		keys = getKeys(enums.PriorityMap)
	case "ownerOf":
		enumMap = convertMap(enums.OwnerOfNames)
		keys = getKeys(enums.OwnerOfNames)
	case "step":
		enumMap = convertMap(enums.StepMap)
		keys = getKeys(enums.StepMap)
	case "sourceLead":
		enumMap = convertMap(enums.SourceLeadMap)
		keys = getKeys(enums.SourceLeadMap)
	case "friendStatus":
		enumMap = convertMap(enums.FriendStatusMap)
		keys = getKeys(enums.FriendStatusMap)
	case "condition":
		enumMap = convertMap(enums.ConditionMap)
		keys = getKeys(enums.ConditionMap)
	case "ruleThen":
		enumMap = convertMap(enums.RuleThenMap)
		keys = getKeys(enums.RuleThenMap)
	case "targetType":
		enumMap = convertMap(enums.TargetTypeMap)
		keys = getKeys(enums.TargetTypeMap)
	case "reportStatus":
		enumMap = convertMap(enums.ReportStatusMap)
		keys = getKeys(enums.ReportStatusMap)
	default:
		return []*sharepb.ItemV3Proto{}, nil
	}

	// Sắp xếp keys theo giá trị tăng dần
	sort.Slice(keys, func(i, j int) bool {
		return getIntValue(keys[i]) < getIntValue(keys[j])
	})

	// Convert sang ItemV3Proto
	result := make([]*sharepb.ItemV3Proto, 0, len(keys))
	for _, key := range keys {
		result = append(result, &sharepb.ItemV3Proto{
			Id:   uint64(getIntValue(key)),
			Name: enumMap[key],
		})
	}

	return result, nil
}

// Helper functions để convert map
func convertMap[T comparable](m map[T]string) map[interface{}]string {
	result := make(map[interface{}]string, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

func getKeys[T comparable](m map[T]string) []interface{} {
	keys := make([]interface{}, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func getIntValue(v interface{}) int64 {
	switch val := v.(type) {
	case int:
		return int64(val)
	case int32:
		return int64(val)
	case int64:
		return val
	case uint:
		return int64(val)
	case uint32:
		return int64(val)
	case uint64:
		return int64(val)
	case enums.ETagContact:
		return int64(val)
	case enums.EStatusContact:
		return int64(val)
	case enums.EVisibility:
		return int64(val)
	case enums.EPriority:
		return int64(val)
	case enums.EOwnerOf:
		return int64(val)
	case enums.EStep:
		return int64(val)
	case enums.ESourceLead:
		return int64(val)
	case enums.FriendStatus:
		return int64(val)
	case enums.ECondition:
		return int64(val)
	case enums.ERuleTrigger:
		return int64(val)
	case enums.TargetType:
		return int64(val)
	case enums.ReportStatus:
		return int64(val)
	default:
		return 0
	}
}