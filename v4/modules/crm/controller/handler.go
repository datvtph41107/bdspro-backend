package controller

import (
	"net/http"

	"userv4/modules/crm/services"
)

// Handler is the modules/crm HTTP entry for core CRM routes (not SEO).
type Handler struct {
	contact      *ContactController
	adminSales   *AdminSalesController
	friend       *FriendController
	friendGroup  *FriendGroupController
	follow       *FollowController
	block        *BlockController
	tag          *TagController
	pipeline     *PipelineController
	stage        *StageController
	lead         *LeadController
	rule         *RuleController
	sharing      *SharingController
	invitation   *InvitationController
	support      *SupportTicketController
}

func NewHandler(
	contactSvc services.ContactService,
	contactUC *services.ContactUsecase,
	adminSales *services.AdminSalesUsecase,
	friend *services.FriendUsecase,
	friendGroup *services.FriendGroupUsecase,
	follow *services.FollowUsecase,
	block *services.BlockUsecase,
	tag *services.TagUsecase,
	pipeline *services.PipelineUsecase,
	stage *services.StageUsecase,
	lead *services.LeadUsecase,
	rule *services.RuleUsecase,
	sharing *services.SharingAccessUsecase,
	invitation *services.InvitationInstallUsecase,
	support *services.SupportTicketUsecase,
) *Handler {
	return &Handler{
		contact:     NewContactController(contactSvc, contactUC),
		adminSales:  NewAdminSalesController(adminSales),
		friend:      NewFriendController(friend),
		friendGroup: NewFriendGroupController(friendGroup),
		follow:      NewFollowController(follow),
		block:       NewBlockController(block),
		tag:         NewTagController(tag),
		pipeline:    NewPipelineController(pipeline),
		stage:       NewStageController(stage),
		lead:        NewLeadController(lead),
		rule:        NewRuleController(rule),
		sharing:     NewSharingController(sharing),
		invitation:  NewInvitationController(invitation),
		support:     NewSupportTicketController(support),
	}
}

// RegisterHTTP mounts all core CRM routes with gateway path parity.
func (h *Handler) RegisterHTTP(mux *http.ServeMux) {
	if h == nil {
		return
	}
	h.contact.RegisterHTTP(mux)
	h.adminSales.RegisterHTTP(mux)
	h.friend.RegisterHTTP(mux)
	h.friendGroup.RegisterHTTP(mux)
	h.follow.RegisterHTTP(mux)
	h.block.RegisterHTTP(mux)
	h.tag.RegisterHTTP(mux)
	h.pipeline.RegisterHTTP(mux)
	h.stage.RegisterHTTP(mux)
	h.lead.RegisterHTTP(mux)
	h.rule.RegisterHTTP(mux)
	h.sharing.RegisterHTTP(mux)
	h.invitation.RegisterHTTP(mux)
	h.support.RegisterHTTP(mux)
}
