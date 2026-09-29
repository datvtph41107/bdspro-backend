package enums

type ERefSourceSys uint32

const (
	ERefSourceSysCountry  ERefSourceSys = 10
	ERefSourceSysInvestor ERefSourceSys = 20
	ERefSourceSysPartner  ERefSourceSys = 30
)

var ERefSourceSystemNames = map[ERefSourceSys]string{
	ERefSourceSysCountry:  "Nhà nước",
	ERefSourceSysInvestor: "Chủ đầu tư",
	ERefSourceSysPartner:  "Đối tác",
}

func (e ERefSourceSys) String() string {
	return ERefSourceSystemNames[e]
}

func (e ERefSourceSys) IsValid() bool {
	return e == ERefSourceSysCountry || e == ERefSourceSysInvestor || e == ERefSourceSysPartner
}
