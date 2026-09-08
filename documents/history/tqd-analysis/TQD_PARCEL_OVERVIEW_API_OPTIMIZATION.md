# TQD Parcel Overview API Optimization

## 1. Đặt Vấn Đề

Khi người dùng chạm vào bản đồ, hệ thống phải trả lời một câu hỏi rất đơn giản:

> "Ở vị trí này, tôi đang xem thửa đất hay vùng quy hoạch nào, và thông tin tổng quan đáng tin cậy nhất là gì?"

Nhưng trước khi tối ưu, luồng `ModalPropertyMap.content.tsx` đang phải gọi nhiều API cho một mục tiêu nghiệp vụ:

1. `GET /v2/tqd/parcels/{id}/info`
2. `GET /v2/tqd/parcels/{id}/overview`
3. `GET /v2/tqd/parcels/{id}/assessment?mode=detail`

Về mặt UI, người dùng chỉ mở popup tổng quan. Về mặt nghiệp vụ, đây là trạng thái free/public hoặc pre-login. Nhưng hệ thống lại kéo cả API `assessment`, vốn được định nghĩa là tầng phân tích chuyên sâu/VIP.

Điều này giống bài toán tàu cao tốc Shinkansen đi vào đường hầm: vấn đề nhìn bên ngoài là tiếng nổ âm thanh, nhưng bản chất kỹ thuật nằm ở cách mũi tàu nén không khí. Các kỹ sư Nhật Bản học từ hình dạng mỏ chim bói cá không phải để "làm đẹp" đầu tàu, mà để giải quyết đúng điểm phát sinh áp lực. Với API parcel cũng vậy: không phải thêm nhiều endpoint hơn, mà phải đưa dữ liệu overview về đúng lõi nghiệp vụ của nó.

## 2. Câu Hỏi Cần Đặt Ra

### 2.1 Người dùng thật sự cần gì ở popup overview?

Popup tổng quan cần trả lời:

- Đây là thửa đất hay vùng quy hoạch?
- Địa chỉ, số tờ, số thửa, diện tích là gì?
- Loại đất/quy hoạch chủ đạo là gì?
- Lớp quy hoạch nào ảnh hưởng chính?
- Tỷ lệ và diện tích overlap là bao nhiêu?
- Có cảnh báo/risk cơ bản không?
- Có đủ dữ liệu spatial để vẽ/zoom/highlight không?

Popup overview chưa cần:

- Conflict chi tiết nhiều lớp.
- Compare lịch sử.
- Historical risk.
- Assessment verdict chuyên sâu.
- Luồng trả phí/VIP.

### 2.2 Vì sao gọi `assessment` trong overview là sai bản chất?

`assessment` là tầng trả lời câu hỏi:

> "Nếu tôi muốn phân tích sâu, ra quyết định, so sánh, đánh giá rủi ro/pháp lý, hệ thống kết luận thế nào?"

Trong khi `overview` là tầng trả lời câu hỏi:

> "Tôi đang nhìn cái gì, loại đất chính là gì, các lớp ảnh hưởng tổng quan ra sao?"

Khi overview phải gọi assessment để lấy `currentPlan` hoặc `landUseGroups`, nghĩa là API overview đang thiếu dữ liệu tổng quan cần thiết. Client phải "mượn" API VIP để vá lỗ hổng free. Đây là dấu hiệu thiết kế sai ranh giới, không phải vấn đề UI.

### 2.3 Dữ liệu nào là bản chất, dữ liệu nào là projection?

Bản chất nghiệp vụ nằm ở pipeline:

1. Lấy parcel.
2. Tìm các region/layer intersect với parcel.
3. Tính overlap area/percent.
4. Chuẩn hóa land use/group/buildability/legal.
5. Rank ra quy hoạch chủ đạo.
6. Tính risk/cảnh báo.
7. Serialize ra response phù hợp từng tầng.

`overview` và `assessment` không nên tự có hai thuật toán riêng. Chúng nên là hai projection từ cùng một pipeline.

