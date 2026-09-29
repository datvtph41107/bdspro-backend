package enums

type ActionPolicy int32

const (
	ActionPolicyAllowDirect    ActionPolicy = 10
	ActionPolicyRequireConfirm ActionPolicy = 20
	ActionPolicyBlock          ActionPolicy = 30
	ActionPolicyRequireReview  ActionPolicy = 40
)

func (e ActionPolicy) String() string {
	switch e {
	case ActionPolicyAllowDirect:
		return "ALLOW_DIRECT"
	case ActionPolicyRequireConfirm:
		return "REQUIRE_CONFIRM"
	case ActionPolicyBlock:
		return "BLOCK"
	case ActionPolicyRequireReview:
		return "REQUIRE_REVIEW"
	default:
		return "UNKNOWN"
	}
}
