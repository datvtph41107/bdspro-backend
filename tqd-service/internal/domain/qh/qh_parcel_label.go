package qh_domain

type QHParcelLabel struct {
	ParcelInfoId uint64
	LabelId      uint64
}

func (QHParcelLabel) TableName() string {
	return "qh_parcel_label"
}
