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
	CrmClient               *client.CrmClient
	// System Config
	// Membership & Plan
	MembershipHandler *handler.MembershipHandler
	PlanHandler       *handler.PlanHandler
	// Internal Service
	InternalHandler       *handler.InternalHandler
	MainAreaHandler       *handler.MainAreaHandler
	PurposeUseHandler     *handler.PurposeUseHandler
	ProfessionHandler     *handler.ProfessionHandler
	TagHandler            *handler.TagHandler
	KYCHandler            *handler.KYCHandler
	PriceTableHandler     *handler.PriceTableHandler
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
	crmClient *client.CrmClient,
	membershipHandler *handler.MembershipHandler,
	planHandler *handler.PlanHandler,
	internalHandler *handler.InternalHandler,
	mainAreaHandler *handler.MainAreaHandler,
	purposeUseHandler *handler.PurposeUseHandler,
	professionHandler *handler.ProfessionHandler,
	tagHandler *handler.TagHandler,
	kycHandler *handler.KYCHandler,
	priceTableHandler *handler.PriceTableHandler,
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
		CrmClient:               crmClient,
		MembershipHandler:       membershipHandler,
		PlanHandler:             planHandler,
		InternalHandler:         internalHandler,
		MainAreaHandler:         mainAreaHandler,
		PurposeUseHandler:       purposeUseHandler,
		ProfessionHandler:       professionHandler,
		TagHandler:              tagHandler,
		KYCHandler:              kycHandler,
		PriceTableHandler:       priceTableHandler,
		UserDashboardStatsJob:   userDashboardStatsJob,
	}
}
