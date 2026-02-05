package events

type MessageCreatedPayload struct {
	MessageID      uint64
	ConversationID uint64
	SenderID       uint64
	Text           string
	Sequence       uint64
	Timestamp      int64
}

// func (uc *MessageUsecase) SendMessage(...) {
//     tx := db.Begin()

//     timeline := insert chat_timelines
//     event := insert chat_events

//     integrationEvent := events.IntegrationEvent{
//         ID: uuid.NewString(),
//         Type: "message.created",
//         AggregateID: timeline.ConversationID,
//         OccurredAt: time.Now().UnixMilli(),
//         Payload: MessageCreatedPayload{...},
//     }

//     tx.Create(&ChatOutbox{
//         AggregateID: timeline.ConversationID,
//         EventType: "message.created",
//         Payload: integrationEvent,
//     })

//     tx.Commit()
// }