Ẩn dụ "đổ đầy thùng gạo/thùng nước" ở đây là: mục tiêu không phải chọn cái xẻng hay vòi nước trước, mà là làm đầy thùng đúng kiểm soát. Trong kỹ thuật, "thùng" là response nghiệp vụ đúng; "xẻng/vòi" là endpoint và query. Nếu endpoint bị chọn sai, dữ liệu có thể vẫn đầy nhưng tràn, thừa, khó kiểm soát, và dễ sai khi vận hành lâu dài.

## 3. Hiện Trạng Trước Khi Tối Ưu

### 3.1 Client gọi thừa API

Trong `react-native-quy-hoach/src/Screens/QHMapScreen/Modal/ModalPropertyMap.content.tsx`, `loadParcelData` gọi đồng thời:

```ts
tqdService.getParcelInfo(id)
tqdService.getParcelOverview(id)
tqdService.getParcelAssessment(id, 'detail')
```

Vấn đề:

- 3 round-trip cho một popup.
- `assessment` bị gọi dù người dùng chưa cần phân tích VIP.
- Risk leak nghiệp vụ: client quen phụ thuộc vào response VIP.
- Latency tăng theo API chậm nhất.
- Khi thay đổi assessment, overview UI có thể vỡ dù nghiệp vụ overview không đổi.

### 3.2 Backend còn hai pipeline gần giống nhau

`GetParcelOverview` và `GetParcelAssessment` đều làm các việc tương tự:

- Lấy rows từ repository.
- Lấy parcel info.
- Lấy legal docs.
- Convert rows sang unified rows.
- Gọi `parcelEngine.Process`.

Điểm lặp này làm phát sinh rủi ro:

- Sửa nguồn dữ liệu overview nhưng quên assessment.
- Sửa cách filter active/as-of-date nhưng chỉ tác động một nhánh.
- Log/error handling khác nhau.
- Khó biết "business rule thật" nằm ở đâu.

### 3.3 Overview dùng row V1 thiếu dữ liệu đất

Trước refactor, overview dùng `GetParcelLayerRows`, trong đó:

- `region_land_use_code` bị để rỗng.
- Tên loại đất lấy từ `qh_labels`.
- Thiếu land use group chuẩn hóa.
- Thiếu `can_build`, priority, build condition.

Trong khi assessment dùng `GetParcelLayerRowsV2`, có join `land_use_codes` và `land_use_groups`.

Kết quả là overview không đủ trả lời "đất gì, nhóm gì, có xây dựng được không" nên client phải gọi assessment.

## 4. Nguyên Tắc Thiết Kế Sau Tối Ưu

### 4.1 Một lõi xử lý, nhiều projection

Pipeline mới:

```text
parcel_id
  -> GetParcelLayerRowsV2
  -> apply preset/as_of_date if needed
  -> GetParcelInfo
  -> Get legal docs
  -> parcelEngine.Process
  -> projection:
       - overview/free
       - assessment/VIP
       - timeline/history
```

Điểm chính:

- `resolveParcelPlanning` là helper dùng chung trong usecase.
- `overview` và `assessment` không tự nhân đôi logic lấy dữ liệu.
- Engine vẫn là nơi quyết định primary planning, land use groups, risk.

### 4.2 Overview phải đủ dùng cho public/free

Overview response được mở rộng để trả về rõ hơn:

- `info`: metadata thửa đất.
- `summary`: thông tin tổng quan có thêm code/color/group/area/canBuild.
- `primaryPlan`: kết luận quy hoạch chủ đạo.
- `landUseGroups`: thống kê nhóm/loại đất theo overlap.
- `layers`: dữ liệu lớp/label/region để UI render chi tiết tổng quan.
- `riskLevel`, `riskScore`, `riskReasons`, `recommendation`: risk cơ bản.
- `conclude_label`: region overlap lớn nhất.

Overview không trả các phần chuyên sâu như compare/historical/conflict detail của assessment.

### 4.3 Assessment giữ vai trò VIP/chuyên sâu

Assessment vẫn tồn tại cho:

- `mode=detail`
- `mode=compare`
- conflict detail
- legal confidence nâng cao
- historical risk
- report/paid workflow sau này

