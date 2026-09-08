package repo

import (
	"bdspro/internal/common"
	"bdspro/internal/domain"
)

type ActionRepo interface {
	common.IBaseRepo[domain.Action, common.DTO]
}
