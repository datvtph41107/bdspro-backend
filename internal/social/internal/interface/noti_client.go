package iusecase

import (
	_enum "common/domain/enum"
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
	CreateToOwner(
		ctx context.Context,
		avatar string,
		title string,
		message []string,
		notificationType _enum.ENotificationType,
		targetId *uint64,
		ownerId uint64,
		ownerOf _enum.EOwnerOf,
		attachData []string,
	) error
}