Nhưng overview không còn phụ thuộc assessment.

## 5. API Trước/Sau

### 5.1 Trước

```text
Modal parcel popup
  -> /parcels/{id}/info
  -> /parcels/{id}/overview
  -> /parcels/{id}/assessment?mode=detail
  -> client tự merge
```

Hệ quả:

- Client biết quá nhiều về backend.
- UI phải merge nhiều response.
- Overview bị trộn với assessment.
- Khó kiểm soát quyền truy cập VIP.

### 5.2 Sau

```text
Modal parcel popup
  -> /parcels/{id}/planning?view=overview
  -> client render từ overview.info + overview.summary + overview.primaryPlan + overview.layers
```

Khi người dùng vào tính năng chuyên sâu:

```text
Detail
  -> /parcels/{id}/planning?view=detail
  -> response.overview + response.assessment

VIP/Assessment
  -> /parcels/{id}/planning?view=assessment
  -> response.assessment

Compare
  -> /parcels/{id}/planning?view=compare
  -> response.assessment
```

Kết quả:

- Popup giảm từ 3 API xuống 1 API.
- Detail giảm từ 3 API xuống 1 API.
- Compare đi qua cùng endpoint canonical, không gọi riêng legacy assessment.
- Backend giảm hai pipeline overview/assessment thành một lõi dùng chung.
- Overview đủ dữ liệu đất/quy hoạch tổng quan.
- Assessment được giữ độc lập để áp quyền/VIP sau này.

### 5.3 API canonical

```proto
rpc GetParcelPlanning(GetParcelPlanningRequest) returns (ParcelPlanningResponse) {
  option (google.api.http) = {
    get: "/v2/tqd/parcels/{parcelId}/planning"
  };
}
```

`view` là projection:

- `overview`: trả `overview`, dùng cho popup/free.
- `detail`: trả cả `overview` và `assessment`, dùng cho màn chi tiết.
- `assessment`: trả `assessment`, dùng cho luồng VIP/report.
- `compare`: trả `assessment` mode compare.

Các endpoint planning cũ đã được loại khỏi contract mới:

- bỏ `/v2/tqd/parcels/{parcelId}/overview`
- bỏ `/v2/tqd/parcels/{parcelId}/assessment`
- frontend không còn wrapper `getParcelOverview/getParcelAssessment/getParcelInfo`

`/v2/tqd/parcels/{parcelId}/info` vẫn có thể giữ cho metadata thuần, nhưng luồng parcel planning không phụ thuộc endpoint này nữa vì `overview.info` đã nằm trong response canonical.

## 6. Mapping Nghiệp Vụ

| Thành phần | Vai trò | Free overview | VIP assessment |
|---|---|---:|---:|
| Parcel info | Định danh thửa, địa chỉ, số tờ/thửa, diện tích | Có | Có thể dùng |
| Layers/labels/regions | Vẽ/hiển thị ảnh hưởng quy hoạch | Có | Có |
| Primary plan | Kết luận chủ đạo | Có | Có |
| Land use groups | Tổng hợp loại đất/nhóm đất | Có ở mức tổng quan | Có chi tiết |
| Basic risk | Cảnh báo nhanh | Có | Có |
| Conflict detail | Phân tích xung đột | Không | Có |
| Compare result | So sánh quy hoạch | Không | Có |
| Historical risk | Xu hướng/lịch sử | Không | Có |

## 7. Dẫn Chứng Từ Code

### 7.1 Client

File:

```text
react-native-quy-hoach/src/Screens/QHMapScreen/Modal/ModalPropertyMap.content.tsx
```

Thay đổi:

