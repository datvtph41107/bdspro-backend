package u_asset_share

type SearchDTO struct {
	Page uint64
	Size uint64
}

type ArchivedDTO struct {
	ID       uint64
	Archived bool
}

type MergeDTO struct {
	IDs []uint64 `binding:"required,notemptyarray"`
}

type SplitDTO struct {
	ID uint64
}

// type
