package services

import (
	"context"
	"strings"

	_dto "common/domain/dto"
	"mdn4/modules/user/entity"
	mapper "mdn4/modules/user/mapper"
	userusecase "mdn4/modules/user/services/user"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
)

// CatalogService owns main-area / purpose-use / profession HTTP surface.
type CatalogService interface {
	GetMainAreas(ctx context.Context) (*userpb.MainAreaPage, error)
	CreateMainArea(ctx context.Context, req *userpb.MainArea) (*userpb.MainArea, error)
	UpdateMainArea(ctx context.Context, req *userpb.MainArea) (*userpb.MainArea, error)
	DeleteMainArea(ctx context.Context, id uint64) (*sharepb.SubmitResponse, error)
	GetMainAreaItems(ctx context.Context, text string, page, size uint32) (*sharepb.ItemV3PageProto, error)

	GetPurposeUses(ctx context.Context) (*userpb.PurposeUsePage, error)
	CreatePurposeUse(ctx context.Context, req *userpb.PurposeUse) (*userpb.PurposeUse, error)
	UpdatePurposeUse(ctx context.Context, req *userpb.PurposeUse) (*userpb.PurposeUse, error)
	DeletePurposeUse(ctx context.Context, id uint64) (*sharepb.SubmitResponse, error)
	GetPurposeUseItems(ctx context.Context) (*sharepb.ItemV3PageProto, error)

	ListProfessions(ctx context.Context, page, size uint32) (*userpb.ProfessionPage, error)
	CreateProfession(ctx context.Context, req *userpb.Profession) (*userpb.Profession, error)
	UpdateProfession(ctx context.Context, req *userpb.Profession) (*userpb.Profession, error)
	DeleteProfession(ctx context.Context, id uint64) (*sharepb.SubmitResponse, error)
	GetProfessionItems(ctx context.Context, page, size uint32) (*userpb.ProfessionPage, error)
}

// LocalCatalogService adapts catalog usecases into modules/user.
// @bind: mdn4/modules/user/services.CatalogService
type LocalCatalogService struct {
	mainArea   IAdminMainAreaUsecase
	mainMapper *mapper.MainAreaMapper

	purpose       userusecase.IAdminPurposeUseUsecase
	purposeMapper *mapper.PurposeUseMapper

	profession       userusecase.IAdminProfessionUsecase
	professionMapper *mapper.ProfessionMapper
}

func NewLocalCatalogService(
	mainArea IAdminMainAreaUsecase,
	mainMapper *mapper.MainAreaMapper,
	purpose userusecase.IAdminPurposeUseUsecase,
	purposeMapper *mapper.PurposeUseMapper,
	profession userusecase.IAdminProfessionUsecase,
	professionMapper *mapper.ProfessionMapper,
) *LocalCatalogService {
	return &LocalCatalogService{
		mainArea:         mainArea,
		mainMapper:       mainMapper,
		purpose:          purpose,
		purposeMapper:    purposeMapper,
		profession:       profession,
		professionMapper: professionMapper,
	}
}

var _ CatalogService = (*LocalCatalogService)(nil)

func (s *LocalCatalogService) GetMainAreas(ctx context.Context) (*userpb.MainAreaPage, error) {
	entities, err := s.mainArea.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	data := s.mainMapper.EntitiesToProtos(entities)
	return &userpb.MainAreaPage{Data: data, Total: uint32(len(data))}, nil
}

func (s *LocalCatalogService) CreateMainArea(ctx context.Context, req *userpb.MainArea) (*userpb.MainArea, error) {
	result, err := s.mainArea.Create(ctx, s.mainMapper.ProtoToEntity(req))
	if err != nil {
		return nil, err
	}
	return s.mainMapper.EntityToProto(result), nil
}

func (s *LocalCatalogService) UpdateMainArea(ctx context.Context, req *userpb.MainArea) (*userpb.MainArea, error) {
	entity, err := s.mainArea.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	entity.Name = req.GetName()
	if req.Description != nil {
		entity.Description = req.GetDescription()
	}
	updated, err := s.mainArea.Update(ctx, req.GetId(), entity)
	if err != nil {
		return nil, err
	}
	return s.mainMapper.EntityToProto(updated), nil
}

func (s *LocalCatalogService) DeleteMainArea(ctx context.Context, id uint64) (*sharepb.SubmitResponse, error) {
	if _, err := s.mainArea.Delete(ctx, id); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: id, Message: "Xóa nhãn khu vực hoạt động thành công"}, nil
}

