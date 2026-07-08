package handler

import (
	_dto "common/domain/dto"
	"context"
	crmpb "pb/types/crm"

	"crm/infra/mapper"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/usecase"
)

type InvitationInstallService struct {
	crmpb.UnimplementedInvitationInstallServiceServer
	uc            *usecase.InvitationInstallUsecase
	ContactMapper *mapper.ContactMapper
}

func NewInvitationInstallService(
	uc *usecase.InvitationInstallUsecase,
	contactMapper *mapper.ContactMapper,
) *InvitationInstallService {
	return &InvitationInstallService{uc: uc, ContactMapper: contactMapper}
}

// @Summary Gửi lời mời
// @Description Gửi lời mời
// @Tags Lời mời
// @Accept json
// @Produce json
// @Param request body crmpb.SendInvitationRequest true "Request"
// @Security BearerAuth
// @Router /invitation/send [post]
func (s *InvitationInstallService) SendInvitation(ctx context.Context, req *crmpb.SendInvitationRequest) (*crmpb.SendInvitationResponse, error) {
	entity := &domain.AppInvitedEntity{
		Content: req.Content,
		Link:    req.Link,
		Channel: req.Channel,
	}
	entities, err := s.uc.SendInvitation(ctx, req.ContactIds, entity)
	if err != nil {
		return nil, err
	}
	items := make([]*crmpb.AppInvited, len(entities))
	for i, entity := range entities {
		items[i] = &crmpb.AppInvited{
			Id:        entity.ID,
			Content:   entity.Content,
			Link:      entity.Link,
			Channel:   entity.Channel,
			ContactId: entity.ContactID,
		}
	}
	return &crmpb.SendInvitationResponse{
		Data: items,
	}, nil
}

// @Summary Danh sách lời mời
// @Description Danh sách lời mời
// @Tags Lời mời
// @Accept json
// @Produce json
// @Param request query crmpb.ListInvitationRequest true "Request"
// @Security BearerAuth
// @Router /invitation/list [get]
func (s *InvitationInstallService) ListInvitation(ctx context.Context, req *crmpb.ListInvitationRequest) (*crmpb.ListInvitationResponse, error) {
	query := dto.ListInvitationRequest{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	}
	entities, total, err := s.uc.ListInvitation(ctx, &query)
	if err != nil {
		return nil, err
	}
	items := make([]*crmpb.AppInvited, len(entities))
	for i, entity := range entities {
		items[i] = &crmpb.AppInvited{
			Id:        entity.ID,
			Content:   entity.Content,
			Link:      entity.Link,
			Channel:   entity.Channel,
			ContactId: entity.ContactID,
			Contact:   s.ContactMapper.ContactToPb(entity.Contact),
		}
	}
	return &crmpb.ListInvitationResponse{
		Data:  items,
		Total: int32(total),
	}, nil
}