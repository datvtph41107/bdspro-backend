package usecase

import (
	base_enum "base/enum"
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"crm/internal"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/interface/provider"
	"crm/internal/repo"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ContactUsecase struct {
	repo               repo.ContactRepo
	friendRepo         repo.FriendRepo
	followRepo         repo.FollowRepo
	leadUsecase        *LeadUsecase // forward reference
	contactProductRepo repo.ContactProductRepo
	userClient         provider.UserClient
	bdsproProvider     provider.BdsproProvider
	ownerUsecase       *OwnerUsecase
	tagUsecase         *TagUsecase

	notificationClient provider.NotificationProvider
	transaction        provider.ITransaction
}

func NewContactUsecase(
	repo repo.ContactRepo,
	friendRepo repo.FriendRepo,
	followRepo repo.FollowRepo,
	contactProductRepo repo.ContactProductRepo,
	userClient provider.UserClient,
	bdsproProvider provider.BdsproProvider,
	ownerUsecase *OwnerUsecase,
	tagUsecase *TagUsecase,

	notificationClient provider.NotificationProvider,
	transaction provider.ITransaction,
) *ContactUsecase {
	return &ContactUsecase{
		repo:               repo,
		friendRepo:         friendRepo,
		followRepo:         followRepo,
		contactProductRepo: contactProductRepo,
		userClient:         userClient,
		bdsproProvider:     bdsproProvider,
		ownerUsecase:       ownerUsecase,
		tagUsecase:         tagUsecase,
		notificationClient: notificationClient,
		transaction:        transaction,
	}
}

func (u *ContactUsecase) GetInterestedContactsByProduct(ctx context.Context, req *dto.ContactProductInterestedFilter) ([]*dto.ContactProductInterestedDTO, int64, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, 0, _errors.ReturnError(_errors.AuthenticationRequired, _errors.WithPublicMessage("Unauthorized"), _errors.WithLegacyCode(401))
	}

	items, total, err := u.repo.GetInterestedContactsByProduct(
		ctx,
		&userID,
		req,
	)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (u *ContactUsecase) PinContact(ctx context.Context, req *dto.PinRequestDTO) error {
	meta, ok := enums.PinTargetMetaMap[req.Target]
	if !ok {
		return errors.New("unsupported pin target")
	}

	var now *time.Time
	if req.Action != "" && req.Action == "pin" {
		t := time.Now()
		now = &t
	}

	return u.repo.UpdatePinnedAtDynamic(
		ctx,
		meta,
		req.TargetIDs,
		now,
	)
}

// updateContactCountForProducts cập nhật ContactCount cho các products bị ảnh hưởng
func (s *ContactUsecase) updateContactCountForProducts(c context.Context, productIds []uint64) {
	fmt.Println("updateContactCountForProducts", productIds)
	if len(productIds) == 0 || s.bdsproProvider == nil {
		return
	}

	// Lấy ContactCount cho các products
	counts, err := s.contactProductRepo.GetContactCountsByProductIds(c, productIds)
	if err != nil {
		fmt.Println("error", err)
		// Log error nhưng không fail transaction
		return
	}

	fmt.Println("counts", counts)

	// Update ContactCount cho từng product
	go func() {
		cloneCtx := _utils.CloneContext(c)
		for productID, count := range counts {
			fmt.Println("productID", productID, "count", count)
			_ = s.bdsproProvider.UpdateContactCount(cloneCtx, productID, count)
			fmt.Println("UpdateContactCount success", productID, count)

		}
	}()
}

func (s *ContactUsecase) CheckPermission(c context.Context, contactId uint64) (*domain.ContactEntity, error) {
	entity, err := s.repo.GetByID(c, contactId)
	if err != nil {
		return nil, err
	}
	if entity.OwnerOf == base_enum.EOwnerOfMember {
		profileId := _utils.GetProfileIdWithContext(c)
		if profileId != entity.OwnerID {
			return nil, _errors.ReturnError(service.ContactAccessDenied)
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

func (s *ContactUsecase) GetListContactByProductId(c context.Context, productId uint64, dto dto.ContactSearchDTO) ([]domain.ContactItem, int64, error) {
	// Lấy ownerId từ context (thường là organizationId)
	ownerId := _utils.GetProfileIdWithContext(c)
	if ownerId == 0 {
		return nil, 0, _errors.ReturnError(service.ContactNotFound, _errors.WithPublicMessage("Không tìm thấy thông tin liên hệ"), _errors.WithLegacyCode(400))
	}
	ownerType := base_enum.EOwnerOfMember
	contacts, total, err := s.repo.GetListContactByProductId(c, productId, ownerId, ownerType, dto)
	if err != nil {
		return nil, 0, err
	}
	return contacts, total, nil
}

func (s *ContactUsecase) UpdateProductIds(c context.Context, contactId uint64, productIds []uint64) error {
	// Kiểm tra quyền truy cập contact
	_, err := s.CheckPermission(c, contactId)
	if err != nil {
		return err
	}

	// Sync products vào bảng tb_contact_product
	return s.contactProductRepo.Bulk(c, productIds, contactId)
}

func (s *ContactUsecase) Create(ctx context.Context, entity *domain.ContactEntity) (*domain.ContactEntity, error) {
	var result *domain.ContactEntity
	err := s.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		ownerId, err := s.ownerUsecase.GetOwnerInfo(txCtx, entity.OwnerOf, entity.OwnerID)
		if err != nil {
			return err
		}
		entity.OwnerID = ownerId

		// Kiểm tra ProfileID đã tồn tại
		if entity.ProfileID != nil && *entity.ProfileID != 0 {
			existing, err := s.repo.GetByProfileID(txCtx, *entity.ProfileID, ownerId, entity.OwnerOf)
			if err != nil {
				return err
			}
			if existing != nil {
				return _errors.ReturnError(service.ContactAlreadyExists)
			}
		}

		// Kiểm tra phone và lấy ProfileID từ user service
		if entity.Phone != "" {
			profiles, err := s.userClient.GetProfileByPhones(txCtx, []string{entity.Phone})
			if err == nil && len(profiles) > 0 {
				if profile, exists := profiles[entity.Phone]; exists {
					if profile.ProfileId == ownerId {
						return _errors.ReturnError(service.SelfContactNotAllowed)
					}
					entity.ProfileID = &profile.ProfileId
				}
			}
		}

		// Tạo contact
		created, err := s.repo.Create(txCtx, entity)
		if err != nil {
			return err
		}
		result = created

		// Xử lý tags
		if len(entity.Tags) > 0 && result.ID > 0 {
			if err := s.tagUsecase.ReplaceContactTags(txCtx, result.ID, entity.Tags); err != nil {
				return err
			}
			tags, _ := s.tagUsecase.ListTagsByContactID(txCtx, result.ID)
			result.Tags = tags
		}

		// Xử lý products
		if len(entity.ProductIDs) > 0 && result.ID > 0 {
			if err := s.contactProductRepo.Bulk(txCtx, entity.ProductIDs, result.ID); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(entity.ProductIDs) > 0 {
		s.updateContactCountForProducts(ctx, entity.ProductIDs)
	}

	// Tạo notification
	go func() {
		s.notificationClient.CreateNotification(ctx,
			result.Avatar,
			"Tạo liên hệ",
			[]string{result.FullName},
			_enum.NotificationContactCreate,
			&result.ID,
			result.OwnerID,
			_enum.EOwnerOfMember,
			[]string{strconv.FormatUint(result.ID, 10)},
		)
	}()

	return result, nil
}

func (s *ContactUsecase) Update(c context.Context, id uint64, entity *domain.ContactEntity) (*domain.ContactEntity, error) {
	_, err := s.CheckPermission(c, id)
	if err != nil {
		return nil, err
	}

	result, err := s.repo.Update(c, id, entity)
	if err != nil {
		return nil, err
	}

	// Xử lý tags nếu có
	if len(entity.Tags) > 0 {
		if err := s.tagUsecase.ReplaceContactTags(c, id, entity.Tags); err != nil {
			return nil, err
		}
		// Load lại tags để trả về
		tags, _ := s.tagUsecase.ListTagsByContactID(c, id)
		result.Tags = tags
	} else {
		// Nếu không có tags trong request, vẫn load tags hiện có
		tags, _ := s.tagUsecase.ListTagsByContactID(c, id)
		result.Tags = tags
	}

	// Xử lý products nếu có
	if len(entity.ProductIDs) > 0 {
		// Lấy productIds bị ảnh hưởng (cũ và mới)
		affectedProductIds, _ := s.contactProductRepo.GetAffectedProductIds(c, id, entity.ProductIDs)
		// Sync products vào bảng tb_contact_product
		if err := s.contactProductRepo.Bulk(c, entity.ProductIDs, id); err != nil {
			return nil, err
		}
		// Cập nhật ContactCount cho các products bị ảnh hưởng
		s.updateContactCountForProducts(c, affectedProductIds)
	}

	return result, nil
}

func (s *ContactUsecase) Delete(c context.Context, id uint64) error {
	_, err := s.CheckPermission(c, id)
	if err != nil {
		return err
	}

	// Lấy productIds trước khi xóa để cập nhật ContactCount
	productIds, _ := s.contactProductRepo.GetProductIds(c, id)

	// Xóa contact
	err = s.repo.Delete(c, id)
	if err != nil {
		return err
	}

	// Cập nhật ContactCount cho các products bị ảnh hưởng
	if len(productIds) > 0 {
		s.updateContactCountForProducts(c, productIds)
	}

	return nil
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

	// Load tags
	tags, _ := s.tagUsecase.ListTagsByContactID(c, id)
	entity.Tags = tags

	// Load products
	productIds, _ := s.contactProductRepo.GetProductIds(c, id)
	entity.ProductIDs = productIds

	return entity, nil
}

func (s *ContactUsecase) GetContactProducts(c context.Context, contactId uint64, pagable _dto.Pagable) ([]uint64, int64, error) {
	// Kiểm tra quyền truy cập contact
	_, err := s.CheckPermission(c, contactId)
	if err != nil {
		return nil, 0, err
	}

	// Lấy product IDs với phân trang
	productIds, total, err := s.contactProductRepo.GetProductIdsWithPagination(c, contactId, pagable)
	if err != nil {
		return nil, 0, err
	}

	return productIds, total, nil
}

func (s *ContactUsecase) RelationShip(c context.Context, profileId uint64) *_dto.RelationShipDTO {
	currentUserId := _utils.GetProfileIdWithContext(c)
	var relationShip *domain.FriendItem
	var following bool
	if currentUserId != 0 {
		relationShip, _ = s.friendRepo.GetRelationShip(c, currentUserId, profileId)
		following, _ = s.followRepo.GetFollowing(c, currentUserId, profileId)
	}

	numFollower, numFollowing, friendNumber, _ := s.followRepo.GetFollowInfoCount(c, profileId)

	return &_dto.RelationShipDTO{
		FriendStatus: &_dto.FriendItemDTO{
			ID:         relationShip.ID,
			ReceiverID: relationShip.ReceiverID,
			CreatedBy:  relationShip.CreatedBy,
			Status:     uint32(relationShip.Status), // TODO: convert to uint32
		},
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

// SearchByPhone tìm kiếm liên hệ theo số điện thoại
// Trả về contact và existed=true nếu tìm thấy, existed=false nếu không tìm thấy
func (s *ContactUsecase) SearchByPhone(c context.Context, ownerType base_enum.EOwnerOf, ownerId uint64, phone string) (*domain.ContactEntity, bool, error) {
	realOwnerId, err := s.ownerUsecase.GetOwnerInfo(c, ownerType, ownerId)
	if err != nil {
		return nil, false, err
	}

	contact, err := s.repo.GetByPhone(c, realOwnerId, ownerType, phone)
	if err != nil {
		return nil, false, err
	}

	if contact == nil {
		return nil, false, nil
	}

	return contact, true, nil
}

// CreateOrRestore tạo contact mới hoặc khôi phục contact đã bị xóa
// Nếu contact đã tồn tại và đang bị xóa (deleted_at IS NOT NULL), sẽ khôi phục và cập nhật thông tin
func (s *ContactUsecase) CreateOrRestore(c context.Context, entity *domain.ContactEntity) (*domain.ContactEntity, error) {
	ownerId, err := s.ownerUsecase.GetOwnerInfo(c, entity.OwnerOf, entity.OwnerID)
	if err != nil {
		return nil, err
	}
	entity.OwnerID = ownerId

	// Kiểm tra contact đã tồn tại chưa (bao gồm cả đã bị xóa)
	if entity.Phone != "" {
		existedContact, err := s.repo.GetByPhoneIncludingDeleted(c, ownerId, entity.OwnerOf, entity.Phone)
		if err != nil {
			return nil, err
		}

		// Nếu contact đã tồn tại và đang bị xóa, khôi phục và cập nhật
		if existedContact != nil && existedContact.ID > 0 {
			if existedContact.DeletedAt != nil {
				// Khôi phục contact đã bị xóa
				if err := s.repo.Restore(c, existedContact.ID); err != nil {
					return nil, err
				}
				// Cập nhật thông tin contact
				entity.ID = existedContact.ID
				result, err := s.repo.Update(c, existedContact.ID, entity)
				if err != nil {
					return nil, err
				}

				// Xử lý tags nếu có
				if len(entity.Tags) > 0 {
					if err := s.tagUsecase.ReplaceContactTags(c, result.ID, entity.Tags); err != nil {
						return nil, err
					}
					tags, _ := s.tagUsecase.ListTagsByContactID(c, result.ID)
					result.Tags = tags
				}

				// Xử lý products nếu có
				if len(entity.ProductIDs) > 0 {
					if err := s.contactProductRepo.Bulk(c, entity.ProductIDs, result.ID); err != nil {
						return nil, err
					}
					// Cập nhật ContactCount cho các products
					s.updateContactCountForProducts(c, entity.ProductIDs)
				}

				return result, nil
			}
			// Nếu contact đã tồn tại và chưa bị xóa, trả về lỗi hoặc contact hiện tại
			return existedContact, nil
		}
	}

	// Kiểm tra số điện thoại trong profile service
	if entity.Phone != "" {
		profiles, err := s.userClient.GetProfileByPhones(c, []string{entity.Phone})
		if err == nil && len(profiles) > 0 {
			if profile, exists := profiles[entity.Phone]; exists {
				if profile.ProfileId == ownerId {
					return nil, _errors.ReturnError(service.SelfContactNotAllowed)
				}
				entity.ProfileID = &profile.ProfileId
			}
		}
	}

	// Tạo contact mới
	result, err := s.repo.Create(c, entity)
	if err != nil {
		return nil, err
	}

	// Xử lý tags nếu có
	if len(entity.Tags) > 0 && result.ID > 0 {
		if err := s.tagUsecase.ReplaceContactTags(c, result.ID, entity.Tags); err != nil {
			return nil, err
		}
		tags, _ := s.tagUsecase.ListTagsByContactID(c, result.ID)
		result.Tags = tags
	}

	// Xử lý products nếu có
	if len(entity.ProductIDs) > 0 && result.ID > 0 {
		if err := s.contactProductRepo.Bulk(c, entity.ProductIDs, result.ID); err != nil {
			return nil, err
		}
		// Cập nhật ContactCount cho các products
		s.updateContactCountForProducts(c, entity.ProductIDs)
	}

	// Tạo thông báo khi tạo liên hệ thành công
	go func() {
		s.notificationClient.CreateNotification(c,
			result.Avatar,
			"Tạo liên hệ",
			[]string{result.FullName},
			_enum.NotificationContactCreate,
			&result.ID,
			ownerId,
			_enum.EOwnerOfMember,
			[]string{strconv.FormatUint(result.ID, 10)},
		)
	}()

	return result, nil
}
