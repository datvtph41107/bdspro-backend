package read

import (
	"fmt"
	"strconv"

	"user/internal/models"
	"user/internal/modules/subscription/domain"
)

func projectAggregate(
	legacy models.Membership,
	aggregate domain.Aggregate,
) (*models.Membership, error) {
	if aggregate.SubjectKind != models.SubscriptionSubjectProfile ||
		aggregate.SubjectID != strconv.FormatUint(legacy.ProfileID, 10) {
		return nil, fmt.Errorf("subscription subject does not match legacy membership")
	}
	if aggregate.CreatedAt.IsZero() || aggregate.UpdatedAt.IsZero() {
		return nil, fmt.Errorf("subscription audit timestamps are missing")
	}
	result := copyMembership(&legacy)
	result.PaidAt = copyTime(&aggregate.StartedAt)
	result.FinishedDate = copyTime(&aggregate.CurrentPeriodEnd)
	result.CreatedAt = copyTime(&aggregate.CreatedAt)
	result.UpdatedAt = copyTime(&aggregate.UpdatedAt)
	return result, nil
}
