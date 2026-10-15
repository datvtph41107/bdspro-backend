package crm

import (
	"context"

	dashboard "base/dashboard"
	_db "common/db"

	"github.com/google/wire"
	"gorm.io/gorm"

	repoiface "userv4/internal/interface/repo"
	"userv4/modules/crm/client"
	"userv4/modules/crm/controller"
	"userv4/modules/crm/mapper"
	crmpostgres "userv4/modules/crm/postgres"
	"userv4/modules/crm/provider"
	repo "userv4/modules/crm/repository"
	"userv4/modules/crm/services"
	"userv4/modules/crm/validator"
)

// ProviderSet wires CRM core HTTP (repos + usecases + NewHandler) on CrmDB.
// Providers take *CrmDB (not bare *gorm.DB) to avoid colliding with the main DB in wire.
var ProviderSet = wire.NewSet(
	ProvideCrmNoopTransaction,
	provider.NewNoopNotification,
	ProvideCRMUserClient,
	provider.NewNoopBdspro,
	provider.NewAllowAllPermission,
	wire.Bind(new(provider.ITransaction), new(*provider.NoopTransaction)),
	wire.Bind(new(provider.NotificationProvider), new(*provider.NoopNotification)),
	wire.Bind(new(provider.UserClient), new(*provider.LocalUserClient)),
	wire.Bind(new(provider.BdsproProvider), new(*provider.NoopBdspro)),
	wire.Bind(new(services.PermissionUsecase), new(*provider.AllowAllPermission)),

	ProvideContactRepo,
	ProvideContactProductRepo,
	ProvideFriendRepo,
	ProvideFriendGroupRepo,
	ProvideFollowRepo,
	ProvideBlockRepo,
	ProvideLeadRepo,
	ProvideStageRepo,
	ProvidePipelineRepo,
	ProvideRuleRepo,
	ProvideTagRepo,
	ProvideSharingRepo,
	ProvideInvitationRepo,
	ProvideSupportTicketRepo,
	ProvideAdminOpportunityEventRepo,
	ProvideDocumentRepo,
	ProvideProductCareRepo,

	// Phase-1 contact repository used by LocalContactService
	repo.NewContactPostgres,
	services.NewLocalContactService,
	wire.Bind(new(services.ContactService), new(*services.LocalContactService)),

	services.NewOwnerUsecase,
	services.NewTagUsecase,
	services.NewContactUsecase,
	services.NewBlockUsecase,
	services.NewFollowUsecase,
	services.NewFriendUsecase,
	services.NewFriendGroupUsecase,
	services.NewPipelineUsecase,
	services.NewStageUsecase,
	services.NewLeadUsecase,
	services.NewRuleUsecase,
	services.NewSharingAccessUsecase,
	services.NewInvitationInstallUsecase,
	services.NewSupportTicketUsecase,
	services.NewAdminSalesUsecase,

	ProvideRateRepo,
	ProvideRegistryRepo,
	ProvideReportRepo,
	ProvideReportReasonRepo,
	ProvideAppointmentRepo,
	ProvideAppointmentStatsJob,
	client.NewUserClient,
	services.NewRegistryUsecase,
	services.NewRateUsecase,
	services.NewReportUsecase,
	services.NewReportReasonUsecase,
	services.NewEnumUsecase,
	services.NewAppointmentUsecase,
	mapper.NewRateMapper,
	mapper.NewReportMapper,
	mapper.NewAppointmentTransformer,
	validator.NewAppointmentValidator,
	controller.NewFeedbackController,
	controller.NewAppointmentController,
	controller.NewEnumController,

	controller.NewHandler,
)

func ProvideCrmNoopTransaction() *provider.NoopTransaction {
	return provider.NewNoopTransaction()
}

func ProvideCRMUserClient(profiles repoiface.IProfileRepo) *provider.LocalUserClient {
	return provider.NewLocalUserClient(profiles)
}

