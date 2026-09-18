package rabbitmq

import (
	"common/logging"
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	paymentcompleted "notification/internal/usecase/eventing/paymentcompleted"
)

const (
	defaultConsumerRetryMin = time.Second
	defaultConsumerRetryMax = 30 * time.Second
)

// PaymentCompletedSupervisor owns the RabbitMQ connection and channel used by
// the payment event consumer. A broker restart is an expected infrastructure
// failure: the durable queue and Inbox remain the source of recovery, while
// this actor reconnects without terminating the Notification API process.
type PaymentCompletedSupervisor struct {
	url            string
	connectionName string
	exchange       string
	service        *paymentcompleted.Service
	retryMin       time.Duration
	retryMax       time.Duration
}

func NewPaymentCompletedSupervisor(
	url string,
	connectionName string,
	exchange string,
	service *paymentcompleted.Service,
) (*PaymentCompletedSupervisor, error) {
	if strings.TrimSpace(url) == "" || strings.TrimSpace(exchange) == "" || service == nil {
		return nil, errors.New("payment completed supervisor missing dependency")
	}
	return &PaymentCompletedSupervisor{
		url:            strings.TrimSpace(url),
		connectionName: strings.TrimSpace(connectionName),
		exchange:       strings.TrimSpace(exchange),
		service:        service,
		retryMin:       defaultConsumerRetryMin,
		retryMax:       defaultConsumerRetryMax,
	}, nil
}

func (s *PaymentCompletedSupervisor) Run(ctx context.Context) error {
	retryDelay := s.retryMin
	for {
		if err := ctx.Err(); err != nil {
			return nil
		}

		connection, err := OpenConnection(s.url, s.connectionName)
		if err == nil {
			var consumer *PaymentCompletedConsumer
			consumer, err = NewPaymentCompletedConsumer(connection, s.exchange, s.service)
			if err == nil {
				// A successfully established consumer resets any startup backoff.
				retryDelay = s.retryMin
				err = consumer.Run(ctx)
				_ = consumer.Close()
			}
			_ = connection.Close()
		}
		if ctx.Err() != nil {
			return nil
		}

		logging.WithComponent(ctx, "payment-consumer").Warn(
			"notification payment consumer reconnecting",
			slog.Duration("retry_in", retryDelay),
			slog.Any("error", err),
		)
		timer := time.NewTimer(retryDelay)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			return nil
		case <-timer.C:
		}
		retryDelay = nextConsumerRetryDelay(retryDelay, s.retryMax)
	}
}

func nextConsumerRetryDelay(current, maximum time.Duration) time.Duration {
	if current <= 0 {
		return defaultConsumerRetryMin
	}
	if maximum <= 0 || current >= maximum {
		return maximum
	}
	next := current * 2
	if next < current || next > maximum {
		return maximum
	}
	return next
}
