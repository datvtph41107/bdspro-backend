//go:build wireinject
// +build wireinject

package wire

import (
	_db "common/db"
	common_injection "common/injection"
	_provider "common/provider"
	common_utils "common/utils"
	"github.com/google/wire"
	social_infra_client "social/infra/client"
	social_infra_impl "social/infra/impl"
	social_infra_mapper "social/infra/mapper"
	social_infra_postgre "social/infra/postgre"
	social_infra_rpc "social/infra/rpc"
	social_infra_service "social/infra/service"
	social_initial "social/initial"
	social_internal_interface "social/internal/interface"
	social_internal_repo "social/internal/repo"
	social_internal_usecase "social/internal/usecase"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	_db.NewTransactionRepo,
	_provider.SyncProviderSet,
	wire.Bind(new(social_internal_interface.BdsproClient), new(*social_infra_client.BdsproClient)),
	wire.Bind(new(social_internal_interface.ChatClient), new(*social_infra_client.ChatClient)),
	wire.Bind(new(social_internal_interface.ITransaction), new(*social_infra_postgre.PostgreTransaction)),
	wire.Bind(new(social_internal_interface.NotiClient), new(*social_infra_client.NotificationClient)),
	wire.Bind(new(social_internal_interface.PermissionUsecase), new(*social_infra_impl.PermissionImpl)),
	wire.Bind(new(social_internal_interface.UserClient), new(*social_infra_client.UserClient)),
	wire.Bind(new(social_internal_repo.CommentRepo), new(*social_infra_postgre.PostgreComment)),
	wire.Bind(new(social_internal_repo.FriendTagRepo), new(*social_infra_postgre.PostgreFriendTag)),
	wire.Bind(new(social_internal_repo.LikeRepo), new(*social_infra_postgre.PostgreLike)),
	wire.Bind(new(social_internal_repo.NewsFeedMediaRepo), new(*social_infra_postgre.PostgreNewsFeedMedia)),
	wire.Bind(new(social_internal_repo.NewsFeedRepo), new(*social_infra_postgre.PostgreNewsFeed)),
	wire.Bind(new(social_internal_repo.NewsFeedShareRepo), new(*social_infra_postgre.PostgreNewsFeedShare)),
	wire.Bind(new(social_internal_repo.ReportReasonRepo), new(*social_infra_postgre.PostgreReportReason)),
	wire.Bind(new(social_internal_repo.ReportRepo), new(*social_infra_postgre.PostgreReport)),
	common_injection.NewHttpClient,
	common_utils.NewSyncUtil,
	social_initial.NewApp,
	social_infra_client.NewBdsproClient,
	social_infra_rpc.NewBdsproRPCClient,
	social_infra_client.NewChatClient,
	social_infra_service.NewCommentService,
	social_internal_usecase.NewCommentUsecase,
	social_internal_usecase.NewFriendTagUsecase,
	social_infra_postgre.NewLikeRepo,
	social_infra_service.NewLikeService,
	social_internal_usecase.NewLikeUsecase,
	social_infra_mapper.NewNewsFeedMapper,
	social_infra_postgre.NewNewsFeedOfGroupPostgres,
	social_infra_postgre.NewNewsFeedOfUserPostgres,
	social_infra_service.NewNewsFeedService,
	social_internal_usecase.NewNewsFeedShareUsecase,
	social_internal_usecase.NewNewsFeedUsecase,
	social_infra_client.NewNotificationClient,
	social_infra_rpc.NewNotificationRPCClient,
	social_internal_usecase.NewNumberUpdateUsecase,
	social_infra_impl.NewPermissionImpl,
	social_infra_postgre.NewPostgreComment,
	social_infra_postgre.NewPostgreFriendTag,
	social_infra_postgre.NewPostgreNewsFeed,
	social_infra_postgre.NewPostgreNewsFeedMedia,
	social_infra_postgre.NewPostgreNewsFeedShare,
	social_infra_postgre.NewPostgreReport,
	social_infra_postgre.NewPostgreReportReason,
	social_infra_postgre.NewPostgreTransaction,
	social_infra_mapper.NewReportMapper,
	social_internal_usecase.NewReportReasonUsecase,
	social_infra_service.NewReportService,
	social_internal_usecase.NewReportUsecase,
	social_infra_client.NewUserClient,
	social_infra_rpc.NewUserRPCClient,
	social_internal_usecase.NewUserUsecase,
)

func InitializeApp() (*social_initial.InitialApp, func(), error) {
	wire.Build(wireSet)
	return nil, nil, nil
}
