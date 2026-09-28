package services

import (
	"context"
	"strconv"

	hubpb "pb/types/hub"
	"userv4/modules/hub/dto"
	"userv4/modules/hub/entity"
	"userv4/modules/hub/repository"
)

type LocalLocationService struct {
	repo repository.LocationRepository
}

func NewLocalLocationService(repo repository.LocationRepository) *LocalLocationService {
	return &LocalLocationService{repo: repo}
}

func (s *LocalLocationService) GetProvinces(ctx context.Context, req dto.GetProvincesRequest) (*hubpb.GetProvincesResponse, error) {
	list, err := s.repo.ListProvinces(ctx, req)
	if err != nil {
		return nil, err
	}
	data := make([]*hubpb.Province, 0, len(list))
	for _, p := range list {
		data = append(data, mapProvinceV1(p))
	}
	return &hubpb.GetProvincesResponse{Data: data}, nil
}

func (s *LocalLocationService) GetDistricts(ctx context.Context, req dto.GetDistrictsRequest) (*hubpb.GetDistrictsResponse, error) {
	list, err := s.repo.ListDistricts(ctx, req)
	if err != nil {
		return nil, err
	}
	data := make([]*hubpb.District, 0, len(list))
	for _, d := range list {
		data = append(data, mapDistrictV1(d))
	}
	return &hubpb.GetDistrictsResponse{Data: data}, nil
}

func (s *LocalLocationService) GetWards(ctx context.Context, req dto.GetWardsRequest) (*hubpb.GetWardsResponse, error) {
	list, err := s.repo.ListWards(ctx, req)
	if err != nil {
		return nil, err
	}
	data := make([]*hubpb.Ward, 0, len(list))
	for _, w := range list {
		data = append(data, mapWardV1(w))
	}
	return &hubpb.GetWardsResponse{Data: data}, nil
}

func (s *LocalLocationService) SearchLocation(ctx context.Context, req dto.SearchLocationRequest) (*hubpb.SearchLocationResponse, error) {
	list, err := s.repo.Search(ctx, req.Keyword)
	if err != nil {
		return nil, err
	}
	data := make([]*hubpb.LocationSearchResult, 0, len(list))
	for i := range list {
		data = append(data, &hubpb.LocationSearchResult{
			DistrictId:   list[i].DistrictID,
			DistrictName: list[i].DistrictName,
			DistrictType: list[i].DistrictType,
			ProvinceId:   list[i].ProvinceID,
			ProvinceName: list[i].ProvinceName,
			ProvinceType: list[i].ProvinceType,
			FullAddress:  list[i].FullAddress,
		})
	}
	return &hubpb.SearchLocationResponse{Data: data}, nil
}

func mapProvinceV1(p *entity.Province) *hubpb.Province {
	if p == nil {
		return nil
	}
	out := &hubpb.Province{
		Id: p.ID, Name: p.Name, Type: p.Type, TypeText: p.TypeText, Slug: p.Slug,
	}
	if p.CreatedAt != nil {
		out.CreatedAt = strconv.FormatInt(p.CreatedAt.UnixMilli(), 10)
	}
	if p.UpdatedAt != nil {
		out.UpdatedAt = strconv.FormatInt(p.UpdatedAt.UnixMilli(), 10)
	}
	return out
}

func mapDistrictV1(d *entity.District) *hubpb.District {
	if d == nil {
		return nil
	}
	return &hubpb.District{
		Id: d.ID, Name: d.Name, ProvinceId: d.ProvinceID, Type: d.Type, TypeText: d.TypeText,
	}
}

func mapWardV1(w *entity.Ward) *hubpb.Ward {
	if w == nil {
		return nil
	}
	return &hubpb.Ward{
		Id: w.ID, Name: w.Name, DistrictId: w.DistrictID, Type: w.Type, TypeText: w.TypeText,
	}
}
