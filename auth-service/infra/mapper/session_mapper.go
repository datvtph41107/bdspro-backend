package mapper

import (
	"auth/internal/domain"
	"auth/internal/dto"
	_utils "common/utils"
	authpb "pb/types/auth"
)

// SessionMapper định nghĩa mapper cho session
type SessionMapper struct {
}

// NewSessionMapper tạo mới SessionMapper
func NewSessionMapper() *SessionMapper {
	return &SessionMapper{}
}

// MapSessionToLoginHistoryItem chuyển đổi UserSessionEntity sang LoginHistoryItem
func (m *SessionMapper) MapSessionToLoginHistoryItem(session *domain.UserSessionEntity) *dto.LoginHistoryItem {
	if session == nil {
		return nil
	}

	createdDate := ""
	if session.CreatedDate != nil {
		createdDate = _utils.FormatTimeToString(session.CreatedDate)
	}

	finishedDate := ""
	if session.FinishedDate != nil {
		finishedDate = _utils.FormatTimeToString(session.FinishedDate)
	}

	return &dto.LoginHistoryItem{
		SessionID:    session.SessionID,
		AuthID:       session.AuthID,
		Platform:     session.Platform,
		Version:      session.Version,
		OS:           session.OS,
		DeviceName:   session.DeviceName,
		CreatedDate:  createdDate,
		FinishedDate: finishedDate,
		IPRequest:    session.IPRequest,
		TotalRequest: session.TotalRequest,
		UserAgent:    session.UserAgent,
		Activate:     session.Activate,
	}
}

// MapSessionToLoginHistoryItemPb chuyển đổi UserSessionEntity sang LoginHistoryItem protobuf
func (m *SessionMapper) MapSessionToLoginHistoryItemPb(session *domain.UserSessionEntity) *authpb.LoginHistoryItem {
	if session == nil {
		return nil
	}

	createdDate := ""
	if session.CreatedDate != nil {
		createdDate = _utils.FormatTimeToString(session.CreatedDate)
	}

	finishedDate := ""
	if session.FinishedDate != nil {
		finishedDate = _utils.FormatTimeToString(session.FinishedDate)
	}

	return &authpb.LoginHistoryItem{
		SessionId:    session.SessionID,
		AuthId:       session.AuthID,
		Platform:     session.Platform,
		Version:      session.Version,
		Os:           session.OS,
		DeviceName:   session.DeviceName,
		CreatedDate:  createdDate,
		FinishedDate: finishedDate,
		IpRequest:    session.IPRequest,
		TotalRequest: session.TotalRequest,
		UserAgent:    session.UserAgent,
		Activate:     session.Activate,
	}
}

// MapSessionsToLoginHistoryList chuyển đổi danh sách UserSessionEntity sang danh sách LoginHistoryItem
func (m *SessionMapper) MapSessionsToLoginHistoryList(sessions []*domain.UserSessionEntity) []dto.LoginHistoryItem {
	if sessions == nil {
		return []dto.LoginHistoryItem{}
	}

	items := make([]dto.LoginHistoryItem, 0, len(sessions))
	for _, session := range sessions {
		if item := m.MapSessionToLoginHistoryItem(session); item != nil {
			items = append(items, *item)
		}
	}

	return items
}

// MapSessionsToLoginHistoryListPb chuyển đổi danh sách UserSessionEntity sang danh sách LoginHistoryItem protobuf
func (m *SessionMapper) MapSessionsToLoginHistoryListPb(sessions []*domain.UserSessionEntity) []*authpb.LoginHistoryItem {
	if sessions == nil {
		return []*authpb.LoginHistoryItem{}
	}

	items := make([]*authpb.LoginHistoryItem, 0, len(sessions))
	for _, session := range sessions {
		if item := m.MapSessionToLoginHistoryItemPb(session); item != nil {
			items = append(items, item)
		}
	}

	return items
}
