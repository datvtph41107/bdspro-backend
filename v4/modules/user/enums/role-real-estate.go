package enums

type ERoleRealEstate int

const (
	ERoleRealEstateBroker    ERoleRealEstate = 10
	ERoleRealEstateAppraiser ERoleRealEstate = 20
	ERoleRealEstateLegal     ERoleRealEstate = 30
	ERoleRealEstateInvestor  ERoleRealEstate = 40
	ERoleRealEstateCustomer  ERoleRealEstate = 50
	ERoleRealEstateManager   ERoleRealEstate = 60
)

var RoleRealEstateMap = map[ERoleRealEstate]string{
	ERoleRealEstateBroker:    "Môi giới",
	ERoleRealEstateAppraiser: "Định giá",
	ERoleRealEstateLegal:     "Pháp lý",
	ERoleRealEstateInvestor:  "Nhà đầu tư",
	ERoleRealEstateCustomer:  "Khách hàng",
	ERoleRealEstateManager:   "Doanh nghiệp",
}

func (r ERoleRealEstate) IsValid() bool {
	return r == ERoleRealEstateBroker ||
		r == ERoleRealEstateAppraiser ||
		r == ERoleRealEstateLegal ||
		r == ERoleRealEstateInvestor ||
		r == ERoleRealEstateCustomer ||
		r == ERoleRealEstateManager
}
