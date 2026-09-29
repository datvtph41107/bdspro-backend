package events

type IntegrationEvent struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	AggregateID uint64 `json:"aggregate_id"`
	OccurredAt  int64  `json:"occurred_at"`
	Payload     any    `json:"payload"`
}
