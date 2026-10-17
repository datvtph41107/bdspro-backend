package controller

import (
	"net/http"

	_utils "common/utils"
	"userv4/modules/bdspro/mapper"
	"userv4/modules/bdspro/services"
)

// Handler is the modules/bdspro HTTP entry for product-tab and related routes.
type Handler struct {
	sync            *SyncController
	deal            *DealController
	dealInternalNote *DealInternalNoteController
	dealInvitation  *DealInvitationController
	dealCost        *DealCostController
	dealContract    *DealContractController
	investment      *InvestmentController
	costType        *CostTypeController
	product         *ProductController
	property        *PropertyController
	public          *PublicController
	productNote     *ProductNoteController
	appointment     *AppointmentController
	propertyHistory *PropertyHistoryController
	catalog         *CatalogController
	projectBuild    *ProjectBuildController
	adminProject    *AdminProjectController
	distribute      *DistributeController
}

func NewHandler(
	syncSvc services.SyncService,
	dealSvc services.DealService,
	dealInternalNoteSvc services.DealInternalNoteService,
	dealCommissionSvc services.DealCommissionService,
	dealInvitationSvc services.DealInvitationService,
	dealCostSvc services.DealCostService,
	contractDealSvc services.ContractDealService,
	investmentSvc services.InvestmentService,
	costTypeSvc services.CostTypeService,
	productSvc services.ProductService,
	propertySvc services.PropertyService,
	publicSvc services.PublicService,
	noteSvc services.ProductNoteService,
	appointmentSvc services.AppointmentService,
	historySvc services.PropertyHistoryService,
	catalogSvc services.CatalogService,
	projectBuildSvc services.ProjectBuildService,
	adminProjectSvc services.AdminProjectService,
	distributionSvc services.DistributionService,
	distributionMapper *mapper.DistributionMapper,
	syncUtil *_utils.SyncUtil,
) *Handler {
	return &Handler{
		sync:            NewSyncController(syncSvc),
		deal:             NewDealController(dealSvc, dealCommissionSvc),
		dealInternalNote: NewDealInternalNoteController(dealInternalNoteSvc),
		dealInvitation:   NewDealInvitationController(dealInvitationSvc),
		dealCost:        NewDealCostController(dealCostSvc),
		dealContract:    NewDealContractController(dealSvc, contractDealSvc),
		investment:      NewInvestmentController(investmentSvc),
		costType:        NewCostTypeController(costTypeSvc),
		product:         NewProductController(productSvc, dealSvc),
		property:        NewPropertyController(propertySvc),
		public:          NewPublicController(publicSvc),
		productNote:     NewProductNoteController(noteSvc),
		appointment:     NewAppointmentController(appointmentSvc),
		propertyHistory: NewPropertyHistoryController(historySvc),
		catalog:         NewCatalogController(catalogSvc),
		projectBuild:    NewProjectBuildController(projectBuildSvc),
		adminProject:    NewAdminProjectController(adminProjectSvc),
		distribute:      NewDistributeController(distributionSvc, distributionMapper, syncUtil),
	}
}

// RegisterHTTP mounts product-tab routes with gateway path parity.
func (h *Handler) RegisterHTTP(mux *http.ServeMux) {
	h.sync.RegisterHTTP(mux)
	h.deal.RegisterHTTP(mux)
	h.dealInternalNote.RegisterHTTP(mux)
	h.dealInvitation.RegisterHTTP(mux)
	h.dealCost.RegisterHTTP(mux)
	h.dealContract.RegisterHTTP(mux)
	h.investment.RegisterHTTP(mux)
	h.costType.RegisterHTTP(mux)
	h.product.RegisterHTTP(mux)
	h.property.RegisterHTTP(mux)
	h.public.RegisterHTTP(mux)
	h.productNote.RegisterHTTP(mux)
	h.appointment.RegisterHTTP(mux)
	h.propertyHistory.RegisterHTTP(mux)
	h.catalog.RegisterHTTP(mux)
	h.projectBuild.RegisterHTTP(mux)
	h.adminProject.RegisterHTTP(mux)
	h.distribute.RegisterHTTP(mux)
}
