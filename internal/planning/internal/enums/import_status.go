package enums

type LayerImportStatus uint32

const (
	LayerImportStatusProcessing = 10 // Processing
	LayerImportStatusDone       = 20 // Done
	LayerImportStatusFailed     = 30 // Failed
)

var LayerImportStatusMap = map[LayerImportStatus]string{
	LayerImportStatusProcessing: "Processing",
	LayerImportStatusDone:       "Done",
	LayerImportStatusFailed:     "Failed",
}

func (e LayerImportStatus) String() string {
	return LayerImportStatusMap[e]
}

func (e LayerImportStatus) IsValid() bool {
	return e == LayerImportStatusProcessing || e == LayerImportStatusDone || e == LayerImportStatusFailed
}
