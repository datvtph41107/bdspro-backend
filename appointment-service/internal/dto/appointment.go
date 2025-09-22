package dto

import (
	bdspropb "pb/types/bdspro"
	organizationpb "pb/types/organization"
	sharepb "pb/types/shared"

	"github.com/tuandq2112/appointment-service/internal/domain/entity"
)

type AppointmentDTO struct {
	Appointment *entity.Appointment
	Deal        *organizationpb.GroupDeal
	Products    []*bdspropb.ProductAttachment
	// DealTransactions []*organizationpb.DealTransactionResponse
	Participants []*sharepb.ProfileItem
}
