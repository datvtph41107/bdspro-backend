# Chuẩn phát triển chức năng có quota

Tài liệu này là quy ước bắt buộc khi thêm hoặc thay đổi một chức năng có giá
trị thương mại. Mục tiêu là để dev mới có thể lần theo một request từ client
đến kết quả mà không phải đoán owner hoặc tự tạo thêm nguồn sự thật.

## Hồ sơ trước khi viết code

Mỗi chức năng phải ghi rõ: giá trị người dùng nhận được, `operation` duy nhất,
`feature`, `meter`, số unit mỗi lần, chu kỳ áp dụng, điều kiện thành công,
điều kiện thất bại, `Idempotency-Key`, kết quả trả về và owner của dữ liệu.

## Trách nhiệm theo service

- User Service quyết định catalog, phiên bản gói, subscription và entitlement.
- Payment Service tạo order, đối soát số tiền và xử lý webhook.
- TQD Service kiểm tra quyền từ User, giữ quota, ghi usage event và tạo kết quả.
- File Service lưu file và cung cấp file effect đã tạo.
- Frontend gửi command, hiển thị trạng thái canonical; không tự cấp quyền hoặc
  tự giảm quota.

## Quy tắc code

1. Không kiểm tra tên gói (`pro`, `vip`) trong nghiệp vụ; chỉ kiểm tra
   `operation` và kết quả entitlement.
2. Mọi thao tác tính quota phải có idempotency. Report trả lại resource đã
   được chấp nhận; lượt đọc đồng bộ có thể đọc lại dữ liệu hiện tại nhưng không
   tạo thêm usage cho cùng command. Không dùng một key cho hai ý định khác nhau.
3. Điểm ghi nhận phụ thuộc hợp đồng chức năng:
   - Đọc đồng bộ: handler tạo response thành công rồi database ghi usage.
   - Report: transaction commit Report + Job + Usage (nếu tính quota). PDF có
     thể vẫn đang chờ worker; tạo lại job không tiêu thụ thêm lượt.
   Không thể suy ra client đã nhận được response từ việc server commit.
4. Reservation thuộc command, không thuộc riêng một HTTP request. Khi callback
   lỗi, database timeout hoặc mất phản hồi, giữ lease để cleanup kiểm tra usage
   bền vững. Hủy ngay có thể mở lại quota trong lúc retry khác đang commit.
   `quota.AcceptanceContext` ràng buộc deadline trước khi lease hết hạn. Handler
   và store phải tôn trọng context; cleanup cần chạy để giải phóng lease lỗi.
5. Sau khi usage đã được xác nhận, lỗi cập nhật Redis là lỗi projection cần
   log/reconcile; không biến kết quả nghiệp vụ thành thất bại. Insert gặp usage
   cũ phải kiểm tra reservation sở hữu event trước khi commit Redis.
6. Comment tiếng Việt phải giải thích owner, bất biến nghiệp vụ và lý do chọn
   transaction; không viết comment lặp lại tên biến.
7. Thêm operation mới phải cập nhật registry, catalog seed, API contract, UI,
   test thành công/thất bại/retry và tài liệu vận hành trong cùng change.

## Kiểm tra khi review

Reviewer phải trả lời được: dữ liệu này thuộc service nào, retry có trừ hai
lần không, thất bại có làm mất quota không, plan version nào được lưu làm bằng
chứng, và người dùng nhận được kết quả nào sau khi thanh toán.


## Tiêu chí test có giá trị

Test phải thể hiện lỗi quan sát được: response thành công nhưng Redis lỗi,
insert đã commit nhưng mất phản hồi, hai retry cùng lease và một request lỗi,
command cũ được gửi sang kỳ mới, entitlement sai bị hiểu thành không tính quota.
Không dùng test chép lại danh sách constant để kết luận đã đồng bộ contract;
đối chiếu annotation protobuf với usage projection. File migration lịch sử
không thay đổi theo registry mới; kiểm tra checksum seed bảo vệ dữ liệu seed
của đúng phiên bản đó.

## Giới hạn cần nhớ

Read adapter hiện chưa lưu request hash và response snapshot. Vì vậy nó chưa
bảo đảm cùng key với payload khác trả conflict, cũng chưa bảo đảm replay trả
cùng nội dung byte. Đây là việc phải hoàn thiện trước khi coi read idempotency
là hợp đồng command đầy đủ. Không gọi retry đơn thuần là exactly-once delivery.

## Read command và retry

Read command phải bind `Subject + Operation + CommandKey` với hash request do
TQD tính từ protobuf deterministic. Bind xảy ra trước quota reservation và trước
business callback. Payload khác cùng key trả conflict. Retry cùng payload có thể
chạy lại read để lấy dữ liệu mới nhưng không được tạo usage event mới; nếu cần
cùng byte phải lưu response snapshot với retention được phê duyệt.

Frontend chỉ tạo key mới cho một interaction mới. `requestData` tự giữ key trong
các lần retry transport/auth; không sinh lại key trong interceptor. Preview công
khai dùng endpoint không tính quota, thao tác người dùng xác nhận mới dùng
operation có entitlement.
