package usecase

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	"fmt"
	bdspropb "pb/types/bdspro"

	sharepb "pb/types/shared"

	"crm/infra/client"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/interface/provider"
	"crm/internal/repo"
)

type AppointmentUsecase interface {
	CreateAppointment(ctx context.Context, appointment *domain.Appointment) (*domain.Appointment, error)
	GetAppointment(ctx context.Context, id uint32) (*dto.AppointmentDTO, error)
	UpdateAppointment(ctx context.Context, appointment *domain.Appointment) (*domain.Appointment, error)
	DeleteAppointment(ctx context.Context, id uint32) error
	ListAppointments(ctx context.Context, filter map[string]any, pagable _dto.Pagable) ([]*dto.AppointmentDTO, int64, error)
	CountCurrent(ctx context.Context) (int64, error)
	GetAppointmentsByProduct(ctx context.Context, productId uint64, pagable _dto.Pagable) ([]*dto.AppointmentDTO, int64, error)
	GetAppointmentsByDeal(ctx context.Context, dealId uint32, page, size int) ([]*dto.AppointmentDTO, int64, error)
	UpdateAppointmentStatus(ctx context.Context, id uint32, action string) (*domain.Appointment, error)
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

// uniqueUserIds loại bỏ các userId trùng lặp và trả về danh sách unique
func (uc *appointmentUsecase) uniqueUserIds(userIds []uint64) []uint64 {
	seen := make(map[uint64]bool)
	unique := make([]uint64, 0)

	for _, id := range userIds {
		if id > 0 && !seen[id] {
			seen[id] = true
			unique = append(unique, id)
		}
	}

	return unique
}

func (uc *appointmentUsecase) CreateAppointment(ctx context.Context, appointment *domain.Appointment) (*domain.Appointment, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return nil, fmt.Errorf("profileId not found in context")
	}

	// Lấy danh sách participantIds từ appointment
	participantIds := make([]uint64, 0)
	if appointment.ParticipantIds != nil {
		for _, id := range appointment.ParticipantIds {
			if id > 0 {
				participantIds = append(participantIds, uint64(id))
			}
		}
	}

	// Thêm userId của user request vào danh sách
	participantIds = append(participantIds, profileId)

	// Loại bỏ trùng lặp và unique danh sách
	participantIds = uc.uniqueUserIds(participantIds)

	// Tạo lại AppointmentParticipants từ danh sách userId đã unique
	appointmentParticipants := make([]domain.AppointmentParticipant, 0, len(participantIds))
	for _, userId := range participantIds {
		role := domain.AppointmentParticipantRoleGuest
		status := domain.AppointmentParticipantStatusPending

		// User request sẽ có role "host" và status "confirmed"
		if userId == profileId {
			role = domain.AppointmentParticipantRoleHost
			status = domain.AppointmentParticipantStatusConfirmed
		}

		appointmentParticipants = append(appointmentParticipants, domain.AppointmentParticipant{
			UserId:    uint32(userId),
			ContactID: nil, // Không dùng contact nữa
			Role:      role,
			Status:    status,
		})
	}

	// Cập nhật lại ParticipantIds và AppointmentParticipants
	// Convert []uint64 sang pq.Int64Array
	participantIdsInt64 := make([]int64, len(participantIds))
	for i, id := range participantIds {
		participantIdsInt64[i] = int64(id)
	}
	appointment.ParticipantIds = participantIdsInt64
	appointment.AppointmentParticipants = appointmentParticipants

	appointment, err := uc.appointmentRepo.Create(ctx, appointment)
	if err != nil {
		return nil, err
	}

	// Cập nhật ScheduleCount cho tất cả products trong appointment
	uc.updateScheduleCountsForProducts(ctx, appointment.ProductIds)

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
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return nil, fmt.Errorf("profileId not found in context")
	}

	// Lấy danh sách participantIds từ appointment
	participantIds := make([]uint64, 0)
	if appointment.ParticipantIds != nil {
		for _, id := range appointment.ParticipantIds {
			if id > 0 {
				participantIds = append(participantIds, uint64(id))
			}
		}
	}

	// Thêm userId của user request vào danh sách
	participantIds = append(participantIds, profileId)

	// Loại bỏ trùng lặp và unique danh sách
	participantIds = uc.uniqueUserIds(participantIds)

	// Tạo lại AppointmentParticipants từ danh sách userId đã unique
	appointmentParticipants := make([]domain.AppointmentParticipant, 0, len(participantIds))
	for _, userId := range participantIds {
		role := domain.AppointmentParticipantRoleGuest
		status := domain.AppointmentParticipantStatusPending

		// User request sẽ có role "host" và status "confirmed"
		if userId == profileId {
			role = domain.AppointmentParticipantRoleHost
			status = domain.AppointmentParticipantStatusConfirmed
		}

		appointmentParticipants = append(appointmentParticipants, domain.AppointmentParticipant{
			UserId:    uint32(userId),
			ContactID: nil, // Không dùng contact nữa
			Role:      role,
			Status:    status,
		})
	}

	// Cập nhật lại ParticipantIds và AppointmentParticipants
	// Convert []uint64 sang pq.Int64Array
	participantIdsInt64 := make([]int64, len(participantIds))
	for i, id := range participantIds {
		participantIdsInt64[i] = int64(id)
	}
	appointment.ParticipantIds = participantIdsInt64
	appointment.AppointmentParticipants = appointmentParticipants

	appointment, err := uc.appointmentRepo.Update(ctx, appointment)
	if err != nil {
		return nil, err
	}

	// Cập nhật ScheduleCount cho tất cả products trong appointment
	uc.updateScheduleCountsForProducts(ctx, appointment.ProductIds)

	return appointment, nil
}

