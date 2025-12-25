package usecase

import (
	base_enum "base/enum"
	_utils "common/utils"
	"context"
	"fmt"
	bdspropb "pb/types/bdspro"

	sharepb "pb/types/shared"

	"crm/infra/client"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/interface/provider"
	"crm/internal/repo"
)

type AppointmentUsecase interface {
	CreateAppointment(ctx context.Context, appointment *domain.Appointment) (*domain.Appointment, error)
	GetAppointment(ctx context.Context, id uint32) (*dto.AppointmentDTO, error)
	UpdateAppointment(ctx context.Context, appointment *domain.Appointment) (*domain.Appointment, error)
	DeleteAppointment(ctx context.Context, id uint32) error
	ListAppointments(ctx context.Context, filter map[string]any, page, size int) ([]*dto.AppointmentDTO, int64, error)
	CountCurrent(ctx context.Context) (int64, error)
	GetAppointmentsByProduct(ctx context.Context, productId uint64, page, size int) ([]*dto.AppointmentDTO, int64, error)
	GetAppointmentsByDeal(ctx context.Context, dealId uint32, page, size int) ([]*dto.AppointmentDTO, int64, error)
}

type appointmentUsecase struct {
	appointmentRepo repo.AppointmentRepository
	contactRepo     repo.ContactRepo
	// organizationClient client.OrganizationClient
	bdsproClient provider.BdsproProvider
	userClient   *client.UserClient
}

func NewAppointmentUsecase(
	appointmentRepo repo.AppointmentRepository,
	contactRepo repo.ContactRepo,
	// organizationClient client.OrganizationClient,
	bdsproClient provider.BdsproProvider,
	userClient *client.UserClient,
) AppointmentUsecase {
	return &appointmentUsecase{
		appointmentRepo: appointmentRepo,
		contactRepo:     contactRepo,
		// organizationClient: organizationClient,
		bdsproClient: bdsproClient,
		userClient:   userClient,
	}
}

func (uc *appointmentUsecase) CreateAppointment(ctx context.Context, appointment *domain.Appointment) (*domain.Appointment, error) {
	// Map ParticipantIds (profileId) -> ContactID
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return nil, fmt.Errorf("profileId not found in context")
	}

	// Map participantIds (profileId) to contactIds
	participantContactMap := make(map[int64]uint64) // map[profileId]contactId
	for _, participantProfileId := range appointment.ParticipantIds {
		// Lấy hoặc tạo contact từ profileId
		contact, err := uc.contactRepo.GetOrCreateByProfileID(
			ctx,
			uint64(participantProfileId),
			profileId,
			base_enum.EOwnerOfMember,
			"", // fullName - sẽ được update sau
			"", // phone - sẽ được update sau
			"", // avatar - sẽ được update sau
		)
		if err != nil {
			// Log error nhưng tiếp tục với các participant khác
			continue
		}
		if contact != nil && contact.ID > 0 {
			participantContactMap[participantProfileId] = contact.ID
		}
	}

	// Lưu map vào context để repository sử dụng
	ctx = context.WithValue(ctx, "participantContactMap", participantContactMap)

	appointment, err := uc.appointmentRepo.Create(ctx, appointment)
	if err != nil {
		return nil, err
	}
	return appointment, nil
}

func (uc *appointmentUsecase) GetAppointment(ctx context.Context, id uint32) (*dto.AppointmentDTO, error) {
	appointment, err := uc.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.mapAppointmentDTO(ctx, appointment)
}

