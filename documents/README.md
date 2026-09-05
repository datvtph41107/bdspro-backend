# Documentation ownership

Tài liệu vận hành canonical cho developer là:

- [`../README.md`](../README.md): entrypoint và public command API.
- [`../CONTRIBUTING.md`](../CONTRIBUTING.md): nguyên tắc thay đổi code.
- [`DEVELOPER-OPERATING-GUIDE.md`](DEVELOPER-OPERATING-GUIDE.md): workspace,
  runtime, log và daily loop.

Các tài liệu còn lại trong thư mục này được giữ làm lịch sử thiết kế và đối
chiếu migration. Chúng không phải runbook và có thể chứa port, service,
`ENV_RUNTIME`, Docker command hoặc deployment topology đã cũ. Không chạy lệnh
từ tài liệu lịch sử nếu chưa đối chiếu với root `Makefile` và `compose.yaml`.

Quy tắc cập nhật: một hành vi developer/runtime mới phải được thêm vào một
trong ba tài liệu canonical ở trên; không tạo thêm quick-start song song.
