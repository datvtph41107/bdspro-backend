package services

import (
	"context"
	"strconv"
	"strings"

	_errors "common/errors"
	hubpb "pb/types/hub"
	"userv4/modules/hub/dto"
	"userv4/modules/hub/entity"
	"userv4/modules/hub/repository"
)

type LocalLocationV2Service struct {
	repo repository.LocationV2Repository
}

func NewLocalLocationV2Service(repo repository.LocationV2Repository) *LocalLocationV2Service {
	return &LocalLocationV2Service{repo: repo}
}

func (s *LocalLocationV2Service) GetProvincesV2(ctx context.Context, req dto.GetProvincesV2Request) (*hubpb.GetProvincesV2Response, error) {
	list, total, err := s.repo.ListProvinces(ctx, req.Keyword, int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}
	data := make([]*hubpb.ProvinceV2, 0, len(list))
	for i := range list {
		data = append(data, mapProvinceV2(&list[i]))
	}
	return &hubpb.GetProvincesV2Response{Data: data, Total: total}, nil
}

func (s *LocalLocationV2Service) GetProvinceV2ByCode(ctx context.Context, code int32, req dto.GetProvinceV2ByCodeRequest) (*hubpb.GetProvinceV2ByCodeResponse, error) {
	row, err := s.repo.GetProvinceByCode(ctx, strconv.Itoa(int(code)), req.IncludeWards)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, _errors.NotFoundException("province not found")
	}
	data := &hubpb.ProvinceV2WithWards{Province: mapProvinceV2(row)}
	if req.IncludeWards {
		wards := make([]*hubpb.WardV2, 0, len(row.Wards))
		for i := range row.Wards {
			wards = append(wards, mapWardV2(&row.Wards[i]))
		}
		data.Wards = wards
	}
	return &hubpb.GetProvinceV2ByCodeResponse{Data: data}, nil
}

func (s *LocalLocationV2Service) GetWardsV2(ctx context.Context, req dto.GetWardsV2Request) (*hubpb.GetWardsV2Response, error) {
	var provinceCode *string
	if req.ProvinceID > 0 {
		p, err := s.repo.GetProvinceByID(ctx, req.ProvinceID)
		if err != nil {
			return nil, err
		}
		if p != nil {
			code := p.Code
			provinceCode = &code
		}
	}
	list, total, err := s.repo.ListWards(ctx, provinceCode, req.Keyword, int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}
	data := make([]*hubpb.WardV2, 0, len(list))
	for i := range list {
		data = append(data, mapWardV2(&list[i]))
	}
	return &hubpb.GetWardsV2Response{Data: data, Total: total}, nil
}

func (s *LocalLocationV2Service) GetWardV2ByCode(ctx context.Context, code int32) (*hubpb.GetWardV2ByCodeResponse, error) {
	row, err := s.repo.GetWardByCode(ctx, strconv.Itoa(int(code)))
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, _errors.NotFoundException("ward not found")
	}
	return &hubpb.GetWardV2ByCodeResponse{Data: mapWardV2(row)}, nil
}

func (s *LocalLocationV2Service) SearchLocationV2(ctx context.Context, req dto.SearchLocationV2Request) (*hubpb.SearchLocationV2Response, error) {
	list, total, err := s.repo.Search(ctx, req.Keyword, int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}
	data := make([]*hubpb.LocationSearchResultV2, 0, len(list))
	for i := range list {
		data = append(data, &hubpb.LocationSearchResultV2{
			WardId:       list[i].WardID,
			WardName:     list[i].WardName,
			WardType:     list[i].WardType,
			ProvinceId:   list[i].ProvinceID,
			ProvinceName: list[i].ProvinceName,
			ProvinceType: list[i].ProvinceType,
			FullAddress:  list[i].FullAddress,
		})
	}
	return &hubpb.SearchLocationV2Response{Data: data, Total: total}, nil
}

func mapProvinceV2(p *entity.ProvinceV2) *hubpb.ProvinceV2 {
	if p == nil {
		return nil
	}
	out := &hubpb.ProvinceV2{
		Id:           p.ID,
		Name:         p.Name,
		Code:         parseInt32Code(p.Code),
		Codename:     p.Codename,
		DivisionType: p.DivisionType,
		PhoneCode:    parseInt32Code(p.PhoneCode),
	}
	if p.CreatedAt != nil {
		out.CreatedAt = strconv.FormatInt(p.CreatedAt.UnixMilli(), 10)
	}
	if p.UpdatedAt != nil {
		out.UpdatedAt = strconv.FormatInt(p.UpdatedAt.UnixMilli(), 10)
	}
	return out
}

func mapWardV2(w *entity.WardV2) *hubpb.WardV2 {
	if w == nil {
		return nil
	}
	out := &hubpb.WardV2{
		Id:            w.ID,
		Name:          w.Name,
		Code:          parseInt32Code(w.Code),
		Codename:      w.Codename,
		DivisionType:  w.DivisionType,
		ShortCodename: w.ShortCodename,
		ProvinceCode:  parseInt32Code(w.ProvinceCode),
	}
	if w.CreatedAt != nil {
		out.CreatedAt = strconv.FormatInt(w.CreatedAt.UnixMilli(), 10)
	}
	if w.UpdatedAt != nil {
		out.UpdatedAt = strconv.FormatInt(w.UpdatedAt.UnixMilli(), 10)
	}
	return out
}

func parseInt32Code(v string) int32 {
	n, _ := strconv.Atoi(strings.TrimSpace(v))
	return int32(n)
}
