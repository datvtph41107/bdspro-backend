package services

import (
	"context"
	"time"

	base_enum "base/enum"
	_errors "common/errors"
	_routes "common/routes"
	_utils "common/utils"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	"mdn4/internal/interface/providers"
	"mdn4/modules/crm/dto"
	"mdn4/modules/crm/entity"
	"mdn4/modules/crm/repository"
)

const contactListUserMaxRecords = 20

// LocalContactService serves contact list/user from the CRM DSN + SyncUtil.
// @bind: mdn4/internal/interface/providers.OriginIDProvider
// @bind: mdn4/modules/crm/services.ContactService
type LocalContactService struct {
	repo    repository.ContactRepository
	sync    *_utils.SyncUtil
	profile providers.ProfileProvider
}

func NewLocalContactService(
	repo repository.ContactRepository,
	sync *_utils.SyncUtil,
	profile providers.ProfileProvider,
) *LocalContactService {
	return &LocalContactService{repo: repo, sync: sync, profile: profile}
}

func (s *LocalContactService) ListUser(ctx context.Context, _ uint64, req dto.ContactSearchRequest) (*crmpb.ContactListDTO, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, _errors.UnauthorizedException()
	}

	size := req.Size
	if size == 0 || size > contactListUserMaxRecords {
		size = contactListUserMaxRecords
	}

	var syncKey string
	if s.sync != nil {
		syncKey = s.sync.GetKey(ctx, _utils.SyncKeyContactMe, 0)
		if !s.sync.HasUpdated(ctx, syncKey, req.Timestamp) && req.Page == 0 {
			return &crmpb.ContactListDTO{
				Timestamp: s.sync.GetTimestamp(ctx, syncKey),
			}, nil
		}
	}

	items, total, err := s.repo.SearchByOwner(ctx, profileID, base_enum.EOwnerOfMember, req.Page, size, req.Text)
	if err != nil {
		return nil, err
	}

	out := make([]*sharepb.ContactDTO, 0, len(items))
	for _, item := range items {
		detail, err := s.detailForList(ctx, profileID, item)
		if err != nil {
			return nil, err
		}
		out = append(out, detail)
	}
	s.enrichProfiles(ctx, out)

	resp := &crmpb.ContactListDTO{
		Data:  out,
		Total: uint32(total),
	}
	if s.sync != nil && syncKey != "" && req.Page == 0 {
		ts := time.Now().UnixMilli()
		_ = s.sync.PutTimestamp(ctx, syncKey, ts)
		resp.Timestamp = ts
	}
	return resp, nil
}

// CreateContact inserts a personal contact (ownerOf=member) and bumps SyncKeyContactMe.
// Returns originID (from origin_profile), contactID, profileID.
func (s *LocalContactService) CreateContact(ctx context.Context, fullName, phone, email, avatar string, linkedProfileID uint64) (uint64, uint64, uint64, error) {
	ownerID := _utils.GetProfileIdWithContext(ctx)
	if ownerID == 0 {
		return 0, 0, 0, _errors.UnauthorizedException()
	}

	contact := &entity.Contact{
		FullName: fullName,
		Phone:    phone,
		Email:    email,
		Avatar:   avatar,
		OwnerID:  ownerID,
		OwnerOf:  base_enum.EOwnerOfMember,
	}
	if linkedProfileID > 0 {
		contact.ProfileID = &linkedProfileID
	}
	if err := s.repo.Create(ctx, contact); err != nil {
		return 0, 0, 0, err
	}

	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyContactMe, 0)
		_ = s.sync.PutTimestamp(ctx, key, 0)
	}

	var originID uint64
	if linkedProfileID > 0 {
		originID, _ = s.repo.GetOriginIDByProfileID(ctx, linkedProfileID)
	}

	var retProfileID uint64
	if contact.ProfileID != nil {
		retProfileID = *contact.ProfileID
	}

	return originID, contact.ID, retProfileID, nil
}

// FindContactIDByPhone finds the contact id owned by the current user for the given phone.
func (s *LocalContactService) FindContactIDByPhone(ctx context.Context, phone string) (uint64, error) {
	ownerID := _utils.GetProfileIdWithContext(ctx)
	if ownerID == 0 {
		return 0, _errors.UnauthorizedException()
	}
	return s.repo.FindIDByPhone(ctx, ownerID, base_enum.EOwnerOfMember, phone)
}

