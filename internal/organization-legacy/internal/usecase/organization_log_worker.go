package usecase

import (
	"context"

	"organization/env"
	"organization/internal/domain/entity"
)

type OrganizationLogEvent struct {
	OrganizationId uint32
	ActorId        uint32
	LogType        string
	LogData        string
}

type OrganizationLogWorker struct {
	logUsecase OrganizationLogActivityUsecase
	queue      chan OrganizationLogEvent
}

func NewOrganizationLogWorker(logUsecase OrganizationLogActivityUsecase) *OrganizationLogWorker {
	worker := &OrganizationLogWorker{
		logUsecase: logUsecase,
		queue:      make(chan OrganizationLogEvent, env.LOG_WORKER_QUEUE_SIZE),
	}

	go worker.start()
	return worker
}

func (w *OrganizationLogWorker) start() {
	for event := range w.queue {
		log := &entity.OrganizationLogActivity{
			OrganizationId: event.OrganizationId,
			ActorId:        event.ActorId,
			LogType:        event.LogType,
			LogData:        event.LogData,
		}
		_, _ = w.logUsecase.CreateOrganizationLogActivity(context.Background(), log)
	}
}

func (w *OrganizationLogWorker) Push(event OrganizationLogEvent) {
	select {
	case w.queue <- event:
	default:
		// Queue full, drop or log error
	}
}
