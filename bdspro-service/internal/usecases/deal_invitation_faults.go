package usecases

import (
	"bdspro/internal/domain"
	"common/fault"
)

// DealInvitationAlreadyExistsCode is the stable application identity for an
// invitation conflict. Protocol mappings consume this identity; messages do
// not define the business semantics.
const DealInvitationAlreadyExistsCode = "bdspro.deal_invitation.already_exists"

const dealInvitationAlreadyExistsMessage = "invitation already exists for this user"

func dealInvitationAlreadyExistsFault() error {
	return fault.Wrap(
		domain.ErrDealInvitationAlreadyExists,
		fault.KindConflict,
		DealInvitationAlreadyExistsCode,
		dealInvitationAlreadyExistsMessage,
	)
}
