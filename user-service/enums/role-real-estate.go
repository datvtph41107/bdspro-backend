package enums

type ERoleRealEstate int

const (
	ERoleRealEstateBroker    ERoleRealEstate = 10
	ERoleRealEstateAppraiser ERoleRealEstate = 20
	ERoleRealEstateLegal     ERoleRealEstate = 30
)

var RoleRealEstateMap = map[ERoleRealEstate]string{
	ERoleRealEstateBroker:    "Môi giới",
	ERoleRealEstateAppraiser: "Định giá",
	ERoleRealEstateLegal:     "Pháp lý",
}