func (uc *appointmentUsecase) UpdateAppointment(ctx context.Context, appointment *domain.Appointment) (*domain.Appointment, error) {
	// Map ParticipantIds (profileId) -> ContactID
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return nil, fmt.Errorf("profileId not found in context")
	}

	// Map participantIds (profileId) to contactIds
	participantContactMap := make(map[int64]uint64) // map[profileId]contactId
	for _, participantProfileId := range appointment.ParticipantIds {
		// Lấy hoặc tạo contact từ profileId
		contact, err := uc.contactRepo.GetOrCreateByProfileID(
			ctx,
			uint64(participantProfileId),
			profileId,
			base_enum.EOwnerOfMember,
			"", // fullName - sẽ được update sau
			"", // phone - sẽ được update sau
			"", // avatar - sẽ được update sau
		)
		if err != nil {
			// Log error nhưng tiếp tục với các participant khác
			continue
		}
		if contact != nil && contact.ID > 0 {
			participantContactMap[participantProfileId] = contact.ID
		}
	}

	// Lưu map vào context để repository sử dụng
	ctx = context.WithValue(ctx, "participantContactMap", participantContactMap)

	appointment, err := uc.appointmentRepo.Update(ctx, appointment)
	if err != nil {
		return nil, err
	}
	return appointment, nil
}