func ProvideContactRepo(crm *gorm.DB) repo.ContactRepo {
	return crmpostgres.NewContactRepo(crm)
}

func ProvideContactProductRepo(crm *gorm.DB) repo.ContactProductRepo {
	return crmpostgres.NewContactProductPostgre(crm)
}

func ProvideFriendRepo(crm *gorm.DB) repo.FriendRepo {
	return crmpostgres.NewFriendRepo(crm)
}

func ProvideFriendGroupRepo(crm *gorm.DB) repo.FriendGroupRepo {
	return crmpostgres.NewGroupRepo(crm)
}

func ProvideFollowRepo(crm *gorm.DB) repo.FollowRepo {
	return crmpostgres.NewFollowPostgre(crm)
}

func ProvideBlockRepo(crm *gorm.DB) repo.BlockRepo {
	return crmpostgres.NewBlockRepo(crm)
}

func ProvideLeadRepo(crm *gorm.DB) repo.LeadRepo {
	return crmpostgres.NewPostgreCustomer(crm)
}

func ProvideStageRepo(crm *gorm.DB) repo.StageRepo {
	return crmpostgres.NewPostgreStage(crm)
}

func ProvidePipelineRepo(crm *gorm.DB) repo.PipelineRepo {
	return crmpostgres.NewPostgrePipeline(crm)
}

func ProvideRuleRepo(crm *gorm.DB) repo.RuleRepo {
	return crmpostgres.NewPostgreRule(crm)
}

func ProvideTagRepo(crm *gorm.DB) repo.TagRepo {
	return crmpostgres.NewTagPostgres(_db.NewTransactionRepo(crm))
}

func ProvideSharingRepo(crm *gorm.DB) repo.SharingAccessRepo {
	return crmpostgres.NewSharingAccessPostgre(crm)
}

func ProvideInvitationRepo(crm *gorm.DB) repo.InvitationInstallRepo {
	return crmpostgres.NewInvitationInstallPostgre(crm)
}

func ProvideSupportTicketRepo(crm *gorm.DB) repo.SupportTicketRepo {
	return crmpostgres.NewSupportTicketPostgres(crm)
}

func ProvideAdminOpportunityEventRepo(crm *gorm.DB) repo.AdminOpportunityEventRepo {
	return crmpostgres.NewAdminOpportunityEventPostgres(crm)
}

func ProvideDocumentRepo(crm *gorm.DB) repo.DocumentRepo {
	return crmpostgres.NewDocumentPostgre(crm)
}

func ProvideProductCareRepo(crm *gorm.DB) repo.ProductCareRepo {
	return crmpostgres.NewProductCarePostgre(crm)
}

func ProvideRateRepo(crm *gorm.DB) repo.RateRepo {
	return crmpostgres.NewRatePostgres(crm)
}

func ProvideRegistryRepo(crm *gorm.DB) repo.RegistryRepo {
	return crmpostgres.NewRegistryPostgres(crm)
}

func ProvideReportRepo(crm *gorm.DB) repo.ReportRepo {
	return crmpostgres.NewReportPostgres(crm)
}

func ProvideReportReasonRepo(crm *gorm.DB) repo.ReportReasonRepo {
	return crmpostgres.NewReportReasonPostgres(crm)
}

func ProvideAppointmentRepo(crm *gorm.DB) repo.AppointmentRepository {
	return crmpostgres.NewAppointmentPostgresRepository(crm)
}

// ProvideAppointmentStatsJob reads appointment_stats snapshots; the daily writer is not scheduled in v4.
func ProvideAppointmentStatsJob(crm *gorm.DB, appointments services.AppointmentUsecase) dashboard.DashboardStatsJob {
	return dashboard.NewDashboardStatsJob(crm, "appointment_stats", func(ctx context.Context) int64 {
		count, err := appointments.CountCurrent(ctx)
		if err != nil {
			return 0
		}
		return count
	})
}
