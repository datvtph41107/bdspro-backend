package initial

import (
	"crm/infra/handler"
	"crm/internal/usecase"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
)

type InitialApp struct {
	RuleUsecase              *usecase.RuleUsecase
	RuleEventUsecase         *usecase.RuleEventUsecase
	PipelineService          *handler.PipelineService
	ContactService           *handler.ContactService
	LeadService              *handler.LeadService
	StageService             *handler.StageService
	RuleService              *handler.RuleService
	FollowService            *handler.FollowService
	BlockService             *handler.BlockService
	FriendGroupService       *handler.FriendGroupService
	FriendService            *handler.FriendService
	SharingAccessService     *handler.SharingAccessService
	InvitationInstallService *handler.InvitationInstallService
	CrmInternalService       *handler.CrmInternalService
	EnumService              *handler.EnumService
	// Feedback handlers
	RateHandler         *handler.RateHandler
	ReportHandler       *handler.ReportHandler
	ReportAdminHandler  *handler.ReportAdminHandler
	ReportReasonHandler *handler.ReportReasonHandler
	RegistryHandler     *handler.RegistryHandler
	AppointmentService  *handler.AppointmentHandler
	Logger              *flogging.FabricLogger
}

func NewInitialApp(
	ruleUsecase *usecase.RuleUsecase,
	ruleEventUsecase *usecase.RuleEventUsecase,
	pipelineService *handler.PipelineService,
	contactService *handler.ContactService,
	customerService *handler.LeadService,
	stageService *handler.StageService,
	ruleService *handler.RuleService,
	followService *handler.FollowService,
	blockService *handler.BlockService,
	friendGroupService *handler.FriendGroupService,
	friendService *handler.FriendService,
	sharingAccessService *handler.SharingAccessService,
	invitationInstallService *handler.InvitationInstallService,
	crmInternalService *handler.CrmInternalService,
	enumService *handler.EnumService,
	// Feedback handlers
	rateHandler *handler.RateHandler,
	reportHandler *handler.ReportHandler,
	reportAdminHandler *handler.ReportAdminHandler,
	reportReasonHandler *handler.ReportReasonHandler,
	registryHandler *handler.RegistryHandler,
	appointmentService *handler.AppointmentHandler,
) *InitialApp {
	return &InitialApp{
		RuleUsecase:              ruleUsecase,
		RuleEventUsecase:         ruleEventUsecase,
		PipelineService:          pipelineService,
		ContactService:           contactService,
		LeadService:              customerService,
		StageService:             stageService,
		RuleService:              ruleService,
		FollowService:            followService,
		BlockService:             blockService,
		FriendGroupService:       friendGroupService,
		FriendService:            friendService,
		SharingAccessService:     sharingAccessService,
		InvitationInstallService: invitationInstallService,
		CrmInternalService:       crmInternalService,
		EnumService:              enumService,
		// Feedback handlers
		RateHandler:         rateHandler,
		ReportHandler:       reportHandler,
		ReportAdminHandler:  reportAdminHandler,
		ReportReasonHandler: reportReasonHandler,
		RegistryHandler:     registryHandler,
		AppointmentService:  appointmentService,
		Logger:              flogging.MustGetLogger(viper.GetString("server.name")),
	}
}
