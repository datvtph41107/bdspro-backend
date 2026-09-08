package enums

type ScopeApply uint32

const (
	ScopeApplyBDSPro   ScopeApply = 10
	ScopeApplyQHPro    ScopeApply = 20
	ScopeApplyAPI      ScopeApply = 30
	ScopeApplyInternal ScopeApply = 40
)

var ScopeApplyNames = map[ScopeApply]string{
	ScopeApplyBDSPro:   "BDSPro",
	ScopeApplyQHPro:    "QHPro",
	ScopeApplyAPI:      "API",
	ScopeApplyInternal: "Nội bộ",
}

func (e ScopeApply) String() string {
	switch e {
	case ScopeApplyBDSPro:
		return "BDSPro"
	case ScopeApplyQHPro:
		return "QHPro"
	case ScopeApplyAPI:
		return "API"
	case ScopeApplyInternal:
		return "Nội bộ"
	}
	return "Unknown"
}

func (e ScopeApply) IsValid() bool {
	switch e {
	case ScopeApplyBDSPro, ScopeApplyQHPro, ScopeApplyAPI, ScopeApplyInternal:
		return true
	}
	return false
}
