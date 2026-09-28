package services

import (
	"context"
	"strconv"

	_errors "common/errors"
	_utils "common/utils"
	hubpb "pb/types/hub"
	"userv4/modules/hub/dto"
	"userv4/modules/hub/entity"
	"userv4/modules/hub/repository"
)

type LocalUserGuideService struct {
	repo repository.UserGuideRepository
	sync *_utils.SyncUtil
}

func NewLocalUserGuideService(repo repository.UserGuideRepository, sync *_utils.SyncUtil) *LocalUserGuideService {
	return &LocalUserGuideService{repo: repo, sync: sync}
}

func (s *LocalUserGuideService) GetUserGuideSimple(ctx context.Context, req dto.GetUserGuideSimpleRequest) (*hubpb.GetUserGuideSimpleResponse, error) {
	list, total, err := s.repo.ListSimple(ctx, req)
	if err != nil {
		return nil, err
	}
	data := make([]*hubpb.UserGuideSimple, 0, len(list))
	for _, g := range list {
		data = append(data, &hubpb.UserGuideSimple{
			Id: g.ID, Title: g.Title, Description: g.Description,
		})
	}
	return &hubpb.GetUserGuideSimpleResponse{Data: data, Total: total}, nil
}

func (s *LocalUserGuideService) GetUserGuide(ctx context.Context, id uint64) (*hubpb.GetUserGuideResponse, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, _errors.NotFoundException("user guide not found")
	}
	return &hubpb.GetUserGuideResponse{Data: mapUserGuideDetail(row)}, nil
}

func (s *LocalUserGuideService) GetUserGuideByKey(ctx context.Context, key string) (*hubpb.GetUserGuideByKeyResponse, error) {
	row, err := s.repo.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, _errors.NotFoundException("user guide not found")
	}
	return &hubpb.GetUserGuideByKeyResponse{Data: mapUserGuideDetail(row)}, nil
}

func (s *LocalUserGuideService) GetAllUserGuidesByKey(ctx context.Context, req dto.GetAllUserGuidesByKeyRequest) (*hubpb.GetAllUserGuidesByKeyResponse, error) {
	ts := repository.EnsureUserGuideSync(ctx, s.sync)
	if s.sync != nil && req.Timestamp > 0 && !s.sync.HasUpdated(ctx, _utils.SyncKeyUserGuideByKey, req.Timestamp) {
		return &hubpb.GetAllUserGuidesByKeyResponse{Timestamp: ts, Unchanged: true}, nil
	}
	list, err := s.repo.ListAllByKey(ctx)
	if err != nil {
		return nil, err
	}
	data := make([]*hubpb.UserGuide, 0, len(list))
	for _, g := range list {
		data = append(data, mapUserGuide(g))
	}
	return &hubpb.GetAllUserGuidesByKeyResponse{Data: data, Timestamp: ts}, nil
}

func (s *LocalUserGuideService) GetUserGuideGroups(ctx context.Context) (*hubpb.GetUserGuideGroupsResponse, error) {
	groups, err := s.repo.ListGroups(ctx)
	if err != nil {
		return nil, err
	}
	data := make([]*hubpb.UserGuideGroup, 0, len(groups))
	for _, g := range groups {
		data = append(data, &hubpb.UserGuideGroup{GroupKey: g.GroupKey, Count: g.Count})
	}
	return &hubpb.GetUserGuideGroupsResponse{Data: data}, nil
}

func mapUserGuideDetail(g *entity.UserGuide) *hubpb.UserGuideDetail {
	if g == nil {
		return nil
	}
	out := &hubpb.UserGuideDetail{
		Id: g.ID, Title: g.Title, Description: g.Description,
		GroupKey: g.GroupKey, Key: g.Key, Mode: hubpb.UserGuideMode(g.Mode),
	}
	if g.CreatedAt != nil {
		out.CreatedAt = strconv.FormatInt(g.CreatedAt.UnixMilli(), 10)
	}
	if g.UpdatedAt != nil {
		out.UpdatedAt = strconv.FormatInt(g.UpdatedAt.UnixMilli(), 10)
	}
	for _, step := range g.Steps {
		out.Steps = append(out.Steps, &hubpb.UserGuideStep{
			Id: step.ID, UserGuideId: step.UserGuideID,
			StepOrder: int32(step.StepOrder), Image: step.Image, Content: step.Content,
		})
	}
	return out
}

func mapUserGuide(g *entity.UserGuide) *hubpb.UserGuide {
	if g == nil {
		return nil
	}
	out := &hubpb.UserGuide{
		Id: g.ID, Title: g.Title, Description: g.Description,
		GroupKey: g.GroupKey, Key: g.Key, Mode: hubpb.UserGuideMode(g.Mode),
	}
	if g.CreatedAt != nil {
		out.CreatedAt = strconv.FormatInt(g.CreatedAt.UnixMilli(), 10)
	}
	if g.UpdatedAt != nil {
		out.UpdatedAt = strconv.FormatInt(g.UpdatedAt.UnixMilli(), 10)
	}
	for _, step := range g.Steps {
		out.Steps = append(out.Steps, &hubpb.UserGuideStep{
			Id: step.ID, UserGuideId: step.UserGuideID,
			StepOrder: int32(step.StepOrder), Image: step.Image, Content: step.Content,
		})
	}
	return out
}
