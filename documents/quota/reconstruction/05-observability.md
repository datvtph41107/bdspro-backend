# Observability tối thiểu

Mục tiêu là debug được một action, không tạo hệ thống metric phức tạp.

## Gateway

- request_id
- operation_id
- idempotency_key_set (không log raw key)
- route
- caller type
- transport outcome
- duration

## TQD report

- request_id / operation_id
- business_operation
- profile_id khi an toàn cho internal log
- subject type/id
- access allowed + limit/period
- reservation_id
- quota used/reserved/remaining
- report_id
- outcome accepted/existing
- usage_key

## Worker

- job_id
- report_id
- attempts
- result completed/retry/failed
- duration

## Cleanup / usage check

Chỉ log khi có action hoặc lỗi:

- reservation committed/canceled bởi cleanup;
- Redis used được sửa theo DB usage;
- số record đã xử lý;
- lỗi DB/Redis.

Không log JWT, raw API key, cookie, file bytes hoặc full Idempotency-Key.
