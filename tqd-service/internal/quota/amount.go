package quota

import "common/businessoperation"

/**
 * AmountForOperation trả số quota cần giữ cho một lần chạy operation.
 *
 * Cost của quota nằm ở G5, không đặt trong business operation hay plan access.
 */
func AmountForOperation(operation businessoperation.Code) (int64, bool) {
	switch operation {
	case businessoperation.WorkspaceReportGenerate:
		return 1, true
	default:
		return 0, false
	}
}
