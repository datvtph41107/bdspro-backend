package enums

type ESourceCustomer int

const (
	SourceCustomerPost     ESourceCustomer = 10
	SourceCustomerAdvise   ESourceCustomer = 20
	SourceCustomerContact  ESourceCustomer = 30
	SourceCustomerZalo     ESourceCustomer = 40
	SourceCustomerFacebook ESourceCustomer = 50
	SourceCustomerOther    ESourceCustomer = 60
)

var SourceCustomerMap = map[ESourceCustomer]string{
	SourceCustomerPost:     "Bài viết",
	SourceCustomerAdvise:   "Tư vấn",
	SourceCustomerContact:  "Danh bạ",
	SourceCustomerZalo:     "Zalo",
	SourceCustomerFacebook: "Facebook",
	SourceCustomerOther:    "Khác",
}
