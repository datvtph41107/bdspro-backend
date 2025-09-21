package initial

import (
	"user/infra/handler"
	"user/infra/routes"
	"user/infra/services"
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
	ProfileService      *services.GrpcProfileService
	AuthService         *services.GrpcAuthService
	RateService         *services.RateService
	// Admin Services
	AdminHandler *handler.AdminHandler
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
	profileService *services.GrpcProfileService,
	grpcAuthService *services.GrpcAuthService,
	rateService *services.RateService,
	adminHandler *handler.AdminHandler,
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
	}
}
