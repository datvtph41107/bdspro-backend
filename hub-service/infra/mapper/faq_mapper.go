package mapper

import (
	"strings"
	"unicode/utf8"

	_utils "common/utils"
	"hub/internal/domain"
	hubpb "pb/types/hub"
)

type FAQMapper struct{}

func NewFAQMapper() *FAQMapper {
	return &FAQMapper{}
}

// SanitizeUTF8 ensures string is valid UTF-8 by replacing invalid sequences
func (m *FAQMapper) SanitizeUTF8(s string) string {
	if s == "" {
		return s
	}
	// ToValidUTF8 replaces invalid UTF-8 byte sequences with empty string (removes invalid bytes)
	// This ensures protobuf can serialize the string without errors
	if !utf8.ValidString(s) {
		return strings.ToValidUTF8(s, "")
	}
	return s
}

// TruncateAnswerSafe safely truncates answer at rune boundaries (UTF-8 safe)
// The total length including "..." will not exceed maxLen
func (m *FAQMapper) TruncateAnswerSafe(s string, maxLen int) string {
	if s == "" || maxLen <= 0 {
		return m.SanitizeUTF8(s)
	}
	// Sanitize first to ensure valid UTF-8
	sanitized := m.SanitizeUTF8(s)
	// Convert to runes to work with characters, not bytes
	runes := []rune(sanitized)

	// If string is already within limit, return as is
	if len(runes) <= maxLen {
		return sanitized
	}

	// If maxLen is too small to fit "...", just truncate without ellipsis
	ellipsisLen := 3 // Length of "..."
	if maxLen <= ellipsisLen {
		return string(runes[:maxLen])
	}

	// Truncate at (maxLen - ellipsisLen) to ensure total length = maxLen
	truncateLen := maxLen - ellipsisLen
	return string(runes[:truncateLen]) + "..."
}

// sanitizeUTF8 is internal method for use within mapper
func (m *FAQMapper) sanitizeUTF8(s string) string {
	return m.SanitizeUTF8(s)
}

// CreateRequestToEntity converts CreateFAQRequest to FAQEntity
func (m *FAQMapper) CreateRequestToEntity(req *hubpb.CreateFAQRequest) *domain.FAQEntity {
	return &domain.FAQEntity{
		Question: req.GetQuestion(),
		Answer:   req.GetAnswer(),
		GroupKey: req.GetGroupKey(),
	}
}

// UpdateRequestToEntity updates FAQEntity with UpdateFAQRequest
func (m *FAQMapper) UpdateRequestToEntity(entity *domain.FAQEntity, req *hubpb.UpdateFAQRequest) {
	if req.GetQuestion() != "" {
		entity.Question = req.GetQuestion()
	}
	if req.GetAnswer() != "" {
		entity.Answer = req.GetAnswer()
	}
	if req.GetGroupKey() != "" {
		entity.GroupKey = req.GetGroupKey()
	}
}

// EntityToDetailProto converts FAQEntity to FAQDetail
func (m *FAQMapper) EntityToDetailProto(entity *domain.FAQEntity) *hubpb.FAQDetail {
	if entity == nil {
		return nil
	}
	return &hubpb.FAQDetail{
		Id:        entity.ID,
		Question:  m.sanitizeUTF8(entity.Question),
		Answer:    m.sanitizeUTF8(entity.Answer),
		GroupKey:  m.sanitizeUTF8(entity.GroupKey),
		CreatedAt: m.sanitizeUTF8(_utils.FormatTimeToString(entity.CreatedAt)),
		UpdatedAt: m.sanitizeUTF8(_utils.FormatTimeToString(entity.UpdatedAt)),
	}
}

// EntitiesToProto converts slice of FAQEntity to slice of FAQ
func (m *FAQMapper) EntitiesToProto(entities []*domain.FAQEntity) []*hubpb.FAQ {
	if entities == nil {
		return []*hubpb.FAQ{}
	}

	result := make([]*hubpb.FAQ, 0, len(entities))
	for _, entity := range entities {
		if entity == nil {
			continue
		}
		result = append(result, &hubpb.FAQ{
			Id:        entity.ID,
			Question:  m.sanitizeUTF8(entity.Question),
			Answer:    m.sanitizeUTF8(entity.Answer),
			GroupKey:  m.sanitizeUTF8(entity.GroupKey),
			CreatedAt: m.sanitizeUTF8(_utils.FormatTimeToString(entity.CreatedAt)),
			UpdatedAt: m.sanitizeUTF8(_utils.FormatTimeToString(entity.UpdatedAt)),
		})
	}
	return result
}
