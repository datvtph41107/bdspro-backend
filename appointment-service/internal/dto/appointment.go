package dto

import (
	pb_bdspro "pb/types/bdspro"
	organizationpb "pb/types/organization"

	"github.com/tuandq2112/appointment-service/internal/domain/entity"
)

type AppointmentDTO struct {
	Appointment *entity.Appointment
	Deal *organizationpb.GroupDeal
	Products []*pb_bdspro.ProductAttachment
	DealTransactions []*organizationpb.DealTransactionResponse
	Participants []*organizationpb.DealInvitationMember
}