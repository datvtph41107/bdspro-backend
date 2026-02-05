package enums

type TimelineType uint32

const (
	TIMELINE_MESSAGE TimelineType = 1
	TIMELINE_EVENT   TimelineType = 2
)

func (e TimelineType) String() string {
	switch e {
	case TIMELINE_MESSAGE:
		return "MESSAGE"
	case TIMELINE_EVENT:
		return "EVENT"
	default:
		return "unknown"
	}
}
