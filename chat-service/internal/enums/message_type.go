package enums

type MessageType uint8

const (
	MessageTypeText   MessageType = 1
	MessageTypeImage  MessageType = 2
	MessageTypeVideo  MessageType = 3
	MessageTypeAudio  MessageType = 4
	MessageTypeFile   MessageType = 5
	MessageTypeSystem MessageType = 100
)

func (t MessageType) String() string {
	switch t {
	case MessageTypeText:
		return "text"
	case MessageTypeImage:
		return "image"
	case MessageTypeVideo:
		return "video"
	case MessageTypeAudio:
		return "audio"
	case MessageTypeFile:
		return "file"
	case MessageTypeSystem:
		return "system"
	default:
		return "unknown"
	}
}
