package dto

import (
	bdspropb "pb/types/bdspro"
	organizationpb "pb/types/organization"
	sharepb "pb/types/shared"

	"appointment/internal/domain/entity"
)

type AppointmentDTO struct {
	Appointment *entity.Appointment
	Deal        *organizationpb.GroupDeal
	Products    []*bdspropb.ProductAttachment
	// DealTransactions []*organizationpb.DealTransactionResponse
	Participants []*sharepb.ProfileItem
}
