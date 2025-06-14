package usecase

import (
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/repo"
	"strings"
)

type ContactUsecase struct {
	repo        repo.ContactRepo
	friendRepo  repo.FriendRepo
	followRepo  repo.FollowRepo
	leadUsecase *LeadUsecase // forward reference
}

func NewContactUsecase(
	repo repo.ContactRepo,
	friendRepo repo.FriendRepo,
	followRepo repo.FollowRepo,
) *ContactUsecase {
	return &ContactUsecase{
		repo:       repo,
		friendRepo: friendRepo,
		followRepo: followRepo,
	}
}

func (s *ContactUsecase) CheckPermission(c context.Context, contactId uint64) (*domain.ContactEntity, error) {
	entity, err := s.repo.GetByID(c, contactId)
	if err != nil {
		return nil, err
	}
	if entity.OwnerType == enums.EOwnerTypeMember {
		profileId := _utils.GetProfileIdWithContext(c)
		if profileId != contactId {
			return nil, &_routes.Except{
				Code:    403,
				Message: "Bạn không có quyền truy cập",
			}
		}
	}
	return entity, nil
}

func (s *ContactUsecase) Search(c context.Context, ownerType enums.EOwnerType, dto dto.ContactSearchDTO) ([]domain.ContactItem, int64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	contacts, total, err := s.repo.Search(c, profileId, ownerType, dto)
	if err != nil {
		return nil, 0, err
	}
	return contacts, total, nil
}

func (s *ContactUsecase) Create(c context.Context, entity *domain.ContactEntity) (*domain.ContactEntity, error) {
	if entity.OwnerType == enums.EOwnerTypeMember {
		profileId := _utils.GetProfileIdWithContext(c)
		entity.OwnerID = profileId
	}
	return s.repo.Create(c, entity)
}

func (s *ContactUsecase) Update(c context.Context, id uint64, entity *domain.ContactEntity) (*domain.ContactEntity, error) {
	_, err := s.CheckPermission(c, id)
	if err != nil {
		return nil, err
	}
	return s.repo.Update(c, entity)
}

func (s *ContactUsecase) Delete(c context.Context, id uint64) error {
	_, err := s.CheckPermission(c, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(c, id)
}

func (s *ContactUsecase) Sync(c context.Context, ownerType enums.EOwnerType, contacts []domain.ContactEntity) ([]domain.ContactEntity, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	for i := 0; i < len(contacts); i++ {
		contacts[i].OwnerID = profileId
		contacts[i].OwnerType = ownerType
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

	if entity.ProfileID != 0 {
		friendStatus, err := s.friendRepo.GetByUserID(c, entity.ProfileID)
		if err != nil {
			return nil, err
		}
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
