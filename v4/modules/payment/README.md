# modules/payment

HTTP catalog owned by the v4 monolith (payment-service schema).

## Routes

- `GET /v2/payment/bank/all` — active banks for client pickers (public, no JWT). Query `timestamp` (ms) for sync cache; unchanged responses set `unchanged: true` with empty `data`. Redis key `pay:bank:all` (`SyncKeyPaymentBankAll`) is exposed in `GET /v2/bdspro/ts` as `bank_all`.

Response shape matches `paymentpb.BankList`: `{ "data": [...], "total": N }` inside the standard envelope (plus top-level `timestamp` when sync is enabled).

## Database

Reads table `banks` via `database.payment.dsn`, falling back to `database.dsn` when empty (same pattern as CRM). Run payment-service migrations on that database if the table is missing.
