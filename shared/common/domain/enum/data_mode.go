package _enum

type EDataMode string

const (
	EDataModePartial EDataMode = "partial"
	EDataModeFull    EDataMode = "full"
)

var DataModeNames = map[EDataMode]string{
	EDataModePartial: "Partial",
	EDataModeFull:    "Full",
}
