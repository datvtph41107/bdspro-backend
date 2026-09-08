package enums

type ELineageStatus uint32

const (
	LineageDraft      ELineageStatus = 10 // user tạo
	LineagePending    ELineageStatus = 20 // chờ verify
	LineageActive     ELineageStatus = 30 // bản chính
	LineageSuperseded ELineageStatus = 40 // bị thay thế
	LineageArchived   ELineageStatus = 50 // lưu trữ
	LineageRejected   ELineageStatus = 60 // bị từ chối
)
