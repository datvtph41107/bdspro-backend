package usecase

import (
	"context"

	"organization/env"
	"organization/internal/domain/entity"
)

type LogEvent struct {
	GroupId uint32
	ActorId uint32
	LogType string
	LogData string
}

type LogWorker struct {
	logUsecase GroupLogActivityUsecase
	queue      chan LogEvent
}

func NewLogWorker(logUsecase GroupLogActivityUsecase) *LogWorker {
	worker := &LogWorker{
		logUsecase: logUsecase,
		queue:      make(chan LogEvent, env.LOG_WORKER_QUEUE_SIZE),
	}

	go worker.start()
	return worker
}

func (w *LogWorker) start() {
	for event := range w.queue {
		log := &entity.GroupLogActivity{
			GroupId: event.GroupId,
			ActorId: event.ActorId,
			LogType: event.LogType,
			LogData: event.LogData,
		}
		_, _ = w.logUsecase.CreateGroupLogActivity(context.Background(), log)
	}
}

func (w *LogWorker) Push(event LogEvent) {
	select {
	case w.queue <- event:
	default:
		// Queue full, drop or log error
	}
}
