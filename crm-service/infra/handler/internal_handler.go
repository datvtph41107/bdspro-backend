package handler

import (
	_utils "common/utils"
	"context"
	"crm/infra/mapper"
	"crm/internal/dto"
	"crm/internal/repo"
	"crm/internal/usecase"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
	"strings"
)

type CrmInternalService struct {
	crmpb.UnimplementedCrmInternalServiceServer
	FriendRepo      repo.FriendRepo
	FriendMapper    *mapper.FriendMapper
	FriendUsecase   *usecase.FriendUsecase
	RegistryUsecase *usecase.RegistryUsecase
	RateUsecase     *usecase.RateUsecase
	ContactRepo     repo.ContactRepo
	ContactMapper   *mapper.ContactMapper
}

func NewCrmInternalService(
	friendRepo repo.FriendRepo,
	friendMapper *mapper.FriendMapper,
	friendUsecase *usecase.FriendUsecase,
	registryUsecase *usecase.RegistryUsecase,
	rateUsecase *usecase.RateUsecase,
	contactRepo repo.ContactRepo,
	contactMapper *mapper.ContactMapper,
) *CrmInternalService {
	return &CrmInternalService{
		FriendRepo:      friendRepo,
		FriendMapper:    friendMapper,
		FriendUsecase:   friendUsecase,
		RegistryUsecase: registryUsecase,
		RateUsecase:     rateUsecase,
		ContactRepo:     contactRepo,
		ContactMapper:   contactMapper,
	}
}

func (s *CrmInternalService) FriendStatusByProfileIds(ctx context.Context, req *crmpb.FriendStatusByProfileIdsRequest) (*crmpb.FriendStatusByProfileIdsResponse, error) {
	friendStatuses, err := s.FriendRepo.GetFriendStatusByProfileIds(ctx, req.CurrentProfileId, req.ProfileIds)
	if err != nil {
		return nil, err
	}

	friendStatusesPb := s.FriendMapper.FriendItemToPbs(friendStatuses)
	return &crmpb.FriendStatusByProfileIdsResponse{FriendStatus: friendStatusesPb}, nil
}

func (s *CrmInternalService) GetFriendInfoByProfileId(ctx context.Context, req *sharepb.IdRequest) (*sharepb.ContactInfoV3Proto, error) {
	result, err := s.FriendUsecase.FriendInfoByProfileId(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &sharepb.ContactInfoV3Proto{
		NumFollower:      result.NumFollower,
		NumFollowing:     result.NumFollowing,
		NumFriend:        result.NumFriend,
		ContactId:        result.ContactId,
		BlockId:          result.BlockId,
		FollowId:         result.FollowId,
		FriendCommons:    result.FriendCommons,
		NumFriendCommons: result.NumFriendCommons,
	}, nil
}

// RegisterEntity registers a new entity in the registry
func (h *CrmInternalService) RegisterEntity(ctx context.Context, req *crmpb.FeedbackRegistry) (*sharepb.SubmitResponse, error) {
	// Check if entity already exists
	existingEntity, err := h.RegistryUsecase.GetRegistryByName(ctx, req.RegistryName)
	if err == nil && existingEntity != nil {
		// Entity already exists, return success without creating new one
		return &sharepb.SubmitResponse{
			Message: "Entity already exists",
		}, nil
	}

	// Convert protobuf to DTO
	cleanName := strings.TrimSpace(strings.ToLower(req.RegistryName))
	createReq := &dto.RegistryCreateRequest{
		RegistryKey:  req.RegistryKey,
		RegistryName: cleanName,
		Port:         int(req.Port),
		Service:      req.Service,
	}

	// Create registry entry
	_, err = h.RegistryUsecase.CreateRegistry(ctx, createReq)
	if err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Message: "Entity registered successfully",
	}, nil
}

func (h *CrmInternalService) GetRateStats(ctx context.Context, req *sharepb.IdRequest) (*sharepb.RateStatsV3Proto, error) {
	rateStats, err := h.RateUsecase.GetRateStats(ctx, req.Id, "user")
	if err != nil {
		return nil, err
	}
	return &sharepb.RateStatsV3Proto{
		TotalReviews: uint64(rateStats.TotalReviews),
		AverageScore: float32(rateStats.AverageScore),
		Star1Count:   uint64(rateStats.Star1Count),
		Star2Count:   uint64(rateStats.Star2Count),
		Star3Count:   uint64(rateStats.Star3Count),
		Star4Count:   uint64(rateStats.Star4Count),
		Star5Count:   uint64(rateStats.Star5Count),
	}, nil
}

func (h *CrmInternalService) GetContactById(ctx context.Context, req *sharepb.IdRequest) (*sharepb.ContactDTO, error) {
	contact, err := h.ContactRepo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	result := &sharepb.ContactDTO{
		Id:        contact.ID,
		FullName:  contact.FullName,
		Phone:     contact.Phone,
		Email:     contact.Email,
		Zalo:      contact.Zalo,
		Company:   contact.Company,
		Note:      contact.Note,
		Birthday:  _utils.FormatTimeToString(contact.Birthday),
		OwnerId:   contact.OwnerID,
		OwnerOf:   int32(contact.OwnerOf),
		ProfileId: contact.ProfileID,
		CreatedAt: _utils.FormatTimeToString(contact.CreatedAt),
		CreatedBy: contact.CreatedBy,
	}
	if contact.FriendStatus != nil {
		result.FriendStatus = &sharepb.FriendItem{
			Id:         contact.FriendStatus.ID,
			Status:     uint32(contact.FriendStatus.Status),
			CreatedBy:  contact.FriendStatus.CreatedBy,
			ReceiverId: contact.FriendStatus.ReceiverID,
		}
	}

	return result, nil
}
