package rabbitmq

import (
	"fmt"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// OpenConnection is a process-composition adapter. The caller owns Close().
// Business/eventing usecases never dial RabbitMQ directly.
func OpenConnection(url, connectionName string) (*amqp.Connection, error) {
	url = strings.TrimSpace(url)
	if url == "" {
		return nil, fmt.Errorf("RabbitMQ URL is required")
	}
	name := strings.TrimSpace(connectionName)
	if name == "" {
		name = "notification-event-consumer"
	}
	conn, err := amqp.DialConfig(url, amqp.Config{
		Heartbeat: 10 * time.Second,
		Properties: amqp.Table{
			"connection_name": name,
			"service":         "notification-service",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("connect RabbitMQ: %w", err)
	}
	return conn, nil
}
