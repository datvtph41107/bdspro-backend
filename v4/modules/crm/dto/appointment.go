package dto

import (
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
	"time"

	domain "userv4/modules/crm/entity"
)

type AppointmentDTO struct {
	Appointment      *domain.Appointment
	Deal             *bdspropb.GroupDeal
	Products         []*bdspropb.ProductAttachment
	Participants     []*sharepb.ProfileItem
	DealContract     *bdspropb.DealContractItem
	ProductUpdatedAt *time.Time `gorm:"column:product_updated_at;->;-:migration" json:"productUpdatedAt"`
}