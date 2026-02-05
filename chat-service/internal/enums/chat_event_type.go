package enums

type ChatEventType uint32

const (
	EVENT_MEMBER_JOIN          ChatEventType = 1
	EVENT_MEMBER_LEAVE         ChatEventType = 2
	EVENT_MEMBER_KICK          ChatEventType = 3
	EVENT_ROLE_CHANGE          ChatEventType = 4
	EVENT_RENAME_GROUP         ChatEventType = 5
	EVENT_AVATAR_CHANGE        ChatEventType = 6
	EVENT_GROUP_SETTING_CHANGE ChatEventType = 7
)

const (
	EVENT_MESSAGE_CREATE ChatEventType = 100
	EVENT_MESSAGE_EDIT   ChatEventType = 101
	EVENT_MESSAGE_DELETE ChatEventType = 102 // soft delete (self)
	EVENT_MESSAGE_REVOKE ChatEventType = 103 // revoke for everyone
)

const (
	EVENT_MESSAGE_PIN   ChatEventType = 110
	EVENT_MESSAGE_UNPIN ChatEventType = 111
)

const (
	EVENT_MESSAGE_REPLY   ChatEventType = 120
	EVENT_MESSAGE_FORWARD ChatEventType = 121
)

func (e ChatEventType) String() string {
	switch e {

	// System
	case EVENT_MEMBER_JOIN:
		return "member_join"
	case EVENT_MEMBER_LEAVE:
		return "member_leave"
	case EVENT_MEMBER_KICK:
		return "member_kick"
	case EVENT_ROLE_CHANGE:
		return "role_change"
	case EVENT_RENAME_GROUP:
		return "rename_group"
	case EVENT_AVATAR_CHANGE:
		return "avatar_change"
	case EVENT_GROUP_SETTING_CHANGE:
		return "group_setting_change"

	// Message
	case EVENT_MESSAGE_CREATE:
		return "message_create"
	case EVENT_MESSAGE_EDIT:
		return "message_edit"
	case EVENT_MESSAGE_DELETE:
		return "message_delete"
	case EVENT_MESSAGE_REVOKE:
		return "message_revoke"

	// Metadata
	case EVENT_MESSAGE_PIN:
		return "message_pin"
	case EVENT_MESSAGE_UNPIN:
		return "message_unpin"

	// Reply / Forward
	case EVENT_MESSAGE_REPLY:
		return "message_reply"
	case EVENT_MESSAGE_FORWARD:
		return "message_forward"

	default:
		return "unknown"
	}
}
