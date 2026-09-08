package enums

type ERecordType string

const (
	ERecordTypeProduct ERecordType = "product"
	ERecordTypeAsset   ERecordType = "asset"
)

var RecordTypeMap = map[ERecordType]string{
	ERecordTypeProduct: "Sản phẩm",
	ERecordTypeAsset:   "Tài sản",
}
