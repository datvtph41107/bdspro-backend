package enums

type EPurposeUsed uint32

const (
	EPurposeUsedODT EPurposeUsed = 10
	EPurposeUsedONT EPurposeUsed = 20
	EPurposeUsedCLN EPurposeUsed = 30
	EPurposeUsedLUC EPurposeUsed = 40
	EPurposeUsedNTS EPurposeUsed = 50
	EPurposeUsedHNK EPurposeUsed = 60
)

var EPurposeUsedNames = map[EPurposeUsed]string{
	EPurposeUsedODT: "ODT", // đất ở đô thị
	EPurposeUsedONT: "ONT", // đất ở nông thôn
	EPurposeUsedCLN: "CLN", // trồng cây lâu năm
	EPurposeUsedLUC: "LUC", // đất lúa
	EPurposeUsedHNK: "HNK", // đất trồng cây hằng năm
	EPurposeUsedNTS: "NTS", // đất nuôi trồng thủy sản
}

func (e EPurposeUsed) String() string {
	return EPurposeUsedNames[e]
}

func (e EPurposeUsed) IsValid() bool {
	return e == EPurposeUsedODT || e == EPurposeUsedONT ||
		e == EPurposeUsedCLN || e == EPurposeUsedLUC ||
		e == EPurposeUsedHNK || e == EPurposeUsedNTS
}
