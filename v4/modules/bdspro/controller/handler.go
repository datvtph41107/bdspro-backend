package controller

import (
	"net/http"

	"userv4/modules/bdspro/services"
)

// Handler is the modules/bdspro HTTP entry for product-tab and related routes.
type Handler struct {
	sync            *SyncController
	deal            *DealController
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
}

func NewHandler(
	syncSvc services.SyncService,
	dealSvc services.DealService,
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
) *Handler {
	return &Handler{
		sync:            NewSyncController(syncSvc),
		deal:            NewDealController(dealSvc),
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
	}
}

// RegisterHTTP mounts product-tab routes with gateway path parity.
func (h *Handler) RegisterHTTP(mux *http.ServeMux) {
	h.sync.RegisterHTTP(mux)
	h.deal.RegisterHTTP(mux)
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
}
