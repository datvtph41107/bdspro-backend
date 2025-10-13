package initial

import (
	"assistant/infra/handler"
)

type InitialApp struct {
	AssistantHandler *handler.AssistantHandler
}

func NewInitialApp(assistantHandler *handler.AssistantHandler) *InitialApp {
	return &InitialApp{
		AssistantHandler: assistantHandler,
	}
}
