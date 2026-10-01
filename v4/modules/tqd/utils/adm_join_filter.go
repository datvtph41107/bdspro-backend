package utils

import (
	"fmt"
	"strconv"

	_utils "common/utils"
)

// PlanningProjectAdmJoinClause lọc đồ án theo composite admId qua bảng nối:
// tỉnh → qh_planning_project_provinces, phường → qh_planning_project_wards.
// projectIDExpr là cột id đồ án, ví dụ "p.id" hoặc "qh_planning_projects.id".
func PlanningProjectAdmJoinClause(admID uint64, projectIDExpr string) (clause string, arg string, ok bool) {
	if admID == 0 || projectIDExpr == "" {
		return "", "", false
	}
	sourceID, kind, err := _utils.DecodeLocationID(admID)
	if err != nil {
		return "", "", false
	}
	idStr := strconv.FormatUint(sourceID, 10)
	switch kind {
	case _utils.LocationIDKindProvince:
		return fmt.Sprintf(`EXISTS (
			SELECT 1 FROM qh_planning_project_provinces ppp
			WHERE ppp.planning_project_id = %s
			  AND ppp.deleted_at IS NULL
			  AND CAST(ppp.province_id AS text) = ?
		)`, projectIDExpr), idStr, true
	case _utils.LocationIDKindWard:
		return fmt.Sprintf(`EXISTS (
			SELECT 1 FROM qh_planning_project_wards ppw
			WHERE ppw.planning_project_id = %s
			  AND ppw.deleted_at IS NULL
			  AND CAST(ppw.ward_id AS text) = ?
		)`, projectIDExpr), idStr, true
	default:
		return "", "", false
	}
}
