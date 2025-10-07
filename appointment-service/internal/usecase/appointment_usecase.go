package usecase

import (
	"context"
	bdspropb "pb/types/bdspro"
	organizationpb "pb/types/organization"

	sharepb "pb/types/shared"

	"appointment/internal/domain/entity"
	"appointment/internal/domain/repository"
	"appointment/internal/dto"
	"appointment/internal/infrastructure/client"
)

type AppointmentUsecase interface {
	CreateAppointment(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error)
	GetAppointment(ctx context.Context, id uint32) (*dto.AppointmentDTO, error)
	UpdateAppointment(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error)
	DeleteAppointment(ctx context.Context, id uint32) error
	ListAppointments(ctx context.Context, filter map[string]any, page, size int) ([]*dto.AppointmentDTO, int64, error)
	CountCurrent(ctx context.Context) (int64, error)
}

type appointmentUsecase struct {
	appointmentRepo    repository.AppointmentRepository
	organizationClient client.OrganizationClient
	bdsproClient       client.BdsproClient
	userClient         *client.UserClient
}

func NewAppointmentUsecase(
	appointmentRepo repository.AppointmentRepository,
	organizationClient client.OrganizationClient,
	bdsproClient client.BdsproClient,
	userClient *client.UserClient,
) AppointmentUsecase {
	return &appointmentUsecase{
		appointmentRepo:    appointmentRepo,
		organizationClient: organizationClient,
		bdsproClient:       bdsproClient,
		userClient:         userClient,
	}
}

func (uc *appointmentUsecase) CreateAppointment(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error) {
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

func (uc *appointmentUsecase) UpdateAppointment(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error) {
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
	appointmentsDTO := make([]*dto.AppointmentDTO, 0)
	for _, appointment := range appointments {
		appointmentDTO, _ := uc.mapAppointmentDTO(ctx, appointment)
		// if err != nil {
		// 	return nil, 0, err
		// }
		appointmentsDTO = append(appointmentsDTO, appointmentDTO)
	}
	return appointmentsDTO, total, nil
}

func (uc *appointmentUsecase) mapAppointmentDTO(ctx context.Context, appointment *entity.Appointment) (*dto.AppointmentDTO, error) {
	// Convert IDs once to avoid repeated conversions
	dealId := uint64(appointment.DealId)
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
		deal     *organizationpb.GroupDeal
		products *bdspropb.ProductAttachmentResponse
		// dealTransactions []*organizationpb.DealTransactionResponse
		participants []*sharepb.ProfileItem
		err          error
	}

	ch := make(chan result, 4)

	// Concurrent API calls
	go func() {
		deal, err := uc.organizationClient.GetDealById(ctx, dealId)
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
		ch <- result{participants: participants, err: err}
	}()

	// Collect results
	var deal *organizationpb.GroupDeal
	var products *bdspropb.ProductAttachmentResponse
	// var dealTransactions []*organizationpb.DealTransactionResponse
	var participants []*sharepb.ProfileItem

	for i := 0; i < 4; i++ {
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
		if res.participants != nil {
			participants = res.participants
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

	return out, nil
}

func (uc *appointmentUsecase) CountCurrent(ctx context.Context) (int64, error) {
	count, err := uc.appointmentRepo.CountCurrent(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}
