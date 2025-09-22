package initial

import (
	"user/infra/handler"
	"user/infra/routes"
)

type InitialApp struct {
	Auth          *routes.AuthRouter
	OAuthRouter   *routes.OAuthRouter
	ProfileRouter *routes.ProfileRouter
	BlockRouter   *routes.BlockRouter
	FollowRouter  *routes.FollowRouter
	FriendRouter  *routes.FriendRouter
	GroupRouter   *routes.GroupRouter
	RateRouter    *routes.RateRouter
	// V2 Handlers (protobuf-based)
	ReportHandler       *handler.ReportHandler
	ReportReasonHandler *handler.ReportReasonHandler
	ReportProofHandler  *handler.ReportProofHandler
	ProfileService      *handler.GrpcProfileService
	AuthService         *handler.GrpcAuthService
	RateService         *handler.RateService
	// Admin Services
	AdminHandler        *handler.AdminHandler
	BookmarkUserHandler *handler.BookmarkUserHandler
	// GroupRouter   *routes.GroupRouter
}

func NewInitialApp(authService *routes.AuthRouter,
	oAuthRouter *routes.OAuthRouter,
	profile *routes.ProfileRouter,
	block *routes.BlockRouter,
	follow *routes.FollowRouter,
	friend *routes.FriendRouter,
	group *routes.GroupRouter,
	rate *routes.RateRouter,
	reportHandler *handler.ReportHandler,
	reportReasonHandler *handler.ReportReasonHandler,
	reportProofHandler *handler.ReportProofHandler,
	profileService *handler.GrpcProfileService,
	grpcAuthService *handler.GrpcAuthService,
	rateService *handler.RateService,
	adminHandler *handler.AdminHandler,
	bookmarkUserHandler *handler.BookmarkUserHandler,
) *InitialApp {
	return &InitialApp{
		Auth:                authService,
		OAuthRouter:         oAuthRouter,
		ProfileRouter:       profile,
		BlockRouter:         block,
		FollowRouter:        follow,
		FriendRouter:        friend,
		GroupRouter:         group,
		RateRouter:          rate,
		ReportHandler:       reportHandler,
		ReportReasonHandler: reportReasonHandler,
		ReportProofHandler:  reportProofHandler,
		ProfileService:      profileService,
		AuthService:         grpcAuthService,
		RateService:         rateService,
		AdminHandler:        adminHandler,
		BookmarkUserHandler: bookmarkUserHandler,
	}
}
