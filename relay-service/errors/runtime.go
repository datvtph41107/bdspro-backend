package errors

import (
	"errors"

	"google.golang.org/protobuf/protoadapt"
)

func ConversationNotFound() error {
	return NewError(int32(ConversationNotFoundCode), errors.New("Không tìm thấy conversation (ID không tồn tại)."))
}

func NotParticipant() error {
	return NewError(int32(NotParticipantCode), errors.New("User không thuộc conversation => không có quyền xem/gửi tin."))
}

func ForwardForbidden() error {
	return NewError(int32(ForwardForbiddenCode), errors.New("Tổ chức bật forbidForward=true => cấm forward tin ra ngoài."))
}

func FileTooLarge() error {
	return NewError(int32(FileTooLargeCode), errors.New("Dung lượng file vượt quá giới hạn cấu hình (VD: >10MB)."))
}

func UserBannedChat() error {
	return NewError(int32(UserBannedChatCode), errors.New("User bị ban chat (admin cấm)."))
}

func InvalidMessageType() error {
	return NewError(int32(InvalidMessageTypeCode), errors.New("content_type không hợp lệ (VD: server không hỗ trợ)."))
}

func CannotRecall() error {
	return NewError(int32(CannotRecallCode), errors.New("User không thể thu hồi tin vì quá thời gian, hoặc không phải chủ tin."))
}

func AlreadyInConversation() error {
	return NewError(int32(AlreadyInConversationCode), errors.New("Khi thêm user vào group, user đã ở trong group đó."))
}

func MaxParticipantsReached() error {
	return NewError(int32(MaxParticipantsReachedCode), errors.New("Nhóm (mini-team) giới hạn 10 người, đã full."))
}

func InvalidConversationType() error {
	return NewError(int32(InvalidConversationTypeCode), errors.New("Hệ thống không hỗ trợ type conversation này."))
}

func MessageNotFound() error {
	return NewError(int32(MessageNotFoundCode), errors.New("Không tìm thấy message (ID không tồn tại)."))
}

func SenderNotAdmin() error {
	return NewError(int32(SenderNotAdminCode), errors.New("User không phải admin của conversation => không có quyền thực hiện hành động này."))
}

func MessageNotInConversation() error {
	return NewError(int32(MessageNotFoundCode), errors.New("Message không thuộc conversation => không thể thực hiện hành động này."))
}

func PermissionDenied() error {
	return NewError(int32(PermissionDeniedCode), errors.New("User không có quyền thực hiện hành động này."))
}

func NotMessageOwner() error {
	return NewError(int32(NotMessageOwnerCode), errors.New("User không phải chủ sở hữu của message => không thể thực hiện hành động này."))
}

func ConversationAlreadyExists(details ...protoadapt.MessageV1) error {
	return NewError(int32(ConversationAlreadyExistsCode), errors.New("Conversation đã tồn tại."), details...)
}
