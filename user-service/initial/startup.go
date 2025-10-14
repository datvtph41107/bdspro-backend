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
	AdminHandler            *handler.AdminHandler
	AdminUserProfileHandler *handler.AdminUserProfileHandler
	BookmarkUserHandler     *handler.BookmarkUserHandler
	FeedbackClient          *client.FeedbackClient
	// System Config
	SystemConfigHandler *handler.SystemConfigHandler
	SystemConfigUsecase *usecases.SystemConfigUsecase
	// Membership & Plan
	MembershipHandler *handler.MembershipHandler
	PlanHandler       *handler.PlanHandler
	// Internal Service
	InternalHandler *handler.InternalHandler
}

func NewInitialApp(
	profile *routes.ProfileRouter,
	block *routes.BlockRouter,
	follow *routes.FollowRouter,
	friend *routes.FriendRouter,
	profileService *handler.GrpcProfileService,
	adminHandler *handler.AdminHandler,
	adminUserProfileHandler *handler.AdminUserProfileHandler,
	bookmarkUserHandler *handler.BookmarkUserHandler,
	feedbackClient *client.FeedbackClient,
	systemConfigHandler *handler.SystemConfigHandler,
	systemConfigUsecase *usecases.SystemConfigUsecase,
	membershipHandler *handler.MembershipHandler,
	planHandler *handler.PlanHandler,
	internalHandler *handler.InternalHandler,
) *InitialApp {
	return &InitialApp{
		ProfileRouter:           profile,
		BlockRouter:             block,
		FollowRouter:            follow,
		FriendRouter:            friend,
		ProfileService:          profileService,
		AdminHandler:            adminHandler,
		AdminUserProfileHandler: adminUserProfileHandler,
		BookmarkUserHandler:     bookmarkUserHandler,
		FeedbackClient:          feedbackClient,
		SystemConfigHandler:     systemConfigHandler,
		SystemConfigUsecase:     systemConfigUsecase,
		MembershipHandler:       membershipHandler,
		PlanHandler:             planHandler,
		InternalHandler:         internalHandler,
	}
}