func (uc *appointmentUsecase) DeleteAppointment(ctx context.Context, id uint32) error {
	err := uc.appointmentRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (uc *appointmentUsecase) ListAppointments(ctx context.Context, filter map[string]any, page, size int) ([]*dto.AppointmentDTO, int64, error) {
	appointments, total, err := uc.appointmentRepo.List(ctx, filter, page, size)
	if err != nil {
		return nil, 0, err
	}
	appointmentsDTO, err := uc.mapAppointmentDTOs(ctx, appointments)
	if err != nil {
		return nil, 0, err
	}

	// for _, appointment := range appointments {
	// 	appointmentDTO, _ := uc.mapAppointmentDTO(ctx, appointment)
	// 	// if err != nil {
	// 	// 	return nil, 0, err
	// 	// }
	// 	appointmentsDTO = append(appointmentsDTO, appointmentDTO)
	// }
	return appointmentsDTO, total, nil
}

func (uc *appointmentUsecase) mapAppointmentDTOs(ctx context.Context, appointments []*domain.Appointment) ([]*dto.AppointmentDTO, error) {
	if len(appointments) == 0 {
		return []*dto.AppointmentDTO{}, nil
	}

	// Collect all unique IDs
	productIdsSet := make(map[uint64]struct{})
	participantIdsSet := make(map[uint64]struct{})
	dealIdsSet := make(map[uint64]struct{})

	for _, appointment := range appointments {
		for _, id := range appointment.ProductIds {
			productIdsSet[uint64(id)] = struct{}{}
		}
		for _, participant := range appointment.AppointmentParticipants {
			if participant.Contact != nil && participant.Contact.ProfileID != nil {
				participantIdsSet[*participant.Contact.ProfileID] = struct{}{}
			}
		}
		if appointment.DealId > 0 {
			dealIdsSet[uint64(appointment.DealId)] = struct{}{}
		}
	}

	// Convert sets to slices
	productIds := make([]uint64, 0, len(productIdsSet))
	for id := range productIdsSet {
		productIds = append(productIds, id)
	}

	participantIds := make([]uint64, 0, len(participantIdsSet))
	for id := range participantIdsSet {
		participantIds = append(participantIds, id)
	}

	dealIds := make([]uint64, 0, len(dealIdsSet))
	for id := range dealIdsSet {
		dealIds = append(dealIds, id)
	}

	// Concurrent API calls
	type result struct {
		deals        map[uint64]*bdspropb.GroupDeal
		products     map[uint64]*bdspropb.ProductAttachment
		participants map[uint64]*sharepb.ProfileItem
		err          error
	}

	ch := make(chan result, 3)

	// Fetch deals (gọi song song cho từng deal)
	go func() {
		dealsMap := make(map[uint64]*bdspropb.GroupDeal)
		deals, err := uc.bdsproClient.GetDealsByIds(ctx, dealIds)
		if err == nil && deals != nil {
			for _, deal := range deals {
				if deal != nil && deal.Id > 0 {
					dealsMap[uint64(deal.Id)] = deal
				}
			}
		}
		ch <- result{deals: dealsMap, err: nil}
	}()

	// Fetch products
	go func() {
		productsResp, err := uc.bdsproClient.GetProductAttachmentByIds(ctx, productIds)
		productsMap := make(map[uint64]*bdspropb.ProductAttachment)
		if err == nil && productsResp != nil && productsResp.Data != nil {
			for _, product := range productsResp.Data {
				if product != nil && product.Id > 0 {
					productsMap[product.Id] = product
				}
			}
		}
		ch <- result{products: productsMap, err: err}
	}()

	// Fetch participants
	go func() {
		participants, err := uc.userClient.GetProfileByIds(ctx, participantIds)
		participantsMap := make(map[uint64]*sharepb.ProfileItem)
		if err == nil && participants != nil {
			for _, participant := range participants {
				if participant != nil && participant.Id > 0 {
					participantsMap[participant.Id] = participant
				}
			}
		}
		ch <- result{participants: participantsMap, err: err}
	}()

	// Collect results
	var dealsMap map[uint64]*bdspropb.GroupDeal
	var productsMap map[uint64]*bdspropb.ProductAttachment
	var participantsMap map[uint64]*sharepb.ProfileItem

	for i := 0; i < 3; i++ {
		res := <-ch
		if res.deals != nil {
			dealsMap = res.deals
		}
		if res.products != nil {
			productsMap = res.products
		}
		if res.participants != nil {
			participantsMap = res.participants
		}
	}

	// Map appointments to DTOs
	appointmentsDTO := make([]*dto.AppointmentDTO, 0, len(appointments))
	for _, appointment := range appointments {
		// Get deal for this appointment
		var deal *bdspropb.GroupDeal
		if appointment.DealId > 0 && dealsMap != nil {
			deal = dealsMap[uint64(appointment.DealId)]
		}

		// Get products for this appointment
		products := make([]*bdspropb.ProductAttachment, 0)
		if productsMap != nil {
			for _, productId := range appointment.ProductIds {
				if product, ok := productsMap[uint64(productId)]; ok {
					products = append(products, product)
				}
			}
		}

		// Get participants for this appointment
		participants := make([]*sharepb.ProfileItem, 0)
		if participantsMap != nil {
			for _, participant := range appointment.AppointmentParticipants {
				if participant.Contact != nil && participant.Contact.ProfileID != nil {
					if participantItem, ok := participantsMap[*participant.Contact.ProfileID]; ok {
						participants = append(participants, participantItem)
					} else {
						participants = append(participants, &sharepb.ProfileItem{
							Id:       *participant.Contact.ProfileID,
							FullName: participant.Contact.FullName,
							Avatar:   participant.Contact.Avatar,
							Phone:    participant.Contact.Phone,
						})
					}
				}
			}
		}

		appointmentDTO := &dto.AppointmentDTO{
			Appointment:  appointment,
			Deal:         deal,
			Products:     products,
			Participants: participants,
		}

		appointmentsDTO = append(appointmentsDTO, appointmentDTO)
	}

	return appointmentsDTO, nil
}

func (uc *appointmentUsecase) mapAppointmentDTO(ctx context.Context, appointment *domain.Appointment) (*dto.AppointmentDTO, error) {
	// Convert IDs once to avoid repeated conversions
	dealId := uint64(appointment.DealId)
	fmt.Println("appointment.ProductIds", appointment.ProductIds)
	productIds := make([]uint64, len(appointment.ProductIds))
	participantIds := make([]uint64, len(appointment.ParticipantIds))

	for i, id := range appointment.ProductIds {
		productIds[i] = uint64(id)
	}
	for i, id := range appointment.ParticipantIds {
		participantIds[i] = uint64(id)
	}

	// Make concurrent API calls for better performance
	type result struct {
		deal     *bdspropb.GroupDeal
		products *bdspropb.ProductAttachmentResponse
		// dealTransactions []*bdspropb.DealTransactionResponse
		participantsMap map[uint64]*sharepb.ProfileItem
		err             error
	}

	ch := make(chan result, 3)

	// Concurrent API calls
	go func() {
		deal, err := uc.bdsproClient.GetDealById(ctx, dealId)
		ch <- result{deal: deal, err: err}
	}()

	go func() {
		products, err := uc.bdsproClient.GetProductAttachmentByIds(ctx, productIds)
		ch <- result{products: products, err: err}
	}()

	// go func() {
	// 	dealTransactions, err := uc.organizationClient.GetDealTransactionsByIds(ctx, []uint64{dealId})
	// 	ch <- result{dealTransactions: dealTransactions, err: err}
	// }()

	go func() {
		participants, err := uc.userClient.GetProfileByIds(ctx, participantIds)
		participantsMap := make(map[uint64]*sharepb.ProfileItem)
		if err == nil && participants != nil {
			for _, participant := range participants {
				if participant != nil && participant.Id > 0 {
					participantsMap[participant.Id] = participant
				}
			}
		}
		ch <- result{participantsMap: participantsMap, err: err}
	}()

	// Collect results
	var deal *bdspropb.GroupDeal
	var products *bdspropb.ProductAttachmentResponse
	// var dealTransactions []*bdspropb.DealTransactionResponse
	var participants []*sharepb.ProfileItem

	for i := 0; i < 3; i++ {
		res := <-ch
		// if res.err != nil {
		// 	return nil, res.err
		// }
		if res.deal != nil {
			deal = res.deal
		}
		if res.products != nil {
			products = res.products
		}
		// if res.dealTransactions != nil {
		// 	dealTransactions = res.dealTransactions
		// }
		if res.participantsMap != nil {
			for _, participant := range appointment.AppointmentParticipants {
				if participant.Contact != nil && participant.Contact.ProfileID != nil {
					if participantItem, ok := res.participantsMap[*participant.Contact.ProfileID]; ok {
						participants = append(participants, participantItem)
					} else {
						participants = append(participants, &sharepb.ProfileItem{
							Id:       *participant.Contact.ProfileID,
							FullName: participant.Contact.FullName,
							Avatar:   participant.Contact.Avatar,
							Phone:    participant.Contact.Phone,
						})
					}
				}
			}
		}
	}

	// Build participants array efficiently
	participantsArray := make([]*sharepb.ProfileItem, 0, len(appointment.ParticipantIds))
	for _, participant := range participants {
		participantsArray = append(participantsArray, participant)
	}

	out := &dto.AppointmentDTO{
		Appointment: appointment,
		Deal:        deal,
		// DealTransactions: dealTransactions,
		Participants: participantsArray,
	}

	if products != nil {
		out.Products = products.Data
	}
	fmt.Println("out.Products", out.Products)
	return out, nil
}

func (uc *appointmentUsecase) CountCurrent(ctx context.Context) (int64, error) {
	count, err := uc.appointmentRepo.CountCurrent(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (uc *appointmentUsecase) GetAppointmentsByProduct(ctx context.Context, productId uint64, page, size int) ([]*dto.AppointmentDTO, int64, error) {
	appointments, total, err := uc.appointmentRepo.GetByProductId(ctx, productId, page, size)
	if err != nil {
		return nil, 0, err
	}
	appointmentsDTO, err := uc.mapAppointmentDTOs(ctx, appointments)
	if err != nil {
		return nil, 0, err
	}
	// appointmentsDTO := make([]*dto.AppointmentDTO, 0)
	// for _, appointment := range appointments {
	// 	appointmentDTO, _ := uc.mapAppointmentDTO(ctx, appointment)
	// 	appointmentsDTO = append(appointmentsDTO, appointmentDTO)
	// }
	return appointmentsDTO, total, nil
}

func (uc *appointmentUsecase) GetAppointmentsByDeal(ctx context.Context, dealId uint32, page, size int) ([]*dto.AppointmentDTO, int64, error) {
	appointments, total, err := uc.appointmentRepo.GetByDealId(ctx, dealId, page, size)
	if err != nil {
		return nil, 0, err
	}
	appointmentsDTO := make([]*dto.AppointmentDTO, 0)
	for _, appointment := range appointments {
		appointmentDTO, _ := uc.mapAppointmentDTO(ctx, appointment)
		appointmentsDTO = append(appointmentsDTO, appointmentDTO)
	}
	return appointmentsDTO, total, nil
}
