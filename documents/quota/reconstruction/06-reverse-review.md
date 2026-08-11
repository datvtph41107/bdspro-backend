# Reverse review

Sau khi forward smoke pass, kiểm ngược để chắc rằng source có thể truy vết.

## Từ usage event

```text
quota_usage_events.usage_key
  ↓ command_key
user_reported.client_request_id
  ↓ job_id
report_jobs
  ↓ operation / operation_id
workspace.report.generate
```

Sau đó dùng subject + period của usage event để tìm subscription/plan access đã cấp limit.

## Từ report

```text
user_reported.id
  ↓
report_jobs.report_id
  ↓
operation_id + command_key
  ↓
request logs
  ↓
gateway Request-ID / frontend action
```

## Câu hỏi nghiệm thu

- Có biết request nào tạo state này không?
- Có biết Actor nào thực hiện không?
- Có biết operation nào bị charge không?
- Có biết subscription/plan nào cho phép không?
- Có biết quota đã reserve/commit/cancel thế nào không?
- Có biết DB commit xảy ra trước hay sau failure không?
- Retry có thể tạo state/usage thứ hai không?

Nếu một câu không trả lời được bằng source + log + DB/Redis key, chưa retire legacy.
