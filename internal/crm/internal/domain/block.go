package domain

type BlockEntity struct {
	ProfileId uint64
	BlockedId uint64
}

func (BlockEntity) TableName() string {
	return "block"
}