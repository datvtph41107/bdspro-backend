package enums

// ERelationType định nghĩa loại relation trong PropertyRelation
type ERelationType uint32

const (
	ERelationTypeProduct ERelationType = 10 // Product
	ERelationTypeAsset   ERelationType = 20 // Asset
)

// String trả về tên của relation type
func (e ERelationType) String() string {
	switch e {
	case ERelationTypeProduct:
		return "Sản phẩm"
	case ERelationTypeAsset:
		return "Tài sản"
	default:
		return "unknown"
	}
}

// IsValid kiểm tra relation type có hợp lệ không
func (e ERelationType) IsValid() bool {
	return e == ERelationTypeProduct || e == ERelationTypeAsset
}
