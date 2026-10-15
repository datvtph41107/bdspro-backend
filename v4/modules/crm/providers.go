package crm

import (
	"context"

	dashboard "base/dashboard"
	_db "common/db"

	"github.com/google/wire"

	infra_postgres "userv4/infra/postgres"
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

func ProvideContactRepo(crm *infra_postgres.CrmDB) repo.ContactRepo {
	return crmpostgres.NewContactRepo(crm.DB)
}

func ProvideContactProductRepo(crm *infra_postgres.CrmDB) repo.ContactProductRepo {
	return crmpostgres.NewContactProductPostgre(crm.DB)
}

func ProvideFriendRepo(crm *infra_postgres.CrmDB) repo.FriendRepo {
	return crmpostgres.NewFriendRepo(crm.DB)
}

func ProvideFriendGroupRepo(crm *infra_postgres.CrmDB) repo.FriendGroupRepo {
	return crmpostgres.NewGroupRepo(crm.DB)
}

func ProvideFollowRepo(crm *infra_postgres.CrmDB) repo.FollowRepo {
	return crmpostgres.NewFollowPostgre(crm.DB)
}

func ProvideBlockRepo(crm *infra_postgres.CrmDB) repo.BlockRepo {
	return crmpostgres.NewBlockRepo(crm.DB)
}

func ProvideLeadRepo(crm *infra_postgres.CrmDB) repo.LeadRepo {
	return crmpostgres.NewPostgreCustomer(crm.DB)
}

func ProvideStageRepo(crm *infra_postgres.CrmDB) repo.StageRepo {
	return crmpostgres.NewPostgreStage(crm.DB)
}

func ProvidePipelineRepo(crm *infra_postgres.CrmDB) repo.PipelineRepo {
	return crmpostgres.NewPostgrePipeline(crm.DB)
}

func ProvideRuleRepo(crm *infra_postgres.CrmDB) repo.RuleRepo {
	return crmpostgres.NewPostgreRule(crm.DB)
}

func ProvideTagRepo(crm *infra_postgres.CrmDB) repo.TagRepo {
	return crmpostgres.NewTagPostgres(_db.NewTransactionRepo(crm.DB))
}

func ProvideSharingRepo(crm *infra_postgres.CrmDB) repo.SharingAccessRepo {
	return crmpostgres.NewSharingAccessPostgre(crm.DB)
}

func ProvideInvitationRepo(crm *infra_postgres.CrmDB) repo.InvitationInstallRepo {
	return crmpostgres.NewInvitationInstallPostgre(crm.DB)
}

func ProvideSupportTicketRepo(crm *infra_postgres.CrmDB) repo.SupportTicketRepo {
	return crmpostgres.NewSupportTicketPostgres(crm.DB)
}

func ProvideAdminOpportunityEventRepo(crm *infra_postgres.CrmDB) repo.AdminOpportunityEventRepo {
	return crmpostgres.NewAdminOpportunityEventPostgres(crm.DB)
}

func ProvideDocumentRepo(crm *infra_postgres.CrmDB) repo.DocumentRepo {
	return crmpostgres.NewDocumentPostgre(crm.DB)
}

func ProvideProductCareRepo(crm *infra_postgres.CrmDB) repo.ProductCareRepo {
	return crmpostgres.NewProductCarePostgre(crm.DB)
}

func ProvideRateRepo(crm *infra_postgres.CrmDB) repo.RateRepo {
	return crmpostgres.NewRatePostgres(crm.DB)
}

func ProvideRegistryRepo(crm *infra_postgres.CrmDB) repo.RegistryRepo {
	return crmpostgres.NewRegistryPostgres(crm.DB)
}

func ProvideReportRepo(crm *infra_postgres.CrmDB) repo.ReportRepo {
	return crmpostgres.NewReportPostgres(crm.DB)
}

func ProvideReportReasonRepo(crm *infra_postgres.CrmDB) repo.ReportReasonRepo {
	return crmpostgres.NewReportReasonPostgres(crm.DB)
}

func ProvideAppointmentRepo(crm *infra_postgres.CrmDB) repo.AppointmentRepository {
	return crmpostgres.NewAppointmentPostgresRepository(crm.DB)
}

// ProvideAppointmentStatsJob reads appointment_stats snapshots; the daily writer is not scheduled in v4.
func ProvideAppointmentStatsJob(crm *infra_postgres.CrmDB, appointments services.AppointmentUsecase) dashboard.DashboardStatsJob {
	return dashboard.NewDashboardStatsJob(crm.DB, "appointment_stats", func(ctx context.Context) int64 {
		count, err := appointments.CountCurrent(ctx)
		if err != nil {
			return 0
		}
		return count
	})
}