func (uc *appointmentUsecase) DeleteAppointment(ctx context.Context, id uint32) error {
	// Lấy appointment trước khi xóa để biết ProductIds
	appointment, err := uc.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = uc.appointmentRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Cập nhật ScheduleCount cho tất cả products trong appointment
	uc.updateScheduleCountsForProducts(ctx, appointment.ProductIds)

	return nil
}

func (uc *appointmentUsecase) ListAppointments(ctx context.Context, filter map[string]any, pagable _dto.Pagable) ([]*dto.AppointmentDTO, int64, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return nil, 0, fmt.Errorf("profileId not found in context")
	}
	appointments, total, err := uc.appointmentRepo.List(ctx, profileId, filter, pagable)
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
			if participant.UserId > 0 {
				participantIdsSet[uint64(participant.UserId)] = struct{}{}
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
				if participant.UserId > 0 {
					if participantItem, ok := participantsMap[uint64(participant.UserId)]; ok {
						participants = append(participants, participantItem)
					} else if participant.Contact != nil {
						participants = append(participants, &sharepb.ProfileItem{
							Id:       uint64(participant.UserId),
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
	participantIds := make([]uint64, len(appointment.AppointmentParticipants))

	for i, id := range appointment.ProductIds {
		productIds[i] = uint64(id)
	}
	for i, participant := range appointment.AppointmentParticipants {
		participantIds[i] = uint64(participant.UserId)
	}

	// Make concurrent API calls for better performance
	type result struct {
		deal            *bdspropb.GroupDeal
		products        *bdspropb.ProductAttachmentResponse
		dealContracts   []*bdspropb.DealContractItem
		participantsMap map[uint64]*sharepb.ProfileItem
		err             error
	}

	// Determine number of goroutines
	numGoroutines := 3
	if appointment.TransactionId > 0 {
		numGoroutines = 4
	}
	ch := make(chan result, numGoroutines)

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

	// Get DealContract if TransactionId exists
	if appointment.TransactionId > 0 {
		go func() {
			transactionIds := []uint64{appointment.TransactionId}
			dealContracts, err := uc.bdsproClient.GetDealContractsByIds(ctx, transactionIds)
			ch <- result{dealContracts: dealContracts, err: err}
		}()
	}

	// Collect results
	var deal *bdspropb.GroupDeal
	var products *bdspropb.ProductAttachmentResponse
	var dealContract *bdspropb.DealContractItem
	var participants []*sharepb.ProfileItem

	for i := 0; i < numGoroutines; i++ {
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
		if len(res.dealContracts) > 0 {
			dealContract = res.dealContracts[0]
		}
		if res.participantsMap != nil {
			for _, participant := range appointment.AppointmentParticipants {
				if participant.UserId != 0 {
					if participantItem, ok := res.participantsMap[uint64(participant.UserId)]; ok && participantItem != nil {
						participants = append(participants, participantItem)
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
		Appointment:  appointment,
		Deal:         deal,
		Participants: participantsArray,
		DealContract: dealContract,
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

func (uc *appointmentUsecase) GetAppointmentsByProduct(ctx context.Context, productId uint64, pagable _dto.Pagable) ([]*dto.AppointmentDTO, int64, error) {
	appointments, total, err := uc.appointmentRepo.GetByProductId(ctx, productId, int(pagable.GetPage()), int(pagable.GetSize()))
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

func (uc *appointmentUsecase) UpdateAppointmentStatus(ctx context.Context, id uint32, action string) (*domain.Appointment, error) {
	// Lấy appointment hiện tại
	appointment, err := uc.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("appointment not found: %v", err)
	}

	// Xác định status mới dựa trên action
	var newStatus enums.EAppointmentStatus
	switch action {
	case "approve":
		newStatus = enums.AppointmentStatusConfirmed
	case "cancel":
		newStatus = enums.AppointmentStatusCancelled
	default:
		return nil, fmt.Errorf("invalid action: %s. Action must be 'approve' or 'cancel'", action)
	}

	// Kiểm tra status hiện tại có thể chuyển sang status mới không
	if appointment.Status != enums.AppointmentStatusPending {
		return nil, fmt.Errorf("appointment status must be pending to update. Current status: %d", appointment.Status)
	}

	// Update status
	err = uc.appointmentRepo.UpdateStatus(ctx, id, newStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to update appointment status: %v", err)
	}

	// Lấy lại appointment đã update
	updatedAppointment, err := uc.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated appointment: %v", err)
	}

	return updatedAppointment, nil
}

// updateScheduleCountsForProducts cập nhật ScheduleCount cho tất cả products
func (uc *appointmentUsecase) updateScheduleCountsForProducts(ctx context.Context, productIds []int64) {
	// Cập nhật ScheduleCount cho mỗi product
	for _, productID := range productIds {
		if productID <= 0 {
			continue
		}
		// Đếm số appointment của product
		count, err := uc.appointmentRepo.CountAppointmentsByProductID(ctx, uint64(productID))
		if err != nil {
			// Log error nhưng không fail toàn bộ operation
			continue
		}

		// Gọi bdspro-service để cập nhật ScheduleCount
		err = uc.bdsproClient.UpdateScheduleCount(ctx, uint64(productID), count)
		if err != nil {
			// Log error nhưng không fail toàn bộ operation
			continue
		}
	}
}