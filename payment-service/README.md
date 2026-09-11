# payment-service

## Role

Order, payment attempt/provider evidence, settlement, fulfillment và durable Outbox.

## Runtime status

Canonical default Compose service. API + owned fulfillment/outbox actors chạy trong owner process mặc định.

## Entrypoint / lifecycle

`main.go` -> Cobra `grpc` -> `cmd/grpc/runtime.go`.

`make server` chạy `go run . grpc`.

## Developer commands

```bash
make help
make test
make build
make server
make dev
```

Nếu thay schema:

```bash
make migrateup
make migrate-version
make new_migration name=<schema_change>
```

Root có thể gọi `make migrate service=payment-service`.

## Dependencies / durable state

PostgreSQL `payment_service`, RabbitMQ, User/Auth/Notification RPC. Canonical migration: `migrate/`.

## SePay hosted checkout

Payment supports two compatible SePay ingress paths:

- `POST /v2/payment/sepay/webhook`: legacy bank-transaction webhook,
  authenticated by `SEPAY_API_KEY`.
- `POST /v2/payment/sepay/ipn`: hosted Payment Gateway IPN for VietQR,
  NAPAS and international cards, authenticated by `X-Secret-Key`.

Hosted checkout is enabled only when these operator-owned values are complete:

```dotenv
SEPAY_GATEWAY_ENVIRONMENT=sandbox
SEPAY_BANK_TRANSFER_FLOW=hosted_checkout
SEPAY_MERCHANT_ID=
SEPAY_MERCHANT_SECRET_KEY=
SEPAY_IPN_SECRET_KEY=
SEPAY_CHECKOUT_BRIDGE_URL=http://127.0.0.1:8000/v2/payment/sepay/checkout
SEPAY_SUCCESS_URL=
SEPAY_ERROR_URL=
SEPAY_CANCEL_URL=
```

`SEPAY_CHECKOUT_URL` is optional; the service selects SePay's canonical
sandbox/production checkout host from `SEPAY_GATEWAY_ENVIRONMENT`. Never put
merchant secrets in React, React Native, protobuf payloads or Git.

SePay requires an HTML `POST` for checkout. Payment signs the canonical field
sequence and returns the Gateway bridge URL; Gateway renders an auto-submitting
form to the allowlisted SePay host. The bridge receives signed fields but never
the merchant secret. Production must set `SEPAY_CHECKOUT_BRIDGE_URL` to the
public HTTPS Gateway URL.

The three browser return URLs are optional. Leave them empty for localhost and
return to the already-open Client tab after a sandbox payment; IPN still
updates durable state. Production must use public URLs owned by the deployed
Client, never localhost.

The browser redirect is not money truth. Only an authenticated `ORDER_PAID`
IPN whose captured/approved VND amount matches frozen Order terms can create a
Settlement and Fulfillment. Provider card detail fields are discarded at the
Gateway boundary and are not retained by QHPRO.

## Bank transfer operating modes

Payment owns one explicit switch; code must not guess from `localhost` or from
the presence of a bank account:

```dotenv
# Hosted Payment Gateway sandbox (optional, separate from Bank Webhooks).
SEPAY_GATEWAY_ENVIRONMENT=sandbox
SEPAY_BANK_TRANSFER_FLOW=hosted_checkout
```

This mode keeps the hosted SePay page. To show the QR inside QHPRO instead,
select `vietqr` and configure Bank Webhooks in either Test mode or Live.

### Inline QR with Bank Webhooks Test mode

Create a bank account in **my.sepay.vn → Test mode → Tài khoản ngân hàng**.
Run the VS Code task **BDSPro · Payment · Configure Webhook Test** and enter
that account's bank BIN and number. The task changes only the payment settings
in `payment-service/.env`, preserves the webhook key (generates one if empty),
and disables card checkout. It does not register the account or webhook on SePay.
The equivalent configuration is:

```dotenv
SEPAY_GATEWAY_ENVIRONMENT=sandbox
SEPAY_BANK_TRANSFER_FLOW=vietqr
SEPAY_WEBHOOK_ENVIRONMENT=test
SEPAY_BANK_TRANSFER_ENABLED=true
SEPAY_CARD_ENABLED=false
SEPAY_QR_BANK_BIN=<six-digit-bank-BIN>
SEPAY_QR_ACCOUNT_NUMBER=<account-created-in-Test-mode>
SEPAY_API_KEY=<webhook-key>
```

Restart Payment and Gateway, then create a new payment attempt. Previously
created attempts retain their original QR or hosted URL for idempotent replay.
The app displays the QR with a Test mode notice; **do not scan it to send real
money**. Test mode does not change the bank-transfer data encoded in a VietQR.
In SePay Test mode, create a JSON webhook for incoming funds using API Key
authentication and the public Gateway URL below. Simulate an incoming transaction
with exactly the order amount and reference. A successful callback returns
`{"success":true}`; the app polls until settlement and fulfillment activate the plan.

Test and Live must use separate application databases and webhook secrets.
The webhook payload itself does not identify its environment. The configuration
therefore requires the explicit pairs `sandbox + test` or `production + live`
when using `vietqr`; it never infers Test mode from a local hostname.

Official guide: [SePay Webhooks Test mode](https://developer.sepay.vn/vi/sepay-webhooks/test-mode/bat-dau-nhanh).

### Inline QR with a live bank account

For a linked real bank account, operators may instead select an inline VietQR:

```dotenv
SEPAY_GATEWAY_ENVIRONMENT=production
SEPAY_BANK_TRANSFER_FLOW=vietqr
SEPAY_WEBHOOK_ENVIRONMENT=live
SEPAY_QR_BANK_BIN=970422
SEPAY_QR_ACCOUNT_NUMBER=
SEPAY_API_KEY=
```

The BIN, account, amount and order reference are rendered to a PNG inside
Payment Service and returned as a data URI; no transaction data is sent to a
public QR-image API. Live QR represents a real transfer even on a developer machine.

For this flow, create the SePay bank webhook at
`https://<public-gateway>/v2/payment/sepay/webhook`, choose JSON and API Key,
and enter the same raw key stored in `SEPAY_API_KEY`. SePay sends
`Authorization: Apikey <key>`; do not include the `Apikey ` prefix in the
environment variable. HTTPS, reference/amount matching, durable provider
transaction dedupe and settlement remain mandatory. HMAC-SHA256 is the next
production hardening step because it requires verification against the raw
HTTP body at the Gateway boundary.

Neither QR display, browser return nor a client-side polling response activates
a subscription. Provider settlement requires an authenticated webhook whose
reference and frozen VND amount match the Order. The separate audited operator
confirmation command remains available for manual reconciliation.

## Read source from here

`cmd/grpc/runtime.go`, `internal/`, `infra/`, `config/`, `migrate/`. Theo flow settlement/fulfillment/outbox thay vì học layer trước.

Đọc theo flow **entrypoint -> config -> business/usecase -> store/client -> actor -> test**.

## Repository rules

- Root orchestration không sở hữu business của service này.
- Cross-service behavior đi qua contract/client; không đọc database owner khác.
- Worker/actor là component của owner trước khi là process/container riêng.

Xem [`../documents/SERVICE-MAP.md`](../documents/SERVICE-MAP.md) và [`../documents/DEVELOPER-OPERATING-GUIDE.md`](../documents/DEVELOPER-OPERATING-GUIDE.md).
