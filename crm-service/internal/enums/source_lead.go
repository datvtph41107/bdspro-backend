package enums

type ESourceLead int

const (
	SourceLeadPost     ESourceLead = 10
	SourceLeadAdvise   ESourceLead = 20
	SourceLeadContact  ESourceLead = 30
	SourceLeadZalo     ESourceLead = 40
	SourceLeadFacebook ESourceLead = 50
	SourceLeadOther    ESourceLead = 60
)

var SourceLeadMap = map[ESourceLead]string{
	SourceLeadPost:     "Bài viết",
	SourceLeadAdvise:   "Tư vấn",
	SourceLeadContact:  "Danh bạ",
	SourceLeadZalo:     "Zalo",
	SourceLeadFacebook: "Facebook",
	SourceLeadOther:    "Khác",
}