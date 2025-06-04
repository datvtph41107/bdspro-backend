package repo

import (
	"bdspro/internal/common"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
)

type AssetShareRepo interface {
	common.IBaseRepo[domain.AssetShare, *dto.AssetShareSearchDTO]
}
