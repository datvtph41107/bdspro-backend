package enums

type PinTarget uint32

const (
	PinTargetUnspecified    PinTarget = 0
	PinTargetContact        PinTarget = 1
	PinTargetContactProduct PinTarget = 2
)

var PinTargetRequest = map[PinTarget]string{
	PinTargetContact:        "Pin liên hệ",
	PinTargetContactProduct: "Pin khách hàng quan tâm",
}

type PinTargetMeta struct {
	TableName string
	// primary key columns (1 hoặc nhiều)
	IDColumns []string
	// ownership
	OwnerIDColumn string
	OwnerOfColumn string
	// pin field
	PinColumn string
}

var PinTargetMetaMap = map[PinTarget]PinTargetMeta{

	PinTargetContact: {
		TableName: "tb_contact",
		IDColumns: []string{"id"},
		// OwnerIDColumn: "owner_id",
		// OwnerOfColumn: "owner_of",
		PinColumn: "priority_pin_at",
	},

	PinTargetContactProduct: {
		TableName: "tb_contact_product",
		IDColumns: []string{"contact_id", "product_id"},
		// OwnerIDColumn: "owner_id", // nếu có
		// OwnerOfColumn: "owner_of", // nếu có
		PinColumn: "priority_pin_at",
	},
}