package enums

type ESourceLead int

const (
	SourceCustomerPost     ESourceLead = 10
	SourceCustomerAdvise   ESourceLead = 20
	SourceCustomerContact  ESourceLead = 30
	SourceCustomerZalo     ESourceLead = 40
	SourceCustomerFacebook ESourceLead = 50
	SourceCustomerOther    ESourceLead = 60
)

var SourceCustomerMap = map[ESourceLead]string{
	SourceCustomerPost:     "Bài viết",
	SourceCustomerAdvise:   "Tư vấn",
	SourceCustomerContact:  "Danh bạ",
	SourceCustomerZalo:     "Zalo",
	SourceCustomerFacebook: "Facebook",
	SourceCustomerOther:    "Khác",
}