func (s *LocalCatalogService) GetMainAreaItems(ctx context.Context, text string, page, size uint32) (*sharepb.ItemV3PageProto, error) {
	entities, total, err := s.mainArea.Search(ctx, text, &_dto.Pagable{Page: page, Size: size})
	if err != nil {
		return nil, err
	}
	filtered := entities
	searchText := strings.TrimSpace(text)
	if searchText != "" {
		searchTextLower := strings.ToLower(searchText)
		filtered = make([]entity.MainAreaEntity, 0, len(entities))
		for _, entity := range entities {
			nameMatch := strings.Contains(strings.ToLower(entity.Name), searchTextLower)
			descMatch := entity.Description != "" && strings.Contains(strings.ToLower(entity.Description), searchTextLower)
			if nameMatch || descMatch {
				filtered = append(filtered, entity)
			}
		}
	}
	return &sharepb.ItemV3PageProto{
		Data:  s.mainMapper.EntitiesToItemProtos(filtered),
		Total: total,
	}, nil
}

func (s *LocalCatalogService) GetPurposeUses(ctx context.Context) (*userpb.PurposeUsePage, error) {
	entities, err := s.purpose.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	data := s.purposeMapper.EntitiesToProtos(entities)
	return &userpb.PurposeUsePage{Data: data, Total: uint32(len(data))}, nil
}

func (s *LocalCatalogService) CreatePurposeUse(ctx context.Context, req *userpb.PurposeUse) (*userpb.PurposeUse, error) {
	result, err := s.purpose.Create(ctx, s.purposeMapper.ProtoToEntity(req))
	if err != nil {
		return nil, err
	}
	return s.purposeMapper.EntityToProto(result), nil
}

func (s *LocalCatalogService) UpdatePurposeUse(ctx context.Context, req *userpb.PurposeUse) (*userpb.PurposeUse, error) {
	entity, err := s.purpose.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	entity.Name = req.GetName()
	if req.Description != nil {
		entity.Description = req.GetDescription()
	}
	updated, err := s.purpose.Update(ctx, req.GetId(), entity)
	if err != nil {
		return nil, err
	}
	return s.purposeMapper.EntityToProto(updated), nil
}

func (s *LocalCatalogService) DeletePurposeUse(ctx context.Context, id uint64) (*sharepb.SubmitResponse, error) {
	if _, err := s.purpose.Delete(ctx, id); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: id, Message: "Xóa mục đích sử dụng thành công"}, nil
}

func (s *LocalCatalogService) GetPurposeUseItems(ctx context.Context) (*sharepb.ItemV3PageProto, error) {
	entities, err := s.purpose.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	data := s.purposeMapper.EntitiesToItemProtos(entities)
	return &sharepb.ItemV3PageProto{Data: data, Total: uint32(len(data))}, nil
}

func (s *LocalCatalogService) ListProfessions(ctx context.Context, page, size uint32) (*userpb.ProfessionPage, error) {
	entities, total, err := s.profession.GetList(ctx, &_dto.Pagable{Page: page, Size: size})
	if err != nil {
		return nil, err
	}
	return &userpb.ProfessionPage{
		Data:  s.professionMapper.EntitiesToProtos(entities),
		Total: uint32(total),
	}, nil
}

func (s *LocalCatalogService) CreateProfession(ctx context.Context, req *userpb.Profession) (*userpb.Profession, error) {
	result, err := s.profession.Create(ctx, s.professionMapper.ProtoToEntity(req))
	if err != nil {
		return nil, err
	}
	return s.professionMapper.EntityToProto(result), nil
}

func (s *LocalCatalogService) UpdateProfession(ctx context.Context, req *userpb.Profession) (*userpb.Profession, error) {
	result, err := s.profession.Update(ctx, req.GetId(), s.professionMapper.ProtoToEntity(req))
	if err != nil {
		return nil, err
	}
	return s.professionMapper.EntityToProto(result), nil
}

func (s *LocalCatalogService) DeleteProfession(ctx context.Context, id uint64) (*sharepb.SubmitResponse, error) {
	if _, err := s.profession.Delete(ctx, id); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: id, Message: "Xóa nghề nghiệp thành công"}, nil
}

func (s *LocalCatalogService) GetProfessionItems(ctx context.Context, page, size uint32) (*userpb.ProfessionPage, error) {
	entities, total, err := s.profession.ListVerified(ctx, &_dto.Pagable{Page: page, Size: size})
	if err != nil {
		return nil, err
	}
	return &userpb.ProfessionPage{
		Data:  s.professionMapper.EntitiesToProtos(entities),
		Total: uint32(total),
	}, nil
}
