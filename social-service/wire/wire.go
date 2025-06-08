//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"social/config"
	social_config "social/config"
	social_infra_impl "social/infra/impl"
	social_infra_postgre "social/infra/postgre"
	social_infra_service "social/infra/service"
	social_internal_interface "social/internal/interface"
	social_internal_repo "social/internal/repo"
	social_internal_usecase "social/internal/usecase"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	wire.Bind(new(social_internal_interface.PermissionUsecase), new(*social_infra_impl.PermissionImpl)),
	wire.Bind(new(social_internal_repo.CommentRepo), new(*social_infra_postgre.PostgreComment)),
	wire.Bind(new(social_internal_repo.NewsFeedRepo), new(*social_infra_postgre.PostgreNewsFeed)),
	wire.Bind(new(social_internal_repo.ReportRepo), new(*social_infra_postgre.PostgreReport)),
	wire.Bind(new(social_internal_repo.ReportReasonRepo), new(*social_infra_postgre.PostgreReportReason)),
	social_infra_service.NewCommentService,
	social_internal_usecase.NewCommentUsecase,
	social_config.NewConfigApp,
	social_infra_service.NewNewsFeedService,
	social_internal_usecase.NewNewsFeedUsecase,
	social_infra_impl.NewPermissionImpl,
	social_infra_postgre.NewPostgreComment,
	social_infra_postgre.NewPostgreNewsFeed,
	social_infra_postgre.NewPostgreReport,
	social_infra_postgre.NewPostgreReportReason,
	social_internal_usecase.NewReportReasonUsecase,
	social_infra_service.NewReportService,
	social_internal_usecase.NewReportUsecase,
)
func InitializeApp() (*config.ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}