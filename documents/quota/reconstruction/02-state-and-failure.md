# State và failure

## Trước ReserveQuota

Không có business state thay đổi.
Lỗi identity/access/quota exhausted trả ngay, không rollback gì.

## Sau ReserveQuota nhưng trước PostgreSQL commit

Redis đang giữ quota tạm.
Nếu build/save transaction lỗi -> `CancelQuota`.

## Sau PostgreSQL commit

Các fact đã bền:

- report tồn tại;
- usage event tồn tại nếu operation có quota;
- report job tồn tại.

Từ đây không rollback report chỉ vì client disconnect, gRPC response fail hoặc Redis response fail.

## Redis commit không chắc chắn

Retry cùng command key trả report durable đã accepted và không reserve/charge lần hai.
Reservation còn `reserved` được `quotacleanup` đối chiếu usage event rồi commit sau khi hết thời gian giữ.
Nếu Redis mất hoặc lệch runtime state, `usagecheck` dùng compare-and-set để đưa `used` về tổng durable usage trong DB mà không ghi đè một commit đồng thời.

## Reservation quá hạn

Không tự hiểu `expired = cancel`.

```text
usage event đúng reservation tồn tại -> CommitQuota
không có usage event đúng reservation -> CancelQuota
```

## Report worker fail

- attempts còn -> job về pending với `available_at` mới;
- quá số lần -> job failed và report failed;
- quota không refund ở đây vì charge point là lúc report job đã được accepted, không phải lúc PDF render xong.
- worker crash để job ở `running`; runner định kỳ release lock stale về `pending` để claim lại.
