package enums

type ConversationType uint32

const (
	ConversationTypeUnknown   ConversationType = 0
	CONVERSATION_TYPE_PRIVATE ConversationType = 1
	CONVERSATION_TYPE_GROUP   ConversationType = 2
	// CONVERSATION_TYPE_CHANNEL ConversationType = 3
	// CONVERSATION_TYPE_COMMUNITY ConversationType = 4
)

func (t ConversationType) String() string {
	switch t {
	case CONVERSATION_TYPE_PRIVATE:
		return "PRIVATE"
	case CONVERSATION_TYPE_GROUP:
		return "GROUP"
	// case CONVERSATION_TYPE_CHANNEL:
	// 	return "CHANNEL"
	// case CONVERSATION_TYPE_COMMUNITY:
	// 	return "COMMUNITY"
	default:
		return "unknown"
	}
}
