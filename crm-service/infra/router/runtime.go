package router

import (
	"github.com/gin-gonic/gin"
)

type RuntimeRouter struct {
	historyRouter  *HistoryRouter
	stageRouter    *StageRouter
	pipelineRouter *PipelineRouter
	ruleRouter     *RuleRouter
	leadRouter     *LeadRouter
}

func NewRuntimeRouter(
	historyRouter *HistoryRouter,
	stageRouter *StageRouter,
	pipelineRouter *PipelineRouter,
	ruleRouter *RuleRouter,
	leadRouter *LeadRouter,
) *RuntimeRouter {
	return &RuntimeRouter{
		historyRouter:  historyRouter,
		stageRouter:    stageRouter,
		pipelineRouter: pipelineRouter,
		ruleRouter:     ruleRouter,
		leadRouter:     leadRouter,
	}
}

// @title CRM API
// @version 1.0
// @description API phần logic chính của CRM. Đã chuyển sang grpc
// @securityDefinitions.apikey BearerAuth
// @in header
// @BasePath /v2/crm
// @name Authorization
func (r *RuntimeRouter) RuntimeRoutes(engine *gin.Engine) {
	group := engine.Group("/crm")

	r.historyRouter.HistoryRoutes(group)
	r.stageRouter.StageRoutes(group)
	r.pipelineRouter.PipelineRoutes(group)
	r.ruleRouter.RuleRoutes(group)
	r.leadRouter.LeadRoutes(group)
}
