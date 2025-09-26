package initial

import (
	"user/infra/client"
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
	// V2 Handlers (protobuf-based)
	ReportHandler       *handler.ReportHandler
	ReportReasonHandler *handler.ReportReasonHandler
	ReportProofHandler  *handler.ReportProofHandler
	ProfileService      *handler.GrpcProfileService
	AuthService         *handler.GrpcAuthService
	// Admin Services
	AdminHandler        *handler.AdminHandler
	BookmarkUserHandler *handler.BookmarkUserHandler
	FeedbackClient      *client.FeedbackClient
	// GroupRouter   *routes.GroupRouter
}

func NewInitialApp(authService *routes.AuthRouter,
	oAuthRouter *routes.OAuthRouter,
	profile *routes.ProfileRouter,
	block *routes.BlockRouter,
	follow *routes.FollowRouter,
	friend *routes.FriendRouter,
	group *routes.GroupRouter,
	reportHandler *handler.ReportHandler,
	reportReasonHandler *handler.ReportReasonHandler,
	reportProofHandler *handler.ReportProofHandler,
	profileService *handler.GrpcProfileService,
	grpcAuthService *handler.GrpcAuthService,
	adminHandler *handler.AdminHandler,
	bookmarkUserHandler *handler.BookmarkUserHandler,
	feedbackClient *client.FeedbackClient,
) *InitialApp {
	return &InitialApp{
		Auth:                authService,
		OAuthRouter:         oAuthRouter,
		ProfileRouter:       profile,
		BlockRouter:         block,
		FollowRouter:        follow,
		FriendRouter:        friend,
		GroupRouter:         group,
		ReportHandler:       reportHandler,
		ReportReasonHandler: reportReasonHandler,
		ReportProofHandler:  reportProofHandler,
		ProfileService:      profileService,
		AuthService:         grpcAuthService,
		AdminHandler:        adminHandler,
		BookmarkUserHandler: bookmarkUserHandler,
		FeedbackClient:      feedbackClient,
	}
}
