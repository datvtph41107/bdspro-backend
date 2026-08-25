package mapper

import (
	"hub/internal/dto"
	hubpb "pb/types/hub"
)

type KeywordSearchMapper struct{}

func NewKeywordSearchMapper() *KeywordSearchMapper {
	return &KeywordSearchMapper{}
}

func (m *KeywordSearchMapper) MapTrackRequest(req *hubpb.TrackKeywordSearchRequest) *dto.TrackKeywordSearchRequest {
	if req == nil {
		return nil
	}
	return &dto.TrackKeywordSearchRequest{
		Keyword: req.GetKeyword(),
		Source:  req.GetSource(),
		Count:   req.GetCount(),
	}
}
