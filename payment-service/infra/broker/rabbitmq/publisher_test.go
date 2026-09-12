package rabbitmq

import (
	"context"
	"errors"
	"testing"
	"time"

	"payment/internal/usecase/outbox"

	"github.com/stretchr/testify/require"
)

func TestNewPublisherRequiresURLAndExchange(t *testing.T) {
	t.Parallel()
	_, err := NewPublisher("", "qhpro.payment.events")
	require.Error(t, err)
	_, err = NewPublisher("amqp://guest:guest@127.0.0.1:5672/", "")
	require.Error(t, err)
}

func TestPublishWithoutBrokerReturnsUnavailable(t *testing.T) {
	t.Parallel()
	publisher, err := NewPublisher("amqp://guest:guest@127.0.0.1:1/", "qhpro.payment.events")
	require.NoError(t, err)
	t.Cleanup(func() { _ = publisher.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = publisher.Publish(ctx, outbox.Message{
		ID: 1, EventID: "evt-1", EventType: "payment.test",
		SchemaVersion: 1, RoutingKey: "payment.test", Payload: []byte(`{}`),
	})
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrPublisherUnavailable), "got %v", err)
}
