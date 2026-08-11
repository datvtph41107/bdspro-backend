# Migration gates

Không bật `QHPRO_REPORT_TARGET_ENABLED=true` cho tới khi các mục dưới đây đều đạt.

1. Chạy `bash scripts/generate-internal-protobuf.sh` để generate `user/access_internal.proto` và `file/internal.proto` chỉ bằng Go/gRPC plugins; không dùng global template có `generate_unbound_methods`.
2. Go toolchain đúng version repo; compile tất cả module liên quan.
3. Review rồi chạy migration V3 `000009_report_quota_target` trong môi trường test/dev; không chạy tự động từ application startup.
4. Xác nhận catalog V2 thực tế có subscription + entitlement cho `workspace.report.generate`.
5. Thiết lập trusted metadata keys/secrets cho Gateway, TQD, user-service, file-service.
6. Redis test/dev sẵn sàng.
7. Test 99/100 concurrent request.
8. Test DB rollback -> reservation canceled.
9. Test DB commit + Redis response failure -> retry không tạo/charge trùng.
10. Test cleanup reservation quá hạn bằng durable usage event.
11. Smoke từ frontend -> gateway -> TQD -> user -> Redis/Postgres.
12. Chỉ sau evidence trên mới bỏ target flag và retire legacy path.

Khi nghiệm thu file gRPC, bật riêng `QHPRO_FILE_GRPC_ENABLED=true`, cấu hình `QHPRO_FILE_GRPC_PORT` nếu không dùng `8203`, rồi test signed upload/read và negative secure-file authorization. Không cần bật file gRPC trước khi renderer có consumer thật.

Renderer không chặn nghiệm thu quota acceptance. Nó chỉ chặn việc report chuyển từ `processing` sang `ready`.
