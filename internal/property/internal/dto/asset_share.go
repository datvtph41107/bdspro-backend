package dto

import _dto "common/domain/dto"

type AssetShareSearchDTO struct {
	_dto.Pagable
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