- `loadParcelData` chỉ gọi `getParcelPlanning(id, 'overview')`.
- Không gọi `getParcelInfo(id)` mặc định vì overview đã có `info`.
- Không gọi `getParcelAssessment(id, 'detail')` trong popup free.
- `mapOverviewToLayers` ưu tiên `overview.layers`, dùng `primaryPlan` để đánh dấu layer chính.
- Nếu backend cũ chưa có `primaryPlan`, client vẫn fallback từ `overview.layers`.
- `DetailScreen` dùng `getParcelPlanning(id, 'detail')` để nhận cả `overview` và `assessment` trong một response.
- Compare dùng `getParcelPlanning(id, 'compare')` và đọc `response.assessment`.
- Service frontend chỉ expose `getParcelPlanning` cho nghiệp vụ planning; wrapper cũ đã bị xóa để tránh màn mới vô tình quay lại endpoint cũ.

### 7.2 Proto contract

File:

```text
shared/protobuf/schema/tqd/parcel.proto
```

Thay đổi:

- `ParcelSummary` thêm các field rõ nghiệp vụ:
  - `dominantLandUseCode`
  - `dominantLandUseColor`
  - `dominantLandUseGroupCode`
  - `dominantLandUseGroupName`
  - `dominantAreaSqm`
  - `dominantCanBuild`
- `ParcelOverviewResponse` thêm:
  - `primaryPlan`
  - `landUseGroups`
- `PlanDetailInfo` thêm:
  - `groupName`
  - `canBuild`
- RPC `GetParcelOverview` và `GetParcelAssessment` bị gỡ khỏi `ParcelService`. `GetParcelPlanning` là contract duy nhất cho planning overview/detail/assessment/compare.

### 7.3 Backend usecase

File:

```text
tqd-service/internal/usecase/parcel_usecase.go
```

Thay đổi:

- Thêm `GetParcelPlanning` làm usecase canonical.
- Thêm `resolveParcelPlanning` làm lõi resolver dùng chung.
- Gỡ `GetParcelOverview`, `GetParcelAssessment`, `GetPublicParcelInfo` khỏi usecase interface để không còn nhiều entrypoint cùng một nghiệp vụ.
- `unifiedToParcelLayersResponse` build `PrimaryPlan`, `LandUseGroups`, summary mở rộng.

### 7.4 Repository SQL

File:

```text
tqd-service/infra/postgres/parcel_postgres.go
```

Thay đổi:

- Overview chuyển sang dùng `GetParcelLayerRowsV2`.
- V2 query bổ sung authority issuing.
- V2 query trả thêm `region_land_use_group` và `region_land_use_color`.

## 8. Vấn Đề Phụ Thuộc Và Không Phụ Thuộc

Bài toán bạn nêu về "nét vẽ phụ thuộc và không phụ thuộc" có thể hiểu như sau:

- Người dùng có quyền vẽ/chọn theo ý định tự nhiên.
- Backend có trách nhiệm chuẩn hóa kết quả để thống kê đúng.
- Hệ thống không nên ép người dùng thay đổi thao tác chỉ vì backend khó tính.
- Nhưng hệ thống phải có chuẩn để tính toán: polygon, parcel, region, overlap, filter.

Trong parcel modal:

- User chỉ chạm parcel.
- UI không nên tự quyết định gọi assessment hay merge ba nguồn.
- Backend phải trả overview đúng chuẩn.

Trong polygon/drawing analysis:

- User vẽ vùng.
- Backend chuẩn hóa polygon, filter theo zoom/layer/land use.
- Kết quả cuối là thống kê đúng, không phải ép người dùng vẽ lại nếu hệ thống có thể xử lý an toàn.

Nguyên tắc vận hành:

```text
Intent của user độc lập với thuật toán nội bộ.
Thuật toán nội bộ phụ thuộc vào chuẩn dữ liệu.
Response public không phụ thuộc vào response VIP.
Assessment phụ thuộc vào cùng core resolver, nhưng không bị overview gọi ngược.
```

## 9. Rủi Ro Nếu Không Tối Ưu

### 9.1 Rủi ro vận hành

- API popup chậm vì phụ thuộc endpoint nặng nhất.
- Tăng tải DB khi mỗi click parcel gọi nhiều query trùng.
- Khó cache vì dữ liệu bị tách thành nhiều endpoint.
- Khó kiểm quota do assessment bị gọi ngầm.

### 9.2 Rủi ro nghiệp vụ

