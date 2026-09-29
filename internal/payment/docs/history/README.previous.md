# Payment Service

Payment is BDSPro's durable commercial execution owner. It owns Order, Attempt, provider evidence, Settlement, Fulfillment, recovery/redrive and Payment Outbox state. User remains owner of Plan/commercial terms, Checkout intent, Subscription and Entitlement.

## Source owners

- `internal/domain/payment`: commercial state/invariants and exact `int64` minor-unit money.
- `internal/usecase/order|attempt|settlement|fulfillment|recovery|outbox`: application owners.
- `internal/domain/wallet`, `internal/usecase/wallet`: live Wallet capability, kept separate from the commercial state machine.
- `internal/domain/bank`, `internal/usecase/bank`: live Bank capability.
- `infra/postgres`: persistence implementations.
- `infra/provider/sepay`: provider protocol verification/normalization.
- `infra/client/user`: User ApplySettlement adapter.
- `infra/broker/rabbitmq`, `worker/outbox.go`: durable Outbox publication owned by the Payment process lifecycle.
- `migrate/`: the only production schema-evolution authority.

Wallet's legacy protobuf doubles are compatibility-boundary values only. Internal Wallet monetary truth is `int64`; fractional wire values are rejected rather than rounded.

## Local

```bash
# payment-service/.env là config process; root .env chỉ dành cho full Compose.
make deps-up
make migrate-up
make build
make server
```

The server process starts the fulfillment worker and Outbox publisher as owned
components. A component failure cancels the process; shutdown waits for the
same cancellation domain and closes RabbitMQ, RPC and database resources.

Repository-level architecture guards protect this ownership from drifting
back to a second standing process.
