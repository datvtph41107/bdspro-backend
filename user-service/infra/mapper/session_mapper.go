package mapper

import (
	_utils "common/utils"
	authpb "pb/types/auth"
	sharepb "pb/types/shared"
	"user/internal/domain/auth"
	"user/internal/dto"
)

// SessionMapper định nghĩa mapper cho session
type SessionMapper struct {
}

// NewSessionMapper tạo mới SessionMapper
func NewSessionMapper() *SessionMapper {
	return &SessionMapper{}
}

// MapSessionToLoginHistoryItem chuyển đổi UserSessionEntity sang LoginHistoryItem
func (m *SessionMapper) MapSessionToLoginHistoryItem(session *auth.UserSessionEntity) *dto.LoginHistoryItem {
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
		DeviceID:     session.DeviceID,
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
func (m *SessionMapper) MapSessionToLoginHistoryItemPb(session *auth.UserSessionEntity) *authpb.LoginHistoryItem {
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
		SessionId:  session.SessionID,
		AuthId:     session.AuthID,
		DeviceName: session.DeviceName,
		Platform:   session.Platform,
		Version:    session.Version,
		Os:         session.OS,
		// DeviceId is currently not part of authpb.LoginHistoryItem
		CreatedDate:  createdDate,
		FinishedDate: finishedDate,
		IpRequest:    session.IPRequest,
		TotalRequest: session.TotalRequest,
		UserAgent:    session.UserAgent,
		Activate:     session.Activate,
	}
}

// MapSessionsToLoginHistoryList chuyển đổi danh sách UserSessionEntity sang danh sách LoginHistoryItem
func (m *SessionMapper) MapSessionsToLoginHistoryList(sessions []*auth.UserSessionEntity) []dto.LoginHistoryItem {
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
func (m *SessionMapper) MapSessionsToLoginHistoryListPb(sessions []*auth.UserSessionEntity) []*authpb.LoginHistoryItem {
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

// MapSessionToSessionV3Proto chuyển đổi UserSessionEntity sang SessionV3Proto
func (m *SessionMapper) MapSessionToSessionV3Proto(session *auth.UserSessionEntity) *sharepb.SessionV3Proto {
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

	lastLogin := ""
	if session.LastLogin != nil {
		lastLogin = _utils.FormatTimeToString(session.LastLogin)
	}

	logoutAt := ""
	if session.LogoutAt != nil {
		logoutAt = _utils.FormatTimeToString(session.LogoutAt)
	}

	lastReq := ""
	if session.LastRequest != nil {
		lastReq = _utils.FormatTimeToString(session.LastRequest)
	}

	return &sharepb.SessionV3Proto{
		SessionId:    session.SessionID,
		AuthId:       session.AuthID,
		DeviceId:     session.DeviceID,
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
		SessionKey:   session.SessionKey,
		ClientId:     session.ClientID,
		LastLogin:    lastLogin,
		LogoutAt:     logoutAt,
		LastReq:      lastReq,
	}
}

// MapSessionsToSessionV3ProtoList chuyển đổi danh sách UserSessionEntity sang danh sách SessionV3Proto
func (m *SessionMapper) MapSessionsToSessionV3ProtoList(sessions []*auth.UserSessionEntity) []*sharepb.SessionV3Proto {
	if sessions == nil {
		return []*sharepb.SessionV3Proto{}
	}

	items := make([]*sharepb.SessionV3Proto, 0, len(sessions))
	for _, session := range sessions {
		if item := m.MapSessionToSessionV3Proto(session); item != nil {
			items = append(items, item)
		}
	}

	return items
}
