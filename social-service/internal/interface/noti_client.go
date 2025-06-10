package iusecase

import (
	"context"
)

type NotiClient interface {
	NotifyToUser(
		ctx context.Context,
		typeNoti int32,
		userId *uint64,
		attachData []string,
		newsFeedOrCommentID uint64,
		isMerge bool,
	) error
}