- Free user có thể nhận dữ liệu vốn thuộc VIP.
- Khi triển khai paywall, phải gỡ phụ thuộc assessment khỏi nhiều UI.
- Overview không còn là contract rõ ràng.

### 9.3 Rủi ro maintainability

- Một rule thay đổi phải sửa nhiều nơi.
- V1/V2 row diverge, dẫn đến overview và assessment kết luận khác nhau.
- Client merge response thủ công, dễ sai field/case.

## 10. Kết Quả Đạt Được

### 10.1 Về API

- Popup parcel từ 3 API về 1 API.
- Detail parcel từ 3 API về 1 API.
- Overview tự đủ dữ liệu tổng quan.
- Assessment không bị gọi thừa.

### 10.2 Về backend

- `GetParcelPlanning` là entrypoint nghiệp vụ chính.
- Overview và assessment dùng chung pipeline `resolveParcelPlanning`.
- Overview dùng row V2 có land use/group/buildability chuẩn hóa.
- Mapper chỉ serialize, không chứa business rule.

### 10.3 Về vận hành lâu dài

- Dễ bật paywall cho assessment.
- Dễ cache overview theo parcel ID.
- Dễ test core resolver độc lập.
- Dễ mở rộng region/parcel/polygon analysis theo cùng mô hình.

## 11. Câu Hỏi Kiểm Tra Sau Refactor

Khi review hoặc phát triển tiếp, nên hỏi:

1. Popup này là free overview hay VIP assessment?
2. Field UI đang cần có thuộc overview hay assessment?
3. Nếu là overview, vì sao overview chưa trả field đó?
4. Dữ liệu này là nguyên liệu DB, kết luận engine, hay projection response?
5. Có đang gọi nhiều endpoint chỉ để merge thành một mô hình UI không?
6. Có endpoint public nào đang phụ thuộc dữ liệu private/VIP không?
7. Nếu đổi rule land use/buildability, cần sửa một nơi hay nhiều nơi?
8. Nếu parcel không intersect region nào, response có rõ là "không có ảnh hưởng" hay bị hiểu là lỗi?
9. Nếu user chưa đăng nhập, response có đủ hiển thị nhưng không lộ phân tích trả phí không?

## 12. Chiến Lược Mở Rộng

### 12.1 Ngắn hạn

- Modal parcel dùng `GetParcelPlanning(view=overview)`.
- Detail dùng `GetParcelPlanning(view=detail)`.
- VIP/compare dùng `GetParcelPlanning(view=assessment|compare)`.
- Không thêm màn mới nào gọi `/overview` hoặc `/assessment`; nếu cần metadata thuần thì tách khỏi planning, không dùng để ghép UI quy hoạch.

### 12.2 Trung hạn

- Thêm cache overview theo parcel ID và version layer.
- Thêm field `accessLevel` hoặc `entitlement` nếu cần phân tầng rõ free/pro/vip.
- Tách `overview` và `assessment` docs trong API gateway.

### 12.3 Dài hạn

- Xây resolver chung cho parcel/region/polygon:

```text
target resolver
  -> spatial resolver
  -> planning resolver
  -> projection resolver
```

- Mỗi use case chỉ chọn projection:
  - map popup
  - detail page
  - report
  - compare
  - SEO source

## 13. Kết Luận

Tối ưu đúng không phải là xóa bớt API một cách cơ học. Tối ưu đúng là trả dữ liệu ở đúng tầng nghiệp vụ.

Với bài toán này:

- `overview` là mũi tàu cần được thiết kế lại để xử lý áp lực popup/free.
- `assessment` là khoang phân tích chuyên sâu, không nên bị kéo vào mọi cú click.
- `parcelEngine` là lõi kỹ thuật, cần là nguồn sự thật duy nhất cho ranking/risk/land use.
- Client nên nhận một response rõ ràng, không phải tự ghép nhiều response để đoán nghiệp vụ.

Khi kiến trúc đi theo hướng này, hệ thống dễ maintain hơn, dễ mở rộng VIP hơn, ít lặp logic hơn, và quan trọng nhất là mỗi endpoint có một lý do tồn tại rõ ràng.
