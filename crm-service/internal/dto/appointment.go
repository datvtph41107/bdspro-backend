package dto

import (
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"

	"crm/internal/domain"
)

type AppointmentDTO struct {
	Appointment *domain.Appointment
	Deal        *bdspropb.GroupDeal
	Products    []*bdspropb.ProductAttachment
	// DealTransactions []*organizationpb.DealTransactionResponse
	Participants []*sharepb.ProfileItem
}
