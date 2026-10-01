package quota

import commonoperation "common/operation"

const (
	AreaAnalysisOperation     commonoperation.Code = "workspace.area.analyze"
	ParcelDetailViewOperation commonoperation.Code = "workspace.parcel.detail.view"
	GeneratedReportOperation  commonoperation.Code = "workspace.report.generate"
)

// profileOperations chỉ sở hữu thứ tự hiển thị usage của TQD. Mã operation
// phải khớp annotation protobuf; User Service sở hữu việc gói nào cho phép dùng
// và tính bao nhiêu unit. Mở map/preview công khai không tiêu thụ lượt.
var profileOperations = [...]commonoperation.Code{
	AreaAnalysisOperation,
	ParcelDetailViewOperation,
	GeneratedReportOperation,
}

// ProfileOperations trả bản sao để caller không thể sửa registry của request
// khác. Test contract đối chiếu tập mã này với annotation RPC trong protobuf.
func ProfileOperations() []commonoperation.Code {
	return append([]commonoperation.Code(nil), profileOperations[:]...)
}
