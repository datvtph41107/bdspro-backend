package enums

type DealMemberType uint

const (
	DealMemberTypeCustomer DealMemberType = 10
	DealMemberTypePartner  DealMemberType = 20
	DealMemberTypeMember   DealMemberType = 30
)

var DealMemberTypeMap = map[DealMemberType]string{
	DealMemberTypeCustomer: "Khách hàng",
	DealMemberTypePartner:  "Đối tác",
	DealMemberTypeMember:   "Thành viên",
}

func (t DealMemberType) IsValid() bool {
	return t == DealMemberTypeCustomer || t == DealMemberTypePartner || t == DealMemberTypeMember
}
