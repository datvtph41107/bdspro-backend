package initial

import (
	"user/infra/client"
	"user/infra/handler"
	"user/infra/routes"
	"user/internal/usecases"
)

type InitialApp struct {
	ProfileRouter  *routes.ProfileRouter
	BlockRouter    *routes.BlockRouter
	FollowRouter   *routes.FollowRouter
	FriendRouter   *routes.FriendRouter
	ProfileService *handler.GrpcProfileService
	ProfileUsecase *usecases.ProfileUsecase
	// Admin Services
	AdminHandler        *handler.AdminHandler
	BookmarkUserHandler *handler.BookmarkUserHandler
	FeedbackClient      *client.FeedbackClient
	// System Config
	SystemConfigHandler *handler.SystemConfigHandler
	SystemConfigUsecase *usecases.SystemConfigUsecase
}

func NewInitialApp(
	profile *routes.ProfileRouter,
	block *routes.BlockRouter,
	follow *routes.FollowRouter,
	friend *routes.FriendRouter,
	profileService *handler.GrpcProfileService,
	adminHandler *handler.AdminHandler,
	bookmarkUserHandler *handler.BookmarkUserHandler,
	feedbackClient *client.FeedbackClient,
	systemConfigHandler *handler.SystemConfigHandler,
	systemConfigUsecase *usecases.SystemConfigUsecase,
) *InitialApp {
	return &InitialApp{
		ProfileRouter:       profile,
		BlockRouter:         block,
		FollowRouter:        follow,
		FriendRouter:        friend,
		ProfileService:      profileService,
		AdminHandler:        adminHandler,
		BookmarkUserHandler: bookmarkUserHandler,
		FeedbackClient:      feedbackClient,
		SystemConfigHandler: systemConfigHandler,
		SystemConfigUsecase: systemConfigUsecase,
	}
}
