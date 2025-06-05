package router

import (
	"github.com/gin-gonic/gin"
)

type RuntimeRouter struct {
	historyRouter  *HistoryRouter
	stageRouter    *StageRouter
	pipelineRouter *PipelineRouter
	ruleRouter     *RuleRouter
	customerRouter *CustomerRouter
}

func NewRuntimeRouter(
	historyRouter *HistoryRouter,
	stageRouter *StageRouter,
	pipelineRouter *PipelineRouter,
	ruleRouter *RuleRouter,
	customerRouter *CustomerRouter,
) *RuntimeRouter {
	return &RuntimeRouter{
		historyRouter:  historyRouter,
		stageRouter:    stageRouter,
		pipelineRouter: pipelineRouter,
		ruleRouter:     ruleRouter,
		customerRouter: customerRouter,
	}
}

func (r *RuntimeRouter) RuntimeRoutes(engine *gin.Engine) {
	group := engine.Group("/crm")

	r.historyRouter.HistoryRoutes(group)
	r.stageRouter.StageRoutes(group)
	r.pipelineRouter.PipelineRoutes(group)
	r.ruleRouter.RuleRoutes(group)
	r.customerRouter.CustomerRoutes(group)
}
