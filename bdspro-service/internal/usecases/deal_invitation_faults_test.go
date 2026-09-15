package usecases

import (
	"errors"
	"testing"

	"bdspro/internal/domain"
	"common/fault"
)

func TestDealInvitationAlreadyExistsFault(t *testing.T) {
	err := dealInvitationAlreadyExistsFault()

	if !errors.Is(
		err,
		domain.ErrDealInvitationAlreadyExists,
	) {
		t.Fatalf(
			"err = %v, want ErrDealInvitationAlreadyExists",
			err,
		)
	}

	failure, ok := fault.As(err)
	if !ok {
		t.Fatalf(
			"err = %T %v, want canonical fault",
			err,
			err,
		)
	}

	if failure.Kind() != fault.KindConflict {
		t.Fatalf(
			"kind = %q, want %q",
			failure.Kind(),
			fault.KindConflict,
		)
	}

	if failure.Code() != DealInvitationAlreadyExistsCode {
		t.Fatalf(
			"code = %q, want %q",
			failure.Code(),
			DealInvitationAlreadyExistsCode,
		)
	}

	if failure.PublicMessage() !=
		dealInvitationAlreadyExistsMessage {
		t.Fatalf(
			"message = %q, want %q",
			failure.PublicMessage(),
			dealInvitationAlreadyExistsMessage,
		)
	}
}
