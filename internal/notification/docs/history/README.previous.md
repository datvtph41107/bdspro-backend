# Notification Service

Notification owns logical notifications and external delivery responsibility. Payment eventing uses a durable responsibility chain:

```text
Payment Outbox -> RabbitMQ -> Notification Inbox -> logical Notification -> DeliveryIntent -> Delivery Worker -> provider
```

## Processes

- `notification-service grpc`: synchronous Notification RPC/API surface.
- `payment-event-worker`: consumes `payment.completed.v1`, validates envelope identity, writes Inbox + logical reaction + DeliveryIntent in one PostgreSQL transaction, then ACKs.
- `delivery-worker`: claims durable DeliveryIntent rows with lease + `claim_version` fencing and performs push delivery.

## Authority

- `migrate/` is the only schema evolution authority.
- RabbitMQ is transport responsibility, not Notification business truth.
- `payment_event_inbox.event_id` is the durable dedupe identity.
- ACK is sent only after Notification DB commit.
- Provider response loss is modeled as `unknown` and can retry with at-least-once external-effect semantics.

## Local operation

```bash
# Makefile tự nạp notification-service/.env, không dùng .env.local.
make doctor
make deps-up
make migrate-up
make generate
make build
make run-grpc
```

Run eventing in separate terminals:

```bash
make run-event-worker
# Firebase credentials are required only for the delivery worker:
make run-delivery-worker
```

The committed repository contains only `config/example.firebase.json`; real Firebase service-account material must be supplied through `NOTIFICATION_FIREBASE_CREDENTIAL_FILE` and must never be committed.
