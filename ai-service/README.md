# ai-service

## Role

Python AI runtime/source module nằm ngoài canonical Go workspace.

## Runtime status

Không thuộc canonical Go workspace/default Compose; module được kiểm syntax/source riêng.

## Entrypoint / lifecycle

Xem Python application entrypoints trong module; source gate mặc định kiểm tra Python syntax.

Không áp dụng Go `service.mk`; dùng module-owned instructions khi cần chạy riêng.

## Developer commands

Module này không dùng canonical Go Make contract. Xem module-owned runtime files và root `make verify` cho source syntax gate.

## Dependencies / durable state

Python requirements/external AI services theo module hiện hữu.

## Read source from here

Python source, `Dockerfile`, module config; không suy diễn Go service contract cho module này.

Đọc theo flow **entrypoint -> config -> business/usecase -> store/client -> actor -> test**.

## Repository rules

- Root orchestration không sở hữu business của service này.
- Cross-service behavior đi qua contract/client; không đọc database owner khác.
- Worker/actor là component của owner trước khi là process/container riêng.

Xem [`../documents/SERVICE-MAP.md`](../documents/SERVICE-MAP.md) và [`../documents/DEVELOPER-OPERATING-GUIDE.md`](../documents/DEVELOPER-OPERATING-GUIDE.md).
