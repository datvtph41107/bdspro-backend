package enums

type SupportTicketStatus uint32

const (
	SupportTicketStatusNew         SupportTicketStatus = 10
	SupportTicketStatusInProgress  SupportTicketStatus = 20
	SupportTicketStatusWaitingUser SupportTicketStatus = 30
	SupportTicketStatusResolved    SupportTicketStatus = 40
	SupportTicketStatusClosed      SupportTicketStatus = 50
	SupportTicketStatusTransferred SupportTicketStatus = 60
)

type SupportTicketPriority uint32

const (
	SupportTicketPriorityLow      SupportTicketPriority = 10
	SupportTicketPriorityNormal   SupportTicketPriority = 20
	SupportTicketPriorityHigh     SupportTicketPriority = 30
	SupportTicketPriorityUrgent   SupportTicketPriority = 40
	SupportTicketPriorityCritical SupportTicketPriority = 50
)

type SupportTicketIssueType uint32

const (
	SupportTicketIssueAccount SupportTicketIssueType = 10
	SupportTicketIssuePackage SupportTicketIssueType = 20
	SupportTicketIssuePayment SupportTicketIssueType = 30
	SupportTicketIssueData    SupportTicketIssueType = 40
	SupportTicketIssueMap     SupportTicketIssueType = 50
	SupportTicketIssueRuntime SupportTicketIssueType = 60
	SupportTicketIssueOther   SupportTicketIssueType = 70
)

type SupportTicketProduct uint32

const (
	SupportTicketProductQHPro   SupportTicketProduct = 10
	SupportTicketProductBDSPro  SupportTicketProduct = 20
	SupportTicketProductAPI     SupportTicketProduct = 30
	SupportTicketProductEcosystem SupportTicketProduct = 40
)

type SupportTicketSource uint32

const (
	SupportTicketSourceAdmin   SupportTicketSource = 10
	SupportTicketSourceApp     SupportTicketSource = 20
	SupportTicketSourceWeb     SupportTicketSource = 30
	SupportTicketSourceHotline SupportTicketSource = 40
	SupportTicketSourceAPI     SupportTicketSource = 50
)

// SupportTicketHandlingTeam — tuyến / bộ phận xử lý (SRS IV.10.7.1 chuyển tuyến).
type SupportTicketHandlingTeam uint32

const (
	SupportTicketTeamCSKH    SupportTicketHandlingTeam = 10
	SupportTicketTeamData    SupportTicketHandlingTeam = 20
	SupportTicketTeamFinance SupportTicketHandlingTeam = 30
	SupportTicketTeamSales   SupportTicketHandlingTeam = 40
	SupportTicketTeamRuntime SupportTicketHandlingTeam = 50
	SupportTicketTeamOther   SupportTicketHandlingTeam = 60
)

func IsValidSupportTicketHandlingTeam(v uint32) bool {
	switch SupportTicketHandlingTeam(v) {
	case SupportTicketTeamCSKH,
		SupportTicketTeamData,
		SupportTicketTeamFinance,
		SupportTicketTeamSales,
		SupportTicketTeamRuntime,
		SupportTicketTeamOther:
		return true
	default:
		return false
	}
}

// ValidSupportTicketTransitions maps from-status -> allowed to-status
var ValidSupportTicketTransitions = map[SupportTicketStatus][]SupportTicketStatus{
	SupportTicketStatusNew:         {SupportTicketStatusInProgress, SupportTicketStatusWaitingUser, SupportTicketStatusTransferred, SupportTicketStatusClosed},
	SupportTicketStatusInProgress:  {SupportTicketStatusWaitingUser, SupportTicketStatusResolved, SupportTicketStatusTransferred, SupportTicketStatusClosed},
	SupportTicketStatusWaitingUser: {SupportTicketStatusInProgress, SupportTicketStatusResolved, SupportTicketStatusClosed},
	SupportTicketStatusResolved:    {SupportTicketStatusClosed, SupportTicketStatusInProgress},
	SupportTicketStatusTransferred: {SupportTicketStatusInProgress, SupportTicketStatusResolved, SupportTicketStatusClosed},
	SupportTicketStatusClosed:      {},
}

func CanTransitionSupportTicket(from, to SupportTicketStatus) bool {
	allowed, ok := ValidSupportTicketTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}
