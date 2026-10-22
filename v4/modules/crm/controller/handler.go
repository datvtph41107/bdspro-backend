package controller

import (
	"net/http"
)

// Handler is the modules/crm HTTP entry for core CRM routes (not SEO).
type Handler struct {
	contact     *ContactController
	adminSales  *AdminSalesController
	friend      *FriendController
	friendGroup *FriendGroupController
	follow      *FollowController
	block       *BlockController
	tag         *TagController
	pipeline    *PipelineController
	stage       *StageController
	lead        *LeadController
	rule        *RuleController
	sharing     *SharingController
	invitation  *InvitationController
	support     *SupportTicketController
	feedback    *FeedbackController
	appointment *AppointmentController
	enum        *EnumController
}

func NewHandler(
	contact *ContactController,
	adminSales *AdminSalesController,
	friend *FriendController,
	friendGroup *FriendGroupController,
	follow *FollowController,
	block *BlockController,
	tag *TagController,
	pipeline *PipelineController,
	stage *StageController,
	lead *LeadController,
	rule *RuleController,
	sharing *SharingController,
	invitation *InvitationController,
	support *SupportTicketController,
	feedback *FeedbackController,
	appointment *AppointmentController,
	enum *EnumController,
) *Handler {
	return &Handler{
		contact:     contact,
		adminSales:  adminSales,
		friend:      friend,
		friendGroup: friendGroup,
		follow:      follow,
		block:       block,
		tag:         tag,
		pipeline:    pipeline,
		stage:       stage,
		lead:        lead,
		rule:        rule,
		sharing:     sharing,
		invitation:  invitation,
		support:     support,
		feedback:    feedback,
		appointment: appointment,
		enum:        enum,
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
	h.feedback.RegisterHTTP(mux)
	h.appointment.RegisterHTTP(mux)
	h.enum.RegisterHTTP(mux)
}