// GetContactProfileID returns the profileId linked to the given contact.
func (s *LocalContactService) GetContactProfileID(ctx context.Context, contactID uint64) (uint64, error) {
	contact, err := s.repo.GetByID(ctx, contactID)
	if err != nil {
		return 0, err
	}
	if contact == nil || contact.ProfileID == nil {
		return 0, nil
	}
	return *contact.ProfileID, nil
}

// GetOriginID returns the origin_id for a profile from origin_profile.
func (s *LocalContactService) GetOriginID(ctx context.Context, profileID uint64) (uint64, error) {
	return s.repo.GetOriginIDByProfileID(ctx, profileID)
}

// DeleteContact soft-deletes a contact and bumps SyncKeyContactMe.
func (s *LocalContactService) DeleteContact(ctx context.Context, contactID uint64) error {
	if err := s.repo.SoftDelete(ctx, contactID); err != nil {
		return err
	}
	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyContactMe, 0)
		_ = s.sync.PutTimestamp(ctx, key, 0)
	}
	return nil
}

func (s *LocalContactService) detailForList(ctx context.Context, currentProfileID uint64, item entity.ContactListItem) (*sharepb.ContactDTO, error) {
	contact, err := s.repo.GetByID(ctx, item.ID)
	if err != nil {
		return nil, err
	}
	if contact == nil {
		return nil, &_routes.Except{Code: 404, Message: "contact not found"}
	}
	if contact.OwnerOf == base_enum.EOwnerOfMember && contact.OwnerID != currentProfileID {
		return nil, &_routes.Except{Code: 403, Message: "Bạn không có quyền truy cập"}
	}

	dto := mapContactToSharePB(contact)

	if contact.ProfileID != nil && *contact.ProfileID > 0 {
		if rel, err := s.repo.GetFriendBetween(ctx, currentProfileID, *contact.ProfileID); err == nil && rel != nil {
			createdBy := rel.CreatedBy
			dto.FriendStatus = &sharepb.FriendItem{
				Id:         rel.ID,
				Status:     rel.Status,
				CreatedBy:  &createdBy,
				ReceiverId: rel.ReceiverID,
			}
			friendID := rel.ID
			dto.FriendId = &friendID
		}
	}

	if productIDs, err := s.repo.ListProductIDs(ctx, contact.ID); err == nil {
		dto.ProductIds = productIDs
	}
	dto.Following = item.Following
	return dto, nil
}

func (s *LocalContactService) enrichProfiles(ctx context.Context, contacts []*sharepb.ContactDTO) {
	if s.profile == nil || len(contacts) == 0 {
		return
	}
	ids := make([]uint64, 0)
	seen := map[uint64]struct{}{}
	for _, c := range contacts {
		if c == nil || c.ProfileId == nil || *c.ProfileId == 0 {
			continue
		}
		id := *c.ProfileId
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return
	}
	m, err := s.profile.GetMapProfileByIds(ctx, ids)
	if err != nil || len(m) == 0 {
		return
	}
	for _, c := range contacts {
		if c == nil || c.ProfileId == nil {
			continue
		}
		p, ok := m[*c.ProfileId]
		if !ok || p == nil {
			continue
		}
		c.ProfileInfo = &sharepb.ProfileItem{
			Id:       p.ProfileID,
			FullName: p.FullName,
			Avatar:   p.Avatar,
			Phone:    p.Phone,
		}
	}
}

func mapContactToSharePB(contact *entity.Contact) *sharepb.ContactDTO {
	if contact == nil {
		return nil
	}
	out := &sharepb.ContactDTO{
		Id:            contact.ID,
		FullName:      contact.FullName,
		Avatar:        contact.Avatar,
		Phone:         contact.Phone,
		OwnerId:       contact.OwnerID,
		Birthday:      _utils.FormatTimeToString(contact.Birthday),
		OwnerOf:       int32(contact.OwnerOf),
		ProfileId:     contact.ProfileID,
		UpdatedAt:     _utils.FormatTimeToString(contact.UpdatedAt),
		CreatedAt:     _utils.FormatTimeToString(contact.CreatedAt),
		CreatedBy:     contact.CreatedBy,
		ContactStatus: 20,
		Email:         contact.Email,
		Zalo:          contact.Zalo,
		Company:       contact.Company,
		Note:          contact.Note,
		LeadId:        contact.LeadID,
	}
	return out
}
