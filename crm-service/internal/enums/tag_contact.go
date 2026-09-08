package enums

type ETagContact int32

const (
	TagContactCustomer ETagContact = 10
	TagContactBroker   ETagContact = 20
	TagContactPartner  ETagContact = 30
	TagContactOwner    ETagContact = 40
)

var TagContactMap = map[ETagContact]string{
	TagContactCustomer: "Khách hàng",
	TagContactBroker:   "Môi giới",
	TagContactPartner:  "Đối tác",
	TagContactOwner:    "Chủ nhà",
}

func (e ETagContact) IsValid() bool {
	switch e {
	case TagContactCustomer, TagContactBroker, TagContactPartner, TagContactOwner:
		return true
	default:
		return false
	}
}