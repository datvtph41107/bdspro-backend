package enums

type SubscriptionStatus uint32

const (
	SubStatusActive  SubscriptionStatus = 10
	SubStatusPaused  SubscriptionStatus = 20
	SubStatusDeleted SubscriptionStatus = 30
)

func (s SubscriptionStatus) String() string {
	switch s {
	case SubStatusActive:
		return "active"
	case SubStatusPaused:
		return "paused"
	case SubStatusDeleted:
		return "deleted"
	}
	return "unknown"
}
