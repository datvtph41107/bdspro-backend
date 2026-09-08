package evaluate

import (
	"common/identity"
	"errors"
	"strconv"

	"user/internal/domain/entitlement"
)

var (
	// ErrActorProfileMissing cho biết request không có profile identity đã xác thực.
	ErrActorProfileMissing = errors.New("actor profile identity is required")
)

// SubjectFromActor là owner duy nhất chuyển identity transport thành billing
// subject. Actor vẫn là người thực hiện (profile); Subject là pool trả tiền và
// tính quota (profile hoặc organization). Report/Job có thể tiếp tục sở hữu
// theo profile, còn usage event dùng subject kết quả này.
func SubjectFromActor(actor identity.Actor) (access.Subject, error) {
	if actor.ProfileID == 0 {
		return access.Subject{}, ErrActorProfileMissing
	}
	if actor.OrganizationID != nil {
		// OrganizationID chỉ được gateway lấy từ token do User Service
		// phát sau khi Organization Service xác minh membership active. Actor.Role
		// là global account role (ví dụ ROLE_USER), không phải organization
		// role; vì vậy không dùng nó làm billing evidence giả.
		return access.Subject{
			Type: access.SubjectOrganization,
			ID:   strconv.FormatUint(*actor.OrganizationID, 10),
		}, nil
	}
	return access.Subject{
		Type: access.SubjectProfile,
		ID:   strconv.FormatUint(actor.ProfileID, 10),
	}, nil
}
