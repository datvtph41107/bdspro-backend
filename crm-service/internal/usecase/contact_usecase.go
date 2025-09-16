package usecase

import (
	base_enum "base/enum"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/interface/provider"
	"crm/internal/repo"
	"strings"
)

type ContactUsecase struct {
	repo         repo.ContactRepo
	friendRepo   repo.FriendRepo
	followRepo   repo.FollowRepo
	leadUsecase  *LeadUsecase // forward reference
	userClient   provider.UserClient
	ownerUsecase *OwnerUsecase
}

func NewContactUsecase(
	repo repo.ContactRepo,
	friendRepo repo.FriendRepo,
	followRepo repo.FollowRepo,
	userClient provider.UserClient,
	ownerUsecase *OwnerUsecase,
) *ContactUsecase {
	return &ContactUsecase{
		repo:         repo,
		friendRepo:   friendRepo,
		followRepo:   followRepo,
		userClient:   userClient,
		ownerUsecase: ownerUsecase,
	}
}

func (s *ContactUsecase) CheckPermission(c context.Context, contactId uint64) (*domain.ContactEntity, error) {
	entity, err := s.repo.GetByID(c, contactId)
	if err != nil {
		return nil, err
	}
	if entity.OwnerOf == base_enum.EOwnerOfMember {
		profileId := _utils.GetProfileIdWithContext(c)
		if profileId != entity.OwnerID {
			return nil, &_routes.Except{
				Code:    403,
				Message: "Bạn không có quyền truy cập",
			}
		}
	}
	return entity, nil
}

func (s *ContactUsecase) Search(c context.Context, ownerType base_enum.EOwnerOf, dto dto.ContactSearchDTO) ([]domain.ContactItem, int64, error) {
	ownerId, err := s.ownerUsecase.GetOwnerInfo(c, ownerType, dto.OwnerId)
	if err != nil {
		return nil, 0, err
	}
	contacts, total, err := s.repo.Search(c, ownerId, ownerType, dto)
	if err != nil {
		return nil, 0, err
	}
	return contacts, total, nil
}

func (s *ContactUsecase) Create(c context.Context, entity *domain.ContactEntity) (*domain.ContactEntity, error) {
	ownerId, err := s.ownerUsecase.GetOwnerInfo(c, entity.OwnerOf, entity.OwnerID)
	if err != nil {
		return nil, err
	}
	entity.OwnerID = ownerId

	// Kiểm tra số điện thoại trong profile service
	if entity.Phone != "" {
		profiles, err := s.userClient.GetProfileByPhones(c, []string{entity.Phone})
		if err == nil && len(profiles) > 0 {
			if profile, exists := profiles[entity.Phone]; exists {
				entity.ProfileID = &profile.ProfileId
			}
		}
	}

	return s.repo.Create(c, entity)
}

func (s *ContactUsecase) Update(c context.Context, id uint64, entity *domain.ContactEntity) (*domain.ContactEntity, error) {
	_, err := s.CheckPermission(c, id)
	if err != nil {
		return nil, err
	}
	return s.repo.Update(c, id, entity)
}

func (s *ContactUsecase) Delete(c context.Context, id uint64) error {
	_, err := s.CheckPermission(c, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(c, id)
}

func (s *ContactUsecase) Sync(c context.Context, ownerType base_enum.EOwnerOf, contacts []domain.ContactEntity) ([]domain.ContactEntity, error) {
	phones := make([]string, len(contacts))
	for i := 0; i < len(contacts); i++ {
		phones[i] = contacts[i].Phone
	}
	profiles, _ := s.userClient.GetProfileByPhones(c, phones)
	if len(profiles) > 0 {
		for i := 0; i < len(contacts); i++ {
			existedProfile, ok := profiles[contacts[i].Phone]
			if ok {
				contacts[i].ProfileID = &existedProfile.ProfileId
			}
		}
	}

	profileId := _utils.GetProfileIdWithContext(c)
	for i := 0; i < len(contacts); i++ {
		contacts[i].OwnerID = profileId
		contacts[i].OwnerOf = ownerType
		contacts[i].Phone = strings.TrimSpace(contacts[i].Phone)
		contacts[i].FullName = strings.TrimSpace(contacts[i].FullName)
	}
	contacts, err := s.repo.Bulk(c, contacts)
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func (s *ContactUsecase) Detail(c context.Context, id uint64) (*domain.ContactEntity, error) {
	entity, err := s.CheckPermission(c, id)
	if err != nil {
		return nil, err
	}

	if entity.ProfileID != nil {
		friendStatus, _ := s.friendRepo.GetByUserID(c, *entity.ProfileID)
		entity.FriendStatus = friendStatus
	}
	return entity, nil
}

func (s *ContactUsecase) RelationShip(c context.Context, profileId uint64) *dto.RelationShipDTO {
	currentUserId := _utils.GetProfileIdWithContext(c)
	var relationShip *domain.FriendItem
	var following bool
	if currentUserId != 0 {
		relationShip, _ = s.friendRepo.GetRelationShip(c, currentUserId, profileId)
		following, _ = s.followRepo.GetFollowing(c, currentUserId, profileId)
	}

	numFollower, numFollowing, friendNumber, _ := s.followRepo.GetFollowInfoCount(c, profileId)

	return &dto.RelationShipDTO{
		FriendStatus: relationShip,
		Following:    following,
		NumFollower:  numFollower,
		NumFollowing: numFollowing,
		NumFriend:    friendNumber,
	}
}

func (s *ContactUsecase) Note(c context.Context, id uint64, note string) error {
	_, err := s.CheckPermission(c, id)
	if err != nil {
		return err
	}
	return s.repo.UpdateNote(c, id, note)
}
