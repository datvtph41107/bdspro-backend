package initial

import (
	"user/infra/client"
	"user/infra/handler"
	"user/infra/routes"
	"user/internal/job"
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
	// Membership & Plan
	MembershipHandler *handler.MembershipHandler
	PlanHandler       *handler.PlanHandler
	// Internal Service
	InternalHandler       *handler.InternalHandler
	MainAreaHandler       *handler.MainAreaHandler
	PurposeUseHandler     *handler.PurposeUseHandler
	UserDashboardStatsJob *job.UserDashboardStatsJob
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
	membershipHandler *handler.MembershipHandler,
	planHandler *handler.PlanHandler,
	internalHandler *handler.InternalHandler,
	mainAreaHandler *handler.MainAreaHandler,
	purposeUseHandler *handler.PurposeUseHandler,
	userDashboardStatsJob *job.UserDashboardStatsJob,
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
		MembershipHandler:       membershipHandler,
		PlanHandler:             planHandler,
		InternalHandler:         internalHandler,
		MainAreaHandler:         mainAreaHandler,
		PurposeUseHandler:       purposeUseHandler,
		UserDashboardStatsJob:   userDashboardStatsJob,
	}
}
