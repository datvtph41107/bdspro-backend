package initial

import (
	"assistant/config"
	"assistant/infra/handler"
)

type InitialApp struct {
	AssistantHandler *handler.AssistantHandler
	Runtime          config.Runtime
}

func NewInitialApp(assistantHandler *handler.AssistantHandler, runtime config.Runtime) *InitialApp {
	return &InitialApp{
		AssistantHandler: assistantHandler,
		Runtime:          runtime,
	}
}
