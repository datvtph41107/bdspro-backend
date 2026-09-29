PHẦN A.  NGUYÊN TẮC SEO QH PRO
A.1	Mục tiêu tài liệu
A.1.1 Mục tiêu SEO của hệ thống
QH Pro là nền tảng dữ liệu số về địa chính, quy hoạch, đất đai và bất động sản. Tài liệu này được xây dựng nhằm đặc tả các yêu cầu SEO cần bổ sung cho hệ thống để bảo đảm dữ liệu và nội dung của QH Pro có khả năng được tìm kiếm, lập chỉ mục, liên kết và khai thác hiệu quả trên các công cụ tìm kiếm và nền tảng AI hiện đại.
Mục tiêu SEO của hệ thống bao gồm:ss
•	Tăng khả năng xuất hiện trên các công cụ tìm kiếm như Google, Bing và các nền tảng tìm kiếm tương tự. 
•	Tăng khả năng được nhận diện, trích xuất và sử dụng bởi các hệ thống AI Search, AI Assistant và các mô hình ngôn ngữ lớn (LLM). 
•	Xây dựng cấu trúc dữ liệu và liên kết ngữ nghĩa (Semantic SEO) phục vụ Knowledge Graph và Entity Search. 
•	Tối ưu khả năng khám phá dữ liệu, điều hướng và liên kết nội bộ giữa các đối tượng dữ liệu trong hệ thống. 
•	Bảo đảm việc triển khai SEO tuân thủ các yêu cầu pháp lý, bảo mật, phân quyền và phạm vi công khai dữ liệu của QH Pro. 
Tài liệu không tập trung vào hoạt động Marketing SEO, Content Marketing hay vận hành SEO, mà tập trung vào các yêu cầu kỹ thuật cần được triển khai trực tiếp trong hệ thống.
 
A.1.2 Phạm vi áp dụng
Tài liệu này áp dụng cho toàn bộ các thành phần thuộc phạm vi công khai hoặc có khả năng công khai của hệ thống QH Pro, bao gồm:
•	Các màn hình tra cứu, hiển thị và khai thác dữ liệu quy hoạch. 
•	Các trang chi tiết đối tượng dữ liệu như thửa đất, khu quy hoạch, đồ án quy hoạch, văn bản pháp lý và các thực thể liên quan. 
•	Các thư viện dữ liệu, kho tri thức quy hoạch và các nội dung phục vụ tìm kiếm. 
•	Các URL, Metadata, Structured Data, Open Graph, Sitemap, AI Summary và các thành phần kỹ thuật SEO khác được sinh ra từ hệ thống. 
Tài liệu không áp dụng cho:
•	Các màn hình quản trị nội bộ. 
•	Các API nội bộ không công khai. 
•	Các chức năng vận hành hệ thống. 
•	Các dữ liệu không đủ điều kiện công khai hoặc không được phép lập chỉ mục. 
•	Các trạng thái giao diện tạm thời, popup kỹ thuật hoặc dữ liệu phát sinh chỉ phục vụ trải nghiệm người dùng. 
A.1.3 Đối tượng sử dụng tài liệu
Tài liệu được xây dựng để phục vụ các nhóm tham gia phát triển và vận hành hệ thống QH Pro, bao gồm:
•	Nhóm Phân tích nghiệp vụ (BA)
•	Xác định phạm vi SEO của từng màn hình và chức năng. 
•	Bảo đảm các yêu cầu SEO được mô tả đầy đủ trong SRS. 
•	Nhóm Thiết kế UI/UX
•	Hiểu các yêu cầu liên quan đến SEO Content Block, Internal Link, Breadcrumb, Structured Content và các thành phần hỗ trợ SEO. 
•	Nhóm Frontend
•	Triển khai Metadata, Canonical, Structured Data, Open Graph, SSR/Hybrid Rendering và các thành phần hiển thị liên quan đến SEO. 
•	Nhóm Backend
•	Xây dựng các dịch vụ sinh URL, Metadata, Sitemap, AI Summary và các SEO Output khác. 
•	Đảm bảo dữ liệu SEO được cung cấp đầy đủ và chính xác cho Frontend. 
•	Nhóm QA
•	Kiểm thử và nghiệm thu các yêu cầu SEO. 
•	Đánh giá tính chính xác, đầy đủ và tuân thủ của SEO Output. 
•	Nhóm Quản lý sản phẩm
•	Kiểm soát định hướng SEO tổng thể. 
•	Đảm bảo việc triển khai SEO thống nhất trên toàn hệ thống. 
 
A.1.4 Quan hệ với bộ SRS QH Pro
Tài liệu này là phụ lục kỹ thuật bổ sung cho bộ tài liệu SRS QH Pro hiện có.
Các nội dung nghiệp vụ, giao diện, luồng xử lý, mô hình dữ liệu, API nghiệp vụ và các yêu cầu hệ thống vẫn được mô tả trong các tài liệu SRS gốc.
Phụ lục SEO chỉ bổ sung các yêu cầu liên quan đến:
•	SEO Scope. 
•	URL và Canonical. 
•	Metadata. 
•	Structured Data. 
•	Open Graph. 
•	AI Summary. 
•	Sitemap. 
•	Internal Link. 
•	Security & Visibility liên quan đến SEO. 
•	Acceptance Criteria và QA Checklist cho SEO. 
Nguyên tắc áp dụng:
•	Không mô tả lại các nội dung đã tồn tại trong SRS. 
•	Không thay thế SRS. 
•	Không sửa đổi nghiệp vụ đã được phê duyệt trong SRS. 
•	Chỉ bổ sung các yêu cầu kỹ thuật cần thiết để triển khai SEO. 
Khi có khác biệt giữa nội dung SRS và Phụ lục SEO, các yêu cầu nghiệp vụ và chức năng trong SRS được ưu tiên làm căn cứ triển khai; Phụ lục SEO chỉ quy định cách thức áp dụng SEO trên các chức năng đó.
 
A.1.5 Quan hệ với Backend / Frontend / Mobile
SEO là yêu cầu xuyên suốt nhiều lớp của hệ thống và không thuộc riêng Backend, Frontend hay Mobile.
•	Backend
Backend chịu trách nhiệm:
•	Sinh URL chuẩn. 
•	Sinh Canonical URL. 
•	Sinh Metadata. 
•	Sinh Structured Data. 
•	Sinh Sitemap. 
•	Sinh AI Summary. 
•	Quản lý trạng thái Index / NoIndex. 
•	Cung cấp SEO Payload cho các ứng dụng sử dụng. 
•	Frontend Web
Frontend Web chịu trách nhiệm:
•	Hiển thị và render các SEO Output do Backend cung cấp. 
•	Triển khai Head Tags. 
•	Triển khai Structured Data. 
•	Triển khai Open Graph. 
•	Triển khai Internal Link và các thành phần hỗ trợ SEO. 
•	Đảm bảo nội dung SEO có thể được công cụ tìm kiếm thu thập và lập chỉ mục. 
•	Mobile App
Ứng dụng Mobile không phải là đối tượng SEO trực tiếp của công cụ tìm kiếm web.
Tuy nhiên Mobile App phải hỗ trợ:
•	Deep Link. 
•	Share Link. 
•	Open Graph. 
•	Đồng bộ URL và Entity với phiên bản Web. 
•	Các cơ chế điều hướng phục vụ Discover, Share và AI Search. 
Mọi yêu cầu SEO được đặc tả trong tài liệu này phải được xem xét trên cả ba lớp Backend, Frontend và Mobile để bảo đảm tính nhất quán của hệ thống.
A.2	Mục tiêu SEO tổng thể
Mục này xác định các mục tiêu SEO cấp hệ thống mà QH Pro cần đạt được. Đây là cơ sở để xây dựng các yêu cầu SEO trong toàn bộ tài liệu, đồng thời là căn cứ đánh giá hiệu quả triển khai SEO của hệ thống.
Các mục tiêu dưới đây không mô tả giải pháp kỹ thuật cụ thể mà tập trung xác định kết quả cuối cùng cần đạt được đối với dữ liệu, nội dung và các thực thể (Entity) được quản lý trong QH Pro.
QH Pro không định hướng SEO theo mô hình Website nội dung hay Website tin tức truyền thống, mà theo mô hình Nền tảng dữ liệu số (Data Platform), trong đó các đối tượng dữ liệu có khả năng được định danh, liên kết, tra cứu, lập chỉ mục và khai thác trên các công cụ tìm kiếm và nền tảng AI.
 
A.2.1 Google Search
Mục tiêu của QH Pro là xây dựng một hệ thống có khả năng được các công cụ tìm kiếm thu thập, lập chỉ mục và hiển thị hiệu quả đối với các dữ liệu và nội dung công khai thuộc lĩnh vực quy hoạch, địa chính, đất đai và bất động sản.
Các nhóm thực thể trọng tâm cần được tối ưu cho Google Search bao gồm:
•	Đơn vị hành chính. 
•	Đơn vị hành chính lịch sử, đơn vị trước sáp nhập, đổi tên hoặc điều chỉnh địa giới. 
•	Thửa đất. 
•	Khu quy hoạch. 
•	Đồ án quy hoạch. 
•	Dự án quy hoạch. 
•	Văn bản pháp lý. 
•	Bản đồ quy hoạch. 
•	Báo cáo và kết quả phân tích. 
Hệ thống phải bảo đảm:
•	URL rõ ràng, ổn định và có ý nghĩa ngữ nghĩa. 
•	Mỗi thực thể quan trọng có khả năng được lập chỉ mục độc lập. 
•	Hạn chế tối đa nội dung trùng lặp. 
•	Duy trì khả năng truy cập đối với các tên gọi lịch sử, địa danh cũ và đơn vị hành chính trước sáp nhập. 
•	Hỗ trợ mở rộng SEO theo dữ liệu được bổ sung trong tương lai. 
Mục tiêu cuối cùng là xây dựng QH Pro trở thành nguồn tham chiếu dữ liệu quy hoạch, địa chính và đất đai có khả năng tiếp cận rộng rãi thông qua các công cụ tìm kiếm.
 
A.2.2 AI Search
Bên cạnh công cụ tìm kiếm truyền thống, QH Pro phải được xây dựng theo hướng sẵn sàng cho các nền tảng AI Search, AI Assistant và các mô hình ngôn ngữ lớn (LLM).
Các nền tảng mục tiêu bao gồm:
•	ChatGPT. 
•	Gemini. 
•	Copilot. 
•	Perplexity. 
•	Các hệ thống AI Search và LLM khác. 
Hệ thống cần bảo đảm:
•	Nội dung có cấu trúc rõ ràng. 
•	Dữ liệu có khả năng đọc hiểu bằng máy. 
•	Các thực thể được định danh nhất quán. 
•	Các quan hệ giữa thực thể được thể hiện rõ ràng. 
•	Thông tin quan trọng có thể được trích xuất, tóm tắt và tham chiếu tự động. 
Mục tiêu là giúp các hệ thống AI có thể:
•	Hiểu đúng các thực thể và dữ liệu của QH Pro. 
•	Hiểu đúng quan hệ giữa địa bàn, thửa đất, khu quy hoạch, đồ án và văn bản pháp lý. 
•	Trích dẫn QH Pro như một nguồn dữ liệu tham khảo. 
•	Trả lời các truy vấn liên quan đến quy hoạch, địa chính và đất đai dựa trên dữ liệu của hệ thống. 
 
A.2.3 Knowledge Graph
QH Pro hướng tới xây dựng một hệ sinh thái dữ liệu có liên kết ngữ nghĩa giữa các thực thể thay vì các trang nội dung rời rạc.
Các nhóm thực thể trọng tâm bao gồm:
•	Đơn vị hành chính. 
•	Đơn vị hành chính lịch sử. 
•	Thửa đất. 
•	Khu quy hoạch. 
•	Dự án quy hoạch. 
•	Đồ án quy hoạch. 
•	Văn bản pháp lý. 
•	Bản đồ quy hoạch. 
•	Báo cáo và kết quả phân tích. 
Giữa các thực thể phải tồn tại các quan hệ rõ ràng như:
•	Thuộc về. 
•	Áp dụng cho. 
•	Nằm trong. 
•	Kế thừa. 
•	Thay thế. 
•	Sáp nhập vào. 
•	Được điều chỉnh bởi. 
•	Được tạo bởi. 
•	Được phê duyệt bởi. 
•	Liên quan. 
Mục tiêu là giúp:
•	Công cụ tìm kiếm hiểu được ngữ cảnh dữ liệu. 
•	Hệ thống AI hiểu được mối quan hệ giữa các thực thể. 
•	Tăng khả năng hiển thị các kết quả giàu ngữ nghĩa (Rich Result). 
•	Hình thành nền tảng dữ liệu phục vụ Semantic Search và AI Search trong dài hạn. 
 
A.2.4 Internal Linking
QH Pro phải được xây dựng như một hệ thống dữ liệu liên kết, trong đó các thực thể có khả năng dẫn chiếu và điều hướng lẫn nhau.
Các liên kết nội bộ cần hình thành giữa:
•	Quốc gia ↔ Tỉnh ↔ Quận/Huyện ↔ Xã/Phường. 
•	Đơn vị hành chính hiện hành ↔ Đơn vị hành chính lịch sử. 
•	Địa bàn ↔ Thửa đất. 
•	Địa bàn ↔ Khu quy hoạch. 
•	Thửa đất ↔ Khu quy hoạch. 
•	Khu quy hoạch ↔ Đồ án quy hoạch. 
•	Đồ án quy hoạch ↔ Văn bản pháp lý. 
•	Văn bản pháp lý ↔ Bản đồ quy hoạch. 
•	Báo cáo ↔ Dữ liệu nguồn. 
•	Các thực thể liên quan khác trong hệ thống. 
Mục tiêu của Internal Linking bao gồm:
•	Tăng khả năng khám phá dữ liệu. 
•	Tăng khả năng thu thập dữ liệu của công cụ tìm kiếm. 
•	Truyền tải ngữ cảnh giữa các thực thể. 
•	Hỗ trợ hình thành Knowledge Graph nội bộ. 
•	Tăng khả năng điều hướng và khai thác dữ liệu của người dùng. 
Nguyên tắc chung là mọi thực thể SEO quan trọng trong hệ thống đều phải có khả năng liên kết tới các thực thể liên quan một cách tự nhiên và có ý nghĩa.
 
A.2.5 Discover & Share
QH Pro phải hỗ trợ khả năng chia sẻ, lan truyền và tái sử dụng dữ liệu trên các nền tảng số.
Các kịch bản trọng tâm bao gồm:
•	Chia sẻ thửa đất. 
•	Chia sẻ khu quy hoạch. 
•	Chia sẻ đồ án quy hoạch. 
•	Chia sẻ văn bản pháp lý. 
•	Chia sẻ bản đồ. 
•	Chia sẻ báo cáo. 
•	Chia sẻ Snapshot. 
•	Chia sẻ kết quả phân tích. 
Hệ thống cần bảo đảm:
•	Mỗi nội dung có URL ổn định và có khả năng tham chiếu lâu dài. 
•	Nội dung chia sẻ có tiêu đề, mô tả và bản xem trước phù hợp. 
•	Hỗ trợ chia sẻ trên công cụ tìm kiếm, mạng xã hội, ứng dụng nhắn tin và nền tảng AI. 
•	Không làm lộ dữ liệu bị hạn chế truy cập hoặc dữ liệu không đủ điều kiện công khai. 
•	Bảo đảm tính nhất quán giữa Web, Mobile và các hình thức chia sẻ khác. 
Mục tiêu cuối cùng là giúp dữ liệu của QH Pro có thể được khám phá, tham chiếu, chia sẻ và tái sử dụng một cách hiệu quả trong toàn bộ hệ sinh thái số.
A.3	Nguyên tắc triển khai SEO
Mục này xác định các nguyên tắc nền tảng phải được áp dụng xuyên suốt trong toàn bộ quá trình thiết kế, phát triển và vận hành SEO của QH Pro.
Các nguyên tắc dưới đây có tính chất bắt buộc và được sử dụng làm căn cứ để:
•	Thiết kế URL. 
•	Thiết kế Metadata. 
•	Thiết kế Structured Data. 
•	Thiết kế Sitemap. 
•	Thiết kế Internal Link. 
•	Thiết kế AI Summary. 
•	Đánh giá khả năng Index / NoIndex. 
•	Đánh giá tính hợp lệ của các yêu cầu SEO trong từng màn hình SRS. 
Khi có nhiều phương án triển khai khác nhau, các nguyên tắc tại mục này được ưu tiên sử dụng để lựa chọn phương án phù hợp.
 
A.3.1 Entity First
QH Pro là nền tảng dữ liệu, không phải website nội dung thông thường. Vì vậy mọi hoạt động SEO phải được xây dựng xoay quanh các thực thể (Entity) của hệ thống thay vì xoay quanh giao diện hoặc từ khóa đơn lẻ.
Các thực thể SEO trọng tâm bao gồm:
•	Đơn vị hành chính. 
•	Đơn vị hành chính lịch sử. 
•	Thửa đất. 
•	Khu quy hoạch. 
•	Đồ án quy hoạch. 
•	Dự án quy hoạch. 
•	Văn bản pháp lý. 
•	Bản đồ quy hoạch. 
•	Báo cáo. 
•	Snapshot. 
Nguyên tắc áp dụng:
•	Mỗi Entity quan trọng phải có khả năng được định danh độc lập. 
•	Mỗi Entity phải có URL tham chiếu ổn định. 
•	Metadata, Structured Data và AI Summary phải được sinh từ Entity. 
•	Internal Link phải được xây dựng dựa trên quan hệ giữa các Entity. 
•	Không xây dựng SEO dựa trên trạng thái giao diện hoặc tổ hợp bộ lọc tạm thời. 
Entity là đơn vị trung tâm của toàn bộ kiến trúc SEO QH Pro.
 
A.3.2 Canonical First
Mỗi nội dung công khai trong hệ thống chỉ được tồn tại một URL chuẩn duy nhất dùng làm nguồn tham chiếu chính thức.
Nguyên tắc này nhằm:
•	Hạn chế trùng lặp nội dung. 
•	Tránh phân tán giá trị SEO. 
•	Tăng tính nhất quán của dữ liệu. 
•	Bảo đảm khả năng tham chiếu lâu dài. 
Các URL phát sinh từ:
•	Bộ lọc. 
•	Tìm kiếm. 
•	Chia sẻ. 
•	Snapshot. 
•	So sánh. 
•	Trạng thái giao diện. 
phải được đánh giá và xác định rõ:
•	Có được Index hay không. 
•	Có Canonical về URL gốc hay không. 
•	Có xuất hiện trong Sitemap hay không. 
Trong mọi trường hợp, Canonical URL phải phản ánh chính xác Entity hoặc nội dung gốc mà hệ thống muốn công cụ tìm kiếm lập chỉ mục.
 
A.3.3 Public / Permission First
SEO chỉ được áp dụng đối với các dữ liệu đủ điều kiện công khai.
Mọi quyết định liên quan đến Index, Metadata, Sitemap hoặc AI Summary phải được kiểm soát bởi trạng thái công khai và phân quyền dữ liệu.
Nguyên tắc áp dụng:
•	Public Data có thể được xem xét SEO. 
•	Restricted Data phải được kiểm soát theo chính sách công bố dữ liệu. 
•	Private Data không được Index. 
•	Nội dung yêu cầu đăng nhập không được xem là nội dung SEO mặc định. 
Việc tối ưu SEO không được làm thay đổi hoặc làm suy giảm các cơ chế bảo mật và phân quyền của hệ thống.
 
A.3.4 NoIndex First
Nguyên tắc mặc định của QH Pro là:
Chỉ Index khi dữ liệu đủ điều kiện.
Không Index khi chưa xác định rõ điều kiện công khai và chất lượng dữ liệu.
Các trường hợp mặc định NoIndex bao gồm:
•	Dữ liệu chưa hoàn chỉnh. 
•	Dữ liệu nháp. 
•	Dữ liệu tạm thời. 
•	Kết quả tìm kiếm phát sinh động. 
•	Trạng thái giao diện tạm thời. 
•	Nội dung nội bộ. 
•	Nội dung yêu cầu đăng nhập. 
•	Nội dung không đủ thông tin để hình thành một thực thể độc lập. 
Việc chuyển từ NoIndex sang Index phải được xác định rõ trong từng đặc tả SEO của màn hình tương ứng.
 
A.3.5 Legal Safe
QH Pro hoạt động trong lĩnh vực quy hoạch, địa chính, đất đai và dữ liệu pháp lý, do đó mọi hoạt động SEO phải tuân thủ yêu cầu pháp lý và trách nhiệm công bố thông tin.
Nguyên tắc áp dụng:
•	Chỉ SEO các dữ liệu được phép công khai. 
•	Không SEO dữ liệu thuộc diện hạn chế truy cập. 
•	Không SEO dữ liệu có rủi ro vi phạm quyền riêng tư. 
•	Không SEO các thông tin nội bộ hoặc thông tin chưa được công bố chính thức. 
•	Không tạo Metadata hoặc AI Summary gây hiểu nhầm về giá trị pháp lý của dữ liệu. 
Trong trường hợp có xung đột giữa mục tiêu SEO và yêu cầu pháp lý, yêu cầu pháp lý luôn được ưu tiên.
 
A.3.6 AI Search Ready
SEO của QH Pro không chỉ phục vụ công cụ tìm kiếm truyền thống mà còn phải phục vụ các hệ thống AI Search và AI Assistant.
Nguyên tắc áp dụng:
•	Dữ liệu phải có cấu trúc rõ ràng. 
•	Các Entity phải được định danh nhất quán. 
•	Các quan hệ giữa Entity phải được thể hiện minh bạch. 
•	Nội dung quan trọng phải có khả năng được trích xuất và tóm tắt. 
•	Hệ thống phải hỗ trợ khả năng đọc hiểu bằng máy. 
Mục tiêu là giúp các nền tảng AI:
•	Hiểu đúng dữ liệu của QH Pro. 
•	Hiểu đúng ngữ cảnh dữ liệu. 
•	Trích dẫn QH Pro như nguồn tham khảo đáng tin cậy. 
 
A.3.7 Internal Link First
Liên kết nội bộ là một thành phần cốt lõi trong kiến trúc SEO của QH Pro.
Mọi Entity quan trọng phải được đặt trong một mạng lưới liên kết có cấu trúc thay vì tồn tại độc lập.
Các liên kết trọng tâm bao gồm:
•	Đơn vị hành chính ↔ Đơn vị hành chính lịch sử. 
•	Địa bàn ↔ Thửa đất. 
•	Địa bàn ↔ Khu quy hoạch. 
•	Thửa đất ↔ Khu quy hoạch. 
•	Khu quy hoạch ↔ Đồ án quy hoạch. 
•	Đồ án ↔ Văn bản pháp lý. 
•	Văn bản ↔ Bản đồ. 
•	Báo cáo ↔ Dữ liệu nguồn. 
Nguyên tắc này nhằm:
•	Tăng khả năng khám phá dữ liệu. 
•	Tăng khả năng thu thập dữ liệu của công cụ tìm kiếm. 
•	Tăng giá trị Semantic SEO. 
•	Hỗ trợ hình thành Knowledge Graph nội bộ. 
 
A.3.8 Không SEO theo trạng thái giao diện tạm thời
SEO phải phản ánh dữ liệu và thực thể của hệ thống, không phản ánh trạng thái sử dụng tức thời của người dùng.
Các trạng thái giao diện sau đây không được xem là đối tượng SEO độc lập:
•	Popup. 
•	Drawer. 
•	Tooltip. 
•	Bộ lọc tạm thời. 
•	Kết quả tìm kiếm tạm thời. 
•	Trạng thái zoom bản đồ. 
•	Trạng thái bật/tắt layer. 
•	Trạng thái lựa chọn đối tượng trên bản đồ. 
•	Trạng thái làm việc cá nhân của người dùng. 
Chỉ khi một nội dung được xác định là một Entity hoặc một URL công khai có giá trị độc lập thì mới được xem xét xây dựng SEO riêng.
Nguyên tắc này giúp bảo đảm:
•	Kiến trúc SEO ổn định. 
•	Hạn chế URL rác. 
•	Hạn chế nội dung trùng lặp. 
•	Duy trì khả năng mở rộng lâu dài của hệ thống.
A.4	SEO Class chuẩn
SEO Class là cơ chế phân loại mức độ ưu tiên SEO đối với từng màn hình, URL và thực thể trong hệ thống QH Pro.
Việc phân loại SEO Class nhằm:
•	Xác định phạm vi Index và NoIndex. 
•	Xác định phạm vi Sitemap. 
•	Xác định mức độ đầu tư SEO. 
•	Xác định phạm vi triển khai Metadata, Structured Data và AI Summary. 
•	Thống nhất cách triển khai SEO trên toàn bộ hệ thống. 
Mỗi màn hình, URL hoặc Entity phải được gán tối thiểu một SEO Class.
Trong trường hợp có nhiều phương án triển khai, SEO Class được sử dụng làm căn cứ ưu tiên.
 
A.4.1 SEO-A – Landing Page trọng tâm
SEO-A là nhóm đối tượng SEO quan trọng nhất của hệ thống.
Đây là các URL và màn hình được ưu tiên cao nhất về:
•	Index. 
•	Sitemap. 
•	Structured Data. 
•	AI Search. 
•	Internal Linking. 
•	Discover. 
•	Đặc điểm
•	Có URL SEO độc lập. 
•	Có giá trị tìm kiếm cao. 
•	Có khả năng thu hút lượng truy cập lớn. 
•	Có nội dung ổn định và có khả năng tồn tại lâu dài. 
•	Đại diện cho các Entity trọng tâm của hệ thống. 
•	Yêu cầu
•	Được Index. 
•	Được đưa vào Sitemap. 
•	Có Canonical URL riêng. 
•	Có Metadata đầy đủ. 
•	Có Structured Data. 
•	Có Open Graph. 
•	Có AI Summary. 
•	Có Internal Link đến các Entity liên quan. 
•	Ví dụ
•	Trang chi tiết đơn vị hành chính. 
•	Trang chi tiết thửa đất. 
•	Trang chi tiết khu quy hoạch. 
•	Trang chi tiết đồ án quy hoạch. 
•	Trang chi tiết văn bản pháp lý. 
•	Trang chi tiết bản đồ quy hoạch. 
•	Trang chi tiết báo cáo. 
Đây là nhóm URL mang lại phần lớn giá trị SEO của QH Pro.
 
A.4.2 SEO-B – Trang điều hướng
SEO-B là nhóm trang hỗ trợ điều hướng và khám phá dữ liệu.
Các trang này có giá trị SEO nhưng không phải đích đến cuối cùng của người dùng.
Mục tiêu chính là:
•	Điều hướng. 
•	Gom nhóm dữ liệu. 
•	Tăng Internal Linking. 
•	Hỗ trợ thu thập dữ liệu. 
•	Đặc điểm
•	Có URL riêng. 
•	Có khả năng được Index. 
•	Chủ yếu hiển thị danh sách hoặc tập hợp dữ liệu. 
•	Dẫn người dùng đến các Entity SEO-A. 
•	Yêu cầu
•	Có thể được Index. 
•	Có thể xuất hiện trong Sitemap. 
•	Có Metadata. 
•	Có Internal Link đầy đủ. 
•	Có thể sử dụng Structured Data phù hợp. 
•	Ví dụ
•	Danh sách đồ án quy hoạch. 
•	Danh sách văn bản pháp lý. 
•	Danh sách khu quy hoạch. 
•	Danh sách bản đồ quy hoạch. 
•	Danh sách đơn vị hành chính. 
•	Trang thư viện quy hoạch. 
SEO-B có vai trò hỗ trợ SEO-A và tăng khả năng khám phá dữ liệu trong hệ thống.
 
A.4.3 SEO-C – Trang SEO có điều kiện
SEO-C là nhóm đối tượng chỉ được SEO khi đáp ứng đủ điều kiện về dữ liệu, nội dung hoặc khả năng công khai.
Không phải mọi URL thuộc nhóm này đều được Index.
•	Đặc điểm
•	Giá trị SEO phụ thuộc vào dữ liệu thực tế. 
•	Có thể chuyển đổi giữa Index và NoIndex. 
•	Có khả năng phát sinh số lượng lớn URL. 
•	Yêu cầu
•	Chỉ được Index khi đủ điều kiện. 
•	Chỉ được đưa vào Sitemap khi đủ điều kiện. 
•	Phải kiểm soát chặt chẽ Canonical. 
•	Phải kiểm soát chất lượng nội dung. 
•	Ví dụ
•	Kết quả tra cứu địa chỉ. 
•	Kết quả tra cứu GPS. 
•	Kết quả tra cứu tờ/thửa. 
•	Trang vùng phân tích. 
•	Trang tổng hợp dữ liệu theo địa bàn. 
•	Trang so sánh dữ liệu. 
Việc chuyển từ NoIndex sang Index phải được quy định cụ thể trong từng đặc tả SEO của màn hình tương ứng.
 
A.4.4 SEO-D – Share / Discover
SEO-D là nhóm URL phục vụ chia sẻ, lan truyền hoặc khám phá dữ liệu nhưng không phải đối tượng SEO trọng tâm.
Các URL này thường được tạo ra để:
•	Chia sẻ nội dung. 
•	Gửi liên kết. 
•	Hiển thị bản xem trước. 
•	Hỗ trợ Discover. 
•	Hỗ trợ AI Search. 
•	Đặc điểm
•	Có URL riêng. 
•	Có Open Graph. 
•	Có Share Preview. 
•	Không phải URL SEO chính thức. 
•	Yêu cầu
•	Có Metadata tối thiểu. 
•	Có Open Graph. 
•	Có Share Preview. 
•	Có thể NoIndex. 
•	Có thể Canonical về URL gốc. 
•	Ví dụ
•	Share Link. 
•	Snapshot Link. 
•	Report Share Link. 
•	Link chia sẻ kết quả phân tích. 
•	Link chia sẻ bản đồ. 
Mục tiêu chính của SEO-D là hỗ trợ chia sẻ và lan truyền dữ liệu, không phải cạnh tranh thứ hạng tìm kiếm.
 
A.4.5 SEO-N – Không SEO
SEO-N là nhóm đối tượng không thuộc phạm vi SEO của hệ thống.
Các URL và màn hình thuộc nhóm này không được Index và không được xem là tài sản SEO.
•	Đặc điểm
•	Không có giá trị SEO. 
•	Không phục vụ tìm kiếm. 
•	Không phục vụ AI Search. 
•	Không phục vụ Discover. 
•	Yêu cầu
•	NoIndex. 
•	Không đưa vào Sitemap. 
•	Không sinh Structured Data. 
•	Không sinh AI Summary. 
•	Không triển khai Internal Link SEO. 
•	Ví dụ
•	Màn hình đăng nhập. 
•	Màn hình đăng ký. 
•	Hồ sơ người dùng. 
•	Quản trị hệ thống. 
•	Cấu hình hệ thống. 
•	API nội bộ. 
•	Popup kỹ thuật. 
•	Drawer kỹ thuật. 
•	Trang lỗi. 
•	Trạng thái giao diện tạm thời. 
Mọi đối tượng không mang giá trị tìm kiếm hoặc không đủ điều kiện công khai mặc định được xếp vào SEO-N cho đến khi có quy định khác.
PHẦN B.  PHẦN B. MA TRẬN RÀ SOÁT SEO TOÀN HỆ THỐNG
Mục tiêu:
Đảm bảo không bỏ sót bất kỳ màn hình, URL hay đối tượng dữ liệu nào cần SEO.
 
B.1	Ma trận Module / Màn hình SRS
Mục đích
Mục này được sử dụng để rà soát toàn bộ các màn hình trong hệ thống QH Pro nhằm xác định:
•	Màn hình nào thuộc phạm vi SEO. 
•	Màn hình nào không thuộc phạm vi SEO. 
•	Entity nào là đối tượng SEO chính. 
•	SEO Class áp dụng. 
•	Loại URL được sử dụng. 
•	Trạng thái Index/NoIndex dự kiến. 
Đây là căn cứ để bảo đảm không bỏ sót bất kỳ màn hình nào cần được bổ sung đặc tả SEO tại Phần D của tài liệu.
 
Nguyên tắc rà soát
Đối với mỗi màn hình trong SRS cần xác định tối thiểu:
Nội dung	Mô tả
Thuộc phạm vi SEO	Có / Không
Entity chính	Entity SEO trọng tâm
Entity phụ	Entity liên quan
SEO Class	SEO-A / SEO-B / SEO-C / SEO-D / SEO-N
URL SEO	Có / Không
Index	Có / Không
Sitemap	Có / Không
Ghi chú	Điều kiện đặc biệt
 
B.1.1 QH.1 – Tra cứu Quy hoạch / Map View Core
Vai trò SEO
QH.1 là module trung tâm của QH Pro.
Module này chứa các chức năng:
•	Tra cứu quy hoạch. 
•	Điều hướng dữ liệu. 
•	Khám phá dữ liệu. 
•	Truy cập các Entity SEO trọng tâm. 
Không phải mọi màn hình trong QH.1 đều được SEO trực tiếp.
Một số màn hình chỉ đóng vai trò điều hướng hoặc hỗ trợ trải nghiệm người dùng.
 
B.1.1.1 QH.1.1 Trang bản đồ quy hoạch
Thuộc tính	Giá trị
Thuộc phạm vi SEO	Có
Entity chính	Planning Region
Entity phụ	Administrative Unit, GIS Layer
SEO Class	SEO-B
URL SEO	Có
Index	Có
Sitemap	Có
Ghi chú	Landing Page tra cứu quy hoạch
 
B.1.1.2 QH.1.2 Thanh tìm kiếm & định vị trên bản đồ
Thuộc tính	Giá trị
Thuộc phạm vi SEO	Không trực tiếp
Entity chính	Không có
Entity phụ	Administrative Unit, Parcel
SEO Class	SEO-N
URL SEO	Không
Index	Không
Sitemap	Không
Ghi chú	Chức năng điều hướng nội bộ
 
B.1.1.3 QH.1.3 Panel lớp dữ liệu
Thuộc tính	Giá trị
Thuộc phạm vi SEO	Không trực tiếp
Entity chính	GIS Layer
Entity phụ	Planning Region
SEO Class	SEO-N
URL SEO	Không
Index	Không
Sitemap	Không
Ghi chú	Thành phần giao diện hỗ trợ tra cứu
 
B.1.1.4 QH.1.4 Popup thông tin nhanh
Thuộc tính	Giá trị
Thuộc phạm vi SEO	Không trực tiếp
Entity chính	Parcel / Planning Region

Entity phụ	Administrative Unit
SEO Class	SEO-N
URL SEO	Không
Index	Không
Sitemap	Không
Ghi chú	Chỉ là điểm truy cập đến màn hình chi tiết
 
B.1.1.5 QH.1.5 Chi tiết thửa đất / vùng quy hoạch / đồ án
Thuộc tính	Giá trị
Thuộc phạm vi SEO	Có
Entity chính	Parcel / Planning Region / Planning Project
Entity phụ	Administrative Unit, Legal Document
SEO Class	SEO-A
URL SEO	Có
Index	Có
Sitemap	Có
Ghi chú	Nhóm URL SEO trọng tâm của hệ thống
 
Kết quả rà soát QH.1
Màn hình	SEO
QH.1.1 Trang bản đồ quy hoạch	Có
QH.1.2 Thanh tìm kiếm & định vị	Không
QH.1.3 Panel lớp dữ liệu	Không
QH.1.4 Popup thông tin nhanh	Không
QH.1.5 Chi tiết thửa đất / vùng quy hoạch / đồ án	Có
Các màn hình cần đặc tả SEO chi tiết tại Phần D
•	D.1.1 QH.1.1 Trang bản đồ quy hoạch 
•	D.1.5 QH.1.5 Chi tiết thửa đất / vùng quy hoạch / đồ án 
Các màn hình không cần đặc tả SEO riêng
•	QH.1.2 
•	QH.1.3 
•	QH.1.4 
Trừ trường hợp phát sinh URL công khai độc lập trong các giai đoạn phát triển tiếp theo.

 
B.1.2 QH.2 – Tra cứu & Định vị
Vai trò SEO
QH.2 là nhóm chức năng giúp người dùng xác định vị trí, tìm kiếm đối tượng và truy cập nhanh đến các thực thể dữ liệu trong hệ thống.
Khác với QH.1 tập trung vào trải nghiệm bản đồ, QH.2 tập trung vào việc tạo điểm vào (Entry Point) cho các truy vấn tìm kiếm thực tế của người dùng.
Đây là nhóm chức năng có giá trị SEO cao vì phần lớn nhu cầu tìm kiếm ngoài Google và AI Search đều xuất phát từ:
•	Tìm theo địa chỉ. 
•	Tìm theo tọa độ. 
•	Tìm theo tờ/thửa. 
•	Tìm theo khu vực quan tâm. 
Tuy nhiên không phải mọi kết quả tìm kiếm đều được phép SEO. Chỉ những kết quả đủ điều kiện hình thành một thực thể hoặc một URL có giá trị độc lập mới được xem xét Index.
 
B.1.2.1 QH.2.1 Tra cứu địa chỉ
Thuộc tính	Giá trị
Thuộc phạm vi SEO	Có
Entity chính	Administrative Unit
Entity phụ	Parcel, Planning Region
SEO Class	SEO-C
URL SEO	Có
Index	Có điều kiện
Sitemap	Có điều kiện
Ghi chú	Nguồn tạo Landing Page địa danh và khu vực
Nhận xét
Đây là một trong những nguồn SEO quan trọng nhất của hệ thống.
Người dùng thường tìm kiếm:
•	Quy hoạch Hà Nội. 
•	Quy hoạch Bắc Ninh. 
•	Quy hoạch phường Tùng Thiện. 
•	Quy hoạch xã Xuân Khanh. 
•	Quy hoạch khu vực Mỹ Đình. 
Các kết quả tra cứu địa chỉ có thể trở thành Landing Page SEO nếu:
•	Địa danh tồn tại. 
•	Có dữ liệu quy hoạch liên quan. 
•	Có nội dung đủ để hình thành trang độc lập. 
 
B.1.2.2 QH.2.2 Tra cứu GPS
Thuộc tính	Giá trị
Thuộc phạm vi SEO	Không trực tiếp
Entity chính	Parcel
Entity phụ	Planning Region
SEO Class	SEO-N
URL SEO	Không
Index	Không
Sitemap	Không
Ghi chú	Chức năng định vị người dùng
Nhận xét
GPS là đầu vào kỹ thuật.
Người dùng tìm:
21.0345,105.8342
không phải mục tiêu SEO.
SEO chỉ phát sinh khi hệ thống xác định được:
GPS
↓
Thửa đất
↓
Khu quy hoạch
↓
Trang chi tiết Entity
Do đó:
•	Kết quả GPS không SEO. 
•	URL GPS không Index. 
•	Không đưa vào Sitemap. 
SEO được chuyển sang trang Entity tương ứng.
 
B.1.2.3 QH.2.3 Tra cứu tờ/thửa
Thuộc tính	Giá trị
Thuộc phạm vi SEO	Có
Entity chính	Parcel
Entity phụ	Administrative Unit, Planning Region
SEO Class	SEO-A
URL SEO	Có
Index	Có
Sitemap	Có
Ghi chú	Một trong các URL SEO trọng tâm của QH Pro
Nhận xét
Đây là nguồn tạo ra các trang:
•	Chi tiết thửa đất. 
•	Thông tin quy hoạch thửa đất. 
•	Thông tin pháp lý liên quan. 
Người dùng thường tìm:
•	Thửa 123 tờ 45. 
•	Quy hoạch thửa đất số 123. 
•	Đất quy hoạch tại thửa X. 
Nếu dữ liệu đủ điều kiện công khai thì đây là nhóm URL có giá trị SEO rất cao.
 
B.1.2.4 QH.2.4 Tra cứu polygon / vùng phân tích
Thuộc tính	Giá trị
Thuộc phạm vi SEO	Có điều kiện
Entity chính	Planning Region
Entity phụ	Parcel, Snapshot
SEO Class	SEO-C
URL SEO	Có điều kiện
Index	Có điều kiện
Sitemap	Có điều kiện
Ghi chú	Chỉ SEO khi vùng phân tích có giá trị độc lập
Nhận xét
Đa số vùng phân tích do người dùng tự vẽ:
•	Không ổn định. 
•	Không có ý nghĩa SEO. 
•	Không nên Index. 
Tuy nhiên một số vùng có thể hình thành Landing Page SEO như:
•	Khu đô thị xác định. 
•	Khu công nghiệp. 
•	Khu chức năng. 
•	Khu quy hoạch trọng điểm. 
•	Khu vực phân tích được công bố công khai. 
Trong các trường hợp này:
•	Có thể sinh URL SEO. 
•	Có thể sinh Metadata. 
•	Có thể Index. 
Nếu chỉ là Polygon tạm thời do người dùng tạo:
•	SEO-N. 
•	NoIndex. 
•	Không Sitemap. 
 
Kết quả rà soát QH.2
Màn hình	SEO
QH.2.1 Tra cứu địa chỉ	Có
QH.2.2 Tra cứu GPS	Không
QH.2.3 Tra cứu tờ/thửa	Có
QH.2.4 Tra cứu polygon / vùng phân tích	Có điều kiện
•	Các màn hình cần đặc tả SEO chi tiết tại Phần D
•	D.2.1 QH.2.1 Tra cứu địa chỉ 
•	D.2.3 QH.2.3 Tra cứu tờ/thửa 
•	D.2.4 QH.2.4 Tra cứu polygon / vùng phân tích 
•	Các màn hình không cần đặc tả SEO riêng
•	QH.2.2 Tra cứu GPS 
SEO của QH.2.2 được kế thừa từ Entity hoặc màn hình đích mà người dùng được chuyển đến sau khi xác định vị trí.

Dưới đây là bảng tổng hợp SEO cho các module từ QH.3 → QH.10, chuẩn hóa để dễ đưa vào tài liệu kiến trúc SEO:
Module	Tên Module	Thuộc phạm vi SEO	SEO Class	Entity chính	Entity phụ	URL SEO	Index	Sitemap	Ghi chú
QH.3	Quản lý & Khai thác GIS Layer	Có điều kiện	SEO-B / SEO-C	GIS Layer	Planning Region, Planning Project	Có điều kiện	Có điều kiện	Có điều kiện	Chỉ áp dụng cho Layer công khai có ý nghĩa độc lập
QH.4	Diễn giải & Phân tích Quy hoạch	Có	SEO-A	Planning Region, Parcel	Planning Project, Legal Document	Có	Có	Có	Nguồn nội dung SEO và AI Search trọng tâm
QH.5	So sánh & Biến động	Có điều kiện	SEO-C	Planning Region	Parcel, Planning Project	Có điều kiện	Có điều kiện	Có điều kiện	Chỉ SEO kết quả biến động có giá trị độc lập
QH.6	Theo dõi & Cảnh báo	Không	SEO-N	Alert	Parcel, Planning Region	Không	Không	Không	Chức năng cá nhân hóa
QH.7	Snapshot & Report	Có điều kiện	SEO-D / SEO-C	Snapshot, Report	Parcel, Planning Region	Có	Có điều kiện	Có điều kiện	Report công khai có thể SEO, Snapshot cá nhân NoIndex
QH.8	Tài khoản & Gói dịch vụ	Không	SEO-N	User	Subscription	Không	Không	Không	Toàn bộ module mặc định NoIndex
QH.9	API & Tích hợp	Không	SEO-N	API Service	Integration	Không	Không	Không	Không SEO API Runtime
QH.10	Quản trị hệ thống	Không	SEO-N	System Object	Admin Object	Không	Không	Không	Khu vực nội bộ hệ thống

Kết quả rà soát tổng hợp
Module	SEO
QH.3 Layer GIS	Có điều kiện
QH.4 Diễn giải & Phân tích	Có
QH.5 So sánh & Biến động	Có điều kiện
QH.6 Theo dõi & Cảnh báo	Không
QH.7 Snapshot & Report	Có điều kiện
QH.8 Tài khoản & Gói dịch vụ	Không
QH.9 API & Tích hợp	Không
QH.10 Quản trị hệ thống	Không
Các module cần đặc tả SEO chi tiết tiếp theo tại Phần D gồm:
•	D.3 QH.3 – Hệ thống Layer GIS 
•	D.4 QH.4 – Diễn giải & Phân tích Quy hoạch 
•	D.5 QH.5 – So sánh & Biến động 
•	D.7 QH.7 – Snapshot & Report 
Các module còn lại mặc định thuộc nhóm SEO-N và không cần đặc tả SEO chuyên sâu.
 
QH.11 – Thư viện số Quy hoạch (Ma trận SEO chi tiết)
Mã	Thành phần	Entity chính	SEO Class	URL SEO	Index	Sitemap	Vai trò
QH.11.1	Workspace thư viện	Planning Library	SEO-B	Có	Có	Có	Landing Page thư viện quy hoạch
QH.11.2	Danh sách đồ án	Planning Project	SEO-B	Có	Có	Có	Danh mục đồ án quy hoạch
QH.11.3	Chi tiết đồ án	Planning Project	SEO-A	Có	Có	Có	URL SEO trọng tâm của hệ thống
QH.11.4	Danh sách văn bản	Legal Document	SEO-B	Có	Có	Có	Danh mục văn bản pháp lý
QH.11.5	Chi tiết văn bản	Legal Document	SEO-A	Có	Có	Có	URL SEO pháp lý trọng tâm
QH.11.6	Danh sách bản đồ	Planning Map	SEO-B	Có	Có	Có	Danh mục bản đồ quy hoạch
QH.11.7	Chi tiết bản đồ	Planning Map	SEO-A	Có	Có	Có	Landing Page bản đồ quy hoạch
QH.11.8	Timeline pháp lý	Planning Project	SEO-C	Có điều kiện	Có điều kiện	Có điều kiện	Hỗ trợ AI Search, Semantic SEO
QH.11.9	Quan hệ phiên bản	Planning Project	SEO-C	Có điều kiện	Có điều kiện	Có điều kiện	Hỗ trợ Historical SEO & Knowledge Graph

Kết luận
QH.11 cùng với:
•	QH.1.5 (Chi tiết thửa đất / vùng quy hoạch) 
•	QH.2.1 (Tra cứu địa chỉ) 
•	QH.2.3 (Tra cứu tờ/thửa) 
•	QH.4 (Diễn giải & Phân tích) 
là 5 cụm SEO trọng tâm nhất của QH Pro, nơi tập trung phần lớn:
•	Landing Page SEO. 
•	AI Search. 
•	Knowledge Graph. 
•	Internal Linking. 
•	Organic Traffic. 
Đây cũng là các nhóm cần được ưu tiên đặc tả SEO chi tiết và triển khai trước trong giai đoạn phát triển hệ thống.

 B.2            Ma trận Entity SEO
Mục đích
Mục này được sử dụng để rà soát và xác định toàn bộ các thực thể (Entity) có khả năng tham gia vào hệ thống SEO của QH Pro.
Khác với B.1 tập trung vào màn hình và chức năng, B.2 tập trung vào dữ liệu và đối tượng nghiệp vụ.
Đây là cơ sở để:
* Xây dựng URL SEO. 
* Xây dựng Metadata. 
* Xây dựng Structured Data. 
* Xây dựng AI Summary. 
* Xây dựng Knowledge Graph. 
* Xây dựng Internal Link. 
* Xây dựng Sitemap. 
Mỗi Entity cần được đánh giá:
* Có phải đối tượng SEO hay không. 
* Có URL SEO riêng hay không. 
* Có khả năng Index hay không. 
* Có vai trò trong AI Search hay không. 
* Có vai trò trong Knowledge Graph hay không. 

Bảng Entity SEO Master
Mã	Entity	SEO Priority	SEO Class	URL SEO	Index	Sitemap	AI Search	Knowledge Graph	Vai trò chính
B.2.1	Administrative Unit	Rất cao	SEO-A	Có	Có	Có	Có	Có	Landing Page địa bàn, KG Hub
B.2.2	Parcel	Rất cao	SEO-A	Có	Có	Có	Có	Có	Thực thể tra cứu trọng tâm
B.2.3	Planning Region	Rất cao	SEO-A	Có	Có	Có	Có	Có	Landing Page quy hoạch
B.2.4	Planning Project	Rất cao	SEO-A	Có	Có	Có	Có	Có	Đồ án quy hoạch
B.2.5	Planning Library	Cao	SEO-B	Có	Có	Có	Có	Có	Điều hướng thư viện
B.2.6	Legal Document	Rất cao	SEO-A	Có	Có	Có	Có	Có	Thực thể pháp lý
B.2.7	GIS Layer	Trung bình	SEO-B / SEO-C	Có điều kiện	Có điều kiện	Có điều kiện	Có	Có	Lớp dữ liệu bản đồ
B.2.8	Report	Cao	SEO-C	Có	Có điều kiện	Có điều kiện	Có	Có	Báo cáo công khai
B.2.9	Snapshot	Thấp	SEO-D	Có	Có điều kiện	Không mặc định	Có điều kiện	Không trọng tâm	Chia sẻ dữ liệu
B.2.10	AI Summary Content	Cao	Kế thừa Entity nguồn	Không độc lập	Không độc lập	Không độc lập	Rất cao	Rất cao	Semantic Layer


B.3	Ma trận URL SEO
Mục đích
Mục này được sử dụng để rà soát toàn bộ các nhóm URL trong hệ thống QH Pro nhằm xác định:
•	URL nào là đối tượng SEO chính. 
•	URL nào được phép Index. 
•	URL nào chỉ phục vụ điều hướng. 
•	URL nào phục vụ chia sẻ. 
•	URL nào không thuộc phạm vi SEO. 
Đây là căn cứ để:
•	Thiết kế URL. 
•	Thiết kế Canonical. 
•	Thiết kế Sitemap. 
•	Kiểm soát Duplicate Content. 
•	Kiểm soát Index / NoIndex. 
•	Xây dựng Internal Linking. 
Nguyên tắc chung:
Không phải mọi URL đều được SEO.
Chỉ những URL đại diện cho một thực thể hoặc nội dung có giá trị độc lập mới được xem xét Index.
 
Mã	URL Type	SEO Class	Priority	Index	Sitemap	Canonical	Vai trò chính
B3.1	Landing URL	SEO-A	Rất cao	Có	Có	Self	Điểm vào SEO chính
B3.2	Detail URL	SEO-A	Rất cao	Có	Có	Self	URL chuẩn của Entity
B3.3	List URL	SEO-B	Cao	Có	Có	Self	Điều hướng & gom nhóm
B3.4	Search Result URL	SEO-C	Trung bình	Có điều kiện	Có điều kiện	Theo quy định	Kết quả tra cứu
B3.5	Compare URL	SEO-C	Trung bình	Có điều kiện	Có điều kiện	Theo quy định	So sánh dữ liệu
B3.6	Report URL	SEO-C	Cao	Có điều kiện	Có điều kiện	Self	Báo cáo công khai
B3.7	Snapshot URL	SEO-D	Thấp	Có điều kiện	Không mặc định	URL nguồn	Lưu trạng thái dữ liệu
B3.8	Share URL	SEO-D	Thấp	Không mặc định	Không	URL gốc	Chia sẻ nội dung
B3.9	API URL	SEO-N	Không	Không	Không	N/A	Runtime API
B3.10	Internal URL	SEO-N	Không	Không	Không	N/A	Chức năng nội bộ
 
Kết luận
Sau khi rà soát toàn bộ URL trong hệ thống QH Pro:
Loại URL	SEO
Landing URL	SEO-A
Detail URL	SEO-A
List URL	SEO-B
Search Result URL	SEO-C
Compare URL	SEO-C
Report URL	SEO-C
Snapshot URL	SEO-D
Share URL	SEO-D
API URL	SEO-N
Internal URL	SEO-N
Đây là ma trận URL chuẩn dùng làm cơ sở cho toàn bộ các đặc tả URL, Canonical, Sitemap và Index Rule tại các phần tiếp theo của tài liệu.

PHẦN C.  BỘ KHUNG SEO CHUẨN BỔ SUNG VÀO SRS
Vai trò
Phần này quy định các tiêu chuẩn dùng chung khi bổ sung SEO vào các màn hình trong bộ SRS QH Pro.
Phần này không mô tả SEO của từng màn hình cụ thể.
Phần này chỉ quy định:
•	Cấu trúc mô tả SEO chuẩn. 
•	Thành phần SEO bắt buộc. 
•	Mẫu Layout SEO chuẩn. 
•	Quy tắc áp dụng thống nhất trong toàn hệ thống. 
Toàn bộ nội dung SEO chi tiết của từng màn hình sẽ được mô tả tại Phần D.
C.1	Bộ khung SEO chuẩn bổ sung vào SRS
Mục đích
Phần này quy định cấu trúc chuẩn để bổ sung nội dung SEO vào từng màn hình trong bộ SRS QH Pro.
Mọi đặc tả SEO tại Phần D phải tuân thủ thống nhất bộ khung này.
Mục tiêu:
•	Đồng nhất cách mô tả. 
•	Tránh bỏ sót yêu cầu SEO. 
•	Giúp Dev triển khai trực tiếp. 
•	Giúp QA nghiệm thu thống nhất. 
Phần C không mô tả SEO của từng màn hình cụ thể.
Phần C chỉ quy định:
•	Viết cái gì. 
•	Viết đến mức nào. 
•	Không viết cái gì. 
 
C.1.1 SEO Scope
Mục đích
Xác định phạm vi SEO của màn hình.
Phải mô tả
•	Có thuộc phạm vi SEO hay không. 
•	Entity SEO chính. 
•	Entity SEO phụ. 
•	SEO Class. 
•	Loại URL SEO. 
Không mô tả
•	Chi tiết Metadata. 
•	Chi tiết Canonical. 
•	Chi tiết Sitemap. 
 
C.1.2 SEO Class
Mục đích
Xác định mức độ ưu tiên SEO.
Phải mô tả
•	SEO-A/B/C/D/N. 
•	Lý do phân loại. 
Không mô tả
•	Metadata chi tiết. 
•	Structured Data chi tiết. 
 
C.1.3 URL & Canonical
Mục đích
Xác định URL SEO chuẩn.
Phải mô tả
•	URL được Index. 
•	URL NoIndex. 
•	Canonical URL. 
•	Duplicate Control. 
Không mô tả
•	Metadata. 
•	Schema. 
 
C.1.4 Metadata
Mục đích
Xác định Metadata phải sinh.
Phải mô tả
•	Title. 
•	Description. 
•	H1. 
•	Robots. 
Không mô tả
•	Cách code Meta Tag. 
 
C.1.5 Structured Data
Mục đích
Xác định Schema áp dụng.
Phải mô tả
•	Schema Type. 
•	Entity Mapping. 
•	Điều kiện áp dụng. 
Không mô tả
•	JSON-LD chi tiết. 
 
C.1.6 Open Graph
Mục đích
Quy định dữ liệu phục vụ chia sẻ.
Phải mô tả
•	OG Title. 
•	OG Description. 
•	OG Image. 
 
C.1.7 AI Summary
Mục đích
Quy định nội dung phục vụ AI Search.
Phải mô tả
•	Có AI Summary hay không. 
•	Dữ liệu nguồn. 
•	Giới hạn dữ liệu công khai. 
 
C.1.8 Sitemap
Mục đích
Quy định việc tham gia Sitemap.
Phải mô tả
•	Sitemap Group. 
•	Điều kiện Index. 
•	Điều kiện loại bỏ. 
 
C.1.9 Internal Link
Mục đích
Quy định liên kết nội bộ.
Phải mô tả
•	Link đến Entity cha. 
•	Link đến Entity con. 
•	Link ngang. 

C.1.10 Security & Visibility
Mục đích
Quy định giới hạn công khai.
Phải mô tả
•	Public. 
•	Restricted. 
•	Private. 
•	Dữ liệu không được SEO. 
 
C.1.11 Cache & Regeneration
Mục đích
Quy định các sự kiện phải cập nhật lại dữ liệu SEO.
Phải mô tả
•	Trigger Event. 
•	SEO Component ảnh hưởng. 
Ví dụ:
Event	Regenerate
parcel.updated	Metadata
legal.updated	AI Summary
visibility.changed	Sitemap
 
C.1.12 Acceptance Criteria
Mục đích
Quy định tiêu chí nghiệm thu SEO.
Phải mô tả
•	URL. 
•	Metadata. 
•	Schema. 
•	Sitemap. 
•	Security. 
Tất cả phải đo kiểm được.
 
C.1.13 QA Checklist
Mục đích
Danh sách kiểm tra cuối cùng.
Phải mô tả
Checklist dạng:
□ URL đúng
□ Canonical đúng
□ Metadata tồn tại
□ Schema hợp lệ
□ Sitemap đúng
□ NoIndex đúng

C.2	Hệ thống trang đích seo (seo landing page system)
C.2.1 Vai trò và phạm vi áp dụng
Mục này quy định hệ thống trang đích SEO chuẩn áp dụng cho toàn bộ QH Pro.
Các trang đích SEO là các trang Web/Public Page được tạo ra nhằm phục vụ:
•	Google Search. 
•	AI Search. 
•	Knowledge Graph. 
•	Internal Linking. 
•	Share Preview. 
•	Traffic Acquisition. 
•	Chuyển đổi người dùng sang môi trường sử dụng chính của QH Pro. 
Các trang này không phải là các màn hình nghiệp vụ chính của hệ thống GIS, nhưng được xây dựng dựa trên dữ liệu, Entity và luồng nghiệp vụ đã được mô tả trong bộ SRS QH Pro.
Mục tiêu của phần này là giúp:
•	Design có cơ sở dựng Mockup thống nhất. 
•	Frontend có cơ sở xây dựng Component và Layout. 
•	Backend có cơ sở cung cấp dữ liệu SEO. 
•	QA có cơ sở nghiệm thu. 
•	Hệ thống SEO được triển khai đồng nhất trên toàn bộ QH Pro. 
Phạm vi áp dụng bao gồm:
•	Trang địa bàn hành chính. 
•	Trang thửa đất. 
•	Trang khu quy hoạch. 
•	Trang đồ án quy hoạch. 
•	Trang văn bản pháp lý. 
•	Trang bản đồ quy hoạch. 
•	Trang báo cáo và phân tích. 
•	Trang Snapshot công khai. 
•	Các Landing Page SEO được sinh tự động từ Entity SEO. 
Không áp dụng cho:
•	Màn hình quản trị. 
•	Màn hình tài khoản người dùng. 
•	Màn hình thanh toán. 
•	Popup runtime. 
•	Drawer runtime. 
•	Layer Toggle. 
•	Filter State. 
•	Search State tạm thời. 
•	API Runtime. 
•	Dữ liệu Private. 
•	Dữ liệu Restricted. 
 
C.2.1.1 Mục tiêu của SEO Landing Page
SEO Landing Page là lớp hiển thị công khai giúp chuyển đổi dữ liệu quy hoạch thành nội dung có thể được:
•	Google thu thập. 
•	AI đọc hiểu. 
•	Người dùng tìm kiếm. 
•	Người dùng chia sẻ. 
•	Hệ thống liên kết nội bộ. 
Mỗi SEO Landing Page phải đạt đồng thời các mục tiêu sau:
•	Mục tiêu SEO
•	Có URL ổn định. 
•	Có khả năng index độc lập. 
•	Có khả năng xuất hiện trên Google Search. 
•	Có khả năng xuất hiện trên AI Search. 
•	Có khả năng được trích dẫn bởi các hệ thống AI. 
•	Mục tiêu dữ liệu
•	Đại diện cho một Entity cụ thể. 
•	Thể hiện được ngữ cảnh dữ liệu. 
•	Liên kết được với các Entity liên quan. 
•	Không tạo nội dung mồ côi (orphan content). 
•	Mục tiêu trải nghiệm
•	Người dùng hiểu nhanh thông tin chính. 
•	Có thể tiếp tục khám phá dữ liệu liên quan. 
•	Có thể chuyển sang môi trường bản đồ hoặc ứng dụng QH Pro. 
•	Mục tiêu chuyển đổi
•	Mở trên bản đồ. 
•	Mở trong ứng dụng. 
•	Đăng ký tài khoản. 
•	Theo dõi dữ liệu. 
•	Chia sẻ dữ liệu. 
•	Tạo báo cáo. 
•	Sử dụng các chức năng nâng cao của QH Pro. 
Nguyên tắc chung:
Mọi SEO Landing Page phải vừa phục vụ SEO, vừa phục vụ điều hướng người dùng vào hệ sinh thái QH Pro.
 
C.2.1.2 Quan hệ với các màn hình nghiệp vụ QH Pro
SEO Landing Page không thay thế các màn hình nghiệp vụ đã được mô tả trong SRS.
Các màn hình nghiệp vụ của QH Pro vẫn là nơi người dùng thao tác chính:
•	Tra cứu bản đồ. 
•	Bật/tắt Layer. 
•	Phân tích quy hoạch. 
•	So sánh dữ liệu. 
•	Theo dõi biến động. 
•	Tạo báo cáo. 
•	Quản lý tài khoản. 
SEO Landing Page chỉ là lớp công khai được sinh ra từ dữ liệu của các màn hình nghiệp vụ đó.
Mối quan hệ:
Dữ liệu GIS
      ↓
Entity SEO
      ↓
SEO Engine
      ↓
SEO Landing Page
      ↓
Google / AI Search
      ↓
Người dùng
      ↓
Mở trên QH Pro
      ↓
Màn hình nghiệp vụ
Nguyên tắc:
•	Không sao chép toàn bộ chức năng nghiệp vụ lên SEO Page. 
•	Không đưa Workspace GIS lên SEO Page. 
•	Không đưa Layer Control lên SEO Page. 
•	Không đưa các trạng thái thao tác tạm thời lên SEO Page. 
•	SEO Page chỉ hiển thị dữ liệu cần thiết để tìm kiếm và khám phá thông tin. 
SEO Landing Page phải đóng vai trò:
Cổng tiếp cận dữ liệu.
Trong khi:
Màn hình nghiệp vụ là môi trường làm việc chính.
 
C.2.1.3 Quan hệ với Phần D – Đặc tả SEO theo màn hình
Phần C.2 quy định:
•	Layout chuẩn. 
•	Component chuẩn. 
•	Cấu trúc trang chuẩn. 
•	Quy tắc hiển thị chuẩn. 
Phần D quy định:
•	SEO áp dụng cho từng màn hình cụ thể. 
•	SEO áp dụng cho từng URL cụ thể. 
•	Metadata cụ thể. 
•	Canonical cụ thể. 
•	Schema cụ thể. 
•	Sitemap cụ thể. 
•	Internal Link cụ thể. 
•	AI Summary cụ thể. 
Mối quan hệ:
Phần C.2
     ↓
SEO Design System

Phần D
     ↓
SEO Implementation Specification
Nguyên tắc:
•	Phần C.2 định nghĩa mẫu. 
•	Phần D áp dụng mẫu vào từng màn hình. 
Không được mô tả trùng lặp giữa hai phần.
Không được định nghĩa lại Layout trong Phần D.
Phần D chỉ được tham chiếu tới Layout đã được định nghĩa tại C.2.
 
C.2.1.4 Nguyên tắc kế thừa Design System QH Pro
Toàn bộ SEO Landing Page phải kế thừa hệ thống nhận diện giao diện của QH Pro.
Bao gồm:
•	Thành phần kế thừa bắt buộc
•	Header. 
•	Footer. 
•	Typography. 
•	Màu sắc thương hiệu. 
•	Icon. 
•	Button. 
•	Card. 
•	Badge. 
•	Table. 
•	Accordion. 
•	Empty State. 
•	CTA Style. 
•	Thành phần điều hướng
•	Logo QH Pro. 
•	Navigation. 
•	Search Entry. 
•	CTA tải App. 
•	CTA đăng nhập. 
•	Thành phần pháp lý
•	Disclaimer. 
•	Nguồn dữ liệu. 
•	Thời gian cập nhật. 
•	Trạng thái dữ liệu. 
Không được:
•	Xây dựng giao diện SEO như website riêng biệt. 
•	Sử dụng bộ màu khác. 
•	Sử dụng hệ thống component khác. 
•	Tạo trải nghiệm tách rời với QH Pro. 
Nguyên tắc:
Người dùng phải nhận biết ngay đây là một phần của hệ sinh thái QH Pro.
 
C.2.1.5 Nguyên tắc HTML-first / SEO-first
Toàn bộ SEO Landing Page phải được xây dựng theo nguyên tắc:
•	HTML-first
Các nội dung quan trọng phải xuất hiện ngay trong HTML render đầu tiên:
•	Title. 
•	Description. 
•	Breadcrumb. 
•	H1. 
•	AI Summary. 
•	Key Facts. 
•	Nội dung chính. 
•	FAQ. 
•	Internal Link. 
•	Schema. 
Không được phụ thuộc hoàn toàn vào JavaScript để sinh nội dung SEO.
•	SEO-first
Thiết kế giao diện phải ưu tiên:
•	Khả năng index. 
•	Khả năng crawl. 
•	Khả năng đọc hiểu của AI. 
•	Khả năng hiển thị Rich Result. 
trước khi tối ưu các hiệu ứng giao diện.
•	Entity-first
Mỗi trang phải đại diện cho một Entity rõ ràng:
•	Administrative Unit. 
•	Parcel. 
•	Planning Region. 
•	Planning Project. 
•	Legal Document. 
•	Planning Map. 
•	Report. 
Không được tạo Landing Page không có Entity trung tâm.
•	Semantic-first
Mỗi trang phải thể hiện rõ:
•	Entity chính. 
•	Entity liên quan. 
•	Quan hệ dữ liệu. 
•	Ngữ cảnh pháp lý. 
•	Ngữ cảnh địa lý. 
•	Performance-first
Trang SEO phải:
•	Render nhanh. 
•	Có HTML đầy đủ. 
•	Hạn chế phụ thuộc JS. 
•	Có khả năng cache. 
•	Có khả năng regenerate tự động. 
Nguyên tắc cuối cùng:
Mọi SEO Landing Page của QH Pro phải được thiết kế để Google, AI và người dùng đều có thể hiểu được nội dung ngay từ lần truy cập đầu tiên.

 
C.2.2 Administrative Unit SEO Layout
C.2.2.1 Vai trò & phạm vi áp dụng
Administrative Unit SEO Layout là mẫu trang đích SEO chuẩn áp dụng cho các Entity địa bàn hành chính trong QH Pro.
Mục tiêu:
•	SEO theo địa bàn. 
•	Tạo Landing Page cho Google Search. 
•	Tạo Landing Page cho AI Search. 
•	Hình thành Knowledge Graph địa lý. 
•	Điều hướng người dùng vào hệ sinh thái QH Pro. 
Áp dụng cho:
•	Tỉnh/Thành phố. 
•	Quận/Huyện. 
•	Xã/Phường. 
•	Đơn vị hành chính lịch sử. 
•	Đơn vị hành chính trước sáp nhập. 
•	Đơn vị hành chính sau sáp nhập. 
 
C.2.2.2 Entity áp dụng
•	Entity chính
Administrative Unit
Bao gồm:
Province
District
Ward
Commune
Historical Administrative Unit
 
•	Entity phụ
Planning Region
Planning Project
Legal Document
Planning Map
Report
 
C.2.2.3 Thành phần tĩnh
Các thành phần cố định trên mọi trang.
•	Header QH Pro
Logo

Menu

Search

CTA tải App

Đăng nhập
 
•	Footer QH Pro
Liên hệ

Chính sách

Điều khoản

API

Tải App

•	Disclaimer
Thông tin mang tính tham khảo.

Người dùng cần đối chiếu với nguồn dữ liệu chính thức khi sử dụng cho mục đích pháp lý hoặc giao dịch.
 
C.2.2.4 Thành phần động
Các thành phần sinh từ SEO Engine.
•	Breadcrumb
Ví dụ
Trang chủ

→ Hà Nội

→ Sơn Tây

→ Tùng Thiện
 
•	H1
Ví dụ
Quy hoạch phường Tùng Thiện, Hà Nội
 
•	AI Summary
Ví dụ
Tóm tắt địa bàn

Tóm tắt quy hoạch

Tóm tắt dữ liệu
 
•	Key Facts
Ví dụ
Tên địa bàn

Loại đơn vị

Mã đơn vị

Tỉnh

Dân số

Diện tích

Ngày cập nhật
 
•	Related Entity
Ví dụ
Đồ án liên quan

Bản đồ liên quan

Văn bản liên quan

Địa bàn lân cận

C.2.2.5 Desktop Layout
+----------------------------------------------------------------------------------+
| HEADER QH PRO                                                                    |
| Logo | Tra cứu | Thư viện | Báo cáo | API | Tải App | Đăng nhập                  |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| BREADCRUMB                                                                       |
| Trang chủ > Hà Nội > Sơn Tây > Tùng Thiện                                        |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| HERO SEO                                                                         |
| +----------------------------------------------------------+------------------+ |
| | H1                                                       | Map Preview      | |
| | Quy hoạch phường Tùng Thiện                              |                  | |
| |                                                          | Bản đồ địa bàn   | |
| | AI Summary                                               |                  | |
| |                                                          | CTA Mở bản đồ    | |
| +----------------------------------------------------------+------------------+ |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| KEY FACTS                                                                        |
| [Loại] [Diện tích] [Dân số] [Tỉnh] [Ngày cập nhật] [Nguồn dữ liệu]              |
+----------------------------------------------------------------------------------+

+---------------------------------------------+------------------------------------+
| MAIN CONTENT                                | RELATED ENTITY                     |
|                                             |                                    |
| Tổng quan địa bàn                           | Địa bàn lân cận                    |
|                                              |                                    |
| Thông tin hành chính                        | Đồ án liên quan                    |
|                                              |                                    |
| Quy hoạch liên quan                         | Văn bản liên quan                  |
|                                              |                                    |
| Bản đồ liên quan                            | Bản đồ liên quan                   |
|                                              |                                    |
| Đơn vị trước/sau sáp nhập                   |                                    |
|                                              |                                    |
| FAQ                                          |                                    |
+---------------------------------------------+------------------------------------+

+----------------------------------------------------------------------------------+
| CTA                                                                              |
| Mở trên bản đồ | Theo dõi | Chia sẻ | Tải App                                   |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| DISCLAIMER                                                                       |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| FOOTER                                                                           |
+----------------------------------------------------------------------------------+
 
C.2.2.6 Responsive / Tablet Layout
+--------------------------------------------------------------+
| HEADER                                                       |
+--------------------------------------------------------------+

| Breadcrumb                                                   |

+--------------------------------------------------------------+
| H1                                                           |
| AI Summary                                                   |
+--------------------------------------------------------------+

| Map Preview                                                  |

+--------------------------------------------------------------+
| Key Facts                                                    |
+--------------------------------------------------------------+

| Tổng quan                                                    |
| Quy hoạch liên quan                                          |
| Đơn vị trước/sau sáp nhập                                    |
| FAQ                                                          |
+--------------------------------------------------------------+

| Related Entity                                               |
+--------------------------------------------------------------+

| CTA                                                          |
+--------------------------------------------------------------+

| Footer                                                       |
+--------------------------------------------------------------+
 
C.2.2.7 Mobile Web / WebView Layout
+--------------------------------------+
| Mobile Header                        |
+--------------------------------------+

| Breadcrumb                           |

+--------------------------------------+
| H1                                   |
+--------------------------------------+

| AI Summary                           |
+--------------------------------------+

| CTA Sticky                           |
| Mở bản đồ                            |
+--------------------------------------+

| Map Preview                          |
+--------------------------------------+

| Key Facts                            |
+--------------------------------------+

| Accordion                            |
| Tổng quan                            |
| Quy hoạch                            |
| Địa bàn liên quan                    |
| FAQ                                  |
+--------------------------------------+

| Related Content                      |
+--------------------------------------+

| Footer                               |
+--------------------------------------+
 
C.2.2.8 Mobile App Layout
+--------------------------------------+
| App Header                           |
+--------------------------------------+

| H1                                   |
| Badge                                |
+--------------------------------------+

| AI Summary Card                      |
+--------------------------------------+

| Mini Map                             |
+--------------------------------------+

| Tab                                  |
| Tổng quan                            |
| Quy hoạch                            |
| Văn bản                              |
| Bản đồ                               |
+--------------------------------------+

| Content                              |
+--------------------------------------+

| Related Entity                       |
+--------------------------------------+

| Bottom CTA                           |
| Mở bản đồ                            |
+--------------------------------------+
 
C.2.2.9 SEO Block bắt buộc
Bắt buộc phải có:
•	Breadcrumb. 
•	H1. 
•	AI Summary. 
•	Key Facts. 
•	Tổng quan địa bàn. 
•	Map Preview. 
•	Related Entity. 
•	CTA. 
•	Disclaimer. 
•	Footer. 
 
C.2.2.10 SEO Block tùy chọn
Có thể bật/tắt:
•	FAQ. 
•	Timeline hành chính. 
•	Dữ liệu dân số. 
•	Dữ liệu kinh tế. 
•	Snapshot liên quan. 
•	Report liên quan. 
 
C.2.2.11 Data Mapping
Block	Nguồn dữ liệu
H1	Administrative Unit
AI Summary	AI Summary Engine
Key Facts	Admin Unit Metadata
Map Preview	GIS Engine
Related Entity	Knowledge Graph
FAQ	SEO Engine

C.2.2.12 Render Rule
Thành phần	Render
Metadata	SSR
H1	SSR
AI Summary	SSR
Key Facts	SSR
FAQ	SSR
Related Entity	SSR/ISR
Map Preview	SSR image
Full Map	CSR
 
C.2.2.13 CTA & Internal Link
CTA chính:
Mở trên bản đồ
CTA phụ:
Xem đồ án liên quan

Xem văn bản liên quan

Xem bản đồ liên quan

Tải App
Internal Link bắt buộc:
Tỉnh
→ Huyện

Huyện
→ Xã

Xã
→ Đồ án

Đồ án
→ Văn bản

Văn bản
→ Bản đồ
 
C.2.2.14 Mapping áp dụng trong QH Pro
Mục SEO	Áp dụng
D.2.1 Tra cứu địa chỉ	✓
D.11.1 Workspace thư viện	✓
D.11.2 Danh sách đồ án	✓
D.11.3 Chi tiết đồ án (địa bàn)	✓
D.11.5 Chi tiết văn bản (địa bàn)	✓
Ghi chú:
Administrative Unit SEO Layout là Layout nền tảng
cho toàn bộ SEO theo địa bàn hành chính của QH Pro.

Mọi Landing Page cấp tỉnh, huyện, xã phải sử dụng
layout này.

 
C.2.3 Parcel SEO Layout
Áp dụng:
•	Thửa đất. 
•	Kết quả tra cứu tờ/thửa. 
•	Kết quả tra cứu GPS đã resolve sang Parcel. 
•	Kết quả tra cứu địa chỉ đã xác định được thửa đất. 
•	Trang Snapshot công khai của thửa đất. 
•	Trang Report công khai của thửa đất. 
 
C.2.3.1 Vai trò & phạm vi áp dụng
Parcel SEO Layout là mẫu trang đích SEO chuẩn dành cho đối tượng Thửa đất (Parcel) trong QH Pro.
Đây là một trong những nhóm Landing Page quan trọng nhất của hệ thống vì phần lớn nhu cầu người dùng cuối đều xoay quanh việc tìm hiểu:
•	Thửa đất nằm ở đâu. 
•	Thửa đất thuộc quy hoạch gì. 
•	Thửa đất có nằm trong khu vực quy hoạch hay không. 
•	Thửa đất chịu tác động của đồ án nào. 
•	Thửa đất liên quan tới các văn bản nào. 
Mục tiêu của Parcel SEO Layout:
•	Tạo Landing Page SEO cho từng thửa đất. 
•	Tạo Landing Page SEO cho kết quả tra cứu tờ/thửa. 
•	Tạo Landing Page SEO cho kết quả GPS. 
•	Tạo Landing Page SEO cho kết quả địa chỉ. 
•	Tạo Knowledge Graph xoay quanh Parcel. 
•	Điều hướng người dùng sang môi trường bản đồ QH Pro. 
Parcel SEO Page không thay thế chức năng tra cứu GIS chuyên sâu.
Parcel SEO Page chỉ cung cấp:
•	Thông tin công khai. 
•	Thông tin tóm tắt. 
•	Thông tin định hướng. 
•	Liên kết tới dữ liệu liên quan. 
 
C.2.3.2 Entity áp dụng
•	Entity chính
Parcel
Bao gồm:
Thửa đất
Tờ bản đồ
Kết quả định vị GPS
Kết quả định vị địa chỉ
 
•	Entity phụ
Administrative Unit

Planning Region

Planning Project

Legal Document

Planning Map

Report

Snapshot
 
C.2.3.3 Thành phần tĩnh
•	Header QH Pro
Logo

Tra cứu

Thư viện quy hoạch

Báo cáo

API

Tải App

Đăng nhập
 
•	Footer QH Pro
Giới thiệu

Điều khoản

Chính sách

API

Liên hệ

Tải App
 
•	Disclaimer
Thông tin quy hoạch chỉ mang tính tham khảo.

Người sử dụng cần đối chiếu với cơ quan nhà nước có thẩm quyền trước khi sử dụng cho mục đích pháp lý hoặc giao dịch.
 
C.2.3.4 Thành phần động
•	Breadcrumb
Ví dụ
Trang chủ

→ Hà Nội

→ Sơn Tây

→ Tùng Thiện

→ Tờ 12

→ Thửa 456
 
•	H1
Ví dụ
Thông tin quy hoạch thửa đất số 456,
tờ bản đồ số 12,
phường Tùng Thiện, Hà Nội

•	AI Summary
Ví dụ
Thửa đất thuộc địa bàn...

Hiện đang chịu tác động của...

Các đồ án quy hoạch liên quan...

Các lưu ý chính...
 
•	Key Facts
Ví dụ
Tờ bản đồ

Số thửa

Địa bàn

Diện tích

Loại đất

Tình trạng dữ liệu

Ngày cập nhật
 
•	Planning Summary
Ví dụ
Loại quy hoạch

Khu quy hoạch liên quan

Mức độ ảnh hưởng

Trạng thái
 
•	Related Entity
Ví dụ
Khu quy hoạch liên quan

Đồ án liên quan

Văn bản liên quan

Bản đồ liên quan

Các thửa lân cận
 
C.2.3.5 Desktop Layout
+----------------------------------------------------------------------------------+
| HEADER QH PRO                                                                    |
| Logo | Tra cứu | Thư viện | Báo cáo | API | Tải App | Đăng nhập                  |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| BREADCRUMB                                                                       |
| Trang chủ > Hà Nội > Sơn Tây > Tùng Thiện > Tờ 12 > Thửa 456                    |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| HERO SEO                                                                         |
| +----------------------------------------------------------+------------------+ |
| | H1                                                       | Parcel Preview   | |
| | Thửa đất 456 - Tờ 12                                     |                  | |
| |                                                          | Mini Map         | |
| | AI Summary                                               |                  | |
| |                                                          | CTA Mở bản đồ    | |
| +----------------------------------------------------------+------------------+ |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| KEY FACTS                                                                        |
| [Tờ] [Thửa] [Địa bàn] [Diện tích] [Loại đất] [Cập nhật]                         |
+----------------------------------------------------------------------------------+

+---------------------------------------------+------------------------------------+
| MAIN CONTENT                                | RELATED ENTITY                     |
|                                             |                                    |
| Tổng quan thửa đất                          | Khu quy hoạch liên quan            |
|                                             |                                    |
| Thông tin quy hoạch                         | Đồ án liên quan                    |
|                                             |                                    |
| Thông tin địa bàn                           | Văn bản liên quan                  |
|                                             |                                    |
| Bản đồ liên quan                            | Thửa lân cận                       |
|                                             |                                    |
| Lưu ý và cảnh báo                           | Snapshot liên quan                 |
|                                             |                                    |
| FAQ                                         |                                    |
+---------------------------------------------+------------------------------------+

+----------------------------------------------------------------------------------+
| CTA                                                                              |
| Mở trên bản đồ | Theo dõi | Chia sẻ | Tạo báo cáo                               |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| DISCLAIMER                                                                       |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| FOOTER                                                                           |
+----------------------------------------------------------------------------------+
 
C.2.3.6 Responsive / Tablet Layout
+--------------------------------------------------------------+
| HEADER                                                       |
+--------------------------------------------------------------+

| Breadcrumb                                                   |

+--------------------------------------------------------------+
| H1                                                           |
| AI Summary                                                   |
+--------------------------------------------------------------+

| Parcel Preview / Mini Map                                   |
+--------------------------------------------------------------+

| Key Facts                                                    |
+--------------------------------------------------------------+

| Tổng quan thửa đất                                           |
| Thông tin quy hoạch                                          |
| Thông tin địa bàn                                            |
| FAQ                                                          |
+--------------------------------------------------------------+

| Related Entity                                               |
+--------------------------------------------------------------+

| CTA                                                          |
+--------------------------------------------------------------+

| Footer                                                       |
+--------------------------------------------------------------+
 
C.2.3.7 Mobile Web / WebView Layout
+--------------------------------------+
| Mobile Header                        |
+--------------------------------------+

| Breadcrumb                           |
+--------------------------------------+

| H1                                   |
+--------------------------------------+

| AI Summary                           |
+--------------------------------------+

| CTA Sticky                           |
| Mở bản đồ                            |
+--------------------------------------+

| Mini Map                             |
+--------------------------------------+

| Key Facts                            |
+--------------------------------------+

| Accordion                            |
| Tổng quan                            |
| Quy hoạch                            |
| Địa bàn                              |
| Văn bản liên quan                    |
| FAQ                                  |
+--------------------------------------+

| Related Content                      |
+--------------------------------------+

| Footer                               |
+--------------------------------------+
 
C.2.3.8 Mobile App Layout
+--------------------------------------+
| App Header                           |
+--------------------------------------+

| H1                                   |
| Badge                                |
+--------------------------------------+

| AI Summary Card                      |
+--------------------------------------+

| Mini Map                             |
+--------------------------------------+

| Tab                                  |
| Tổng quan                            |
| Quy hoạch                            |
| Văn bản                              |
| Bản đồ                               |
+--------------------------------------+

| Content                              |
+--------------------------------------+

| Related Entity                       |
+--------------------------------------+

| Bottom CTA                           |
| Mở bản đồ                            |
+--------------------------------------+
 
C.2.3.9 SEO Block bắt buộc
Bắt buộc phải có:
•	Breadcrumb. 
•	H1. 
•	AI Summary. 
•	Key Facts. 
•	Parcel Preview / Mini Map. 
•	Tổng quan thửa đất. 
•	Thông tin quy hoạch. 
•	Related Entity. 
•	CTA. 
•	Disclaimer. 
•	Footer. 
 
C.2.3.10 SEO Block tùy chọn
Có thể bật/tắt:
•	FAQ. 
•	Snapshot liên quan. 
•	Report liên quan. 
•	Thửa lân cận. 
•	Timeline biến động. 
•	Lịch sử cập nhật dữ liệu. 
•	AI Insight nâng cao. 

C.2.3.11 Data Mapping
Block	Nguồn dữ liệu
H1	Parcel
AI Summary	AI Summary Engine
Key Facts	Parcel Metadata
Parcel Preview	GIS Engine
Planning Summary	Planning Engine
Related Entity	Knowledge Graph
FAQ	SEO Engine
 
C.2.3.12 Render Rule
Thành phần	Render
Metadata	SSR
H1	SSR
AI Summary	SSR
Key Facts	SSR
Planning Summary	SSR
FAQ	SSR
Related Entity	SSR / ISR
Mini Map	SSR Image
Full Map	CSR
 
C.2.3.13 CTA & Internal Link
•	CTA chính
Mở trên bản đồ
 
•	CTA phụ
Xem khu quy hoạch liên quan

Xem đồ án liên quan

Xem văn bản liên quan

Theo dõi biến động

Tạo báo cáo

Tải App
 
•	Internal Link bắt buộc
Thửa đất
→ Địa bàn

Thửa đất
→ Khu quy hoạch

Khu quy hoạch
→ Đồ án

Đồ án
→ Văn bản

Văn bản
→ Bản đồ

Bản đồ
→ Report
 
C.2.3.14 Mapping áp dụng trong QH Pro
Mục SEO	Áp dụng
D.1.5 Chi tiết thửa đất	✓
D.2.2 Tra cứu GPS	✓
D.2.3 Tra cứu tờ/thửa	✓
D.2.1 Tra cứu địa chỉ (khi resolve về thửa)	✓
D.7.1 Snapshot công khai	✓
D.7.3 Report công khai	✓
•	Ghi chú
Parcel SEO Layout là Layout chuẩn cho toàn bộ
SEO liên quan đến thửa đất trong QH Pro.

Mọi Landing Page lấy Parcel làm Entity trung tâm
đều phải sử dụng Layout này.
 
Đây là một trong những layout SEO quan trọng nhất của QH Pro vì nó trực tiếp phục vụ nhu cầu tra cứu quy hoạch theo thửa đất — nhóm truy vấn có giá trị SEO và chuyển đổi người dùng cao nhất.

 
C.2.4 Planning Region SEO Layout
Áp dụng:
•	Khu quy hoạch. 
•	Vùng quy hoạch. 
•	Khu chức năng. 
•	Khu đô thị. 
•	Khu công nghiệp. 
•	Vùng phân tích công khai đã được chuẩn hóa thành Planning Region. 
•	Kết quả tra cứu polygon nếu vùng có giá trị SEO độc lập. 
 
C.2.4.1 Vai trò & phạm vi áp dụng
Planning Region SEO Layout là mẫu trang đích SEO chuẩn dành cho các khu vực/vùng quy hoạch trong QH Pro.
Layout này phục vụ các truy vấn người dùng thường gặp như:
•	Quy hoạch khu vực X. 
•	Quy hoạch khu đô thị X. 
•	Quy hoạch khu công nghiệp X. 
•	Khu vực X thuộc quy hoạch gì. 
•	Bản đồ quy hoạch vùng X. 
•	Đồ án quy hoạch khu X. 
Mục tiêu của Planning Region SEO Layout:
•	Tạo Landing Page SEO cho từng vùng/khu quy hoạch. 
•	Giúp Google và AI hiểu được phạm vi, loại quy hoạch, cơ sở pháp lý và các đối tượng liên quan của vùng. 
•	Liên kết vùng quy hoạch với địa bàn hành chính, đồ án, văn bản, bản đồ, layer và báo cáo. 
•	Dẫn người dùng sang Map View để xem trực tiếp trên bản đồ. 
•	Dẫn người dùng sang Thư viện Quy hoạch để xem hồ sơ, đồ án và văn bản nguồn. 
Planning Region SEO Page không thay thế màn hình phân tích vùng hoặc Map Workspace. Trang này chỉ hiển thị nội dung công khai, có cấu trúc, phục vụ tìm kiếm và điều hướng.
 
C.2.4.2 Entity áp dụng
•	Entity chính
Planning Region
Bao gồm:
Khu quy hoạch
Vùng quy hoạch
Khu chức năng
Khu đô thị
Khu công nghiệp
Vùng phân tích công khai
•	Entity phụ
Administrative Unit
Planning Project
Legal Document
Planning Map
GIS Layer
Parcel
Report
Snapshot
 
C.2.4.3 Thành phần tĩnh
•	Header QH Pro
Logo

Tra cứu

Thư viện quy hoạch

Báo cáo

API

Tải App

Đăng nhập
•	Footer QH Pro
Giới thiệu

Điều khoản

Chính sách

API

Liên hệ

Tải App
•	Disclaimer
Thông tin quy hoạch được tổng hợp từ dữ liệu công khai và dữ liệu được chuẩn hóa trong QH Pro.

Người sử dụng cần đối chiếu với hồ sơ, bản đồ và văn bản pháp lý chính thức trước khi sử dụng cho mục đích pháp lý, đầu tư hoặc giao dịch.
•	CTA mặc định
Mở vùng này trên bản đồ

Xem đồ án liên quan

Xem văn bản pháp lý

Chia sẻ
 
C.2.4.4 Thành phần động
•	Breadcrumb
Ví dụ:
Trang chủ

→ Hà Nội

→ Quy hoạch

→ Khu quy hoạch ven sông Hồng
hoặc:
Trang chủ

→ Bình Dương

→ Khu công nghiệp

→ KCN VSIP
•	H1
Ví dụ:
Quy hoạch Khu đô thị Tây Hồ Tây, Hà Nội
hoặc:
Thông tin quy hoạch Khu công nghiệp VSIP, Bình Dương
•	AI Summary
Nội dung cần tóm tắt:
Vị trí và phạm vi khu vực

Loại quy hoạch chính

Địa bàn hành chính liên quan

Đồ án quy hoạch nguồn

Văn bản pháp lý liên quan

Các lớp bản đồ / layer liên quan

Lưu ý dữ liệu và trạng thái hiệu lực
•	Key Facts
Ví dụ:
Tên khu/vùng

Loại vùng

Địa bàn

Diện tích

Trạng thái quy hoạch

Đồ án nguồn

Số văn bản liên quan

Ngày cập nhật
•	Planning Summary
Ví dụ:
Loại quy hoạch

Chức năng sử dụng chính

Lớp quy hoạch chính

Lớp phụ / overlay

Trạng thái pháp lý

Mức độ tin cậy dữ liệu
•	Related Entity
Ví dụ:
Địa bàn liên quan

Đồ án liên quan

Văn bản liên quan

Bản đồ liên quan

Layer GIS liên quan

Thửa đất liên quan

Báo cáo / Snapshot liên quan
 
C.2.4.5 Desktop Layout
+----------------------------------------------------------------------------------+
| HEADER QH PRO                                                                    |
| Logo | Tra cứu | Thư viện | Báo cáo | API | Tải App | Đăng nhập                  |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| BREADCRUMB                                                                       |
| Trang chủ > Hà Nội > Quy hoạch > Khu quy hoạch ven sông Hồng                    |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| HERO SEO                                                                         |
| +----------------------------------------------------------+------------------+ |
| | H1                                                       | Region Preview   | |
| | Quy hoạch Khu quy hoạch ven sông Hồng                    |                  | |
| |                                                          | Mini Map /       | |
| | Badge: Khu quy hoạch | Đang hiệu lực | Dữ liệu công khai  | Boundary Image   | |
| |                                                          |                  | |
| | AI Summary ngắn                                          | [Mở trên bản đồ] | |
| |                                                          |                  | |
| | CTA: Mở bản đồ | Xem đồ án | Chia sẻ                    |                  | |
| +----------------------------------------------------------+------------------+ |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| KEY FACTS                                                                        |
| [Loại vùng] [Địa bàn] [Diện tích] [Trạng thái] [Đồ án nguồn] [Cập nhật]         |
+----------------------------------------------------------------------------------+

+---------------------------------------------+------------------------------------+
| MAIN CONTENT                                | RELATED ENTITY                     |
|                                             |                                    |
| Tổng quan khu/vùng quy hoạch                | Địa bàn liên quan                  |
|                                             |                                    |
| Phạm vi và ranh giới                        | Đồ án quy hoạch liên quan          |
|                                             |                                    |
| Thông tin quy hoạch chính                   | Văn bản pháp lý liên quan          |
|                                             |                                    |
| Bản đồ / Layer GIS liên quan                | Bản đồ quy hoạch liên quan         |
|                                             |                                    |
| Timeline / Trạng thái hiệu lực nếu có       | Khu vực lân cận                    |
|                                             |                                    |
| Lưu ý và cảnh báo dữ liệu                   | Báo cáo / Snapshot liên quan       |
|                                             |                                    |
| FAQ                                         |                                    |
+---------------------------------------------+------------------------------------+

+----------------------------------------------------------------------------------+
| CTA                                                                              |
| Mở trên bản đồ | Xem trong Thư viện Quy hoạch | Theo dõi | Chia sẻ | Tạo báo cáo   |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| DISCLAIMER / LAST UPDATED                                                        |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| FOOTER QH PRO                                                                    |
+----------------------------------------------------------------------------------+
 
C.2.4.6 Responsive / Tablet Layout
+--------------------------------------------------------------+
| HEADER QH PRO                                                |
+--------------------------------------------------------------+

| Breadcrumb                                                   |
+--------------------------------------------------------------+

| H1                                                           |
| Badge trạng thái                                             |
| AI Summary                                                   |
+--------------------------------------------------------------+

| Region Preview / Mini Map                                    |
| CTA: Mở trên bản đồ                                          |
+--------------------------------------------------------------+

| Key Facts dạng card                                          |
| [Loại vùng] [Địa bàn] [Diện tích] [Trạng thái]               |
+--------------------------------------------------------------+

| Tổng quan khu/vùng quy hoạch                                 |
+--------------------------------------------------------------+

| Phạm vi và ranh giới                                         |
+--------------------------------------------------------------+

| Thông tin quy hoạch chính                                    |
+--------------------------------------------------------------+

| Bản đồ / Layer GIS liên quan                                 |
+--------------------------------------------------------------+

| Related Entity dạng card / carousel                          |
+--------------------------------------------------------------+

| FAQ dạng accordion                                           |
+--------------------------------------------------------------+

| CTA                                                          |
+--------------------------------------------------------------+

| Disclaimer / Footer                                          |
+--------------------------------------------------------------+
 
C.2.4.7 Mobile Web / WebView Layout
+--------------------------------------+
| Mobile Header QH Pro                 |
+--------------------------------------+

| Breadcrumb rút gọn                   |
+--------------------------------------+

| H1                                   |
| Badge: Loại vùng / Trạng thái        |
+--------------------------------------+

| AI Summary ngắn                      |
+--------------------------------------+

| CTA Sticky                           |
| [Mở bản đồ] [Chia sẻ]                |
+--------------------------------------+

| Mini Map / Boundary Preview          |
+--------------------------------------+

| Key Facts                            |
+--------------------------------------+

| Accordion                            |
| - Tổng quan                          |
| - Phạm vi                            |
| - Quy hoạch                          |
| - Đồ án / Văn bản                    |
| - Bản đồ liên quan                   |
| - FAQ                                |
+--------------------------------------+

| Related Content                      |
+--------------------------------------+

| Disclaimer                           |
+--------------------------------------+

| Footer                               |
+--------------------------------------+
 
C.2.4.8 Mobile App Layout
+--------------------------------------+
| App Header                           |
| Back | Tên khu/vùng | Share          |
+--------------------------------------+

| H1 / Region Name                     |
| Badge: Khu quy hoạch / Hiệu lực      |
+--------------------------------------+

| AI Summary Card                      |
+--------------------------------------+

| Mini Map / Boundary Preview          |
+--------------------------------------+

| Tab                                  |
| [Tổng quan] [Quy hoạch] [Hồ sơ]      |
| [Bản đồ] [Liên quan]                 |
+--------------------------------------+

| Content Cards                        |
| - Phạm vi                            |
| - Quy hoạch chính                    |
| - Đồ án nguồn                        |
| - Văn bản                            |
+--------------------------------------+

| Related Entity                       |
+--------------------------------------+

| Bottom CTA                           |
| [Mở bản đồ] [Theo dõi] [Chia sẻ]     |
+--------------------------------------+

C.2.4.9 SEO Block bắt buộc
Bắt buộc phải có:
•	Breadcrumb. 
•	H1. 
•	AI Summary. 
•	Badge loại vùng / trạng thái. 
•	Key Facts. 
•	Region Preview / Mini Map. 
•	Tổng quan khu/vùng quy hoạch. 
•	Phạm vi và ranh giới. 
•	Thông tin quy hoạch chính. 
•	Related Entity. 
•	CTA mở trên bản đồ. 
•	Disclaimer. 
•	Footer. 
 
C.2.4.10 SEO Block tùy chọn
Có thể bật/tắt theo dữ liệu:
•	FAQ. 
•	Timeline pháp lý. 
•	Quan hệ phiên bản. 
•	So sánh vùng trước/sau điều chỉnh. 
•	Danh sách thửa đất liên quan. 
•	Báo cáo liên quan. 
•	Snapshot liên quan. 
•	AI Insight nâng cao. 
•	Cảnh báo xung đột dữ liệu. 
•	Breakdown diện tích theo lớp quy hoạch. 
 
C.2.4.11 Data Mapping
Block	Nguồn dữ liệu
H1	Planning Region
Breadcrumb	Administrative Unit + Planning Region
Badge trạng thái	Planning Region Metadata / Legal Status
AI Summary	AI Summary Engine
Key Facts	Planning Region Metadata
Region Preview	GIS Engine / Map Preview Service
Phạm vi & ranh giới	Geometry / Boundary Metadata
Thông tin quy hoạch chính	Planning Information Model
Đồ án liên quan	Planning Project
Văn bản liên quan	Legal Document
Bản đồ / Layer liên quan	Planning Map / GIS Layer
Related Entity	Knowledge Graph
FAQ	SEO Engine
Disclaimer	Legal Safe Policy
Last Updated	Data Version / Updated At
 
C.2.4.12 Render Rule
Thành phần	Render
Metadata	SSR
Canonical	SSR
Breadcrumb	SSR
H1	SSR
Badge trạng thái	SSR
AI Summary	SSR
Key Facts	SSR
Planning Summary	SSR
Main SEO Content	SSR
FAQ	SSR nếu có
Related Entity	SSR / ISR
Schema JSON-LD	SSR
Region Preview	SSR image hoặc CSR nhẹ
Full Map Viewer	CSR
CTA / Share	CSR được phép
Quy tắc:
Không phụ thuộc vào Full Map Engine để render nội dung SEO chính.

Full Map chỉ được mở khi user bấm CTA “Mở trên bản đồ”.
 
C.2.4.13 CTA & Internal Link
•	CTA chính
Mở vùng này trên bản đồ
•	CTA phụ
Xem đồ án liên quan

Xem văn bản pháp lý

Xem bản đồ quy hoạch

Theo dõi khu vực

Tạo báo cáo

Chia sẻ

Tải App
•	Internal Link bắt buộc
Planning Region
→ Administrative Unit

Planning Region
→ Planning Project

Planning Project
→ Legal Document

Planning Project
→ Planning Map

Planning Region
→ GIS Layer

Planning Region
→ Report / Snapshot
•	Internal Link khuyến nghị
Planning Region
→ Planning Region lân cận

Planning Region
→ Parcel liên quan

Planning Region
→ Compare / biến động nếu có dữ liệu
 
C.2.4.14 Mapping áp dụng trong QH Pro
Mục SEO	Áp dụng
D.1.5 Chi tiết vùng quy hoạch	✓
D.2.4 Tra cứu polygon / vùng phân tích	✓ nếu vùng được public và chuẩn hóa
D.3.1 Layer Public	✓ khi layer đại diện cho một vùng/khu quy hoạch
D.4.1 Màn hình kết quả diễn giải	✓ khi output là vùng/khu quy hoạch
D.4.2 Màn hình phân tích quy hoạch	✓ khi vùng phân tích được public
D.5.1 So sánh quy hoạch	✓ nếu có URL công khai cho vùng
D.5.3 Biến động quy hoạch	✓ nếu biến động gắn với vùng
D.7.1 Snapshot công khai	✓ nếu Snapshot gắn với Planning Region
D.7.3 Report công khai	✓ nếu Report gắn với Planning Region
•	Ghi chú
Planning Region SEO Layout là Layout chuẩn cho toàn bộ
trang SEO lấy khu/vùng quy hoạch làm Entity trung tâm.

Không dùng layout này cho polygon tạm thời do người dùng tự vẽ
nếu polygon chưa được chuẩn hóa thành vùng công khai có giá trị SEO.

 
C.2.5 Planning Project SEO Layout
Áp dụng:
•	Đồ án quy hoạch. 
•	Quy hoạch chung. 
•	Quy hoạch phân khu. 
•	Quy hoạch chi tiết. 
•	Quy hoạch sử dụng đất. 
•	Đồ án điều chỉnh quy hoạch. 
•	Đồ án quy hoạch chuyên ngành có công bố công khai. 
 
C.2.5.1 Vai trò & phạm vi áp dụng
Planning Project SEO Layout là mẫu trang đích SEO chuẩn dành cho các Đồ án Quy hoạch (Planning Project) trong QH Pro.
Đây là một trong những loại Landing Page SEO quan trọng nhất của hệ thống vì phần lớn dữ liệu quy hoạch đều được tổ chức xoay quanh các đồ án quy hoạch được phê duyệt bởi cơ quan nhà nước có thẩm quyền.
Mục tiêu:
•	Tạo Landing Page SEO cho từng đồ án quy hoạch. 
•	Tạo điểm truy cập chính thức tới hồ sơ quy hoạch trên QH Pro. 
•	Tạo Knowledge Graph giữa đồ án, địa bàn, văn bản, bản đồ, layer GIS và các phiên bản điều chỉnh. 
•	Tạo nguồn dữ liệu chuẩn cho Google Search và AI Search. 
•	Dẫn người dùng tới Thư viện số Quy hoạch và Map View. 
Planning Project SEO Page không thay thế màn hình Chi tiết Đồ án trong Thư viện Quy hoạch.
Planning Project SEO Page là lớp Public SEO được tối ưu cho:
•	Google Search. 
•	AI Search. 
•	Discover. 
•	Chia sẻ liên kết. 
•	Internal Linking. 
 
C.2.5.2 Entity áp dụng
•	Entity chính
Planning Project
Bao gồm:
Quy hoạch chung

Quy hoạch phân khu

Quy hoạch chi tiết

Quy hoạch sử dụng đất

Đồ án điều chỉnh

Đồ án chuyên ngành
 
•	Entity phụ
Administrative Unit

Planning Region

Legal Document

Planning Map

GIS Layer

Parcel

Report

Snapshot
 
C.2.5.3 Thành phần tĩnh
•	Header QH Pro
Logo

Tra cứu

Thư viện Quy hoạch

Báo cáo

API

Tải App

Đăng nhập
 
•	Footer QH Pro
Giới thiệu

Liên hệ

Điều khoản

Chính sách

API

Tải App
 
•	Disclaimer
Thông tin quy hoạch được tổng hợp từ hồ sơ, bản đồ và văn bản quy hoạch đã được chuẩn hóa trong QH Pro.

Người sử dụng cần đối chiếu hồ sơ pháp lý chính thức trước khi sử dụng cho mục đích pháp lý hoặc giao dịch.
 
C.2.5.4 Thành phần động
•	Breadcrumb
Ví dụ:
Trang chủ

→ Hà Nội

→ Quy hoạch

→ Quy hoạch phân khu S4

→ Đồ án quy hoạch
 
•	H1
Ví dụ:
Đồ án Quy hoạch phân khu đô thị S4, Hà Nội
 
•	AI Summary
Bao gồm:
Mục tiêu đồ án

Phạm vi nghiên cứu

Địa bàn áp dụng

Loại quy hoạch

Trạng thái hiệu lực

Các thay đổi đáng chú ý

Văn bản pháp lý chính
 
•	Key Facts
Ví dụ:
Tên đồ án

Loại quy hoạch

Cấp quy hoạch

Địa bàn

Cơ quan phê duyệt

Số quyết định

Ngày phê duyệt

Trạng thái hiệu lực

Ngày cập nhật
 
•	Legal Timeline
Ví dụ:
Lập quy hoạch

Thẩm định

Phê duyệt

Điều chỉnh

Bổ sung

Thay thế

Hết hiệu lực
 
•	Related Entity
Ví dụ:
Địa bàn liên quan

Vùng quy hoạch liên quan

Văn bản pháp lý

Bản đồ quy hoạch

Layer GIS

Phiên bản liên quan

Đồ án liên quan
 
C.2.5.5 Desktop Layout
+----------------------------------------------------------------------------------+
| HEADER QH PRO                                                                    |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| BREADCRUMB                                                                       |
| Trang chủ > Hà Nội > Quy hoạch > Quy hoạch phân khu S4                           |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| HERO SEO                                                                         |
| +----------------------------------------------------------+------------------+ |
| | H1                                                       | Project Preview  | |
| | Quy hoạch phân khu S4                                    |                  | |
| |                                                          | Cover / Map      | |
| | Badge: Đang hiệu lực                                     |                  | |
| | Badge: Quy hoạch phân khu                                |                  | |
| |                                                          | CTA Mở đồ án     | |
| | AI Summary                                               |                  | |
| +----------------------------------------------------------+------------------+ |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| KEY FACTS                                                                        |
| [Loại] [Cấp] [Địa bàn] [QĐ phê duyệt] [Hiệu lực] [Cập nhật]                     |
+----------------------------------------------------------------------------------+

+---------------------------------------------+------------------------------------+
| MAIN CONTENT                                | RELATED ENTITY                     |
|                                             |                                    |
| Tổng quan đồ án                             | Địa bàn liên quan                  |
|                                             |                                    |
| Phạm vi nghiên cứu                          | Văn bản liên quan                  |
|                                             |                                    |
| Mục tiêu quy hoạch                          | Bản đồ liên quan                   |
|                                             |                                    |
| Chỉ tiêu chính                              | Layer GIS liên quan                |
|                                             |                                    |
| Timeline pháp lý                            | Phiên bản liên quan                |
|                                             |                                    |
| Hồ sơ & tài liệu                            | Đồ án liên quan                    |
|                                             |                                    |
| FAQ                                         |                                    |
+---------------------------------------------+------------------------------------+

+----------------------------------------------------------------------------------+
| LEGAL TIMELINE                                                                    |
| Draft → Approved → Adjusted → Effective → Replaced                               |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| CTA                                                                              |
| Xem trong thư viện | Mở bản đồ | Tải hồ sơ | Chia sẻ | Theo dõi                  |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| DISCLAIMER                                                                        |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| FOOTER                                                                           |
+----------------------------------------------------------------------------------+
 
C.2.5.6 Responsive / Tablet Layout
+--------------------------------------------------------------+
| HEADER                                                       |
+--------------------------------------------------------------+

| Breadcrumb                                                   |
+--------------------------------------------------------------+

| H1                                                           |
| Badge                                                       |
| AI Summary                                                   |
+--------------------------------------------------------------+

| Project Preview                                              |
+--------------------------------------------------------------+

| Key Facts                                                    |
+--------------------------------------------------------------+

| Tổng quan đồ án                                              |
+--------------------------------------------------------------+

| Mục tiêu quy hoạch                                           |
+--------------------------------------------------------------+

| Timeline                                                     |
+--------------------------------------------------------------+

| Hồ sơ & tài liệu                                              |
+--------------------------------------------------------------+

| Related Entity                                               |
+--------------------------------------------------------------+

| FAQ                                                          |
+--------------------------------------------------------------+

| CTA                                                          |
+--------------------------------------------------------------+
 
C.2.5.7 Mobile Web / WebView Layout
+--------------------------------------+
| Mobile Header                        |
+--------------------------------------+

| Breadcrumb                           |
+--------------------------------------+

| H1                                   |
| Badge                                |
+--------------------------------------+

| AI Summary                           |
+--------------------------------------+

| CTA Sticky                           |
| Mở đồ án                             |
+--------------------------------------+

| Preview                              |
+--------------------------------------+

| Key Facts                            |
+--------------------------------------+

| Accordion                            |
| - Tổng quan                          |
| - Quy hoạch                          |
| - Timeline                           |
| - Hồ sơ                              |
| - Văn bản                            |
| - FAQ                                |
+--------------------------------------+

| Related Content                      |
+--------------------------------------+

| Footer                               |
+--------------------------------------+
 
C.2.5.8 Mobile App Layout
+--------------------------------------+
| App Header                           |
+--------------------------------------+

| H1                                   |
| Badge                                |
+--------------------------------------+

| AI Summary Card                      |
+--------------------------------------+

| Project Preview                      |
+--------------------------------------+

| Tabs                                 |
| Tổng quan                            |
| Hồ sơ                                |
| Bản đồ                               |
| Văn bản                              |
| Liên quan                            |
+--------------------------------------+

| Content                              |
+--------------------------------------+

| Related Entity                       |
+--------------------------------------+

| Bottom CTA                           |
| Xem đồ án                            |
+--------------------------------------+
 
C.2.5.9 SEO Block bắt buộc
Bắt buộc phải có:
•	Breadcrumb. 
•	H1. 
•	AI Summary. 
•	Badge loại quy hoạch. 
•	Badge trạng thái hiệu lực. 
•	Key Facts. 
•	Tổng quan đồ án. 
•	Timeline pháp lý. 
•	Hồ sơ & tài liệu. 
•	Related Entity. 
•	CTA. 
•	Disclaimer. 
•	Footer. 
 
C.2.5.10 SEO Block tùy chọn
Có thể bật/tắt:
•	FAQ. 
•	So sánh phiên bản. 
•	So sánh điều chỉnh. 
•	Chỉ tiêu quy hoạch chi tiết. 
•	AI Insight. 
•	Snapshot liên quan. 
•	Report liên quan. 
•	Bản đồ nổi bật. 
•	Phân tích tác động. 
 
C.2.5.11 Data Mapping
Block	Nguồn dữ liệu
H1	Planning Project
AI Summary	AI Summary Engine
Key Facts	Project Metadata
Timeline	Legal Timeline
Hồ sơ	Project Repository
Văn bản	Legal Document
Bản đồ	Planning Map
Layer GIS	GIS Layer
Related Entity	Knowledge Graph
FAQ	SEO Engine
 
C.2.5.12 Render Rule
Thành phần	Render
Metadata	SSR
Canonical	SSR
H1	SSR
AI Summary	SSR
Key Facts	SSR
Timeline	SSR
FAQ	SSR
Related Entity	SSR / ISR
Schema JSON-LD	SSR
Preview	SSR
Full Viewer	CSR
 
C.2.5.13 CTA & Internal Link
•	CTA chính
Xem đồ án trong Thư viện Quy hoạch
•	CTA phụ
Mở trên bản đồ

Xem văn bản pháp lý

Tải hồ sơ

Theo dõi cập nhật

Tạo báo cáo

Chia sẻ
•	Internal Link bắt buộc
Planning Project
→ Administrative Unit

Planning Project
→ Planning Region

Planning Project
→ Legal Document

Planning Project
→ Planning Map

Planning Project
→ GIS Layer

Planning Project
→ Report

Planning Project
→ Snapshot
•	Internal Link khuyến nghị
Planning Project
→ Phiên bản trước

Planning Project
→ Phiên bản điều chỉnh

Planning Project
→ Đồ án liên quan

C.2.5.14 Mapping áp dụng trong QH Pro
Mục SEO	Áp dụng
D.1.5 Chi tiết đồ án quy hoạch	✓
D.4.1 Kết quả diễn giải quy hoạch	✓
D.4.3 Thông tin quy hoạch	✓
D.11.3 Chi tiết đồ án	✓
D.11.8 Timeline pháp lý	✓
D.11.9 Quan hệ phiên bản	✓
D.7.1 Snapshot công khai	✓
D.7.3 Report công khai	✓
•	Ghi chú
Planning Project SEO Layout là Layout trung tâm
của toàn bộ SEO Thư viện Quy hoạch.

Mọi URL SEO lấy Đồ án Quy hoạch làm Entity trung tâm
đều phải sử dụng Layout này.

Đây là Layout có giá trị SEO cao nhất
trong nhóm dữ liệu quy hoạch của QH Pro.

 
C.2.6 Legal Document SEO Layout
Áp dụng:
•	Quyết định. 
•	Nghị quyết. 
•	Thông báo. 
•	Văn bản điều chỉnh. 
•	Văn bản phê duyệt. 
•	Văn bản chấp thuận. 
•	Văn bản thẩm định. 
•	Văn bản hướng dẫn. 
•	Văn bản pháp lý liên quan đến quy hoạch, đất đai và địa chính. 
 
C.2.6.1 Vai trò & phạm vi áp dụng
Legal Document SEO Layout là mẫu trang đích SEO chuẩn dành cho các Văn bản pháp lý trong QH Pro.
Đây là lớp nội dung có giá trị SEO rất cao vì:
•	Người dùng thường tìm kiếm trực tiếp theo số hiệu văn bản. 
•	Google ưu tiên các nội dung có tính tham chiếu pháp lý. 
•	AI Search thường trích dẫn các văn bản nguồn. 
•	Văn bản là cơ sở pháp lý của toàn bộ hệ thống quy hoạch. 
Mục tiêu:
•	Tạo Landing Page SEO cho từng văn bản. 
•	Tạo nguồn tham chiếu pháp lý cho các đồ án quy hoạch. 
•	Tạo Knowledge Graph giữa văn bản, đồ án, địa bàn và bản đồ. 
•	Hỗ trợ AI Search trích xuất chính xác nguồn dữ liệu. 
•	Tăng độ tin cậy và tính pháp lý của QH Pro. 
Legal Document SEO Page không thay thế kho tài liệu gốc.
Trang này đóng vai trò:
Văn bản pháp lý
↓
Giải thích ngữ cảnh
↓
Liên kết dữ liệu quy hoạch
↓
Điều hướng người dùng
 
C.2.6.2 Entity áp dụng
•	Entity chính
Legal Document
Bao gồm:
Quyết định

Nghị quyết

Thông báo

Văn bản phê duyệt

Văn bản điều chỉnh

Văn bản thẩm định

Văn bản hướng dẫn
 
•	Entity phụ
Planning Project

Administrative Unit

Planning Region

Planning Map

GIS Layer

Parcel

Report

Snapshot
 
C.2.6.3 Thành phần tĩnh
•	Header QH Pro
Logo

Tra cứu

Thư viện Quy hoạch

Báo cáo

API

Tải App

Đăng nhập
 
•	Footer QH Pro
Giới thiệu

Điều khoản

Chính sách

API

Liên hệ

Tải App
 
•	Disclaimer
Văn bản được số hóa và chuẩn hóa từ nguồn công khai hoặc nguồn dữ liệu được tích hợp vào QH Pro.

Người sử dụng cần đối chiếu văn bản gốc được công bố bởi cơ quan có thẩm quyền khi sử dụng cho mục đích pháp lý.
 
C.2.6.4 Thành phần động
•	Breadcrumb
Ví dụ:
Trang chủ

→ Thư viện Quy hoạch

→ Văn bản pháp lý

→ Quyết định 1234/QĐ-UBND
 
•	H1
Ví dụ:
Quyết định số 1234/QĐ-UBND phê duyệt Quy hoạch phân khu S4
 
•	AI Summary
Bao gồm:
Mục đích văn bản

Cơ quan ban hành

Đối tượng áp dụng

Địa bàn áp dụng

Đồ án liên quan

Hiệu lực pháp lý

Các điểm đáng chú ý
 
•	Key Facts
Ví dụ:
Số hiệu

Loại văn bản

Ngày ban hành

Ngày hiệu lực

Cơ quan ban hành

Tình trạng hiệu lực

Địa bàn áp dụng

Ngày cập nhật
 
•	Legal Relationship
Ví dụ:
Văn bản được ban hành để

Phê duyệt đồ án nào

Điều chỉnh văn bản nào

Bị thay thế bởi văn bản nào

Liên quan đến đồ án nào
 
•	Related Entity
Ví dụ:
Đồ án liên quan

Địa bàn liên quan

Bản đồ liên quan

Layer GIS liên quan

Văn bản liên quan

Báo cáo liên quan
 
C.2.6.5 Desktop Layout
+----------------------------------------------------------------------------------+
| HEADER QH PRO                                                                    |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| BREADCRUMB                                                                       |
| Trang chủ > Thư viện > Văn bản > QĐ 1234/QĐ-UBND                                |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| HERO SEO                                                                         |
| +----------------------------------------------------------+------------------+ |
| | H1                                                       | Document Card    | |
| | Quyết định 1234/QĐ-UBND                                  |                  | |
| |                                                          | PDF Preview      | |
| | Badge: Quyết định                                        |                  | |
| | Badge: Còn hiệu lực                                      |                  | |
| |                                                          | Tải PDF          | |
| | AI Summary                                               |                  | |
| +----------------------------------------------------------+------------------+ |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| KEY FACTS                                                                        |
| [Loại] [Số hiệu] [Cơ quan] [Ban hành] [Hiệu lực] [Cập nhật]                     |
+----------------------------------------------------------------------------------+

+---------------------------------------------+------------------------------------+
| MAIN CONTENT                                | RELATED ENTITY                     |
|                                             |                                    |
| Tổng quan văn bản                           | Đồ án liên quan                    |
|                                             |                                    |
| Nội dung chính                              | Văn bản liên quan                  |
|                                             |                                    |
| Phạm vi áp dụng                             | Địa bàn liên quan                  |
|                                             |                                    |
| Quan hệ pháp lý                             | Bản đồ liên quan                   |
|                                             |                                    |
| Timeline hiệu lực                           | Layer GIS liên quan                |
|                                             |                                    |
| FAQ                                         | Báo cáo liên quan                  |
+---------------------------------------------+------------------------------------+

+----------------------------------------------------------------------------------+
| LEGAL RELATIONSHIP                                                               |
| Văn bản gốc → Điều chỉnh → Thay thế → Hiệu lực                                  |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| DOCUMENT DOWNLOAD                                                                |
| PDF | DOC | Nguồn công bố                                                        |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| CTA                                                                              |
| Xem đồ án | Mở bản đồ | Tải PDF | Chia sẻ                                        |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| DISCLAIMER                                                                       |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| FOOTER                                                                           |
+----------------------------------------------------------------------------------+
 
C.2.6.6 Responsive / Tablet Layout
+--------------------------------------------------------------+
| HEADER                                                       |
+--------------------------------------------------------------+

| Breadcrumb                                                   |
+--------------------------------------------------------------+

| H1                                                           |
| Badge                                                       |
| AI Summary                                                   |
+--------------------------------------------------------------+

| Document Preview                                             |
+--------------------------------------------------------------+

| Key Facts                                                    |
+--------------------------------------------------------------+

| Tổng quan                                                    |
+--------------------------------------------------------------+

| Nội dung chính                                               |
+--------------------------------------------------------------+

| Quan hệ pháp lý                                              |
+--------------------------------------------------------------+

| Download                                                     |
+--------------------------------------------------------------+

| Related Entity                                               |
+--------------------------------------------------------------+

| FAQ                                                          |
+--------------------------------------------------------------+

| CTA                                                          |
+--------------------------------------------------------------+
 
C.2.6.7 Mobile Web / WebView Layout
+--------------------------------------+
| Mobile Header                        |
+--------------------------------------+

| Breadcrumb                           |
+--------------------------------------+

| H1                                   |
| Badge                                |
+--------------------------------------+

| AI Summary                           |
+--------------------------------------+

| CTA Sticky                           |
| Tải PDF                              |
+--------------------------------------+

| Document Preview                     |
+--------------------------------------+

| Key Facts                            |
+--------------------------------------+

| Accordion                            |
| - Tổng quan                          |
| - Nội dung chính                     |
| - Quan hệ pháp lý                    |
| - Đồ án liên quan                    |
| - FAQ                                |
+--------------------------------------+

| Related Content                      |
+--------------------------------------+

| Footer                               |
+--------------------------------------+
 
C.2.6.8 Mobile App Layout
+--------------------------------------+
| App Header                           |
+--------------------------------------+

| H1                                   |
| Badge                                |
+--------------------------------------+

| AI Summary Card                      |
+--------------------------------------+

| Document Preview                     |
+--------------------------------------+

| Tabs                                 |
| Tổng quan                            |
| Nội dung                             |
| Đồ án                                |
| Bản đồ                               |
| Liên quan                            |
+--------------------------------------+

| Content                              |
+--------------------------------------+

| Related Entity                       |
+--------------------------------------+

| Bottom CTA                           |
| Tải PDF                              |
+--------------------------------------+
 
C.2.6.9 SEO Block bắt buộc
Bắt buộc phải có:
•	Breadcrumb. 
•	H1. 
•	AI Summary. 
•	Badge loại văn bản. 
•	Badge trạng thái hiệu lực. 
•	Key Facts. 
•	Tổng quan văn bản. 
•	Quan hệ pháp lý. 
•	Download nguồn gốc. 
•	Related Entity. 
•	CTA. 
•	Disclaimer. 
•	Footer. 

C.2.6.10 SEO Block tùy chọn
Có thể bật/tắt:
•	FAQ. 
•	Timeline hiệu lực. 
•	So sánh phiên bản. 
•	So sánh trước/sau điều chỉnh. 
•	AI Insight. 
•	Trích dẫn nổi bật. 
•	Snapshot liên quan. 
•	Report liên quan. 
•	Danh sách điều khoản quan trọng. 
 
C.2.6.11 Data Mapping
Block	Nguồn dữ liệu
H1	Legal Document
AI Summary	AI Summary Engine
Key Facts	Legal Metadata
Timeline	Legal Timeline
Quan hệ pháp lý	Legal Relationship Graph
Download	Document Repository
Đồ án liên quan	Planning Project
Bản đồ liên quan	Planning Map
Related Entity	Knowledge Graph
FAQ	SEO Engine
 
C.2.6.12 Render Rule
Thành phần	Render
Metadata	SSR
Canonical	SSR
H1	SSR
AI Summary	SSR
Key Facts	SSR
Quan hệ pháp lý	SSR
Timeline	SSR
FAQ	SSR
Related Entity	SSR / ISR
Schema JSON-LD	SSR
PDF Preview	SSR
PDF Viewer	CSR
 
C.2.6.13 CTA & Internal Link
•	CTA chính
Xem văn bản gốc
•	CTA phụ
Xem đồ án liên quan

Mở trên bản đồ

Tải PDF

Theo dõi cập nhật

Tạo báo cáo

Chia sẻ
•	Internal Link bắt buộc
Legal Document
→ Planning Project

Planning Project
→ Planning Region

Planning Project
→ Administrative Unit

Legal Document
→ Planning Map

Legal Document
→ Report

Legal Document
→ Snapshot
•	Internal Link khuyến nghị
Legal Document
→ Văn bản bị thay thế

Legal Document
→ Văn bản thay thế

Legal Document
→ Văn bản liên quan

Legal Document
→ Timeline pháp lý
 
C.2.6.14 Mapping áp dụng trong QH Pro
Mục SEO	Áp dụng
D.11.4 Danh sách văn bản	✓
D.11.5 Chi tiết văn bản	✓
D.11.8 Timeline pháp lý	✓
D.11.9 Quan hệ phiên bản	✓
D.4.3 Thông tin quy hoạch	✓
D.7.1 Snapshot công khai	✓
D.7.3 Report công khai	✓
•	Ghi chú
Legal Document SEO Layout là Layout chuẩn
cho toàn bộ trang SEO lấy Văn bản pháp lý
làm Entity trung tâm.

Đây là lớp dữ liệu có độ tin cậy cao nhất
trong hệ sinh thái SEO của QH Pro.

Mọi URL SEO của Quyết định, Nghị quyết,
Thông báo, Văn bản phê duyệt và Văn bản điều chỉnh
đều phải sử dụng Layout này.
 
Lưu ý cho QH Pro: Layout này đặc biệt quan trọng với AI Search. Trong nhiều trường hợp, ChatGPT, Gemini, Perplexity sẽ ưu tiên trích dẫn văn bản pháp lý hơn cả trang quy hoạch. Vì vậy phần AI Summary, Quan hệ pháp lý, Timeline hiệu lực và Liên kết tới đồ án quy hoạch là các block có giá trị SEO và AI Search rất cao, không nên lược bỏ.
 
C.2.7 Planning Map SEO Layout
Áp dụng:
•	Bản đồ quy hoạch. 
•	Layer GIS công khai. 
•	Bản đồ chuyên đề. 
•	Bản đồ quy hoạch sử dụng đất. 
•	Bản đồ quy hoạch xây dựng. 
•	Bản đồ quy hoạch giao thông. 
•	Bản đồ quy hoạch phân khu. 
•	Bản đồ quy hoạch chi tiết. 
•	Bản đồ nguồn được công khai trong Thư viện Quy hoạch. 
 
C.2.7.1 Vai trò & phạm vi áp dụng
Planning Map SEO Layout là mẫu trang đích SEO chuẩn dành cho các bản đồ quy hoạch và layer GIS công khai trong QH Pro.
Đây là loại trang SEO có giá trị đặc biệt vì người dùng thường tìm kiếm trực tiếp các truy vấn như:
•	Bản đồ quy hoạch Hà Nội. 
•	Bản đồ quy hoạch sử dụng đất Sơn Tây. 
•	Bản đồ quy hoạch phân khu S4. 
•	Bản đồ quy hoạch giao thông. 
•	Bản đồ quy hoạch chi tiết 1/500. 
•	Layer quy hoạch sử dụng đất. 
Mục tiêu của Planning Map SEO Layout:
•	Tạo Landing Page SEO cho từng bản đồ quy hoạch công khai. 
•	Giúp Google và AI hiểu bản đồ thuộc đồ án nào, địa bàn nào, loại quy hoạch nào. 
•	Liên kết bản đồ với đồ án, văn bản pháp lý, layer GIS và Map View. 
•	Tạo điểm vào trực tiếp để người dùng mở bản đồ trong QH Pro. 
•	Hỗ trợ SEO hình ảnh, SEO bản đồ và Semantic SEO. 
Planning Map SEO Page không thay thế Map View nghiệp vụ.
Trang này chỉ cung cấp bản xem trước, metadata, mô tả, legend, nguồn dữ liệu và CTA để mở bản đồ đầy đủ trong QH Pro.
 
C.2.7.2 Entity áp dụng
•	Entity chính
Planning Map
Bao gồm:
Bản đồ quy hoạch
Bản đồ chuyên đề
Bản đồ nguồn
Bản đồ PDF
Bản đồ CAD đã chuyển đổi preview
Bản đồ ảnh scan
Bản đồ GIS publish
•	Entity phụ
GIS Layer
Planning Project
Legal Document
Administrative Unit
Planning Region
Parcel
Report
Snapshot
 
C.2.7.3 Thành phần tĩnh
•	Header QH Pro
Logo

Tra cứu

Thư viện Quy hoạch

Báo cáo

API

Tải App

Đăng nhập
•	Footer QH Pro
Giới thiệu

Điều khoản

Chính sách

API

Liên hệ

Tải App
•	Disclaimer
Bản đồ quy hoạch được hiển thị dưới dạng bản xem trước hoặc dữ liệu GIS đã được chuẩn hóa trong QH Pro.

Người sử dụng cần đối chiếu bản đồ gốc, hồ sơ quy hoạch và văn bản pháp lý chính thức trước khi sử dụng cho mục đích pháp lý, đầu tư hoặc giao dịch.
•	CTA mặc định
Mở bản đồ trong QH Pro

Xem đồ án nguồn

Xem văn bản pháp lý

Chia sẻ
 
C.2.7.4 Thành phần động
•	Breadcrumb
Ví dụ:
Trang chủ

→ Thư viện Quy hoạch

→ Bản đồ quy hoạch

→ Bản đồ quy hoạch sử dụng đất Sơn Tây
hoặc:
Trang chủ

→ Hà Nội

→ Quy hoạch phân khu S4

→ Bản đồ quy hoạch sử dụng đất
•	H1
Ví dụ:
Bản đồ quy hoạch sử dụng đất Sơn Tây, Hà Nội
hoặc:
Bản đồ quy hoạch phân khu S4, Hà Nội
•	AI Summary
Bao gồm:
Tên bản đồ

Loại bản đồ

Phạm vi địa lý

Đồ án nguồn

Văn bản pháp lý liên quan

Layer GIS tương ứng

Trạng thái dữ liệu

Lưu ý khi sử dụng bản đồ
•	Key Facts
Ví dụ:
Tên bản đồ

Loại bản đồ

Định dạng nguồn

Địa bàn áp dụng

Đồ án nguồn

Tình trạng dữ liệu

Hệ quy chiếu

Ngày cập nhật
•	Map Metadata
Ví dụ:
Nguồn bản đồ

Tỷ lệ bản đồ

Hệ tọa độ

Định dạng gốc

Trạng thái số hóa

Layer GIS liên kết

Legend / chú giải

Ngày publish
•	Related Entity
Ví dụ:
Đồ án nguồn

Văn bản pháp lý

Layer GIS

Địa bàn liên quan

Khu quy hoạch liên quan

Bản đồ liên quan

Báo cáo / Snapshot liên quan
 
C.2.7.5 Desktop Layout
+----------------------------------------------------------------------------------+
| HEADER QH PRO                                                                    |
| Logo | Tra cứu | Thư viện | Báo cáo | API | Tải App | Đăng nhập                  |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| BREADCRUMB                                                                       |
| Trang chủ > Thư viện > Bản đồ quy hoạch > Bản đồ QHSDĐ Sơn Tây                  |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| HERO SEO                                                                         |
| +----------------------------------------------------------+------------------+ |
| | H1                                                       | Map Preview      | |
| | Bản đồ quy hoạch sử dụng đất Sơn Tây                     |                  | |
| |                                                          | Thumbnail /      | |
| | Badge: Bản đồ quy hoạch                                  | Mini Map         | |
| | Badge: Dữ liệu công khai                                 |                  | |
| |                                                          | [Mở bản đồ]      | |
| | AI Summary                                               |                  | |
| +----------------------------------------------------------+------------------+ |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| KEY FACTS                                                                        |
| [Loại bản đồ] [Địa bàn] [Đồ án nguồn] [Định dạng] [Hệ tọa độ] [Cập nhật]       |
+----------------------------------------------------------------------------------+

+---------------------------------------------+------------------------------------+
| MAIN CONTENT                                | RELATED ENTITY                     |
|                                             |                                    |
| Tổng quan bản đồ                            | Đồ án nguồn                        |
|                                             |                                    |
| Phạm vi địa lý                              | Văn bản pháp lý                    |
|                                             |                                    |
| Metadata bản đồ                             | Layer GIS liên kết                 |
|                                             |                                    |
| Chú giải / Legend                           | Bản đồ liên quan                   |
|                                             |                                    |
| Cách sử dụng bản đồ                         | Địa bàn liên quan                  |
|                                             |                                    |
| FAQ                                         | Báo cáo / Snapshot liên quan       |
+---------------------------------------------+------------------------------------+

+----------------------------------------------------------------------------------+
| MAP PREVIEW EXPANDED                                                              |
| Ảnh bản đồ / Preview GIS / Legend rút gọn                                         |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| CTA                                                                              |
| Mở bản đồ trong QH Pro | Xem đồ án nguồn | Xem văn bản | Chia sẻ | Tải App       |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| DISCLAIMER / LAST UPDATED                                                        |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| FOOTER QH PRO                                                                    |
+----------------------------------------------------------------------------------+
 
C.2.7.6 Responsive / Tablet Layout
+--------------------------------------------------------------+
| HEADER QH PRO                                                |
+--------------------------------------------------------------+

| Breadcrumb                                                   |
+--------------------------------------------------------------+

| H1                                                           |
| Badge: Loại bản đồ / Trạng thái                              |
| AI Summary                                                   |
+--------------------------------------------------------------+

| Map Preview                                                  |
| CTA: Mở bản đồ                                               |
+--------------------------------------------------------------+

| Key Facts dạng card                                          |
| [Loại] [Địa bàn] [Đồ án nguồn] [Định dạng]                  |
+--------------------------------------------------------------+

| Tổng quan bản đồ                                             |
+--------------------------------------------------------------+

| Metadata bản đồ                                              |
+--------------------------------------------------------------+

| Chú giải / Legend                                            |
+--------------------------------------------------------------+

| Related Entity dạng card / carousel                          |
+--------------------------------------------------------------+

| FAQ dạng accordion                                           |
+--------------------------------------------------------------+

| CTA                                                          |
+--------------------------------------------------------------+

| Disclaimer / Footer                                          |
+--------------------------------------------------------------+
 
C.2.7.7 Mobile Web / WebView Layout
+--------------------------------------+
| Mobile Header QH Pro                 |
+--------------------------------------+

| Breadcrumb rút gọn                   |
+--------------------------------------+

| H1                                   |
| Badge: Bản đồ / Layer công khai      |
+--------------------------------------+

| AI Summary ngắn                      |
+--------------------------------------+

| CTA Sticky                           |
| [Mở bản đồ] [Chia sẻ]                |
+--------------------------------------+

| Map Preview                          |
+--------------------------------------+

| Key Facts                            |
+--------------------------------------+

| Accordion                            |
| - Tổng quan                          |
| - Metadata                           |
| - Chú giải                           |
| - Đồ án / Văn bản                    |
| - Bản đồ liên quan                   |
| - FAQ                                |
+--------------------------------------+

| Related Content                      |
+--------------------------------------+

| Disclaimer                           |
+--------------------------------------+

| Footer                               |
+--------------------------------------+
 
C.2.7.8 Mobile App Layout
+--------------------------------------+
| App Header                           |
| Back | Tên bản đồ | Share            |
+--------------------------------------+

| H1 / Map Name                        |
| Badge: Bản đồ quy hoạch              |
+--------------------------------------+

| AI Summary Card                      |
+--------------------------------------+

| Map Preview / Mini Viewer            |
+--------------------------------------+

| Tab                                  |
| [Tổng quan] [Metadata] [Legend]      |
| [Đồ án] [Văn bản]                    |
+--------------------------------------+

| Content Cards                        |
| - Phạm vi                            |
| - Định dạng                          |
| - Hệ tọa độ                          |
| - Trạng thái dữ liệu                 |
+--------------------------------------+

| Related Entity                       |
+--------------------------------------+

| Bottom CTA                           |
| [Mở bản đồ] [Tải / Chia sẻ]          |
+--------------------------------------+
 
C.2.7.9 SEO Block bắt buộc
Bắt buộc phải có:
•	Breadcrumb. 
•	H1. 
•	AI Summary. 
•	Badge loại bản đồ. 
•	Badge trạng thái dữ liệu. 
•	Key Facts. 
•	Map Preview. 
•	Metadata bản đồ. 
•	Đồ án nguồn. 
•	Văn bản pháp lý liên quan. 
•	Related Entity. 
•	CTA mở bản đồ. 
•	Disclaimer. 
•	Footer. 

C.2.7.10 SEO Block tùy chọn
Có thể bật/tắt theo dữ liệu:
•	FAQ. 
•	Legend chi tiết. 
•	Bản đồ preview mở rộng. 
•	Link tải file gốc. 
•	Link xem CAD/PDF. 
•	Layer GIS liên kết. 
•	So sánh bản đồ trước/sau. 
•	Timeline dữ liệu. 
•	Snapshot liên quan. 
•	Report liên quan. 
•	AI Insight nâng cao. 
 
C.2.7.11 Data Mapping
Block	Nguồn dữ liệu
H1	Planning Map
Breadcrumb	Planning Library / Administrative Unit / Planning Project
Badge loại bản đồ	Planning Map Metadata
Badge trạng thái	Publish Status / Legal Status
AI Summary	AI Summary Engine
Key Facts	Planning Map Metadata
Map Preview	Map Preview Service / GIS Engine
Metadata bản đồ	Planning Map Metadata
Legend	Layer Legend / Map Legend
Đồ án nguồn	Planning Project
Văn bản pháp lý	Legal Document
Layer GIS liên kết	GIS Layer
Related Entity	Knowledge Graph
FAQ	SEO Engine
Last Updated	Data Version / Updated At
 
C.2.7.12 Render Rule
Thành phần	Render
Metadata	SSR
Canonical	SSR
Breadcrumb	SSR
H1	SSR
Badge	SSR
AI Summary	SSR
Key Facts	SSR
Main SEO Content	SSR
Metadata bản đồ	SSR
Legend rút gọn	SSR nếu có
Related Entity	SSR / ISR
FAQ	SSR nếu có
Schema JSON-LD	SSR
Map Preview Image	SSR
Full Map Viewer	CSR
File Viewer PDF/CAD	CSR
CTA / Share	CSR được phép
Quy tắc:
Trang SEO bản đồ không được phụ thuộc vào Full Map Viewer
để hiển thị nội dung index chính.

Bản đồ đầy đủ chỉ mở khi user bấm CTA “Mở bản đồ”.
 
C.2.7.13 CTA & Internal Link
•	CTA chính
Mở bản đồ trong QH Pro
•	CTA phụ
Xem đồ án nguồn

Xem văn bản pháp lý

Xem layer GIS

Xem legend

Tải tài liệu nếu được phép

Chia sẻ

Tải App
•	Internal Link bắt buộc
Planning Map
→ Planning Project

Planning Map
→ Legal Document

Planning Map
→ GIS Layer

Planning Map
→ Administrative Unit

Planning Map
→ Planning Region
•	Internal Link khuyến nghị
Planning Map
→ Bản đồ liên quan

Planning Map
→ Phiên bản bản đồ trước/sau

Planning Map
→ Report / Snapshot

Planning Map
→ Tra cứu bản đồ theo địa bàn
 
C.2.7.14 Mapping áp dụng trong QH Pro
Mục SEO	Áp dụng
D.3.1 Layer Public	✓
D.3.2 Layer Metadata	✓
D.3.3 Layer Category	✓ nếu category có URL công khai
D.11.6 Danh sách bản đồ	✓
D.11.7 Chi tiết bản đồ	✓
D.11.3 Chi tiết đồ án	✓ với tab Bản đồ & GIS
D.4.3 Thông tin quy hoạch	✓ khi có bản đồ liên quan
D.7.1 Snapshot công khai	✓ nếu Snapshot lấy bản đồ làm đối tượng trung tâm
D.7.3 Report công khai	✓ nếu Report gắn với bản đồ/layer
•	Ghi chú
Planning Map SEO Layout là Layout chuẩn cho toàn bộ
trang SEO lấy bản đồ quy hoạch hoặc layer GIS công khai
làm Entity trung tâm.

Không áp dụng layout này cho layer nội bộ, layer kỹ thuật,
layer tạm thời hoặc layer không được publish công khai.

 
C.2.8 Report & Analysis SEO Layout
Áp dụng:
•	Báo cáo công khai. 
•	Snapshot công khai. 
•	Kết quả phân tích quy hoạch. 
•	Kết quả diễn giải quy hoạch. 
•	AI Summary Page. 
•	Báo cáo chuyên đề. 
•	Báo cáo tổng hợp theo địa bàn. 
•	Báo cáo tổng hợp theo đồ án. 
•	Báo cáo tổng hợp theo khu vực. 
•	Báo cáo sinh tự động từ AI Engine. 
 
C.2.8.1 Vai trò & phạm vi áp dụng
Report & Analysis SEO Layout là mẫu trang đích SEO chuẩn dành cho các nội dung phân tích, tổng hợp và diễn giải dữ liệu trong QH Pro.
Khác với các Layout trước:
•	Administrative Unit → tập trung vào địa bàn. 
•	Parcel → tập trung vào thửa đất. 
•	Planning Region → tập trung vào vùng quy hoạch. 
•	Planning Project → tập trung vào đồ án. 
•	Legal Document → tập trung vào văn bản. 
•	Planning Map → tập trung vào bản đồ. 
Report & Analysis SEO Layout tập trung vào:
Tri thức

Kết luận

Phân tích

Diễn giải

Báo cáo

Insight
Đây là nhóm Landing Page có giá trị rất cao đối với:
•	Google Search. 
•	AI Search. 
•	Discover. 
•	Social Share. 
•	SEO Content. 
•	Topical Authority. 
Mục tiêu:
•	Chuyển dữ liệu GIS thành nội dung dễ hiểu. 
•	Chuyển dữ liệu quy hoạch thành tri thức. 
•	Tạo nội dung SEO quy mô lớn. 
•	Tạo nguồn dữ liệu AI-friendly. 
•	Tăng lượng truy cập từ các truy vấn dài (Long-tail SEO). 
Ví dụ:
Quy hoạch phường Tùng Thiện có gì đáng chú ý?

Đất ở khu vực nào bị ảnh hưởng bởi quy hoạch?

Phân tích quy hoạch phân khu S4.

Các thay đổi chính của quy hoạch sử dụng đất Sơn Tây.

AI tóm tắt Quy hoạch chung Hà Nội.
 
C.2.8.2 Entity áp dụng
•	Entity chính
Report

Analysis

AI Summary

Snapshot
Bao gồm:
Báo cáo

Báo cáo chuyên đề

Snapshot công khai

Kết quả phân tích

Kết quả diễn giải

AI Summary Page

Insight Page
•	Entity phụ
Administrative Unit

Parcel

Planning Region

Planning Project

Legal Document

Planning Map

GIS Layer
 
C.2.8.3 Thành phần tĩnh
•	Header QH Pro
Logo

Tra cứu

Thư viện Quy hoạch

Báo cáo

API

Tải App

Đăng nhập
•	Footer QH Pro
Giới thiệu

Điều khoản

Chính sách

API

Liên hệ

Tải App
•	Disclaimer
Các nội dung phân tích, diễn giải và AI Summary được sinh từ dữ liệu hiện có trong QH Pro.

Thông tin mang tính hỗ trợ tham khảo, không thay thế kết luận pháp lý hoặc ý kiến chuyên môn của cơ quan nhà nước có thẩm quyền.
 
C.2.8.4 Thành phần động
•	Breadcrumb
Ví dụ:
Trang chủ

→ Báo cáo

→ Hà Nội

→ Phân tích quy hoạch Sơn Tây
hoặc:
Trang chủ

→ AI Summary

→ Quy hoạch phân khu S4
 
•	H1
Ví dụ:
Phân tích quy hoạch khu vực Sơn Tây, Hà Nội
hoặc:
AI Summary Quy hoạch phân khu S4
 
•	AI Summary
Đây là Block trung tâm của Layout.
Bao gồm:
Tóm tắt nhanh

Các điểm nổi bật

Các tác động chính

Các rủi ro cần lưu ý

Các đồ án liên quan

Các khu vực liên quan

Các nguồn dữ liệu tham chiếu
 
•	Key Facts
Ví dụ:
Loại báo cáo

Địa bàn

Đồ án

Phạm vi dữ liệu

Thời điểm tạo

Phiên bản dữ liệu

Ngày cập nhật
 
•	Analysis Result
Ví dụ:
Nhận định chính

Xu hướng

Tác động

So sánh

Kết luận

Khuyến nghị

•	Related Entity
Ví dụ:
Địa bàn liên quan

Đồ án liên quan

Bản đồ liên quan

Văn bản liên quan

Snapshot liên quan

Báo cáo liên quan
 
C.2.8.5 Desktop Layout
+----------------------------------------------------------------------------------+
| HEADER QH PRO                                                                    |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| BREADCRUMB                                                                       |
| Trang chủ > Báo cáo > Hà Nội > Phân tích quy hoạch Sơn Tây                      |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| HERO SEO                                                                         |
| +----------------------------------------------------------+------------------+ |
| | H1                                                       | Cover Image      | |
| | Phân tích quy hoạch Sơn Tây                              |                  | |
| |                                                          | Snapshot Preview | |
| | Badge: Report                                            |                  | |
| | Badge: Public                                            |                  | |
| |                                                          | Share            | |
| | AI Summary                                               |                  | |
| +----------------------------------------------------------+------------------+ |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| KEY FACTS                                                                        |
| [Loại] [Địa bàn] [Đồ án] [Nguồn dữ liệu] [Phiên bản] [Cập nhật]                |
+----------------------------------------------------------------------------------+

+---------------------------------------------+------------------------------------+
| MAIN CONTENT                                | RELATED ENTITY                     |
|                                             |                                    |
| Executive Summary                           | Đồ án liên quan                    |
|                                             |                                    |
| Phân tích chi tiết                          | Văn bản liên quan                  |
|                                             |                                    |
| Nhận định chính                             | Bản đồ liên quan                   |
|                                             |                                    |
| Tác động                                    | Snapshot liên quan                 |
|                                             |                                    |
| Khuyến nghị                                 | Báo cáo liên quan                  |
|                                             |                                    |
| FAQ                                         |                                    |
+---------------------------------------------+------------------------------------+

+----------------------------------------------------------------------------------+
| DATA SOURCES                                                                      |
| Administrative Unit | Parcel | Project | Map | Document                          |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| CTA                                                                              |
| Mở bản đồ | Xem đồ án | Tạo báo cáo mới | Chia sẻ                               |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| DISCLAIMER                                                                       |
+----------------------------------------------------------------------------------+

+----------------------------------------------------------------------------------+
| FOOTER                                                                           |
+----------------------------------------------------------------------------------+
 
C.2.8.6 Responsive / Tablet Layout
+--------------------------------------------------------------+
| HEADER                                                       |
+--------------------------------------------------------------+

| Breadcrumb                                                   |
+--------------------------------------------------------------+

| H1                                                           |
| Badge                                                       |
| AI Summary                                                   |
+--------------------------------------------------------------+

| Snapshot Preview                                             |
+--------------------------------------------------------------+

| Key Facts                                                    |
+--------------------------------------------------------------+

| Executive Summary                                            |
+--------------------------------------------------------------+

| Phân tích                                                    |
+--------------------------------------------------------------+

| Kết luận                                                     |
+--------------------------------------------------------------+

| Related Entity                                               |
+--------------------------------------------------------------+

| FAQ                                                          |
+--------------------------------------------------------------+

| CTA                                                          |
+--------------------------------------------------------------+
 
C.2.8.7 Mobile Web / WebView Layout
+--------------------------------------+
| Mobile Header                        |
+--------------------------------------+

| Breadcrumb                           |
+--------------------------------------+

| H1                                   |
| Badge                                |
+--------------------------------------+

| AI Summary                           |
+--------------------------------------+

| CTA Sticky                           |
| Chia sẻ                              |
+--------------------------------------+

| Snapshot Preview                     |
+--------------------------------------+

| Key Facts                            |
+--------------------------------------+

| Accordion                            |
| - Tóm tắt                            |
| - Phân tích                          |
| - Kết luận                           |
| - Nguồn dữ liệu                      |
| - FAQ                                |
+--------------------------------------+

| Related Content                      |
+--------------------------------------+

| Footer                               |
+--------------------------------------+
 
C.2.8.8 Mobile App Layout
+--------------------------------------+
| App Header                           |
+--------------------------------------+

| H1                                   |
| Badge                                |
+--------------------------------------+

| AI Summary Card                      |
+--------------------------------------+

| Snapshot Preview                     |
+--------------------------------------+

| Tabs                                 |
| Tổng quan                            |
| Phân tích                            |
| Dữ liệu                              |
| Liên quan                            |
+--------------------------------------+

| Content                              |
+--------------------------------------+

| Related Entity                       |
+--------------------------------------+

| Bottom CTA                           |
| Mở bản đồ                            |
+--------------------------------------+
 
C.2.8.9 SEO Block bắt buộc
Bắt buộc phải có:
•	Breadcrumb. 
•	H1. 
•	AI Summary. 
•	Badge loại báo cáo. 
•	Key Facts. 
•	Executive Summary. 
•	Analysis Result. 
•	Data Sources. 
•	Related Entity. 
•	CTA. 
•	Disclaimer. 
•	Footer. 
 
C.2.8.10 SEO Block tùy chọn
Có thể bật/tắt:
•	FAQ. 
•	Biểu đồ. 
•	Timeline. 
•	Compare View. 
•	AI Insight nâng cao. 
•	Snapshot Gallery. 
•	Report Download. 
•	Bảng số liệu. 
•	Video Summary. 
•	Change History. 
 
C.2.8.11 Data Mapping
Block	Nguồn dữ liệu
H1	Report / Analysis
AI Summary	AI Engine
Key Facts	Report Metadata
Executive Summary	Analysis Engine
Analysis Result	Analysis Engine
Data Sources	Data Lineage Engine
Related Entity	Knowledge Graph
FAQ	SEO Engine
Snapshot Preview	Snapshot Service
 
C.2.8.12 Render Rule
Thành phần	Render
Metadata	SSR
Canonical	SSR
H1	SSR
AI Summary	SSR
Executive Summary	SSR
Analysis Result	SSR
Data Sources	SSR
FAQ	SSR
Related Entity	SSR / ISR
Schema JSON-LD	SSR
Snapshot Preview	SSR
Chart Interactive	CSR
Compare Tool	CSR
Quy tắc:
Toàn bộ nội dung phân tích chính phải có mặt trong HTML SSR.

Không được phụ thuộc vào JavaScript để sinh phần nội dung SEO cốt lõi.

C.2.8.13 CTA & Internal Link
•	CTA chính
Mở trên bản đồ
•	CTA phụ
Xem đồ án liên quan

Xem văn bản liên quan

Tạo báo cáo

Tải App

Chia sẻ
•	Internal Link bắt buộc
Report
→ Administrative Unit

Report
→ Planning Project

Report
→ Planning Region

Report
→ Planning Map

Report
→ Legal Document

Report
→ Snapshot
•	Internal Link khuyến nghị
Report
→ Report liên quan

AI Summary
→ Báo cáo gốc

Snapshot
→ Report

Report
→ Compare Page
 
C.2.8.14 Mapping áp dụng trong QH Pro
Mục SEO	Áp dụng
D.4.1 Màn hình kết quả diễn giải	✓
D.4.2 Màn hình phân tích quy hoạch	✓
D.4.3 Màn hình thông tin quy hoạch	✓
D.5.1 So sánh quy hoạch	✓
D.5.2 So sánh dữ liệu	✓
D.5.3 Biến động quy hoạch	✓
D.7.1 Snapshot	✓
D.7.2 Share Snapshot	✓
D.7.3 Report	✓
D.7.4 Share Report	✓
•	Ghi chú
Report & Analysis SEO Layout là Layout trung tâm
cho toàn bộ hệ thống SEO Content của QH Pro.

Đây là nhóm Landing Page có khả năng mở rộng số lượng lớn nhất,
phù hợp cho Google Search, AI Search và Discover.

Mọi URL SEO được sinh từ AI Summary,
Report, Snapshot hoặc Analysis đều phải sử dụng Layout này.
 
•	Lưu ý chiến lược SEO
Trong toàn bộ hệ thống QH Pro:
Administrative Unit
Parcel
Planning Region
Planning Project
Legal Document
Planning Map
là nhóm SEO dữ liệu gốc.
Còn:
Report
Analysis
Snapshot
AI Summary
là nhóm SEO nội dung.
Đây sẽ là nguồn tạo traffic lớn nhất về dài hạn vì có khả năng sinh ra hàng chục nghìn đến hàng triệu trang SEO chất lượng từ cùng một kho dữ liệu quy hoạch.

 
C.2.9 Ma trận Mapping SEO Layout ↔ Entity ↔ Màn hình SRS
•	Vai trò
Mục này đóng vai trò là ma trận liên kết trung tâm giữa:
•	Entity SEO. 
•	SEO Layout. 
•	Màn hình nghiệp vụ trong SRS. 
•	URL SEO. 
•	SEO Class. 
•	Schema Type. 
Mục tiêu:
•	Tránh áp dụng sai Layout. 
•	Tránh tạo URL SEO sai Entity. 
•	Tránh trùng lặp SEO giữa các màn hình. 
•	Giúp Dev xác định chính xác Layout cần sử dụng. 
•	Giúp QA nghiệm thu nhanh. 
•	Giúp Backend sinh URL, Schema và Metadata thống nhất. 
Nguyên tắc:
Một Entity chính
→ Một SEO Layout chuẩn

Một URL SEO
→ Một Layout chuẩn

Một màn hình SRS
→ Có thể sinh nhiều URL SEO

Nhưng mỗi URL SEO
→ Chỉ được dùng một Layout chính
 
C.2.9.1 Mapping Layout ↔ Entity
SEO Layout	Entity chính
Administrative Unit SEO Layout	Administrative Unit
Parcel SEO Layout	Parcel
Planning Region SEO Layout	Planning Region
Planning Project SEO Layout	Planning Project
Legal Document SEO Layout	Legal Document
Planning Map SEO Layout	Planning Map
Report & Analysis SEO Layout	Report / Analysis / Snapshot / AI Summary
•	Quy tắc
Administrative Unit
→ Không dùng Parcel Layout

Parcel
→ Không dùng Planning Project Layout

Planning Project
→ Không dùng Legal Document Layout

Legal Document
→ Không dùng Report Layout
Mỗi Entity chỉ có một Layout chuẩn.
 
C.2.9.2 Mapping Layout ↔ D.x.x
•	Administrative Unit SEO Layout
Màn hình SRS
D.2.1 Tra cứu địa chỉ
D.11.1 Workspace thư viện
D.11.2 Danh sách đồ án
D.11.3 Chi tiết đồ án (theo địa bàn)
D.11.5 Chi tiết văn bản (theo địa bàn)
 
•	Parcel SEO Layout
Màn hình SRS
D.1.5 Chi tiết thửa đất
D.2.2 Tra cứu GPS
D.2.3 Tra cứu tờ/thửa
D.2.1 Tra cứu địa chỉ (khi resolve sang Parcel)
D.7.1 Snapshot
D.7.3 Report
 
•	Planning Region SEO Layout
Màn hình SRS
D.1.5 Chi tiết vùng quy hoạch
D.2.4 Tra cứu polygon
D.3.1 Layer Public
D.4.1 Diễn giải quy hoạch
D.4.2 Phân tích quy hoạch
D.5.1 So sánh quy hoạch
D.5.3 Biến động quy hoạch
 
•	Planning Project SEO Layout
Màn hình SRS
D.1.5 Chi tiết đồ án
D.4.1 Diễn giải quy hoạch
D.4.3 Thông tin quy hoạch
D.11.3 Chi tiết đồ án
D.11.8 Timeline pháp lý
D.11.9 Quan hệ phiên bản
 
•	Legal Document SEO Layout
Màn hình SRS
D.11.4 Danh sách văn bản
D.11.5 Chi tiết văn bản
D.11.8 Timeline pháp lý
D.11.9 Quan hệ phiên bản
D.4.3 Thông tin quy hoạch
 
•	Planning Map SEO Layout
Màn hình SRS
D.3.1 Layer Public
D.3.2 Layer Metadata
D.3.3 Layer Category
D.11.6 Danh sách bản đồ
D.11.7 Chi tiết bản đồ
 
•	Report & Analysis SEO Layout
Màn hình SRS
D.4.1 Kết quả diễn giải
D.4.2 Kết quả phân tích
D.5.1 So sánh quy hoạch
D.5.2 So sánh dữ liệu
D.5.3 Biến động
D.7.1 Snapshot
D.7.2 Share Snapshot
D.7.3 Report
D.7.4 Share Report
 
C.2.9.3 Mapping Layout ↔ URL SEO
SEO Layout	URL Pattern
Administrative Unit	/dia-ban/{slug}
Parcel	/thua-dat/{slug}
Planning Region	/vung-quy-hoach/{slug}
Planning Project	/do-an-quy-hoach/{slug}
Legal Document	/van-ban/{slug}
Planning Map	/ban-do/{slug}
Report & Analysis	/bao-cao/{slug}
 
•	Ví dụ
/dia-ban/ha-noi

/dia-ban/ha-noi/son-tay

/thua-dat/ha-noi-son-tay-to-12-thua-456

/vung-quy-hoach/ven-song-hong

/do-an-quy-hoach/quy-hoach-phan-khu-s4

/van-ban/1234-qd-ubnd

/ban-do/qhsdd-son-tay

/bao-cao/phan-tich-quy-hoach-son-tay
 
•	URL không được index
/map

/search

/compare

/report-builder

/snapshot-builder

/api

/admin

/account

/pricing

/login
 
C.2.9.4 Mapping Layout ↔ SEO Class
SEO Layout	SEO Class
Administrative Unit	SEO-A
Parcel	SEO-A
Planning Region	SEO-A
Planning Project	SEO-A
Legal Document	SEO-A
Planning Map	SEO-B
Report & Analysis	SEO-B / SEO-C
 
•	Giải thích
•	SEO-A
Là tài sản SEO cốt lõi.
Ví dụ:
Địa bàn

Thửa đất

Đồ án

Văn bản
Luôn ưu tiên index.
 
•	SEO-B
Trang điều hướng.
Ví dụ:
Bản đồ

Danh sách

Thư viện
Có thể index.
 
•	SEO-C
Trang sinh động.
Ví dụ:
AI Summary

Snapshot

Report

Analysis
Chỉ index khi đạt điều kiện chất lượng.
 
•	SEO-D
Share URL
Không đưa sitemap.
 
•	SEO-N
Admin

Account

API

Internal Tool
Không SEO.
 
C.2.9.5 Mapping Layout ↔ Schema Type
SEO Layout	Schema Type	
Administrative Unit	Place	
Parcel	Place	
Planning Region	Place	
Planning Project	CreativeWork	
Legal Document	Legislation	
Planning Map	Dataset	
Report & Analysis	Report	
 
•	Schema bổ sung
•	Administrative Unit
Place

AdministrativeArea
 
•	Parcel
Place

GeoShape
 
•	Planning Region
Place

GeoShape
 
•	Planning Project
CreativeWork

Dataset
 
•	Legal Document
Legislation
 
•	Planning Map
Dataset
 
•	Report & Analysis
Report

Article
 
•	Nguyên tắc cuối cùng
Entity
↓
Layout
↓
URL
↓
SEO Class
↓
Schema
↓
SRS Screen
Đây là chuỗi Mapping chuẩn bắt buộc áp dụng trong toàn bộ hệ thống QH Pro.
Mọi URL SEO được tạo mới trong tương lai phải được đối chiếu với ma trận này trước khi triển khai để bảo đảm:
•	Không trùng lặp SEO. 
•	Không sai Entity. 
•	Không sai Layout. 
•	Không sai Schema. 
•	Không sai Sitemap. 
•	Không sai Canonical. 
•	Không làm suy giảm chất lượng SEO tổng thể của hệ thống.
C.3	SEO Reference Library
C.3.0  Vai trò
Là thư viện tiêu chuẩn SEO dùng chung cho toàn bộ hệ thống QH Pro.
Mục này quy định các nguyên tắc SEO cấp hệ thống nhằm bảo đảm tính thống nhất giữa các màn hình, các URL, các Entity và các module của QH Pro.
Mục tiêu:
•	Thống nhất cách triển khai SEO trong toàn hệ thống. 
•	Tránh mỗi màn hình triển khai SEO theo một cách khác nhau. 
•	Làm căn cứ tham chiếu khi xây dựng nội dung SEO tại Phần D. 
•	Giảm trùng lặp nội dung giữa các màn hình. 
•	Giảm khối lượng bảo trì khi thay đổi tiêu chuẩn SEO trong tương lai. 
•	Hỗ trợ Dev, QA và Product cùng sử dụng một bộ quy tắc thống nhất. 
Phần này không mô tả SEO của từng màn hình cụ thể.
Các yêu cầu triển khai trực tiếp cho từng màn hình sẽ được mô tả tại:
PHẦN D. NỘI DUNG SEO BỔ SUNG CHO TỪNG MÀN HÌNH SRS
Trong trường hợp có khác biệt giữa:
C.3 SEO Reference Library
và:
D.x.x Nội dung SEO của màn hình
thì:
Ưu tiên áp dụng D.x.x

C.3.1 URL Rule
•	Vai trò
Quy định nguyên tắc thiết kế URL SEO trong toàn bộ hệ thống.
•	Mục tiêu
Đảm bảo URL:
•	Dễ hiểu. 
•	Dễ đọc. 
•	Có ý nghĩa ngữ nghĩa. 
•	Ổn định lâu dài. 
•	Hỗ trợ Google Search. 
•	Hỗ trợ AI Search. 
•	Hỗ trợ Internal Linking. 
•	Hỗ trợ Knowledge Graph. 
•	Nguyên tắc
Một Entity
→ Một URL chuẩn

Một URL chuẩn
→ Một Entity
URL phải được sinh từ Entity.
Không sinh URL SEO từ:
Filter

Sort

Zoom

Map State

Popup Runtime

Trạng thái giao diện tạm thời
•	Đối tượng được phép tạo URL SEO
Administrative Unit

Parcel

Planning Region

Planning Project

Legal Document

Planning Map

Report

Snapshot

AI Summary (nếu có Landing Page riêng)
•	Đối tượng không thuộc phạm vi URL SEO
Login

Register

Profile

Workspace

Admin

Search Runtime

Filter Runtime

API Runtime

Map Runtime State
 
C.3.2 Metadata Rule
•	Vai trò
Quy định nguyên tắc Metadata dùng chung.
•	Mục tiêu
Đảm bảo mọi Landing Page SEO đều có Metadata đầy đủ, nhất quán và có khả năng được Search Engine hiểu chính xác.
•	Thành phần Metadata tối thiểu
Page Title

Meta Description

H1

Open Graph Title

Open Graph Description

Open Graph Image

AI Snippet

Last Updated
•	Nguyên tắc
Metadata phải được sinh từ dữ liệu hệ thống.
Không hard-code Metadata trong Frontend.
Metadata phải phản ánh:
Entity

Địa bàn

Ngữ cảnh quy hoạch

Trạng thái dữ liệu
•	Không được
Metadata rỗng

Metadata trùng lặp

Metadata không liên quan

Nhồi nhét từ khóa
 
C.3.3 C.3.3 Canonical Rule
•	Vai trò
Quy định URL chuẩn của hệ thống.
•	Mục tiêu
Ngăn chặn:
Duplicate Content

Duplicate URL

Keyword Cannibalization
•	Nguyên tắc
Một Entity
→ Một Canonical URL
Mọi URL phụ phải:
Canonical
hoặc
Redirect
về URL chuẩn.
•	Đối tượng thường phát sinh URL phụ
Alias URL

Historical URL

Search URL

Share URL

Filter URL

Sort URL
 
C.3.4 C.3.4 Schema Rule
•	Vai trò
Quy định nguyên tắc Structured Data.
•	Mục tiêu
Giúp:
Google Search

AI Search

Knowledge Graph

Semantic Search
hiểu đúng dữ liệu QH Pro.
•	Nguyên tắc
Mọi Landing Page SEO phải có Schema phù hợp với Entity chính.
Schema phải:
JSON-LD

Machine Readable

SSR Rendered
•	Schema phải thể hiện
Entity

Relationship

Administrative Context

Legal Context

Source Context
 
C.3.5 C.3.5 Sitemap Rule
•	Vai trò
Quy định nguyên tắc xây dựng Sitemap.
•	Mục tiêu
Kiểm soát:
Discovery

Crawl

Index
•	Nguyên tắc
Chỉ URL:
Public

Canonical

Indexable
được phép xuất hiện trong Sitemap.
•	Không đưa vào Sitemap
Private URL

Draft URL

Technical URL

Search URL

Share URL

Admin URL
 
C.3.6 C.3.6 Index Rule
•	Vai trò
Quy định nguyên tắc Index và NoIndex.
•	Mục tiêu
Đảm bảo chỉ các URL có giá trị SEO thực sự được phép xuất hiện trên Search Engine.
•	Được phép Index
URL:
Public

Có Metadata

Có Content

Có Schema

Có Canonical

Có Internal Link
•	Bắt buộc NoIndex
Private

Restricted

Draft

Search URL

Technical URL

Share URL

Admin URL

Temporary URL
•	Nguyên tắc
Chất lượng trước số lượng.
 
C.3.7 C.3.7 Internal Linking Rule
•	Vai trò
Quy định nguyên tắc liên kết nội bộ.
•	Mục tiêu
Xây dựng mạng lưới dữ liệu liên kết xuyên suốt trong toàn hệ thống.
•	Thành phần bắt buộc
Mọi Landing Page SEO phải có:
Breadcrumb

Related Entity

Related Content
•	Không được tồn tại
Orphan Page
•	Nguyên tắc
Mọi Entity phải có khả năng liên kết tới:
Entity cha

Entity con

Entity liên quan

Nguồn dữ liệu liên quan
nếu tồn tại.
 
C.3.8 C.3.8 AI Search Rule
•	Vai trò
Quy định nguyên tắc tối ưu dữ liệu cho AI Search.
•	Mục tiêu
Giúp:
ChatGPT

Gemini

Copilot

Perplexity

Claude

Các hệ thống AI Search khác
có khả năng hiểu và trích dẫn dữ liệu QH Pro.
•	Thành phần bắt buộc
Đối với Landing Page SEO-A và SEO-B:
AI Summary

AI Snippet

Structured Data

Citation Source
•	Nguyên tắc
Thông tin quan trọng phải tồn tại dưới dạng:
HTML

Text

Schema

Metadata
Không được chỉ tồn tại dưới dạng:
Image

Map Tile

Canvas

Popup Runtime
 
C.3.9 Cache & Regeneration Rule
•	Vai trò
Quy định nguyên tắc cập nhật dữ liệu SEO.
•	Mục tiêu
Đảm bảo dữ liệu SEO luôn mới và phản ánh đúng trạng thái dữ liệu thực tế.
•	Nguyên tắc
Entity thay đổi
→ Regenerate phần bị ảnh hưởng
Không được:
Rebuild toàn bộ hệ thống
cho các thay đổi thông thường.
•	Ưu tiên
Event Driven

Incremental Regeneration

Queue Processing
 
C.3.10 Knowledge Graph Rule
•	Vai trò
Quy định nguyên tắc xây dựng mạng lưới dữ liệu ngữ nghĩa của QH Pro.
•	Mục tiêu
Biến QH Pro thành hệ thống dữ liệu có khả năng được Search Engine và AI hiểu như một Knowledge Graph.

•	Entity lõi
Administrative Unit

Parcel

Planning Region

Planning Project

Legal Document

Planning Map

Report

Snapshot
•	Nguyên tắc
Mọi Entity phải xác định được:
Nó là gì

Thuộc về đâu

Liên quan tới cái gì

Áp dụng cho đâu

Được tạo bởi ai

Được phê duyệt bởi ai
•	Relationship phải được thể hiện thông qua
Schema

Internal Link

AI Summary

Related Entity

Metadata
 
C.3.11 Acceptance Rule
•	Vai trò
Là tiêu chuẩn nghiệm thu SEO cấp hệ thống.
•	Một Landing Page SEO được coi là đạt yêu cầu khi
✓ Có URL chuẩn

✓ Có Metadata

✓ Có Canonical

✓ Có Schema

✓ Có Sitemap

✓ Có Index Rule phù hợp

✓ Có Internal Link

✓ Có AI Summary

✓ Có Related Entity

✓ Có Related Content

✓ Có Knowledge Graph Relationship
•	Một Landing Page SEO không đạt yêu cầu khi
✗ Không có Canonical

✗ Không có Metadata

✗ Không có Schema

✗ Không có Internal Link

✗ Không có AI Summary

✗ Không có Entity Relationship

✗ Không có Sitemap

✗ Bị Duplicate Content
 
•	Ghi chú triển khai
Phần C.3 là thư viện tiêu chuẩn SEO dùng chung của toàn bộ hệ thống.
Mọi yêu cầu SEO cụ thể cho từng màn hình, từng URL, từng Entity và từng luồng nghiệp vụ phải được mô tả tại:
PHẦN D. NỘI DUNG SEO BỔ SUNG CHO TỪNG MÀN HÌNH SRS
Dev, Design và QA triển khai trực tiếp theo Phần D; Phần C.3 chỉ đóng vai trò tiêu chuẩn tham chiếu để bảo đảm tính nhất quán trên toàn bộ hệ thống QH Pro.

PHẦN D.  NỘI DUNG SEO BỔ SUNG CHO TỪNG MÀN HÌNH SRS
Đây là phần Dev sử dụng trực tiếp.
Mỗi mục dưới đây sẽ được soạn đầy đủ theo bộ khung tại Phần C và sau đó copy ngược vào SRS gốc.
D.0	Nguyên tắc sử dụng Phần D
•	Vai trò
Phần D là phần triển khai SEO theo từng màn hình, từng module và từng đối tượng dữ liệu trong bộ SRS QH Pro.
Đây là phần tài liệu được sử dụng trực tiếp bởi:
•	Product Owner. 
•	UI/UX Designer. 
•	Frontend Developer. 
•	Backend Developer. 
•	QA/QC. 
•	SEO Team. 
Mục tiêu của Phần D là chuyển hóa các nguyên tắc SEO cấp hệ thống thành các yêu cầu triển khai cụ thể cho từng màn hình của QH Pro.
 
•	Quan hệ giữa Phần C và Phần D
Phần D không lặp lại toàn bộ các quy tắc SEO dùng chung đã được quy định tại Phần C.
Cụ thể:
•	C.2 – SEO Layout Template
Quy định:
•	Khung Layout chuẩn. 
•	Thành phần giao diện SEO chuẩn. 
•	Bố cục nội dung SEO chuẩn. 
•	Responsive Rule. 
•	Mobile Rule. 
•	Desktop Rule. 
Mỗi màn hình tại Phần D phải lựa chọn và áp dụng một hoặc nhiều Layout Template phù hợp từ C.2.
 
•	C.3 – SEO Reference Library
Quy định:
•	URL Rule. 
•	Metadata Rule. 
•	Canonical Rule. 
•	Schema Rule. 
•	Sitemap Rule. 
•	Index Rule. 
•	Internal Linking Rule. 
•	AI Search Rule. 
•	Cache & Regeneration Rule. 
•	Knowledge Graph Rule. 
•	Acceptance Rule. 
Các quy tắc này được coi là tiêu chuẩn mặc định của toàn hệ thống.
Mọi màn hình trong Phần D đều mặc định áp dụng các quy tắc tại C.3, trừ khi có quy định khác được nêu rõ trong từng mục D.x.x.
 
•	Nguyên tắc mô tả tại Phần D
Mỗi mục D.x.x chỉ mô tả:
•	Các yêu cầu SEO đặc thù của màn hình
Bao gồm:
•	SEO Scope. 
•	SEO Class. 
•	Landing Page hay không. 
•	Entity SEO sử dụng. 
•	Layout Template áp dụng. 
•	URL Pattern đặc thù. 
•	Metadata đặc thù. 
•	Schema đặc thù. 
•	AI Summary đặc thù. 
•	Internal Link đặc thù. 
•	Sitemap đặc thù. 
•	Index / NoIndex đặc thù. 
•	Acceptance Criteria đặc thù. 
 
•	Không lặp lại các quy tắc dùng chung
Không mô tả lại:
•	URL Rule chung. 
•	Metadata Rule chung. 
•	Canonical Rule chung. 
•	Schema Rule chung. 
•	Sitemap Rule chung. 
•	Internal Linking Rule chung. 
•	AI Search Rule chung. 
•	Cache Rule chung. 
Các nội dung này đã được quy định tại C.3.
 
•	Thứ tự đọc tài liệu khi triển khai
Đối với Design:
C.2 SEO Layout Template
↓
D.x.x tương ứng
 
Đối với Dev:
C.3 SEO Reference Library
↓
D.x.x tương ứng
 
Đối với QA:
C.3 SEO Reference Library
↓
D.x.x tương ứng
↓
Phần E Checklist & Nghiệm thu
 
•	Nguyên tắc ưu tiên
Trong trường hợp có khác biệt giữa:
C.2 SEO Layout Template
C.3 SEO Reference Library
và:
D.x.x
thì ưu tiên áp dụng:
D.x.x
vì đây là phần mô tả SEO cụ thể cho từng màn hình nghiệp vụ.
 
•	Mục tiêu cuối cùng của Phần D
Bảo đảm mỗi màn hình thuộc phạm vi SEO của QH Pro đều được mô tả đầy đủ và nhất quán về:
•	Giao diện SEO. 
•	Nội dung SEO. 
•	URL SEO. 
•	Metadata. 
•	Structured Data. 
•	AI Search. 
•	Internal Linking. 
•	Indexing. 
•	Knowledge Graph. 
Từ đó giúp Design, Dev và QA có thể triển khai, kiểm thử và vận hành SEO một cách thống nhất trên toàn bộ hệ thống QH Pro.


 
D.1	QH.1 – Tra cứu Quy hoạch / Map View Core
D.1.1 QH.1.1 Trang bản đồ quy hoạch
16. SEO
16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes
SEO Class	SEO-B
Entity chính	Planning Region
Entity phụ	Administrative Unit, GIS Layer, Planning Project
URL Type	Landing URL
Mục tiêu SEO	Trang điều hướng và điểm vào chính của hệ thống tra cứu quy hoạch
•	Ghi chú
•	Đây là Landing Page SEO của module tra cứu quy hoạch. 
•	Không phải URL SEO chính của Parcel hoặc Planning Project. 
•	Chức năng chính là điều hướng đến các Entity SEO-A. 
 
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach
•	URL phụ được phép sử dụng
/quy-hoach/{province-slug}
/quy-hoach/{province-slug}/{district-slug}
•	Canonical
•	Trang bản đồ toàn quốc: 
/quy-hoach
•	Trang địa phương: 
/quy-hoach/{administrative-unit-slug}
•	Không tạo URL SEO riêng cho
•	Zoom Level. 
•	Layer Toggle. 
•	Base Map. 
•	Theme. 
•	Sort. 
•	Filter tạm thời. 
•	Trạng thái panel. 
•	Trạng thái popup. 
•	Vị trí camera. 
•	Toạ độ người dùng. 
•	Duplicate Control
Mọi URL chứa:
?zoom=
?lat=
?lng=
?layer=
?theme=
phải canonical về URL chuẩn tương ứng.
 
•	16.3 Metadata
•	Title
Trang toàn quốc:
Bản đồ quy hoạch toàn quốc | QH Pro
Trang địa phương:
Bản đồ quy hoạch {administrative_unit_name} | QH Pro
•	Description
Tra cứu bản đồ quy hoạch, quy hoạch sử dụng đất, quy hoạch xây dựng, quy hoạch giao thông và các lớp dữ liệu quy hoạch liên quan tại {administrative_unit_name}.
•	H1
Trang toàn quốc:
Bản đồ quy hoạch
Trang địa phương:
Bản đồ quy hoạch {administrative_unit_name}
•	Robots
index,follow
•	Trường hợp NoIndex
noindex,nofollow
khi:
•	Không xác định được địa bàn. 
•	Không có dữ liệu công khai. 
•	Trang lỗi. 
•	Trạng thái Restricted. 
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Place
Entity Mapping	Administrative Unit
•	Schema bổ sung
Schema	Điều kiện
Place	Luôn áp dụng
Dataset	Có dữ liệu GIS công khai
CollectionPage	Trang danh mục địa bàn
•	Không inject Schema khi
•	Draft. 
•	Restricted. 
•	Deleted. 
•	Không có dữ liệu địa bàn hợp lệ. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Map Preview của địa bàn
OG Type	website
•	Share Preview
Ưu tiên hiển thị:
•	Tên địa bàn. 
•	Ảnh bản đồ. 
•	Số lượng lớp dữ liệu công khai. 
 
•	16.6 AI Summary
•	Có AI Summary
Yes
•	Dữ liệu nguồn
•	Administrative Unit. 
•	Planning Region. 
•	Planning Project. 
•	GIS Layer công khai. 
•	Nội dung AI Summary
•	Mô tả địa bàn. 
•	Tình trạng dữ liệu quy hoạch hiện có. 
•	Các nhóm quy hoạch đang áp dụng. 
•	Các đồ án quy hoạch nổi bật. 
•	Không public
•	Layer nội bộ. 
•	Layer Restricted. 
•	Metadata kỹ thuật. 
•	Dữ liệu Partner Only. 
•	Dữ liệu chưa công bố. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	planning-map.xml
•	Tham gia Sitemap khi
•	Public. 
•	Active. 
•	Có dữ liệu quy hoạch. 
•	Có Administrative Unit hợp lệ. 
•	Loại khỏi Sitemap khi
•	Draft. 
•	Restricted. 
•	Deleted. 
•	Empty Dataset. 
 
•	16.8 Internal Link
•	Link đến
•	Trang đơn vị hành chính. 
•	Trang chi tiết thửa đất. 
•	Trang chi tiết khu quy hoạch. 
•	Trang chi tiết đồ án quy hoạch. 
•	Trang thư viện quy hoạch. 
•	Trang văn bản pháp lý. 
•	Link nhận từ
•	Trang chủ. 
•	Tra cứu địa chỉ. 
•	Tra cứu GPS. 
•	Tra cứu tờ/thửa. 
•	Chi tiết đồ án. 
•	Báo cáo. 
•	Snapshot. 
 
•	16.9 Security & Visibility
•	Public
•	Metadata. 
•	AI Summary. 
•	Preview bản đồ. 
•	Administrative Unit. 
•	Danh sách lớp dữ liệu công khai. 
•	Restricted
•	Layer trả phí. 
•	Layer đối tác. 
•	Layer chưa công bố. 
•	Dữ liệu nội bộ. 
•	Private
•	Thông tin quản trị. 
•	Log hệ thống. 
•	Cấu hình kỹ thuật. 
•	Permission Rule
Không được hiển thị trong:
•	Metadata. 
•	Open Graph. 
•	Structured Data. 
•	AI Summary. 
đối với dữ liệu Restricted hoặc Private.
 
•	16.10 Cache & Regeneration
Event	Regenerate
administrative_unit.updated	Metadata
planning_region.updated	Metadata
planning_project.updated	AI Summary
gis_layer.updated	AI Summary, Open Graph
visibility.changed	Sitemap, Robots
map_preview.updated	Open Graph
legal_status.changed	Metadata
 
•	16.11 Acceptance Criteria
•	AC-D1.1-001
Trang bản đồ quy hoạch có URL chuẩn duy nhất.
•	AC-D1.1-002
Tất cả URL chứa trạng thái bản đồ động đều canonical về URL chuẩn.
•	AC-D1.1-003
Metadata được sinh đầy đủ theo địa bàn.
•	AC-D1.1-004
Structured Data hợp lệ.
•	AC-D1.1-005
Open Graph hiển thị đúng Preview bản đồ.
•	AC-D1.1-006
AI Summary được sinh từ dữ liệu công khai.
•	AC-D1.1-007
Trang đủ điều kiện xuất hiện trong Sitemap.
•	AC-D1.1-008
Không lộ dữ liệu Restricted hoặc Private.
 
•	16.12 QA Checklist
•	URL
□ URL đúng chuẩn
□ URL duy nhất
□ Canonical đúng
□ Không index URL chứa map state
 
•	Metadata
□ Title tồn tại
□ Description tồn tại
□ H1 tồn tại
□ Robots đúng
 
•	Structured Data
□ Schema hợp lệ
□ Entity Mapping đúng
□ Không inject khi Restricted
 
•	Open Graph
□ OG Title đúng
□ OG Description đúng
□ OG Image đúng
□ Share Preview hiển thị đúng
 
•	AI Summary
□ Có AI Summary
□ Nội dung đúng địa bàn
□ Không lộ dữ liệu Restricted
 
•	Sitemap
□ Có trong planning-map.xml
□ Loại khỏi Sitemap đúng điều kiện
 
•	Security
□ Không lộ Layer Restricted
□ Không lộ Layer Private
□ Không lộ Metadata nội bộ
□ Permission hoạt động đúng

D.1.2 QH.1.2 Thanh tìm kiếm & định vị trên bản đồ
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	No
SEO Class	SEO-N
Entity chính	Không có
Entity phụ	Administrative Unit, Parcel, Planning Region
URL Type	Internal URL
Mục tiêu SEO	Điều hướng người dùng đến các Entity SEO khác
•	Ghi chú
•	Đây là thành phần hỗ trợ tìm kiếm và điều hướng. 
•	Không phải đối tượng SEO độc lập. 
•	Không tạo Landing Page SEO. 
•	Không tham gia Sitemap. 
•	SEO được kế thừa từ Entity đích mà người dùng được chuyển tới. 
 
•	16.2 URL & Canonical
•	URL chuẩn
Không áp dụng.
•	URL phụ
Không áp dụng.
•	URL NoIndex
Toàn bộ URL phát sinh từ trạng thái tìm kiếm tạm thời.
Ví dụ:
?search=
?q=
?suggest=
?lat=
?lng=
?keyword=
•	Canonical
Canonical về URL của màn hình cha:
/quy-hoach
hoặc URL Entity đích sau khi người dùng lựa chọn kết quả.
•	Không tạo URL SEO riêng cho
•	Keyword tìm kiếm. 
•	Từ khóa gợi ý. 
•	Tọa độ định vị. 
•	Kết quả autocomplete. 
•	Trạng thái tìm kiếm tạm thời. 
 
•	16.3 Metadata
•	Title
Không sinh Metadata riêng.

•	Description
Không sinh Metadata riêng.
•	H1
Kế thừa từ màn hình cha.
•	Robots
noindex,nofollow
đối với mọi trạng thái tìm kiếm nội bộ.
•	Ghi chú
Thanh tìm kiếm không được phép sinh Title hoặc Description độc lập.
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Không áp dụng
Entity Mapping	Không áp dụng
•	Điều kiện áp dụng
Không áp dụng.
•	Điều kiện loại bỏ
Luôn loại bỏ.
•	Ghi chú
Structured Data chỉ được sinh tại Entity đích.
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Không áp dụng
OG Description	Không áp dụng
OG Image	Không áp dụng
•	Share Preview
Không hỗ trợ.
•	Ghi chú
Thanh tìm kiếm không phải đối tượng được chia sẻ.
 
•	16.6 AI Summary
•	Có AI Summary
No
•	Dữ liệu nguồn
Không áp dụng.
•	Nội dung public
Không áp dụng.
•	Nội dung không public
Toàn bộ:
•	Keyword tìm kiếm. 
•	Lịch sử tìm kiếm. 
•	Gợi ý cá nhân hóa. 
•	Dữ liệu vị trí người dùng. 
•	Ghi chú
AI Summary chỉ tồn tại ở Entity đích.
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	Không áp dụng
•	Tham gia Sitemap khi
Không có trường hợp nào.
•	Loại khỏi Sitemap khi
Luôn loại khỏi Sitemap.
 
•	16.8 Internal Link
•	Link đến
•	Administrative Unit. 
•	Parcel. 
•	Planning Region. 
•	Planning Project. 
•	Legal Document. 
•	Link nhận từ
•	QH.1.1 Trang bản đồ quy hoạch. 
•	Menu điều hướng. 
•	Trang chủ. 
•	Ghi chú
Thanh tìm kiếm đóng vai trò Entry Point nội bộ.
 
•	16.9 Security & Visibility
•	Public
•	Danh sách gợi ý công khai. 
•	Administrative Unit công khai. 
•	Restricted
•	Dữ liệu vị trí chính xác của người dùng. 
•	Gợi ý cá nhân hóa. 
•	Private
•	Lịch sử tìm kiếm. 
•	Hành vi người dùng. 
•	Nhật ký định vị. 
•	Permission Rule
Không được đưa các dữ liệu Restricted hoặc Private vào:
•	Metadata. 
•	Open Graph. 
•	AI Summary. 
•	Structured Data. 
 
•	16.10 Cache & Regeneration
Event	Regenerate
administrative_unit.updated	Search Index
parcel.updated	Search Index
planning_region.updated	Search Index
planning_project.updated	Search Index
•	Ghi chú
Không phát sinh Metadata Regeneration.
Không phát sinh Sitemap Regeneration.
 
•	16.11 Acceptance Criteria
•	AC-D1.2-001
Không sinh URL SEO cho từ khóa tìm kiếm.
•	AC-D1.2-002
Không sinh Metadata độc lập.
•	AC-D1.2-003
Không sinh Structured Data.
•	AC-D1.2-004
Không tham gia Sitemap.
•	AC-D1.2-005
Kết quả tìm kiếm điều hướng đúng đến Entity SEO đích.
•	AC-D1.2-006
Không lộ dữ liệu vị trí hoặc lịch sử tìm kiếm của người dùng.
 
•	16.12 QA Checklist
•	URL
□ Không sinh URL SEO từ keyword
□ Không sinh URL SEO từ autocomplete
□ Canonical đúng về URL cha hoặc Entity đích
 
•	Metadata
□ Không sinh Metadata riêng
□ Không sinh H1 riêng
 
•	Structured Data
□ Không inject Schema
□ Không inject Entity Mapping
 
•	Open Graph
□ Không sinh OG riêng
□ Không sinh Share Preview riêng
 
•	Sitemap
□ Không xuất hiện trong Sitemap
□ Không tạo URL Index mới
 
•	Security
□ Không lộ lịch sử tìm kiếm
□ Không lộ vị trí người dùng
□ Không lộ dữ liệu Restricted
□ Không lộ dữ liệu Private
D.1.3 QH.1.3 Panel lớp dữ liệu
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	No
SEO Class	SEO-N
Entity chính	GIS Layer
Entity phụ	Planning Region, Planning Project, Administrative Unit
URL Type	Internal URL
Mục tiêu SEO	Hỗ trợ người dùng lựa chọn và hiển thị dữ liệu quy hoạch trên bản đồ
•	Ghi chú
•	Panel lớp dữ liệu là thành phần giao diện hỗ trợ tra cứu. 
•	Không phải Landing Page. 
•	Không phải Entity SEO độc lập. 
•	Không tham gia SEO trực tiếp. 
•	Giá trị SEO được chuyển giao cho các trang chi tiết Layer, Khu quy hoạch hoặc Đồ án quy hoạch nếu tồn tại URL SEO tương ứng. 
 
•	16.2 URL & Canonical
•	URL chuẩn
Không áp dụng.
•	URL phụ
Không áp dụng.
•	URL NoIndex
Toàn bộ trạng thái phát sinh từ:
?layer=
?layers=
?group=
?category=
?overlay=
?showLayer=
?hideLayer=
?legend=
•	Canonical
Canonical về:
/quy-hoach
hoặc URL Entity đích tương ứng.
•	Không tạo URL SEO riêng cho
•	Layer Toggle. 
•	Layer Group. 
•	Layer Category. 
•	Layer Visibility. 
•	Layer Order. 
•	Layer Opacity. 
•	Legend State. 
•	Duplicate Control
Mọi trạng thái thay đổi lớp dữ liệu không được tạo URL SEO độc lập.
 
•	16.3 Metadata
•	Title
Không sinh Metadata riêng.
•	Description
Không sinh Metadata riêng.
•	H1
Kế thừa từ màn hình cha.
•	Robots
noindex,nofollow
•	Ghi chú
Panel lớp dữ liệu không phải đối tượng SEO độc lập nên không được sinh Metadata riêng.
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Không áp dụng
Entity Mapping	Không áp dụng
•	Điều kiện áp dụng
Không áp dụng.
•	Điều kiện loại bỏ
Luôn loại bỏ.
•	Ghi chú
Schema chỉ được áp dụng tại:
•	Trang chi tiết Layer. 
•	Trang chi tiết Đồ án. 
•	Trang chi tiết Khu quy hoạch. 
•	Trang chi tiết Bản đồ. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Không áp dụng
OG Description	Không áp dụng
OG Image	Không áp dụng
•	Share Preview
Không hỗ trợ.
•	Ghi chú
Không cho phép chia sẻ trạng thái Panel lớp dữ liệu như một đối tượng SEO.
 
•	16.6 AI Summary
•	Có AI Summary
No
•	Dữ liệu nguồn
Không áp dụng.
•	Nội dung public
Không áp dụng.
•	Nội dung không public
•	Danh sách Layer Restricted. 
•	Layer nội bộ. 
•	Layer trả phí. 
•	Layer đối tác. 
•	Metadata kỹ thuật. 
•	Quy tắc hiển thị hệ thống. 
•	Ghi chú
AI Summary chỉ tồn tại ở Entity đích, không tồn tại ở Panel lớp dữ liệu.
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	Không áp dụng
•	Tham gia Sitemap khi
Không có trường hợp nào.


•	Loại khỏi Sitemap khi
Luôn loại khỏi Sitemap.
•	Ghi chú
Panel lớp dữ liệu không được xuất hiện trong bất kỳ Sitemap nào.
 
•	16.8 Internal Link
•	Link đến
•	Chi tiết Layer GIS. 
•	Chi tiết Khu quy hoạch. 
•	Chi tiết Đồ án quy hoạch. 
•	Chi tiết Bản đồ quy hoạch. 
•	Chi tiết Đơn vị hành chính. 
•	Link nhận từ
•	Trang bản đồ quy hoạch. 
•	Kết quả tìm kiếm. 
•	Chi tiết Đồ án. 
•	Chi tiết Bản đồ. 
•	Ghi chú
Panel lớp dữ liệu chỉ đóng vai trò điều hướng nội bộ.
 
•	16.9 Security & Visibility
•	Public
•	Tên Layer công khai. 
•	Nhóm Layer công khai. 
•	Chú giải công khai. 
•	Restricted
•	Layer trả phí. 
•	Layer đối tác. 
•	Layer chưa công bố. 
•	Layer đang kiểm duyệt. 
•	Private
•	Layer nội bộ. 
•	Layer kỹ thuật. 
•	Metadata vận hành. 
•	Cấu hình hiển thị. 
•	Permission Rule
Không được đưa dữ liệu Restricted hoặc Private vào:
•	Metadata. 
•	Open Graph. 
•	AI Summary. 
•	Structured Data. 
•	Sitemap. 
 
•	16.10 Cache & Regeneration
Event	Regenerate
layer.created	Layer Index
layer.updated	Layer Index
layer.deleted	Layer Index
layer.visibility.changed	Permission Cache
layer.category.changed	Search Cache
•	Ghi chú
Không phát sinh:
•	Metadata Regeneration. 
•	Schema Regeneration. 
•	Sitemap Regeneration. 
đối với Panel lớp dữ liệu.
 
•	16.11 Acceptance Criteria
•	AC-D1.3-001
Không sinh URL SEO cho trạng thái Layer.
•	AC-D1.3-002
Không sinh Metadata riêng.
•	AC-D1.3-003
Không sinh Structured Data.
•	AC-D1.3-004
Không tham gia Sitemap.
•	AC-D1.3-005
Không sinh Open Graph.
•	AC-D1.3-006
Không sinh AI Summary.
•	AC-D1.3-007
Không lộ Layer Restricted hoặc Layer Private.
•	AC-D1.3-008
Các liên kết điều hướng từ Panel dẫn đúng đến Entity SEO tương ứng.
 
•	16.12 QA Checklist
•	URL
□ Không sinh URL SEO cho Layer Toggle
□ Không sinh URL SEO cho Layer State
□ Canonical đúng về URL cha
 
•	Metadata
□ Không sinh Metadata riêng
□ Không sinh H1 riêng
□ Robots đúng
 
•	Structured Data
□ Không inject Schema
□ Không inject Entity Mapping
 
•	Open Graph
□ Không sinh OG Title
□ Không sinh OG Description
□ Không sinh OG Image
 
•	AI Summary
□ Không sinh AI Summary
□ Không công khai Layer Restricted
 
•	Sitemap
□ Không xuất hiện trong Sitemap
□ Không tạo URL Index mới
 
•	Security
□ Không lộ Layer trả phí
□ Không lộ Layer đối tác
□ Không lộ Layer nội bộ
□ Permission hoạt động đúng
 
•	Navigation
□ Điều hướng đúng tới Đồ án quy hoạch
□ Điều hướng đúng tới Bản đồ quy hoạch
□ Điều hướng đúng tới Khu quy hoạch
□ Điều hướng đúng tới Layer công khai tương ứng
•	Ghi chú triển khai
QH.1.3 là màn hình SEO-N, tương tự QH.1.2. Dev không cần triển khai các thành phần SEO (Metadata, Schema, Sitemap, AI Summary) cho màn hình này. Trọng tâm chỉ là:
•	Không tạo URL SEO sai. 
•	Không tạo nội dung trùng lặp. 
•	Không làm lộ Layer bị hạn chế. 
•	Điều hướng đúng tới các Entity SEO-A của hệ thống.

D.1.4 QH.1.4 Popup thông tin nhanh
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	No
SEO Class	SEO-N
Entity chính	Parcel / Planning Region / Planning Project
Entity phụ	Administrative Unit
URL Type	Internal UI Component
Mục tiêu SEO	Điều hướng người dùng đến trang chi tiết Entity SEO
•	Ghi chú
•	Popup thông tin nhanh chỉ là thành phần giao diện. 
•	Không phải Landing Page. 
•	Không phải Entity SEO độc lập. 
•	Không được Index. 
•	Không tham gia Sitemap. 
•	SEO được chuyển tiếp tới màn hình chi tiết tương ứng. 
 
•	16.2 URL & Canonical
•	URL chuẩn
Không áp dụng.
•	URL phụ
Không áp dụng.
•	URL NoIndex
Toàn bộ trạng thái Popup.
Ví dụ:
?popup=parcel
?popup=planning-region
?popup=project
?selected=
?feature=
•	Canonical
Canonical về:
•	Trang bản đồ quy hoạch. 
•	Hoặc URL chi tiết Entity tương ứng. 
•	Không tạo URL SEO riêng cho
•	Popup mở. 
•	Popup đóng. 
•	Feature được chọn. 
•	Highlight trạng thái. 
•	Hover trạng thái. 
•	Duplicate Control
Không cho phép Popup sinh URL SEO riêng.
 
•	16.3 Metadata
•	Title
Không sinh Metadata riêng.
•	Description
Không sinh Metadata riêng.
•	H1
Không sinh H1 riêng.
•	Robots
noindex,nofollow
•	Ghi chú
Popup không phải trang SEO.
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Không áp dụng
Entity Mapping	Không áp dụng
•	Điều kiện áp dụng
Không áp dụng.
•	Điều kiện loại bỏ
Luôn loại bỏ.
•	Ghi chú
Schema chỉ được sinh tại trang chi tiết Entity.
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Không áp dụng
OG Description	Không áp dụng
OG Image	Không áp dụng
•	Share Preview
Không hỗ trợ.
•	Ghi chú
Popup không được chia sẻ độc lập.
 
•	16.6 AI Summary
•	Có AI Summary
No
•	Dữ liệu nguồn
Không áp dụng.
•	Nội dung public
Không áp dụng.
•	Nội dung không public
•	Nội dung Popup. 
•	Dữ liệu tạm thời. 
•	Dữ liệu Restricted. 
•	Ghi chú
AI Summary chỉ xuất hiện tại màn hình chi tiết.
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	Không áp dụng
•	Tham gia Sitemap khi
Không có.
•	Loại khỏi Sitemap khi
Luôn loại bỏ.
 
•	16.8 Internal Link
•	Link đến
•	Chi tiết thửa đất. 
•	Chi tiết khu quy hoạch. 
•	Chi tiết đồ án quy hoạch. 
•	Chi tiết văn bản pháp lý. 
•	Chi tiết bản đồ quy hoạch. 
•	Link nhận từ
•	Trang bản đồ quy hoạch. 
•	Kết quả tìm kiếm. 
•	Layer GIS. 
•	Ghi chú
Popup đóng vai trò cầu nối điều hướng.
 
•	16.9 Security & Visibility
•	Public
•	Thông tin tóm tắt công khai. 
•	Tên đối tượng. 
•	Trạng thái công khai. 
•	Restricted
•	Dữ liệu trả phí. 
•	Dữ liệu đối tác. 
•	Dữ liệu hạn chế. 
•	Private
•	Ghi chú nội bộ. 
•	Thông tin quản trị. 
•	Permission Rule
Không được đưa dữ liệu Restricted hoặc Private vào Popup công khai.
 
•	16.10 Cache & Regeneration
Event	Regenerate
parcel.updated	Popup Cache
planning_region.updated	Popup Cache
planning_project.updated	Popup Cache
visibility.changed	Popup Cache
•	Ghi chú
Không phát sinh Metadata, Schema hoặc Sitemap Regeneration.
 
•	16.11 Acceptance Criteria
•	AC-D1.4-001
Không sinh URL SEO cho Popup.
•	AC-D1.4-002
Không sinh Metadata riêng.
•	AC-D1.4-003
Không sinh Structured Data.
•	AC-D1.4-004
Không tham gia Sitemap.
•	AC-D1.4-005
Điều hướng đúng tới Entity SEO đích.
•	AC-D1.4-006
Không hiển thị dữ liệu Restricted hoặc Private.
 
•	16.12 QA Checklist
•	URL
□ Không sinh URL SEO
□ Canonical đúng về URL cha hoặc Entity đích
•	Metadata
□ Không sinh Metadata riêng
□ Không sinh H1 riêng
•	Structured Data
□ Không inject Schema
□ Không inject Entity Mapping
•	Open Graph
□ Không sinh OG
•	Sitemap
□ Không xuất hiện trong Sitemap
•	Security
□ Không lộ dữ liệu Restricted
□ Không lộ dữ liệu Private
•	Navigation
□ Điều hướng đúng tới trang chi tiết
 
 
D.1.5 QH.1.5 Chi tiết thửa đất / vùng quy hoạch / đồ án
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes
SEO Class	SEO-A
Entity chính	Parcel / Planning Region / Planning Project
Entity phụ	Administrative Unit, Legal Document, GIS Layer
URL Type	Detail URL
Mục tiêu SEO	URL SEO trọng tâm của QH Pro, phục vụ Google Search, AI Search, Knowledge Graph và Internal Linking
•	Ghi chú
Đây là nhóm URL SEO quan trọng nhất của toàn bộ hệ thống.
Mọi đối tượng:
•	Thửa đất. 
•	Khu quy hoạch. 
•	Đồ án quy hoạch. 
đều phải có URL SEO chuẩn riêng.
 
•	16.2 URL & Canonical
•	URL chuẩn
•	Thửa đất
/quy-hoach/thua-dat/{parcel-slug}
•	Vùng quy hoạch
/quy-hoach/khu-vuc/{planning-region-slug}
•	Đồ án quy hoạch
/quy-hoach/do-an/{planning-project-slug}
•	URL phụ
Cho phép URL lịch sử phục vụ redirect.
•	Canonical
Luôn trỏ về URL chuẩn của Entity.
•	Không tạo URL SEO riêng cho
•	Tab. 
•	Layer. 

•	Bộ lọc. 
•	Chế độ xem. 
•	Bản đồ nền. 
•	Trạng thái mở rộng. 
•	Duplicate Control
Mỗi Entity chỉ có một URL SEO chuẩn.
 
•	16.3 Metadata
•	Title
•	Parcel
{parcel_name} | Thông tin quy hoạch và địa chính | QH Pro
•	Planning Region
{planning_region_name} | Thông tin quy hoạch | QH Pro
•	Planning Project
{planning_project_name} | Đồ án quy hoạch | QH Pro
•	Description
Thông tin quy hoạch, vị trí, hiện trạng, pháp lý, dữ liệu liên quan và các đồ án áp dụng cho {entity_name}.
•	H1
{entity_name}
•	Robots
index,follow
•	NoIndex khi
•	Draft. 
•	Deleted. 
•	Restricted. 
•	Không có dữ liệu công khai. 
•	Geometry không hợp lệ. 
 
•	16.4 Structured Data
Entity	Schema
Parcel	Place
Planning Region	Place
Planning Project	CreativeWork
Legal Document	Legislation
Administrative Unit	AdministrativeArea
•	Không inject Schema khi
•	Draft. 
•	Restricted. 
•	Deleted. 
•	Missing Geometry. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Map Preview
OG Type	article
•	Share Preview
Hiển thị:
•	Tên đối tượng. 
•	Địa bàn. 
•	Bản đồ Preview. 
•	Trạng thái quy hoạch. 
 
•	16.6 AI Summary
•	Có AI Summary
Yes
•	Dữ liệu nguồn
•	Parcel. 
•	Planning Region. 
•	Planning Project. 
•	Administrative Unit. 
•	Legal Document. 
•	GIS Layer. 
•	Nội dung AI Summary
•	Tóm tắt quy hoạch. 
•	Tóm tắt pháp lý. 
•	Tóm tắt vị trí. 
•	Quan hệ với các đồ án khác. 
•	Các cảnh báo quan trọng. 
•	Không public
•	Internal Note. 
•	Restricted Data. 
•	Private Data. 
•	Partner Data. 
•	Dữ liệu chưa công bố. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Parcel	parcel.xml
Planning Region	planning-region.xml
Planning Project	planning-project.xml
•	Tham gia khi
•	Public. 
•	Active. 
•	Có dữ liệu hợp lệ. 
•	Có URL chuẩn. 
•	Loại khỏi Sitemap khi
•	Draft. 
•	Deleted. 
•	Restricted. 
•	Inactive. 
 
•	16.8 Internal Link
•	Link đến
•	Administrative Unit. 
•	Parcel liên quan. 
•	Planning Region liên quan. 
•	Planning Project liên quan. 
•	Legal Document liên quan. 
•	Bản đồ quy hoạch liên quan. 
•	Báo cáo liên quan. 
•	Link nhận từ
•	Trang bản đồ quy hoạch. 
•	Tra cứu địa chỉ. 
•	Tra cứu GPS. 
•	Tra cứu tờ/thửa. 
•	Thư viện quy hoạch. 
•	Báo cáo. 
•	Snapshot. 
•	AI Search Result. 
•	Ghi chú
Đây là Node trung tâm của Knowledge Graph.
 
•	16.9 Security & Visibility
•	Public
•	Metadata. 
•	AI Summary. 
•	Bản đồ Preview. 
•	Thông tin quy hoạch công khai. 
•	Văn bản công khai. 
•	Restricted
•	Dữ liệu trả phí. 
•	Dữ liệu đối tác. 
•	Dữ liệu hạn chế. 
•	Private
•	Ghi chú nội bộ. 
•	Nhật ký xử lý. 
•	Metadata vận hành. 
•	Permission Rule
Không được đưa dữ liệu Restricted hoặc Private vào:
•	Metadata. 
•	Open Graph. 
•	Schema. 
•	AI Summary. 
•	Sitemap. 
 
•	16.10 Cache & Regeneration
Event	Regenerate
parcel.updated	Metadata, Schema
planning_region.updated	Metadata, Schema
planning_project.updated	Metadata, AI Summary
legal_document.updated	AI Summary
geometry.updated	Schema, Open Graph
visibility.changed	Sitemap, Robots
administrative_unit.updated	Metadata
 
•	16.11 Acceptance Criteria
•	AC-D1.5-001
Mỗi Entity có đúng một URL chuẩn.
•	AC-D1.5-002
Canonical luôn trỏ về URL chuẩn.
•	AC-D1.5-003
Metadata được sinh đầy đủ.
•	AC-D1.5-004
Structured Data hợp lệ.
•	AC-D1.5-005
Open Graph hiển thị đúng Preview.
•	AC-D1.5-006
AI Summary được sinh đúng dữ liệu công khai.
•	AC-D1.5-007
Xuất hiện trong Sitemap đúng nhóm.
•	AC-D1.5-008
Không lộ dữ liệu Restricted hoặc Private.
•	AC-D1.5-009
Internal Link hoạt động đầy đủ.
 
•	16.12 QA Checklist
•	URL
□ URL đúng chuẩn
□ URL duy nhất
□ Canonical đúng
□ Không có URL trùng lặp
 
•	Metadata
□ Title tồn tại
□ Description tồn tại
□ H1 tồn tại
□ Robots đúng
 
•	Structured Data
□ Schema hợp lệ
□ Entity Mapping đúng
□ Không inject khi Restricted
 
•	Open Graph
□ OG Title đúng
□ OG Description đúng
□ OG Image đúng
□ Share Preview đúng
 
•	AI Summary
□ Có AI Summary
□ Nội dung đúng dữ liệu nguồn
□ Không lộ dữ liệu Restricted
 
•	Sitemap
□ Xuất hiện trong Sitemap đúng nhóm
□ Loại khỏi Sitemap đúng điều kiện
 
•	Security
□ Không lộ dữ liệu Restricted
□ Không lộ dữ liệu Private
□ Permission hoạt động đúng
 
•	Knowledge Graph
□ Có liên kết tới Administrative Unit
□ Có liên kết tới Planning Region
□ Có liên kết tới Planning Project
□ Có liên kết tới Legal Document
□ Có liên kết tới các Entity liên quan
•	Ghi chú triển khai
D.1.5 là một trong các URL SEO-A quan trọng nhất của QH Pro và phải được ưu tiên triển khai đầy đủ toàn bộ:
•	URL & Canonical. 
•	Metadata. 
•	Structured Data. 
•	Open Graph. 
•	AI Summary. 
•	Sitemap. 
•	Internal Linking. 
•	Knowledge Graph. 
•	Security Control.

 
D.2	QH.2 – Tra cứu & Định vị
D.2.1 QH.2.1 Tra cứu địa chỉ
16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes
SEO Class	SEO-B / SEO-C
Entity chính	Administrative Unit
Entity phụ	Parcel, Planning Region, Planning Project
URL Type	Search / Landing URL
Mục tiêu SEO	Hỗ trợ tra cứu theo địa danh, điều hướng tới các trang địa bàn và Entity SEO-A
•	Ghi chú
•	Tra cứu địa chỉ là điểm vào quan trọng của SEO địa phương. 
•	Không SEO mọi truy vấn tìm kiếm tự do. 
•	Chỉ SEO các địa danh hoặc khu vực có dữ liệu ổn định, công khai và đủ giá trị lập chỉ mục. 
 
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/{administrative-unit-slug}
•	URL tra cứu không Index
/quy-hoach/tim-kiem?q={keyword}
/quy-hoach/dia-chi?keyword={keyword}
•	Canonical
•	Nếu kết quả xác định được địa bàn: canonical về URL địa bàn. 
•	Nếu kết quả xác định được Parcel: canonical về URL chi tiết thửa đất. 
•	Nếu chỉ là query tạm thời: canonical về /quy-hoach. 
•	Không tạo URL SEO riêng cho
•	Từ khóa tìm kiếm tự do. 
•	Autocomplete. 
•	Suggestion. 
•	Search session. 
•	Kết quả không xác định được Entity. 
 
•	16.3 Metadata
•	Title
Quy hoạch {administrative_unit_name} | QH Pro
•	Description
Tra cứu thông tin quy hoạch, bản đồ quy hoạch, thửa đất, khu quy hoạch và các đồ án liên quan tại {administrative_unit_name}.
•	H1
Quy hoạch {administrative_unit_name}
•	Robots
index,follow
•	NoIndex khi
•	Không xác định được địa bàn hợp lệ. 
•	Kết quả chỉ là query tạm thời. 
•	Không có dữ liệu công khai. 
•	Địa bàn bị Restricted. 
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Place / AdministrativeArea
Entity Mapping	Administrative Unit
•	Schema bổ sung
Schema	Điều kiện
CollectionPage	Trang danh sách kết quả địa bàn
Dataset	Có dữ liệu quy hoạch công khai
•	Không inject Schema nếu
•	Địa bàn không hợp lệ. 
•	Dữ liệu Restricted. 
•	Không có dữ liệu công khai. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Map Preview của địa bàn
OG Type	website
•	Share Preview
Hiển thị:
•	Tên địa bàn. 
•	Bản đồ Preview. 
•	Nhóm dữ liệu quy hoạch hiện có. 
 
•	16.6 AI Summary
•	Có AI Summary
Yes, với địa bàn đủ điều kiện công khai.
•	Dữ liệu nguồn
•	Administrative Unit. 
•	Planning Region. 
•	Planning Project. 
•	GIS Layer công khai. 
•	Legal Document công khai. 
•	Nội dung AI Summary
•	Tóm tắt địa bàn. 
•	Các nhóm quy hoạch hiện có. 
•	Đồ án quy hoạch liên quan. 
•	Văn bản pháp lý nổi bật. 
•	Liên kết đến các Entity chính. 
•	Không public
•	Search history. 
•	Query cá nhân hóa. 
•	Dữ liệu Restricted. 
•	Layer nội bộ. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	administrative-unit.xml
•	Tham gia khi
•	Địa bàn Public. 

•	Có dữ liệu quy hoạch. 
•	Có URL chuẩn. 
•	Có Metadata hợp lệ. 
•	Loại bỏ khi
•	Restricted. 
•	Deleted. 
•	Không có dữ liệu công khai. 
•	Không xác định được Entity. 
 
•	16.8 Internal Link
•	Link đến
•	Trang đơn vị hành chính cấp cha. 
•	Trang đơn vị hành chính cấp con. 
•	Thửa đất liên quan. 
•	Khu quy hoạch liên quan. 
•	Đồ án quy hoạch liên quan. 
•	Văn bản pháp lý liên quan. 
•	Link nhận từ
•	Trang bản đồ quy hoạch. 
•	Thanh tìm kiếm. 
•	Chi tiết thửa đất. 
•	Chi tiết đồ án. 
•	Thư viện quy hoạch. 
 
•	16.9 Security & Visibility
•	Public
•	Tên địa bàn. 
•	Bản đồ Preview. 
•	Metadata địa bàn. 
•	Dữ liệu quy hoạch công khai. 
•	Restricted
•	Layer hạn chế. 
•	Dữ liệu đối tác. 
•	Dữ liệu trả phí. 
•	Private
•	Query cá nhân. 
•	Lịch sử tìm kiếm. 
•	Vị trí người dùng. 
•	Permission Rule
Không đưa dữ liệu Restricted / Private vào Metadata, Schema, OG, AI Summary hoặc Sitemap.
 
•	16.10 Cache & Regeneration
Event	Regenerate
administrative_unit.updated	Metadata, Sitemap
administrative_unit.merged	Canonical, Redirect, Metadata
administrative_unit.renamed	Metadata, Canonical
planning_region.updated	AI Summary
planning_project.updated	AI Summary
visibility.changed	Sitemap, Robots
map_preview.updated	Open Graph
 
•	16.11 Acceptance Criteria
•	AC-D2.1-001
Địa bàn hợp lệ có URL SEO chuẩn.
•	AC-D2.1-002
Query tìm kiếm tạm thời không được Index.
•	AC-D2.1-003
Canonical trỏ đúng về URL địa bàn hoặc Entity đích.
•	AC-D2.1-004
Metadata sinh đầy đủ theo địa bàn.
•	AC-D2.1-005
AI Summary chỉ dùng dữ liệu công khai.
•	AC-D2.1-006
Địa bàn đủ điều kiện xuất hiện trong Sitemap.
•	AC-D2.1-007
Hỗ trợ tên địa bàn cũ, trước sáp nhập và sau sáp nhập.
 
•	16.12 QA Checklist
•	URL
□ URL địa bàn đúng chuẩn
□ Query tạm thời NoIndex
□ Canonical đúng Entity đích
□ Legacy URL redirect/canonical đúng
•	Metadata
□ Title đúng địa bàn
□ Description đúng địa bàn
□ H1 đúng địa bàn
□ Robots đúng
•	Structured Data
□ Schema AdministrativeArea hợp lệ
□ Entity Mapping đúng
•	Open Graph
□ OG Title đúng
□ OG Image đúng Map Preview
•	AI Summary
□ Nội dung đúng địa bàn
□ Không lộ dữ liệu Restricted
•	Sitemap
□ Có trong administrative-unit.xml khi đủ điều kiện
□ Loại khỏi Sitemap đúng điều kiện
•	Security
□ Không lộ search history
□ Không lộ vị trí cá nhân
□ Không lộ dữ liệu Restricted
 
D.2.2 QH.2.2 Tra cứu GPS
16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	No
SEO Class	SEO-N
Entity chính	Không có
Entity phụ	Parcel, Administrative Unit, Planning Region
URL Type	Internal / Runtime URL
Mục tiêu SEO	Xác định vị trí người dùng và điều hướng đến Entity SEO phù hợp
•	Ghi chú
•	GPS là input kỹ thuật, không phải đối tượng SEO. 
•	Không Index URL chứa tọa độ người dùng. 
•	SEO chỉ phát sinh tại Entity đích sau khi hệ thống xác định vị trí. 
 
•	16.2 URL & Canonical
•	URL chuẩn
Không áp dụng.
•	URL NoIndex
/quy-hoach/gps
/quy-hoach?lat={lat}&lng={lng}
/quy-hoach/vi-tri-hien-tai
•	Canonical
•	Nếu resolve được Parcel: canonical về URL Parcel. 
•	Nếu resolve được Administrative Unit: canonical về URL địa bàn. 
•	Nếu không resolve được Entity: canonical về /quy-hoach. 
•	Không tạo URL SEO riêng cho
•	Tọa độ GPS. 
•	Accuracy. 
•	Device Location. 
•	Radius. 
•	Session. 
•	Tracking State. 
 
•	16.3 Metadata
•	Title
Không sinh Metadata riêng cho GPS runtime.
•	Description
Không sinh Description riêng.
•	H1
Kế thừa từ màn hình cha.
•	Robots
noindex,nofollow
•	Ghi chú
Không được đưa tọa độ người dùng vào Metadata.
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Không áp dụng
Entity Mapping	Không áp dụng
•	Điều kiện áp dụng
Không áp dụng.
•	Điều kiện loại bỏ
Luôn loại bỏ với GPS runtime.
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Không áp dụng
OG Description	Không áp dụng
OG Image	Không áp dụng
•	Share Preview
Không hỗ trợ chia sẻ GPS runtime.
Nếu người dùng chia sẻ kết quả, hệ thống phải chia sẻ URL Entity đích.
 
•	16.6 AI Summary
•	Có AI Summary
No
•	Dữ liệu nguồn

Không áp dụng cho GPS runtime.
•	Không public
•	Tọa độ người dùng. 
•	Vị trí hiện tại. 
•	Lịch sử vị trí. 
•	Thiết bị. 
•	Session. 
•	Ghi chú
AI Summary chỉ sinh tại Entity đích.
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	Không áp dụng
•	Tham gia khi
Không có trường hợp nào.
•	Loại khỏi Sitemap khi
Luôn loại bỏ.
 
•	16.8 Internal Link
•	Link đến
•	Parcel được xác định từ GPS. 
•	Administrative Unit được xác định từ GPS. 
•	Planning Region được xác định từ GPS. 
•	Link nhận từ
•	Trang bản đồ quy hoạch. 
•	Nút định vị. 
•	Mobile App. 
•	Ghi chú
GPS chỉ đóng vai trò điều hướng.
 
•	16.9 Security & Visibility
•	Public
Không có dữ liệu GPS cá nhân nào được public.
•	Restricted
•	Kết quả định vị chính xác. 
•	Dữ liệu vị trí runtime. 
•	Private
•	Tọa độ người dùng. 
•	Lịch sử vị trí. 
•	Session GPS. 
•	Device ID. 
•	Permission Rule
Không đưa dữ liệu GPS cá nhân vào bất kỳ SEO Output nào.
 
•	16.10 Cache & Regeneration
Event	Regenerate
parcel.updated	GPS Resolver Index
administrative_unit.updated	GPS Resolver Index
planning_region.updated	GPS Resolver Index
geometry.updated	GPS Resolver Index
visibility.changed	Permission Cache
•	Ghi chú
Không regenerate Metadata, Schema, OG hoặc Sitemap cho GPS runtime.
 
•	16.11 Acceptance Criteria
•	AC-D2.2-001
Không có GPS runtime URL được Index.
•	AC-D2.2-002
Không đưa tọa độ người dùng vào Metadata, Schema, OG hoặc AI Summary.
•	AC-D2.2-003
Không có GPS URL trong Sitemap.
•	AC-D2.2-004
Kết quả GPS điều hướng đúng tới Entity đích.
•	AC-D2.2-005
Dữ liệu vị trí cá nhân được bảo vệ.
 
•	16.12 QA Checklist
•	URL
□ GPS URL NoIndex
□ Không có URL tọa độ được Index
□ Canonical đúng về Entity đích hoặc /quy-hoach
•	Metadata
□ Không sinh Metadata riêng
□ Không chứa lat/lng trong Metadata
•	Structured Data
□ Không inject Schema
•	Open Graph
□ Không sinh OG riêng
□ Share kết quả dẫn tới Entity đích
•	Sitemap
□ Không xuất hiện trong Sitemap
•	Security
□ Không lộ tọa độ người dùng
□ Không lộ lịch sử vị trí
□ Không lộ Device ID / Session
 
D.2.3 QH.2.3 Tra cứu tờ/thửa
16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes
SEO Class	SEO-A / SEO-C
Entity chính	Parcel
Entity phụ	Administrative Unit, Planning Region
URL Type	Search Result / Detail URL
Mục tiêu SEO	Chuyển truy vấn tờ/thửa thành trang chi tiết thửa đất có khả năng Index
•	Ghi chú
•	Form tra cứu tờ/thửa không phải URL SEO chính. 
•	URL SEO chính là trang chi tiết Parcel. 
•	Search Result chỉ Index khi đủ điều kiện hình thành URL Parcel công khai. 
 
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/thua-dat/{parcel-slug}
•	URL tra cứu
/quy-hoach/tra-cuu-to-thua
•	URL NoIndex
/quy-hoach/tra-cuu-to-thua?to={sheet}&thua={parcel}
/quy-hoach/tim-thua?sheet={sheet}&parcel={parcel}
•	Canonical
•	Nếu tìm được Parcel: canonical về URL chi tiết Parcel. 
•	Nếu không tìm được Parcel: canonical về trang tra cứu tờ/thửa. 
•	Query runtime luôn NoIndex. 
•	Không tạo URL SEO riêng cho
•	Query tờ. 
•	Query thửa. 
•	Kết quả không xác định. 
•	Search session. 
•	Pagination tạm thời. 
 
•	16.3 Metadata
•	Title
Trang tra cứu:
Tra cứu tờ thửa quy hoạch | QH Pro
Trang Parcel:
Thửa đất {parcel_number}, tờ {sheet_number} tại {administrative_unit_name} | QH Pro
•	Description
Tra cứu thông tin quy hoạch, vị trí, địa chính và các dữ liệu liên quan của thửa đất số {parcel_number}, tờ bản đồ {sheet_number} tại {administrative_unit_name}.
•	H1
Thửa đất {parcel_number}, tờ {sheet_number}
•	Robots
•	Trang tra cứu form: 
noindex,follow
•	Trang Parcel đủ điều kiện: 
index,follow
•	NoIndex khi
•	Không tìm được Parcel. 
•	Parcel Restricted. 
•	Parcel không có dữ liệu công khai. 
•	Geometry không hợp lệ. 
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Place
Entity Mapping	Parcel
•	Schema bổ sung
Schema	Điều kiện
Dataset	Có dữ liệu quy hoạch công khai
BreadcrumbList	Có cấu trúc địa bàn đầy đủ
•	Không inject Schema nếu
•	Parcel Restricted. 
•	Parcel Draft. 
•	Missing Geometry. 
•	Không xác định được địa bàn. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Map Preview của Parcel
OG Type	article
•	Share Preview
Hiển thị:
•	Số tờ. 
•	Số thửa. 
•	Địa bàn. 
•	Bản đồ Preview. 
•	Tình trạng dữ liệu quy hoạch công khai. 
 
•	16.6 AI Summary
•	Có AI Summary
Yes, với Parcel đủ điều kiện công khai.
•	Dữ liệu nguồn
•	Parcel. 
•	Administrative Unit. 
•	Planning Region. 
•	Legal Document. 
•	GIS Layer công khai. 
•	Nội dung AI Summary
•	Tóm tắt vị trí. 
•	Tóm tắt quy hoạch. 
•	Tóm tắt pháp lý. 
•	Khu quy hoạch liên quan. 
•	Cảnh báo dữ liệu nếu có. 
•	Không public
•	Chủ sử dụng. 
•	Dữ liệu cá nhân. 
•	Internal Note. 
•	Restricted Data. 
•	Partner Data. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	parcel.xml
•	Tham gia khi
•	Parcel Public. 
•	Active. 
•	Có URL chuẩn. 
•	Có Geometry hợp lệ. 
•	Có Metadata tối thiểu. 
•	Loại bỏ khi
•	Deleted. 
•	Restricted. 
•	Draft. 
•	Missing Geometry. 
•	Không có dữ liệu công khai. 
 
•	16.8 Internal Link
•	Link đến
•	Administrative Unit. 
•	Planning Region. 
•	Planning Project. 
•	Legal Document. 
•	Report liên quan. 
•	Snapshot liên quan. 
•	Link nhận từ
•	Tra cứu tờ/thửa. 
•	Trang bản đồ. 
•	Tra cứu địa chỉ. 
•	Planning Region. 
•	Report. 
•	Snapshot. 
•	Thư viện quy hoạch. 
 
•	16.9 Security & Visibility
•	Public
•	Số tờ. 
•	Số thửa. 
•	Địa bàn. 
•	Geometry công khai. 
•	Thông tin quy hoạch công khai. 
•	Restricted
•	Dữ liệu trả phí. 
•	Dữ liệu đối tác. 
•	Dữ liệu hạn chế. 
•	Private
•	Chủ sử dụng đất. 
•	Thông tin cá nhân. 
•	Ghi chú nội bộ. 
•	Audit Log. 
•	Permission Rule
Không đưa dữ liệu Restricted hoặc Private vào SEO Output.
 
•	16.10 Cache & Regeneration
Event	Regenerate
parcel.updated	Metadata, AI Summary
parcel.geometry.updated	Schema, OG
parcel.visibility.changed	Sitemap, Robots
administrative_unit.updated	Metadata
planning_region.updated	AI Summary
legal_document.updated	AI Summary
 
•	16.11 Acceptance Criteria
•	AC-D2.3-001
Form tra cứu tờ/thửa không được Index như URL kết quả động.
•	AC-D2.3-002
Parcel đủ điều kiện có URL SEO chuẩn.
•	AC-D2.3-003
Canonical của query tra cứu trỏ về Parcel URL khi tìm thấy Parcel.
•	AC-D2.3-004
Metadata chứa đúng số tờ, số thửa và địa bàn.
•	AC-D2.3-005
Structured Data hợp lệ.
•	AC-D2.3-006
Parcel đủ điều kiện xuất hiện trong parcel.xml.
•	AC-D2.3-007
Không lộ dữ liệu cá nhân hoặc Restricted.
 
•	16.12 QA Checklist
•	URL
□ Form tra cứu NoIndex
□ Query runtime NoIndex
□ Parcel URL đúng chuẩn
□ Canonical đúng Parcel
•	Metadata
□ Title có số tờ/số thửa
□ Description đúng địa bàn
□ H1 tồn tại
□ Robots đúng
•	Structured Data
□ Schema Place hợp lệ
□ BreadcrumbList đúng nếu có
•	Open Graph
□ OG Image đúng Map Preview
□ Share Preview không lộ dữ liệu riêng tư
•	AI Summary
□ Có AI Summary cho Parcel Public
□ Không lộ chủ sử dụng / dữ liệu cá nhân
•	Sitemap
□ Parcel đủ điều kiện có trong parcel.xml
□ Parcel Restricted bị loại khỏi Sitemap
•	Security
□ Không lộ dữ liệu cá nhân
□ Không lộ dữ liệu Restricted
□ Permission hoạt động đúng
 
D.2.4 QH.2.4 Tra cứu polygon / vùng phân tích
16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes, có điều kiện
SEO Class	SEO-C / SEO-D
Entity chính	Planning Region / Report / Snapshot
Entity phụ	Parcel, GIS Layer, Legal Document
URL Type	Analysis / Snapshot / Report URL
Mục tiêu SEO	Chuyển vùng phân tích đủ điều kiện thành nội dung công khai có giá trị SEO hoặc Share
•	Ghi chú
•	Polygon runtime không phải đối tượng SEO. 
•	Vùng phân tích chỉ SEO khi được lưu thành Report, Snapshot công khai hoặc Planning Region có giá trị độc lập. 
•	Vùng người dùng vẽ tạm thời mặc định NoIndex. 
 
•	16.2 URL & Canonical
•	URL chuẩn
Với Report công khai:
/quy-hoach/bao-cao/{report-slug}
Với Snapshot công khai:
/quy-hoach/snapshot/{snapshot-slug}
Với Planning Region công khai:
/quy-hoach/khu-vuc/{planning-region-slug}
•	URL NoIndex
/quy-hoach/phan-tich?polygon={geometry}
/quy-hoach/ve-vung?bbox={bbox}
/quy-hoach/analysis-session/{session-id}
•	Canonical
•	Report công khai canonical về Report URL. 
•	Snapshot công khai canonical về Snapshot URL hoặc Entity nguồn. 
•	Polygon runtime canonical về /quy-hoach. 
•	Nếu vùng trùng Planning Region đã có: canonical về Planning Region URL. 
•	Không tạo URL SEO riêng cho
•	Geometry raw. 
•	BBox. 
•	Vertex. 
•	Draw State. 
•	Session. 
•	Temporary Polygon. 
 
•	16.3 Metadata
•	Title
Report:
Báo cáo phân tích quy hoạch {region_name} | QH Pro
Snapshot:
Snapshot quy hoạch {region_name} | QH Pro
Planning Region:
Quy hoạch khu vực {region_name} | QH Pro
•	Description
Phân tích quy hoạch, hiện trạng, lớp dữ liệu liên quan và các thông tin pháp lý trong phạm vi {region_name}.
•	H1
Phân tích quy hoạch {region_name}
•	Robots
•	Runtime Polygon: 
noindex,nofollow
•	Public Report / Region: 
index,follow
•	Share-only Snapshot: 
noindex,follow
•	NoIndex khi
•	Polygon tạm thời. 
•	Private Snapshot. 
•	Restricted Report. 
•	Missing Geometry. 
•	Không có dữ liệu công khai. 
 
•	16.4 Structured Data
Entity	Schema
Planning Region	Place
Report	Report
Snapshot	CreativeWork
GIS Layer	Dataset
•	Không inject Schema nếu
•	Polygon runtime. 
•	Private. 
•	Restricted. 
•	Missing Geometry. 
•	Dữ liệu chưa công bố. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Map Preview của vùng phân tích
OG Type	article
•	Share Preview
Hiển thị:
•	Tên vùng. 
•	Bản đồ Preview. 
•	Tóm tắt kết quả phân tích. 
•	Cảnh báo nếu là Snapshot tham khảo. 
 
•	16.6 AI Summary
•	Có AI Summary
Yes, với Report hoặc Planning Region công khai.
•	Dữ liệu nguồn
•	Geometry vùng. 
•	Planning Region. 
•	Parcel trong vùng. 
•	GIS Layer công khai. 
•	Legal Document. 
•	Planning Project. 
•	Nội dung AI Summary
•	Tóm tắt phạm vi vùng. 
•	Tỷ lệ các lớp quy hoạch chính. 
•	Các cảnh báo quan trọng. 
•	Văn bản pháp lý liên quan. 
•	Đồ án liên quan. 
•	Không public
•	Polygon cá nhân. 
•	Geometry private. 
•	Kết quả phân tích nội bộ. 
•	Restricted Layer. 
•	Dữ liệu người dùng. 

•	16.7 Sitemap
Thành phần	Giá trị
Report	report.xml
Planning Region	planning-region.xml
Snapshot	Không mặc định
•	Tham gia khi
•	Public. 
•	Active. 
•	Có Geometry hợp lệ. 
•	Có Metadata. 
•	Có nội dung phân tích đủ điều kiện. 
•	Loại bỏ khi
•	Runtime. 
•	Private. 
•	Restricted. 
•	Deleted. 
•	Missing Geometry. 
•	Share-only. 
 
•	16.8 Internal Link
•	Link đến
•	Planning Region. 
•	Parcel trong vùng. 
•	Planning Project. 
•	Legal Document. 
•	GIS Layer công khai. 
•	Report liên quan. 
•	Link nhận từ
•	Trang bản đồ. 
•	Tra cứu polygon. 
•	Report. 
•	Snapshot. 
•	Thư viện quy hoạch. 
•	Planning Region liên quan. 
 
•	16.9 Security & Visibility
•	Public
•	Tên vùng công khai. 
•	Bản đồ Preview. 
•	Tóm tắt phân tích. 
•	Layer công khai. 
•	Văn bản công khai. 
•	Restricted
•	Layer trả phí. 
•	Layer đối tác. 
•	Phân tích chuyên sâu yêu cầu quyền. 
•	Private
•	Polygon cá nhân. 
•	Vùng do user lưu riêng. 
•	Session phân tích. 
•	Dữ liệu cá nhân. 
•	Permission Rule
Không đưa dữ liệu Private / Restricted vào Metadata, Schema, OG, AI Summary hoặc Sitemap.
 
•	16.10 Cache & Regeneration
Event	Regenerate
report.updated	Metadata, AI Summary
snapshot.visibility.changed	Robots, Sitemap
region.geometry.updated	Schema, OG
gis_layer.updated	AI Summary
legal_document.updated	AI Summary
visibility.changed	Sitemap, Robots
map_preview.updated	Open Graph
 
•	16.11 Acceptance Criteria
•	AC-D2.4-001
Polygon runtime luôn NoIndex.
•	AC-D2.4-002
Không Index URL chứa raw geometry, bbox hoặc session.
•	AC-D2.4-003
Report công khai có URL SEO chuẩn.
•	AC-D2.4-004
Planning Region công khai có URL SEO chuẩn.
•	AC-D2.4-005
Snapshot cá nhân không xuất hiện trong Sitemap.
•	AC-D2.4-006
AI Summary chỉ dùng dữ liệu công khai.
•	AC-D2.4-007
Không lộ polygon cá nhân hoặc dữ liệu Restricted.
 
•	16.12 QA Checklist
•	URL
□ Polygon runtime NoIndex
□ Raw geometry không Index
□ Report URL đúng chuẩn
□ Region URL đúng chuẩn
□ Canonical đúng Entity nguồn
•	Metadata
□ Title đúng vùng / Report
□ Description đúng phạm vi
□ H1 tồn tại
□ Robots đúng theo trạng thái
•	Structured Data
□ Schema Report hợp lệ
□ Schema Place hợp lệ nếu là Region
□ Không inject Schema cho runtime
•	Open Graph
□ OG Image đúng Map Preview
□ Share Preview không lộ dữ liệu private
•	AI Summary
□ Có AI Summary cho Report công khai
□ Không dùng layer Restricted
□ Không lộ polygon cá nhân
•	Sitemap
□ Report đủ điều kiện có trong report.xml
□ Runtime / Private / Share-only bị loại khỏi Sitemap
•	Security
□ Không lộ Geometry Private
□ Không lộ dữ liệu người dùng
□ Không lộ dữ liệu Restricted
□ Permission hoạt động đúng

 
D.3	QH.3 – Hệ thống Layer GIS
D.3.1 Layer Public
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes, có điều kiện
SEO Class	SEO-B / SEO-C
Entity chính	GIS Layer
Entity phụ	Administrative Unit, Planning Region, Planning Project, Legal Document
URL Type	Layer Detail / Dataset URL
Mục tiêu SEO	SEO các lớp dữ liệu GIS công khai có giá trị độc lập, phục vụ Google Search, AI Search và Knowledge Graph
•	Ghi chú
•	Chỉ SEO các Layer công khai, ổn định, có ngữ nghĩa quy hoạch rõ ràng. 
•	Không SEO Layer kỹ thuật, Layer runtime, Layer private hoặc Layer chưa công bố. 
•	Layer Public đóng vai trò bổ trợ cho các trang SEO-A như Đồ án, Khu quy hoạch, Bản đồ quy hoạch và Đơn vị hành chính. 
 
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/layer/{layer-slug}
•	URL phụ
/quy-hoach/layer/{administrative-unit-slug}/{layer-slug}
•	URL NoIndex
/quy-hoach?layer={layer-id}
/quy-hoach?layers={layer-list}
/quy-hoach/layer-preview/{layer-id}
•	Canonical
•	Layer công khai có URL riêng: canonical về URL Layer chuẩn. 
•	Layer gắn với đồ án cụ thể: canonical về URL chi tiết đồ án nếu Layer không có giá trị độc lập. 
•	Layer chỉ là trạng thái bật/tắt trên bản đồ: canonical về /quy-hoach. 
•	Không tạo URL SEO riêng cho
•	Layer toggle. 
•	Layer opacity. 
•	Layer order. 
•	Layer visibility state. 
•	Layer filter runtime. 
•	Layer preview tạm thời. 
 
•	16.3 Metadata
•	Title
{layer_name} | Lớp dữ liệu quy hoạch | QH Pro
•	Description
Thông tin lớp dữ liệu {layer_name}, phạm vi áp dụng, trạng thái công khai, nguồn dữ liệu và các quy hoạch liên quan trên hệ thống QH Pro.
•	H1
{layer_name}
•	Robots
index,follow
•	NoIndex khi
•	Layer không công khai. 
•	Layer đang Draft. 
•	Layer Restricted. 
•	Layer chỉ phục vụ runtime. 
•	Layer thiếu metadata tối thiểu. 
•	Layer không có phạm vi địa lý hoặc ngữ nghĩa rõ ràng. 
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Dataset
Entity Mapping	GIS Layer
•	Schema bổ sung
Schema	Điều kiện
Place	Layer gắn với địa bàn cụ thể
CreativeWork	Layer là bản đồ / tài liệu quy hoạch số hóa
BreadcrumbList	Có cấu trúc điều hướng đầy đủ
•	Không inject Schema nếu
•	Layer Restricted. 
•	Layer Draft. 
•	Layer thiếu metadata. 
•	Layer thiếu phạm vi áp dụng. 
•	Layer là runtime layer. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Map Preview của Layer
OG Type	article
•	Share Preview
Hiển thị:
•	Tên Layer. 
•	Nhóm Layer. 
•	Địa bàn áp dụng. 
•	Map Preview. 
•	Trạng thái dữ liệu công khai. 
 
•	16.6 AI Summary
•	Có AI Summary
Yes, nếu Layer đủ điều kiện công khai.
•	Dữ liệu nguồn
•	GIS Layer. 
•	Layer Metadata. 
•	Administrative Unit. 
•	Planning Region. 
•	Planning Project. 
•	Legal Document. 
•	Nội dung AI Summary
•	Mô tả Layer. 
•	Phạm vi áp dụng. 
•	Nhóm dữ liệu. 
•	Nguồn dữ liệu. 
•	Trạng thái công khai. 
•	Liên kết tới đồ án / văn bản liên quan. 
•	Không public
•	Internal Layer. 
•	Restricted Layer. 
•	Partner-only Layer. 
•	Metadata kỹ thuật nội bộ. 
•	Rule render / priority nội bộ. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	gis-layer.xml
•	Tham gia khi
•	Public. 
•	Active. 
•	Có Metadata hợp lệ. 
•	Có phạm vi địa lý hoặc phạm vi áp dụng rõ ràng. 
•	Có URL chuẩn. 
•	Loại bỏ khi
•	Draft. 
•	Restricted. 
•	Deleted. 
•	Runtime only. 
•	Internal only. 
•	Missing Metadata. 
 
•	16.8 Internal Link
•	Link đến
•	Administrative Unit. 
•	Planning Region. 
•	Planning Project. 
•	Legal Document. 
•	Planning Map. 
•	Trang bản đồ quy hoạch. 
•	Link nhận từ
•	Trang bản đồ quy hoạch. 
•	Panel lớp dữ liệu. 
•	Chi tiết đồ án. 
•	Chi tiết bản đồ quy hoạch. 
•	Thư viện quy hoạch. 
•	Báo cáo / Snapshot. 
 
•	16.9 Security & Visibility
•	Public
•	Tên Layer. 
•	Mô tả Layer. 
•	Nhóm Layer. 
•	Phạm vi áp dụng. 
•	Map Preview công khai. 
•	Restricted
•	Layer trả phí. 
•	Layer đối tác. 
•	Layer chưa công bố. 
•	Layer có hạn chế truy cập. 
•	Private
•	Layer nội bộ. 
•	Metadata vận hành. 
•	Rule render. 
•	Rule priority. 
•	Log xử lý dữ liệu. 
•	Permission Rule
Không đưa Layer Restricted hoặc Private vào Metadata, Schema, OG, AI Summary hoặc Sitemap.
 
•	16.10 Cache & Regeneration
Event	Regenerate
layer.created	Metadata, Sitemap
layer.updated	Metadata, AI Summary
layer.visibility.changed	Robots, Sitemap
layer.metadata.updated	Metadata, Schema
layer.geometry.updated	Schema, OG
layer.map_preview.updated	Open Graph
legal_document.updated	AI Summary
 
•	16.11 Acceptance Criteria
•	AC-D3.1-001
Layer Public đủ điều kiện có URL SEO chuẩn.
•	AC-D3.1-002
Layer runtime không được Index.
•	AC-D3.1-003
Canonical đúng về URL Layer hoặc URL Entity nguồn.
•	AC-D3.1-004
Metadata sinh đầy đủ từ Layer Metadata.
•	AC-D3.1-005
Structured Data dạng Dataset hợp lệ.
•	AC-D3.1-006
Layer đủ điều kiện xuất hiện trong gis-layer.xml.
•	AC-D3.1-007
Không lộ Layer Restricted hoặc Private.
 
•	16.12 QA Checklist
•	URL
□ URL Layer đúng chuẩn
□ Runtime Layer NoIndex
□ Canonical đúng
□ Không sinh URL từ layer toggle
•	Metadata
□ Title tồn tại
□ Description tồn tại
□ H1 tồn tại
□ Robots đúng
•	Structured Data
□ Schema Dataset hợp lệ
□ Entity Mapping đúng GIS Layer
□ Không inject khi Restricted
•	Open Graph
□ OG Title đúng
□ OG Description đúng
□ OG Image đúng Map Preview
•	AI Summary
□ AI Summary đúng Layer Public
□ Không dùng metadata nội bộ
□ Không lộ Layer Restricted
•	Sitemap
□ Layer đủ điều kiện có trong gis-layer.xml
□ Layer Draft / Restricted bị loại khỏi Sitemap
•	Security
□ Không lộ Layer nội bộ
□ Không lộ rule render
□ Không lộ metadata vận hành
 
D.3.2 Layer Metadata
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes, có điều kiện
SEO Class	SEO-C
Entity chính	GIS Layer Metadata
Entity phụ	GIS Layer, Planning Project, Legal Document
URL Type	Metadata / Dataset URL
Mục tiêu SEO	Chuẩn hóa metadata công khai để hỗ trợ SEO, AI Search, Structured Data và độ tin cậy dữ liệu
•	Ghi chú
•	Layer Metadata không phải lúc nào cũng có trang SEO riêng. 
•	Metadata chủ yếu là dữ liệu nguồn phục vụ SEO Output cho Layer, Đồ án, Bản đồ và Báo cáo. 
•	Chỉ SEO riêng khi metadata có giá trị công khai độc lập. 
 
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/layer/{layer-slug}/metadata
•	URL NoIndex
/quy-hoach/layer-metadata-preview/{metadata-id}
/quy-hoach/layer/{layer-id}/metadata?raw=true
•	Canonical
•	Nếu Metadata có giá trị độc lập: canonical về URL Metadata. 
•	Nếu Metadata chỉ bổ trợ cho Layer: canonical về URL Layer. 
•	Raw metadata luôn NoIndex. 
•	Không tạo URL SEO riêng cho
•	Raw JSON. 
•	Preview nội bộ. 
•	Debug metadata. 
•	Version kỹ thuật. 
•	Metadata chưa công bố. 
 
•	16.3 Metadata
•	Title
Metadata lớp dữ liệu {layer_name} | QH Pro
•	Description
Thông tin metadata của lớp dữ liệu {layer_name}, bao gồm nguồn dữ liệu, phạm vi áp dụng, trạng thái công khai, phiên bản và các tài liệu liên quan.
•	H1
Metadata lớp dữ liệu {layer_name}
•	Robots
index,follow
•	NoIndex khi
•	Metadata nội bộ. 
•	Metadata raw. 
•	Metadata thiếu thông tin bắt buộc. 
•	Layer nguồn không công khai. 
•	Metadata Draft / Restricted. 
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Dataset
Entity Mapping	GIS Layer Metadata
•	Schema bổ sung
Schema	Điều kiện
CreativeWork	Metadata gắn với tài liệu / bản đồ
BreadcrumbList	Có điều hướng đầy đủ
•	Không inject Schema nếu
•	Metadata raw. 
•	Metadata Restricted. 
•	Layer nguồn không Public. 
•	Thiếu trường nguồn dữ liệu hoặc phạm vi áp dụng. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Map Preview của Layer
OG Type	article
•	Share Preview
Hiển thị:
•	Tên Layer. 
•	Nguồn dữ liệu. 
•	Phiên bản dữ liệu. 
•	Phạm vi áp dụng. 
 
•	16.6 AI Summary
•	Có AI Summary
Yes, nếu Metadata công khai.
•	Dữ liệu nguồn
•	Layer Metadata. 
•	GIS Layer. 
•	Legal Document. 
•	Planning Project. 
•	Administrative Unit. 
•	Nội dung AI Summary
•	Nguồn dữ liệu. 
•	Phạm vi áp dụng. 
•	Phiên bản dữ liệu. 
•	Trạng thái công khai. 
•	Độ tin cậy dữ liệu. 
•	Văn bản liên quan. 
•	Không public
•	Internal Metadata. 
•	Import Log. 
•	Validation Error chi tiết. 
•	Source Credential. 
•	Rule kỹ thuật nội bộ. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	gis-layer-metadata.xml
•	Tham gia khi
•	Metadata Public. 
•	Layer nguồn Public. 
•	Có URL chuẩn. 
•	Có Metadata tối thiểu. 
•	Có giá trị tham chiếu độc lập. 
•	Loại bỏ khi
•	Raw. 
•	Draft. 
•	Restricted. 
•	Internal. 
•	Missing Required Metadata. 
•	Layer nguồn không Public. 
 
•	16.8 Internal Link
•	Link đến
•	GIS Layer. 
•	Planning Project. 
•	Legal Document. 
•	Planning Map. 
•	Administrative Unit. 
•	Report sử dụng Layer. 
•	Link nhận từ
•	Trang Layer Public. 
•	Chi tiết bản đồ quy hoạch. 
•	Chi tiết đồ án. 
•	Thư viện quy hoạch. 
•	Report / Snapshot. 
 
•	16.9 Security & Visibility
•	Public
•	Tên nguồn dữ liệu. 
•	Ngày cập nhật. 
•	Phạm vi áp dụng. 
•	Trạng thái công khai. 
•	Phiên bản công khai. 
•	Restricted
•	Metadata đối tác. 
•	Metadata chưa thẩm định. 
•	Metadata có điều kiện sử dụng. 
•	Private
•	Import Log. 
•	Error Log. 
•	Credential. 
•	Internal Mapping. 
•	Data Pipeline Note. 
•	Permission Rule
Không đưa metadata nội bộ hoặc metadata hạn chế vào bất kỳ SEO Output nào.
 
•	16.10 Cache & Regeneration
Event	Regenerate
layer.metadata.updated	Metadata, Schema, AI Summary
layer.source.updated	Metadata, AI Summary
layer.version.updated	Metadata, Sitemap
layer.visibility.changed	Robots, Sitemap
legal_document.updated	AI Summary
validation_status.changed	Robots
 
•	16.11 Acceptance Criteria
•	AC-D3.2-001
Metadata công khai đủ điều kiện có URL hoặc canonical đúng về Layer.
•	AC-D3.2-002
Raw metadata luôn NoIndex.
•	AC-D3.2-003
Metadata không được làm lộ dữ liệu nội bộ.
•	AC-D3.2-004
Structured Data dạng Dataset hợp lệ.
•	AC-D3.2-005
AI Summary chỉ dùng metadata công khai.
•	AC-D3.2-006
Sitemap đúng điều kiện.
 
•	16.12 QA Checklist
•	URL
□ Metadata URL đúng chuẩn
□ Raw metadata NoIndex
□ Canonical đúng về Metadata hoặc Layer
•	Metadata
□ Title tồn tại
□ Description tồn tại
□ H1 tồn tại
□ Robots đúng
•	Structured Data
□ Schema Dataset hợp lệ
□ Không inject khi Layer nguồn Restricted
•	Open Graph
□ OG đúng Layer nguồn
□ Preview không lộ dữ liệu nội bộ
•	AI Summary
□ AI Summary mô tả đúng nguồn dữ liệu
□ Không lộ import log / credential / error log
•	Sitemap
□ Metadata đủ điều kiện có trong Sitemap
□ Raw / Internal bị loại khỏi Sitemap
•	Security
□ Không lộ metadata nội bộ
□ Không lộ pipeline note
□ Không lộ thông tin đối tác hạn chế
 
D.3.3 Layer Category
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes, có điều kiện
SEO Class	SEO-B
Entity chính	Layer Category
Entity phụ	GIS Layer, Administrative Unit, Planning Project
URL Type	List / Category URL
Mục tiêu SEO	Tạo trang danh mục điều hướng các nhóm Layer công khai, hỗ trợ Internal Link và Discoverability
•	Ghi chú
•	Layer Category có thể SEO khi là nhóm dữ liệu công khai có giá trị tìm kiếm. 
•	Không SEO Category nội bộ, category kỹ thuật hoặc category chỉ dùng cấu hình hiển thị. 
•	Đây là nhóm SEO hỗ trợ, không phải SEO đích chính. 
 
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/layer-category/{category-slug}

•	URL phụ
/quy-hoach/layer/{category-slug}
•	URL NoIndex
/quy-hoach?layer_category={category-id}
/quy-hoach/layer-category/{category-slug}?sort={sort}
/quy-hoach/layer-category/{category-slug}?filter={filter}
•	Canonical
•	Category công khai: canonical về URL Category chuẩn. 
•	Category chỉ là filter runtime: canonical về /quy-hoach. 
•	Category gắn với thư viện quy hoạch: có thể canonical về Workspace / danh sách bản đồ nếu không đủ giá trị độc lập. 
•	Không tạo URL SEO riêng cho
•	Sort. 
•	Filter. 
•	Layer toggle. 
•	Runtime category state. 
•	Category ẩn. 
•	Category nội bộ. 
 
•	16.3 Metadata
•	Title
{category_name} | Nhóm lớp dữ liệu quy hoạch | QH Pro
•	Description
Danh sách các lớp dữ liệu thuộc nhóm {category_name}, bao gồm các lớp quy hoạch, địa chính, bản đồ và dữ liệu GIS công khai liên quan.
•	H1
{category_name}
•	Robots
index,follow
•	NoIndex khi
•	Category nội bộ. 
•	Category không có Layer công khai. 
•	Category chỉ dùng cho runtime UI. 
•	Category Draft / Restricted. 
•	Category trùng lặp với danh mục khác. 
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	CollectionPage
Entity Mapping	Layer Category
•	Schema bổ sung
Schema	Điều kiện
Dataset	Category đại diện cho nhóm dữ liệu công khai
BreadcrumbList	Có cấu trúc điều hướng đầy đủ
•	Không inject Schema nếu
•	Category không công khai. 
•	Không có Layer công khai. 
•	Category chỉ là UI group. 
•	Category Restricted. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Map Preview / Category Preview
OG Type	website
•	Share Preview
Hiển thị:
•	Tên Category. 
•	Số lượng Layer công khai. 
•	Nhóm dữ liệu. 
•	Preview bản đồ nếu có. 
 
•	16.6 AI Summary
•	Có AI Summary
Yes, nếu Category công khai và có đủ Layer.
•	Dữ liệu nguồn
•	Layer Category. 
•	GIS Layer công khai. 
•	Administrative Unit. 
•	Planning Project. 
•	Legal Document liên quan. 
•	Nội dung AI Summary
•	Mô tả nhóm Layer. 
•	Các Layer nổi bật. 
•	Phạm vi dữ liệu. 
•	Ứng dụng trong tra cứu quy hoạch. 
•	Liên kết tới các Layer / Đồ án / Bản đồ liên quan. 
•	Không public
•	Category nội bộ. 
•	Category kỹ thuật. 
•	Layer Restricted. 
•	Rule nhóm nội bộ. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	gis-layer-category.xml
•	Tham gia khi
•	Public. 
•	Active. 
•	Có tối thiểu một Layer công khai. 
•	Có URL chuẩn. 
•	Có Metadata hợp lệ. 
•	Loại bỏ khi
•	Draft. 
•	Restricted. 
•	Empty Category. 
•	Internal Category. 
•	Runtime-only Category. 
 
•	16.8 Internal Link
•	Link đến
•	GIS Layer công khai trong Category. 
•	Planning Project liên quan. 
•	Planning Map liên quan. 
•	Administrative Unit liên quan. 
•	Thư viện quy hoạch. 
•	Link nhận từ
•	Trang bản đồ quy hoạch. 
•	Panel lớp dữ liệu. 
•	Workspace thư viện. 
•	Trang danh sách bản đồ. 
•	Chi tiết Layer. 
 
•	16.9 Security & Visibility
•	Public
•	Tên Category. 
•	Mô tả Category. 
•	Danh sách Layer công khai. 
•	Preview công khai. 
•	Restricted
•	Category trả phí. 
•	Category đối tác. 
•	Category có Layer hạn chế. 
•	Private
•	Category nội bộ. 
•	Category kỹ thuật. 
•	Category dùng cho cấu hình hiển thị. 
•	Rule phân nhóm nội bộ. 
•	Permission Rule
Category SEO chỉ được hiển thị Layer công khai. Không được để danh sách công khai chứa Layer Restricted hoặc Private.
 
•	16.10 Cache & Regeneration
Event	Regenerate
category.created	Metadata, Sitemap
category.updated	Metadata, AI Summary
category.visibility.changed	Robots, Sitemap
layer.added_to_category	AI Summary, Sitemap
layer.removed_from_category	AI Summary, Sitemap
layer.visibility.changed	AI Summary, Sitemap
category.preview.updated	Open Graph
 
•	16.11 Acceptance Criteria
•	AC-D3.3-001
Category công khai có URL chuẩn.
•	AC-D3.3-002
Category runtime hoặc nội bộ luôn NoIndex.
•	AC-D3.3-003
Metadata sinh đúng theo Category.
•	AC-D3.3-004
Structured Data dạng CollectionPage hợp lệ.
•	AC-D3.3-005
Sitemap chỉ chứa Category đủ điều kiện.
•	AC-D3.3-006
Không hiển thị Layer Restricted trong Category công khai.
 
•	16.12 QA Checklist
•	URL
□ Category URL đúng chuẩn
□ Filter / sort NoIndex
□ Canonical đúng
•	Metadata
□ Title tồn tại
□ Description tồn tại
□ H1 tồn tại
□ Robots đúng
•	Structured Data
□ CollectionPage hợp lệ
□ Dataset chỉ inject khi đủ điều kiện
•	Open Graph
□ OG Title đúng
□ OG Image đúng Category Preview
•	AI Summary
□ AI Summary mô tả đúng nhóm Layer
□ Không nhắc Layer Restricted
•	Sitemap
□ Category đủ điều kiện có trong gis-layer-category.xml
□ Empty / Internal Category bị loại khỏi Sitemap
•	Security
□ Không lộ Category nội bộ
□ Không lộ Layer Restricted
□ Không lộ rule phân nhóm kỹ thuật
 
D.4	QH.4 – Diễn giải & Phân tích Quy hoạch
D.4.1 Màn hình kết quả diễn giải
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes
SEO Class	SEO-A / SEO-C
Entity chính	Planning Information
Entity phụ	Parcel, Planning Region, Planning Project, Legal Document
URL Type	Detail / Analysis Result URL
Mục tiêu SEO	Biến kết quả diễn giải quy hoạch thành nội dung có khả năng đọc hiểu, lập chỉ mục và phục vụ AI Search
•	Ghi chú
•	Kết quả diễn giải là lớp nội dung có giá trị SEO cao. 
•	Không SEO kết quả tạm thời chưa lưu. 
•	Chỉ SEO kết quả gắn với Entity công khai và đủ dữ liệu. 
 
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/dien-giai/{entity-slug}
•	URL NoIndex
/quy-hoach/dien-giai?selection={id}
/quy-hoach/dien-giai?mode=preview
•	Canonical
•	Nếu diễn giải gắn với Parcel: canonical về URL Parcel. 
•	Nếu gắn với Planning Region: canonical về URL Planning Region. 
•	Nếu gắn với Planning Project: canonical về URL Planning Project. 
•	Nếu là bản diễn giải độc lập được công khai: canonical về URL diễn giải. 
•	Không tạo URL SEO riêng cho
•	Preview mode. 
•	Tab. 
•	Filter. 
•	Layer state. 
•	Selection runtime. 
•	Temporary analysis. 
 
•	16.3 Metadata
•	Title
Diễn giải quy hoạch {entity_name} | QH Pro
•	Description
Tóm tắt và diễn giải thông tin quy hoạch, trạng thái pháp lý, lớp dữ liệu liên quan và các cảnh báo chính đối với {entity_name}.
•	H1
Diễn giải quy hoạch {entity_name}
•	Robots
index,follow
•	NoIndex khi
•	Kết quả tạm thời. 
•	Dữ liệu Restricted. 
•	Không có Entity nguồn hợp lệ. 
•	Không đủ dữ liệu diễn giải. 
•	Draft / Deleted. 
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Article / CreativeWork
Entity Mapping	Planning Information
•	Schema bổ sung
Schema	Điều kiện
Place	Diễn giải gắn với Parcel / Region
Dataset	Có dữ liệu quy hoạch công khai
BreadcrumbList	Có cấu trúc điều hướng
•	Không inject Schema nếu
•	Restricted. 
•	Draft. 
•	Missing Entity. 
•	Nội dung diễn giải chưa đủ điều kiện công khai. 

•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Map Preview của Entity nguồn
OG Type	article
•	Share Preview
Hiển thị:
•	Tên Entity. 
•	Tóm tắt diễn giải. 
•	Map Preview. 
•	Trạng thái dữ liệu. 
 
•	16.6 AI Summary
•	Có AI Summary
Yes
•	Dữ liệu nguồn
•	Planning Information Model. 
•	Parcel. 
•	Planning Region. 
•	Planning Project. 
•	Legal Document. 
•	GIS Layer công khai. 
•	Nội dung AI Summary
•	Tóm tắt lớp quy hoạch chính. 
•	Trạng thái pháp lý tổng hợp. 
•	Ý nghĩa quy hoạch. 
•	Cảnh báo chính. 
•	Văn bản pháp lý liên quan. 
•	Không public
•	Internal Note. 
•	Raw Layer chưa chuẩn hóa. 
•	Restricted Data. 
•	Private Data. 
•	Dữ liệu chưa công bố. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	planning-information.xml
•	Tham gia khi
•	Public. 
•	Active. 
•	Có Entity nguồn hợp lệ. 
•	Có nội dung diễn giải đủ điều kiện. 
•	Có Metadata hợp lệ. 
•	Loại bỏ khi
•	Restricted. 
•	Draft. 
•	Deleted. 
•	Không có Entity nguồn. 
•	Nội dung quá ngắn hoặc không đủ ý nghĩa độc lập. 
 
•	16.8 Internal Link
•	Link đến
•	Entity nguồn. 
•	Administrative Unit. 
•	Planning Region. 
•	Planning Project. 
•	Legal Document. 
•	GIS Layer công khai. 
•	Report liên quan. 
•	Link nhận từ
•	Trang chi tiết thửa đất. 
•	Trang chi tiết khu quy hoạch. 
•	Trang bản đồ. 
•	Kết quả tra cứu. 
•	Report. 
•	Snapshot. 
 
•	16.9 Security & Visibility
•	Public
•	Tóm tắt diễn giải. 
•	Metadata. 
•	AI Summary. 
•	Map Preview. 
•	Văn bản công khai. 
•	Restricted
•	Diễn giải chuyên sâu trả phí. 
•	Dữ liệu đối tác. 
•	Layer hạn chế. 
•	Private
•	Ghi chú nội bộ. 
•	Dữ liệu cá nhân. 
•	Log phân tích. 
•	Permission Rule
Không đưa dữ liệu Restricted hoặc Private vào Metadata, Schema, OG, AI Summary hoặc Sitemap.
 
•	16.10 Cache & Regeneration
Event	Regenerate
planning_information.updated	Metadata, AI Summary
legal_status.changed	Metadata, AI Summary
layer.updated	AI Summary
entity.visibility.changed	Robots, Sitemap
map_preview.updated	Open Graph
legal_document.updated	AI Summary
 
•	16.11 Acceptance Criteria
•	AC-D4.1-001
Kết quả diễn giải công khai có URL chuẩn hoặc canonical đúng về Entity nguồn.
•	AC-D4.1-002
Kết quả tạm thời luôn NoIndex.
•	AC-D4.1-003
Metadata sinh đầy đủ theo Entity nguồn.
•	AC-D4.1-004
Structured Data hợp lệ.
•	AC-D4.1-005
AI Summary chỉ sử dụng dữ liệu công khai.
•	AC-D4.1-006
Sitemap đúng điều kiện.
•	AC-D4.1-007
Không lộ dữ liệu Restricted hoặc Private.
 
•	16.12 QA Checklist
•	URL
□ URL đúng chuẩn
□ Canonical đúng Entity nguồn
□ Runtime URL NoIndex
•	Metadata
□ Title tồn tại
□ Description tồn tại
□ H1 tồn tại
□ Robots đúng
•	Structured Data
□ Schema hợp lệ
□ Entity Mapping đúng
•	Open Graph
□ OG Title đúng
□ OG Description đúng
□ OG Image đúng
•	AI Summary
□ Có AI Summary
□ Nội dung đúng nguồn
□ Không lộ dữ liệu Restricted
•	Sitemap
□ Có trong Sitemap khi đủ điều kiện
□ Loại khỏi Sitemap đúng điều kiện
•	Security
□ Không lộ dữ liệu Restricted
□ Không lộ dữ liệu Private
 
D.4.2 Màn hình phân tích quy hoạch
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes, có điều kiện
SEO Class	SEO-C
Entity chính	Analysis Result / Report
Entity phụ	Parcel, Planning Region, GIS Layer, Legal Document
URL Type	Analysis / Report URL
Mục tiêu SEO	SEO các kết quả phân tích công khai có giá trị độc lập; NoIndex các phiên phân tích tạm thời
•	Ghi chú
•	Phân tích runtime không SEO. 
•	Chỉ SEO khi kết quả phân tích được lưu, công khai và có giá trị tham chiếu độc lập. 
•	Report công khai có thể trở thành SEO-C hoặc SEO-A tùy chất lượng nội dung. 
 
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/phan-tich/{analysis-slug}
hoặc:
/quy-hoach/bao-cao/{report-slug}
•	URL NoIndex
/quy-hoach/phan-tich?session={id}
/quy-hoach/phan-tich?polygon={geometry}
/quy-hoach/phan-tich?layers={layers}
•	Canonical
•	Analysis công khai canonical về URL Analysis. 
•	Report công khai canonical về URL Report. 
•	Analysis runtime canonical về Entity nguồn hoặc /quy-hoach. 
•	Không tạo URL SEO riêng cho
•	Session phân tích. 
•	Layer selection. 
•	Runtime polygon. 
•	Filter. 
•	Compare mode tạm thời. 
 
•	16.3 Metadata
•	Title
Phân tích quy hoạch {analysis_name} | QH Pro
•	Description
Kết quả phân tích quy hoạch, lớp dữ liệu, phạm vi ảnh hưởng, trạng thái pháp lý và các cảnh báo liên quan đối với {analysis_name}.
•	H1
Phân tích quy hoạch {analysis_name}
•	Robots
•	Public Analysis / Report: 
index,follow
•	Runtime / Private: 
noindex,nofollow
•	NoIndex khi
•	Analysis runtime. 
•	Private Report. 
•	Restricted Result. 
•	Missing Geometry. 
•	Không đủ dữ liệu công khai. 
 
•	16.4 Structured Data
Entity	Schema
Analysis Result	Dataset
Report	Report
Planning Region	Place
Legal Document	CreativeWork
•	Không inject Schema nếu
•	Runtime. 
•	Restricted. 
•	Private. 
•	Missing Geometry. 
•	Draft. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Map Preview / Analysis Preview
OG Type	article
•	Share Preview
Hiển thị:
•	Tên phân tích. 
•	Phạm vi phân tích. 
•	Kết quả tóm tắt. 
•	Map Preview. 
 
•	16.6 AI Summary
•	Có AI Summary
Yes, với kết quả phân tích công khai.
•	Dữ liệu nguồn
•	Analysis Result. 
•	Planning Information Model. 
•	GIS Layer công khai. 
•	Legal Document. 
•	Planning Region. 
•	Parcel liên quan. 
•	Nội dung AI Summary
•	Tóm tắt phạm vi phân tích. 
•	Tỷ lệ / mức độ ảnh hưởng. 
•	Lớp quy hoạch chính. 
•	Cảnh báo chính. 
•	Văn bản pháp lý liên quan. 
•	Không public
•	Kết quả phân tích riêng tư. 
•	Polygon cá nhân. 
•	Layer Restricted. 
•	Internal Note. 
•	Dữ liệu trả phí chưa mở công khai. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	analysis-report.xml
•	Tham gia khi
•	Public. 
•	Active. 
•	Có URL chuẩn. 
•	Có nội dung phân tích đủ điều kiện. 
•	Có Metadata hợp lệ. 
•	Loại bỏ khi
•	Runtime. 
•	Restricted. 
•	Private. 
•	Deleted. 
•	Missing Geometry. 
 
•	16.8 Internal Link
•	Link đến
•	Entity nguồn. 
•	Planning Region. 
•	Planning Project. 
•	Legal Document. 
•	GIS Layer công khai. 
•	Report liên quan. 
•	Link nhận từ
•	Trang bản đồ. 
•	Chi tiết thửa đất. 
•	Chi tiết khu quy hoạch. 
•	Snapshot. 
•	Report. 
•	Thư viện quy hoạch. 
 
•	16.9 Security & Visibility
•	Public
•	Metadata. 
•	Tóm tắt phân tích. 
•	Map Preview. 
•	Kết quả công khai. 
•	Restricted
•	Dữ liệu trả phí. 
•	Layer đối tác. 
•	Kết quả phân tích nâng cao. 
•	Private
•	Session phân tích. 
•	Polygon cá nhân. 
•	Lịch sử thao tác. 
•	Ghi chú nội bộ. 
•	Permission Rule
Không đưa dữ liệu Restricted hoặc Private vào SEO Output.
 
•	16.10 Cache & Regeneration
Event	Regenerate
analysis.updated	Metadata, AI Summary
report.updated	Metadata, Sitemap
geometry.updated	Schema, OG
layer.updated	AI Summary
legal_document.updated	AI Summary
visibility.changed	Robots, Sitemap
 
•	16.11 Acceptance Criteria
•	AC-D4.2-001
Analysis runtime luôn NoIndex.
•	AC-D4.2-002
Report / Analysis công khai có URL SEO chuẩn.
•	AC-D4.2-003
Metadata sinh đầy đủ theo Analysis Name.
•	AC-D4.2-004
Structured Data hợp lệ.
•	AC-D4.2-005
AI Summary chỉ dùng dữ liệu công khai.
•	AC-D4.2-006
Không lộ polygon cá nhân hoặc session phân tích.
•	AC-D4.2-007
Sitemap đúng điều kiện.
 
•	16.12 QA Checklist
•	URL
□ Runtime Analysis NoIndex
□ Public Analysis URL đúng chuẩn
□ Canonical đúng
•	Metadata
□ Title tồn tại
□ Description tồn tại
□ H1 tồn tại
□ Robots đúng
•	Structured Data
□ Schema Dataset / Report hợp lệ
□ Không inject khi runtime
•	Open Graph
□ OG Image đúng Preview
□ Share Preview đúng dữ liệu công khai
•	AI Summary
□ Có AI Summary cho Analysis công khai
□ Không dùng layer Restricted
•	Sitemap
□ Public Analysis có trong Sitemap khi đủ điều kiện
□ Runtime bị loại khỏi Sitemap
•	Security
□ Không lộ polygon cá nhân
□ Không lộ session
□ Không lộ dữ liệu Restricted
 
D.4.3 Màn hình thông tin quy hoạch
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes
SEO Class	SEO-A
Entity chính	Planning Information
Entity phụ	Parcel, Planning Region, Planning Project, Legal Document
URL Type	Detail URL
Mục tiêu SEO	Cung cấp trang thông tin quy hoạch chuẩn hóa, có thể Index, phục vụ Google Search và AI Search
•	Ghi chú
•	Đây là màn hình có giá trị SEO cao. 
•	Nội dung phải được sinh từ Planning Information Model. 
•	Không cho phép frontend tự diễn giải từ layer thô để tạo nội dung SEO. 
 
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/thong-tin/{entity-slug}
•	URL phụ
/quy-hoach/thong-tin/{administrative-unit-slug}/{entity-slug}
•	URL NoIndex
/quy-hoach/thong-tin?selection={id}
/quy-hoach/thong-tin?tab={tab}
/quy-hoach/thong-tin?layer={layer}
•	Canonical
•	Nếu thông tin gắn với Parcel: canonical về URL thông tin Parcel. 
•	Nếu gắn với Planning Region: canonical về URL Planning Region. 
•	Nếu gắn với Planning Project: canonical về URL Planning Project. 
•	Không tạo URL SEO riêng cho
•	Tab. 
•	Layer state. 
•	Filter. 
•	Runtime selection. 
•	View mode. 
•	16.3 Metadata
•	Title
Thông tin quy hoạch {entity_name} | QH Pro
•	Description
Thông tin quy hoạch tổng hợp của {entity_name}, bao gồm lớp quy hoạch chính, trạng thái pháp lý, đồ án liên quan, văn bản pháp lý và cảnh báo dữ liệu.
•	H1
Thông tin quy hoạch {entity_name}
•	Robots
index,follow
•	NoIndex khi
•	Không có Entity nguồn. 
•	Thông tin chưa đủ dữ liệu. 
•	Restricted. 
•	Draft. 
•	Deleted. 
•	Dữ liệu không có trạng thái pháp lý rõ ràng. 
 
•	16.4 Structured Data
Entity	Schema
Planning Information	Dataset / CreativeWork
Parcel	Place
Planning Region	Place
Planning Project	CreativeWork
Legal Document	Legislation / CreativeWork
•	Không inject Schema nếu
•	Missing Entity. 
•	Restricted. 
•	Draft. 
•	Missing Legal Status. 
•	Không đủ dữ liệu công khai. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Map Preview của Entity
OG Type	article
•	Share Preview
Hiển thị:
•	Tên Entity. 
•	Tóm tắt thông tin quy hoạch. 
•	Trạng thái pháp lý. 
•	Map Preview. 
 
•	16.6 AI Summary
•	Có AI Summary
Yes
•	Dữ liệu nguồn
•	Planning Information Model. 
•	Legal Status Aggregate. 
•	Primary Layer. 
•	Secondary Layers công khai. 
•	Legal Document. 
•	Planning Project. 
•	Administrative Unit. 
•	Nội dung AI Summary
•	Tóm tắt quy hoạch chính. 
•	Tóm tắt lớp phụ có ý nghĩa. 
•	Trạng thái pháp lý tổng hợp. 
•	Cảnh báo dữ liệu. 
•	Đồ án và văn bản liên quan. 
•	Không public
•	Raw Layer chưa normalize. 
•	Layer Restricted. 
•	Internal Note. 
•	Private Data. 
•	Dữ liệu chưa công bố. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	planning-information.xml
•	Tham gia khi
•	Public. 
•	Active. 
•	Có Entity nguồn. 
•	Có thông tin quy hoạch hợp lệ. 
•	Có Metadata hợp lệ. 
•	Loại bỏ khi
•	Restricted. 
•	Draft. 
•	Deleted. 
•	Missing Entity. 
•	Missing Legal Status. 
•	Không có dữ liệu công khai. 
 
•	16.8 Internal Link
•	Link đến
•	Entity nguồn. 
•	Administrative Unit. 
•	Planning Region. 
•	Planning Project. 
•	Legal Document. 
•	GIS Layer công khai. 
•	Report liên quan. 
•	Snapshot liên quan. 
•	Link nhận từ
•	Trang bản đồ. 
•	Popup nhanh. 
•	Chi tiết thửa đất. 
•	Tra cứu tờ/thửa. 
•	Tra cứu địa chỉ. 
•	Report. 
•	Snapshot. 
•	Thư viện quy hoạch. 
 
•	16.9 Security & Visibility
•	Public
•	Metadata. 
•	Tóm tắt quy hoạch. 
•	Legal Status công khai. 
•	Map Preview. 
•	Văn bản công khai. 
•	Restricted
•	Layer trả phí. 
•	Layer đối tác. 
•	Phân tích chuyên sâu. 
•	Dữ liệu chưa mở công khai. 
•	Private
•	Internal Note. 
•	Raw Engine Output. 
•	Audit Log. 
•	User Context. 
•	Permission Rule
Không expose raw engine output, restricted layer hoặc private data trong bất kỳ SEO Output nào.
 
•	16.10 Cache & Regeneration
Event	Regenerate
planning_information.updated	Metadata, AI Summary
legal_status.changed	Metadata, AI Summary, Robots
primary_layer.changed	Metadata, AI Summary
secondary_layers.changed	AI Summary
legal_document.updated	AI Summary
entity.visibility.changed	Robots, Sitemap
map_preview.updated	Open Graph
 
•	16.11 Acceptance Criteria
•	AC-D4.3-001
Thông tin quy hoạch công khai có URL chuẩn.
•	AC-D4.3-002
Canonical đúng Entity nguồn.
•	AC-D4.3-003
Metadata sinh từ Planning Information Model.
•	AC-D4.3-004
Structured Data hợp lệ.
•	AC-D4.3-005
AI Summary không dùng raw layer chưa chuẩn hóa.
•	AC-D4.3-006
Không lộ dữ liệu Restricted hoặc Private.
•	AC-D4.3-007
Sitemap đúng điều kiện.
 
•	16.12 QA Checklist
•	URL
□ URL đúng chuẩn
□ Canonical đúng
□ Runtime selection NoIndex
•	Metadata
□ Title đúng Entity
□ Description đúng Planning Information
□ H1 tồn tại
□ Robots đúng
•	Structured Data
□ Schema hợp lệ
□ Entity Mapping đúng
□ Không inject khi Missing Legal Status
•	Open Graph
□ OG Title đúng
□ OG Image đúng Map Preview
•	AI Summary
□ Sinh từ Planning Information Model
□ Không dùng raw layer chưa normalize
□ Không lộ Restricted Data
•	Sitemap
□ Có trong planning-information.xml khi đủ điều kiện
□ Loại khỏi Sitemap đúng điều kiện
•	Security
□ Không lộ raw engine output
□ Không lộ internal note
□ Không lộ dữ liệu Restricted
 
D.5	QH.5 – So sánh & Biến động
D.5.1 So sánh quy hoạch
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes, có điều kiện
SEO Class	SEO-C
Entity chính	Planning Comparison
Entity phụ	Planning Region, Planning Project, Legal Document, GIS Layer
URL Type	Compare URL
Mục tiêu SEO	SEO các kết quả so sánh quy hoạch công khai, có giá trị tham chiếu độc lập
•	Ghi chú
•	Không SEO phiên so sánh tạm thời do người dùng tạo. 
•	Chỉ SEO khi kết quả so sánh được lưu, công khai và có đủ dữ liệu giải thích. 
•	Nếu kết quả so sánh chỉ bổ trợ cho một đồ án hoặc vùng quy hoạch, canonical về Entity nguồn. 
 
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/so-sanh/{comparison-slug}
•	URL NoIndex
/quy-hoach/so-sanh?from={version-a}&to={version-b}
/quy-hoach/so-sanh?layers={layer-list}
/quy-hoach/so-sanh?session={session-id}
•	Canonical
•	So sánh công khai có giá trị độc lập: canonical về URL so sánh. 
•	So sánh gắn với đồ án: canonical về URL chi tiết đồ án. 
•	So sánh gắn với vùng quy hoạch: canonical về URL vùng quy hoạch. 
•	So sánh runtime: canonical về /quy-hoach. 
•	Không tạo URL SEO riêng cho
•	Layer state. 
•	Runtime version selection. 
•	Temporary compare session. 
•	Filter. 
•	Sort. 
•	Tab. 
•	Map state. 
 
•	16.3 Metadata
•	Title
So sánh quy hoạch {comparison_name} | QH Pro
•	Description
So sánh các phiên bản, lớp dữ liệu và thay đổi quy hoạch liên quan đến {comparison_name}, bao gồm phạm vi, nội dung khác biệt và văn bản pháp lý liên quan.
•	H1
So sánh quy hoạch {comparison_name}
•	Robots
•	Kết quả công khai đủ điều kiện: 
index,follow
•	Phiên so sánh tạm thời: 
noindex,nofollow
•	NoIndex khi
•	Runtime compare. 
•	Private. 
•	Restricted. 
•	Thiếu dữ liệu nguồn. 
•	Không có khác biệt đáng kể. 
•	Không đủ nội dung mô tả. 
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Dataset / CreativeWork
Entity Mapping	Planning Comparison
•	Schema bổ sung
Schema	Điều kiện
Place	So sánh gắn với vùng quy hoạch
Legislation / CreativeWork	Có văn bản pháp lý liên quan
BreadcrumbList	Có cấu trúc điều hướng
•	Không inject Schema nếu
•	Runtime. 
•	Draft. 
•	Restricted. 
•	Missing Source Version. 
•	Missing Legal Reference. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Compare Preview / Map Diff Preview
OG Type	article
•	Share Preview
Hiển thị:
•	Tên khu vực / đồ án. 
•	Phiên bản so sánh. 
•	Tóm tắt thay đổi chính. 
•	Ảnh bản đồ so sánh. 
 
•	16.6 AI Summary
•	Có AI Summary
Yes, với kết quả so sánh công khai.
•	Dữ liệu nguồn
•	Planning Comparison. 
•	Planning Region. 
•	Planning Project. 
•	GIS Layer công khai. 
•	Legal Document. 
•	Version / Lineage. 
•	Nội dung AI Summary
•	Đối tượng được so sánh. 
•	Phiên bản trước / sau. 
•	Các thay đổi chính. 
•	Ý nghĩa quy hoạch. 
•	Văn bản pháp lý liên quan. 
•	Cảnh báo dữ liệu nếu có. 
•	Không public
•	Layer Restricted. 
•	Dữ liệu chưa công bố. 
•	Internal Note. 
•	Raw diff nội bộ. 
•	Phiên bản đối tác không công khai. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	planning-comparison.xml
•	Tham gia khi
•	Public. 
•	Active. 
•	Có URL chuẩn. 
•	Có dữ liệu so sánh hợp lệ. 
•	Có AI Summary / mô tả đủ điều kiện. 
•	Có Entity nguồn hợp lệ. 
•	Loại bỏ khi
•	Runtime. 
•	Draft. 
•	Restricted. 
•	Private. 
•	Missing Source Version. 
•	Deleted. 
 
•	16.8 Internal Link
•	Link đến
•	Planning Region. 
•	Planning Project. 
•	Legal Document. 
•	GIS Layer công khai. 
•	Timeline pháp lý. 
•	Quan hệ phiên bản. 
•	Report liên quan. 
•	Link nhận từ
•	Chi tiết đồ án. 
•	Chi tiết vùng quy hoạch. 
•	Timeline pháp lý. 
•	Quan hệ phiên bản. 
•	Report. 
•	Snapshot. 
•	Thư viện quy hoạch. 
 
•	16.9 Security & Visibility
•	Public
•	Tên so sánh. 
•	Tóm tắt thay đổi. 
•	Map Diff Preview. 
•	Văn bản công khai. 
•	Layer công khai. 
•	Restricted
•	Layer trả phí. 
•	Dữ liệu đối tác. 
•	Phiên bản chưa công bố. 
•	Kết quả phân tích chuyên sâu. 
•	Private
•	Compare session cá nhân. 
•	Ghi chú nội bộ. 
•	Raw diff nội bộ. 
•	Audit log. 
•	Permission Rule
Không đưa dữ liệu Restricted hoặc Private vào Metadata, Schema, OG, AI Summary hoặc Sitemap.

•	16.10 Cache & Regeneration
Event	Regenerate
comparison.updated	Metadata, AI Summary
planning_version.updated	Metadata, AI Summary
layer.updated	AI Summary, OG
legal_document.updated	AI Summary
visibility.changed	Robots, Sitemap
map_diff_preview.updated	Open Graph
 
•	16.11 Acceptance Criteria
•	AC-D5.1-001
Runtime compare luôn NoIndex.
•	AC-D5.1-002
Kết quả so sánh công khai đủ điều kiện có URL chuẩn.
•	AC-D5.1-003
Canonical đúng về URL so sánh hoặc Entity nguồn.
•	AC-D5.1-004
Metadata mô tả đúng đối tượng và phiên bản so sánh.
•	AC-D5.1-005
AI Summary chỉ dùng dữ liệu công khai.
•	AC-D5.1-006
Sitemap chỉ chứa kết quả so sánh đủ điều kiện.
•	AC-D5.1-007
Không lộ dữ liệu Restricted hoặc Private.
 
•	16.12 QA Checklist
•	URL
□ Runtime compare NoIndex
□ Public Compare URL đúng chuẩn
□ Canonical đúng
□ Không sinh URL từ session / filter / layer state
•	Metadata
□ Title tồn tại
□ Description tồn tại
□ H1 tồn tại
□ Robots đúng
•	Structured Data
□ Schema hợp lệ
□ Entity Mapping đúng Planning Comparison
□ Không inject khi thiếu version nguồn
•	Open Graph
□ OG Title đúng
□ OG Image đúng Map Diff Preview
•	AI Summary
□ Tóm tắt đúng nội dung so sánh
□ Không dùng Layer Restricted
□ Không lộ raw diff nội bộ
•	Sitemap
□ Compare đủ điều kiện có trong planning-comparison.xml
□ Runtime / Private / Restricted bị loại khỏi Sitemap
•	Security
□ Không lộ phiên bản chưa công bố
□ Không lộ dữ liệu đối tác
□ Không lộ compare session cá nhân
 
D.5.2 So sánh dữ liệu
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes, có điều kiện
SEO Class	SEO-C
Entity chính	Data Comparison
Entity phụ	GIS Layer, Planning Region, Parcel, Legal Document
URL Type	Compare / Report URL
Mục tiêu SEO	SEO các kết quả so sánh dữ liệu công khai có giá trị phân tích, không SEO thao tác so sánh runtime
•	Ghi chú
•	So sánh dữ liệu thường là nghiệp vụ phân tích, không phải mọi kết quả đều SEO. 
•	Chỉ SEO khi kết quả được chuẩn hóa thành Report hoặc trang so sánh công khai. 
•	Các thao tác so sánh tạm thời mặc định NoIndex. 
 
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/so-sanh-du-lieu/{comparison-slug}
hoặc:
/quy-hoach/bao-cao/{report-slug}
•	URL NoIndex
/quy-hoach/so-sanh-du-lieu?dataset_a={id}&dataset_b={id}
/quy-hoach/so-sanh-du-lieu?session={session-id}
/quy-hoach/so-sanh-du-lieu?layers={layer-list}
•	Canonical
•	So sánh dữ liệu công khai: canonical về URL Compare. 
•	Nếu xuất thành Report công khai: canonical về URL Report. 
•	Runtime comparison: canonical về /quy-hoach. 
•	Không tạo URL SEO riêng cho
•	Dataset selection runtime. 
•	Layer selection. 
•	Filter. 
•	Sort. 
•	Pagination tạm thời. 
•	Session. 
 
•	16.3 Metadata
•	Title
So sánh dữ liệu quy hoạch {comparison_name} | QH Pro
•	Description
So sánh dữ liệu quy hoạch, lớp GIS, phạm vi áp dụng, nguồn dữ liệu và các khác biệt chính liên quan đến {comparison_name}.
•	H1
So sánh dữ liệu {comparison_name}
•	Robots
•	Public Compare / Report: 
index,follow
•	Runtime: 
noindex,nofollow
•	NoIndex khi
•	Runtime session. 
•	Dữ liệu nguồn Restricted. 
•	Dữ liệu thiếu metadata. 
•	So sánh không có giá trị độc lập. 
•	Dữ liệu chưa công bố. 
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Dataset / Report
Entity Mapping	Data Comparison
•	Schema bổ sung
Schema	Điều kiện
CreativeWork	Có báo cáo giải thích
BreadcrumbList	Có điều hướng đầy đủ
•	Không inject Schema nếu
•	Runtime. 
•	Missing Dataset. 
•	Dataset Restricted. 
•	Missing Metadata. 
•	Private Report. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Data Compare Preview
OG Type	article
•	Share Preview
Hiển thị:
•	Tên bộ dữ liệu so sánh. 
•	Phạm vi so sánh. 
•	Tóm tắt khác biệt. 
•	Preview biểu đồ / bản đồ nếu có. 
 
•	16.6 AI Summary
•	Có AI Summary
Yes, với kết quả công khai.
•	Dữ liệu nguồn
•	Data Comparison. 
•	GIS Layer công khai. 
•	Layer Metadata. 
•	Planning Region. 
•	Legal Document. 
•	Report nếu có. 
•	Nội dung AI Summary
•	Dữ liệu được so sánh. 
•	Nguồn dữ liệu. 
•	Phạm vi so sánh. 
•	Khác biệt chính. 
•	Độ tin cậy / cảnh báo dữ liệu. 
•	Văn bản liên quan. 
•	Không public
•	Raw dataset nội bộ. 
•	Layer Restricted. 
•	Metadata kỹ thuật. 
•	Import log. 
•	Validation error nội bộ. 
•	Partner data chưa công khai. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	data-comparison.xml / report.xml
•	Tham gia khi
•	Public. 
•	Active. 
•	Có URL chuẩn. 
•	Có dữ liệu nguồn công khai. 
•	Có nội dung so sánh đủ điều kiện. 
•	Loại bỏ khi
•	Runtime. 
•	Restricted. 
•	Private. 
•	Missing Dataset. 
•	Missing Metadata. 
•	Deleted. 
 
•	16.8 Internal Link
•	Link đến
•	GIS Layer nguồn. 
•	Layer Metadata. 
•	Planning Region. 
•	Planning Project. 
•	Legal Document. 
•	Report liên quan. 
•	Link nhận từ
•	Layer Public. 
•	Layer Metadata. 
•	Màn hình phân tích quy hoạch. 
•	Report. 
•	Snapshot. 
•	Thư viện quy hoạch. 
 
•	16.9 Security & Visibility
•	Public
•	Tóm tắt so sánh. 
•	Metadata công khai. 
•	Layer công khai. 
•	Preview công khai. 
•	Restricted
•	Layer trả phí. 
•	Dataset đối tác. 
•	Dữ liệu chưa kiểm duyệt. 
•	Báo cáo phân tích chuyên sâu. 
•	Private
•	Raw dataset nội bộ. 
•	Import log. 
•	Validation log. 
•	Compare session cá nhân. 
•	Permission Rule
Không đưa Raw Dataset, Layer Restricted hoặc Metadata nội bộ vào SEO Output.
 
•	16.10 Cache & Regeneration
Event	Regenerate
data_comparison.updated	Metadata, AI Summary
layer.metadata.updated	Metadata, AI Summary
layer.visibility.changed	Robots, Sitemap
dataset.updated	AI Summary
report.updated	Metadata, Sitemap
preview.updated	Open Graph
 
•	16.11 Acceptance Criteria
•	AC-D5.2-001
Runtime data comparison luôn NoIndex.
•	AC-D5.2-002
So sánh dữ liệu công khai đủ điều kiện có URL chuẩn.
•	AC-D5.2-003
Canonical đúng Compare hoặc Report URL.
•	AC-D5.2-004
Metadata sinh từ dữ liệu công khai.
•	AC-D5.2-005
Structured Data hợp lệ.
•	AC-D5.2-006
Không lộ raw dataset hoặc metadata nội bộ.
•	AC-D5.2-007
Sitemap đúng điều kiện.
 
•	16.12 QA Checklist
•	URL
□ Runtime comparison NoIndex
□ Public Compare URL đúng chuẩn
□ Canonical đúng
□ Không Index URL chứa dataset query
•	Metadata
□ Title tồn tại
□ Description tồn tại
□ H1 tồn tại
□ Robots đúng
•	Structured Data
□ Schema Dataset / Report hợp lệ
□ Không inject khi dataset Restricted
•	Open Graph
□ OG Preview đúng dữ liệu công khai
□ Không lộ dataset nội bộ
•	AI Summary
□ Tóm tắt đúng dữ liệu nguồn
□ Không lộ import log / validation log
□ Không dùng Layer Restricted
•	Sitemap
□ Compare / Report đủ điều kiện có trong Sitemap
□ Runtime / Private bị loại khỏi Sitemap
•	Security
□ Không lộ raw dataset
□ Không lộ metadata nội bộ
□ Không lộ dữ liệu đối tác hạn chế
D.5.3 Biến động quy hoạch
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes
SEO Class	SEO-A / SEO-C
Entity chính	Planning Change
Entity phụ	Planning Region, Planning Project, Legal Document, GIS Layer
URL Type	Change / Timeline / Report URL
Mục tiêu SEO	SEO các biến động quy hoạch công khai, phục vụ Google Search, AI Search, Historical SEO và Knowledge Graph
•	Ghi chú
•	Đây là nhóm nội dung có giá trị SEO cao. 
•	Biến động công khai, có căn cứ pháp lý, có tác động đáng kể có thể xếp SEO-A. 
•	Biến động nhỏ hoặc nội dung bổ trợ xếp SEO-C. 
•	Runtime diff hoặc biến động chưa xác minh mặc định NoIndex. 
 
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/bien-dong/{change-slug}
•	URL theo đồ án
/quy-hoach/do-an/{project-slug}/bien-dong/{change-slug}
•	URL NoIndex
/quy-hoach/bien-dong?from={version-a}&to={version-b}
/quy-hoach/bien-dong?layer={layer-id}
/quy-hoach/bien-dong?session={session-id}
•	Canonical
•	Biến động công khai độc lập: canonical về URL biến động. 
•	Biến động là một phần của đồ án: canonical về URL chi tiết đồ án hoặc Timeline pháp lý. 
•	Biến động là kết quả diff runtime: canonical về Entity nguồn. 
•	Không tạo URL SEO riêng cho
•	Runtime diff. 
•	Version selector. 
•	Layer state. 
•	Filter. 
•	Map state. 
•	Preview mode. 
 
•	16.3 Metadata
•	Title
Biến động quy hoạch {change_name} | QH Pro
•	Description
Thông tin biến động quy hoạch {change_name}, bao gồm nội dung thay đổi, phạm vi ảnh hưởng, phiên bản trước sau và văn bản pháp lý liên quan.
•	H1
Biến động quy hoạch {change_name}
•	Robots
•	Biến động công khai: 
index,follow
•	Biến động chưa xác minh / runtime: 
noindex,nofollow
•	NoIndex khi
•	Chưa xác minh. 
•	Không có văn bản pháp lý. 
•	Runtime diff. 
•	Restricted. 
•	Private. 
•	Missing Source Version. 
•	Missing Target Version. 
 
•	16.4 Structured Data
Entity	Schema
Planning Change	CreativeWork / Dataset
Planning Project	CreativeWork
Legal Document	Legislation / CreativeWork
Planning Region	Place
•	Schema bổ sung
Schema	Điều kiện
Event	Nếu biến động được mô tả như sự kiện pháp lý
BreadcrumbList	Có điều hướng đầy đủ
•	Không inject Schema nếu
•	Change chưa xác minh. 
•	Missing Legal Document. 
•	Restricted. 
•	Runtime diff. 
•	Missing Version. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Change Preview / Map Diff Preview
OG Type	article
•	Share Preview
Hiển thị:
•	Tên biến động. 
•	Phạm vi ảnh hưởng. 
•	Phiên bản trước / sau. 
•	Văn bản pháp lý liên quan. 
•	Map Diff Preview. 
 
•	16.6 AI Summary
•	Có AI Summary
Yes
•	Dữ liệu nguồn
•	Planning Change. 
•	Planning Project. 
•	Planning Region. 
•	Legal Document. 
•	GIS Layer công khai. 
•	Version / Lineage. 
•	Nội dung AI Summary
•	Tóm tắt biến động. 
•	Nội dung thay đổi. 
•	Phạm vi ảnh hưởng. 
•	Căn cứ pháp lý. 
•	Phiên bản trước / sau. 
•	Tác động chính. 
•	Cảnh báo dữ liệu nếu có. 
•	Không public
•	Biến động chưa xác minh. 
•	Draft change. 
•	Restricted Layer. 
•	Internal Note. 
•	Raw diff nội bộ. 
•	Dữ liệu chưa công bố. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	planning-change.xml
•	Tham gia khi
•	Public. 
•	Verified. 
•	Có URL chuẩn. 
•	Có căn cứ pháp lý. 
•	Có Entity nguồn. 
•	Có Metadata hợp lệ. 
•	Loại bỏ khi
•	Draft. 
•	Unverified. 
•	Restricted. 
•	Private. 
•	Runtime diff. 
•	Deleted. 
•	Missing Legal Document. 
 
•	16.8 Internal Link
•	Link đến
•	Planning Project. 
•	Planning Region. 
•	Legal Document. 
•	Timeline pháp lý. 
•	Quan hệ phiên bản. 
•	GIS Layer liên quan. 
•	Report liên quan. 
•	Link nhận từ
•	Chi tiết đồ án. 
•	Timeline pháp lý. 
•	Quan hệ phiên bản. 
•	So sánh quy hoạch. 
•	Thư viện quy hoạch. 
•	Report. 
•	Snapshot. 
 
•	16.9 Security & Visibility
•	Public
•	Tên biến động. 
•	Tóm tắt thay đổi. 
•	Văn bản pháp lý công khai. 
•	Map Diff Preview. 
•	Phiên bản công khai. 
•	Restricted
•	Biến động chưa công bố. 
•	Layer trả phí. 
•	Dữ liệu đối tác. 
•	Phân tích tác động chuyên sâu. 
•	Private
•	Internal Note. 
•	Raw diff. 
•	Audit Log. 
•	Draft version. 
•	Permission Rule
Không đưa biến động chưa xác minh hoặc dữ liệu Restricted / Private vào SEO Output.
 
•	16.10 Cache & Regeneration
Event	Regenerate
planning_change.created	Metadata, Sitemap
planning_change.updated	Metadata, AI Summary
planning_change.verified	Robots, Sitemap
legal_document.updated	AI Summary
version_lineage.updated	AI Summary, Internal Link
map_diff_preview.updated	Open Graph
visibility.changed	Robots, Sitemap
 
•	16.11 Acceptance Criteria
•	AC-D5.3-001
Biến động công khai, đã xác minh có URL SEO chuẩn.
•	AC-D5.3-002
Biến động runtime hoặc chưa xác minh luôn NoIndex.
•	AC-D5.3-003
Canonical đúng về URL biến động hoặc Entity nguồn.
•	AC-D5.3-004
Metadata mô tả đúng nội dung biến động.
•	AC-D5.3-005
AI Summary có căn cứ pháp lý và không dùng dữ liệu chưa công bố.
•	AC-D5.3-006
Biến động đủ điều kiện xuất hiện trong planning-change.xml.
•	AC-D5.3-007
Không lộ raw diff hoặc draft version.
 
•	16.12 QA Checklist
•	URL
□ URL biến động đúng chuẩn
□ Runtime diff NoIndex
□ Canonical đúng
□ Không sinh URL từ version selector
•	Metadata
□ Title tồn tại
□ Description nêu đúng biến động
□ H1 tồn tại
□ Robots đúng
•	Structured Data
□ Schema hợp lệ
□ Không inject nếu thiếu Legal Document
□ Entity Mapping đúng Planning Change
•	Open Graph
□ OG Image đúng Map Diff Preview
□ Share Preview đúng dữ liệu công khai
•	AI Summary
□ Tóm tắt đúng biến động
□ Có căn cứ pháp lý
□ Không lộ dữ liệu chưa công bố
□ Không lộ raw diff nội bộ
•	Sitemap
□ Biến động verified có trong planning-change.xml
□ Draft / Unverified / Restricted bị loại khỏi Sitemap
•	Security
□ Không lộ draft version
□ Không lộ internal note
□ Không lộ Layer Restricted
□ Permission hoạt động đúng

 
D.6	QH.6 – Theo dõi & Cảnh báo
D.6.1 Danh sách theo dõi
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	No
SEO Class	SEO-N
Entity chính	Watchlist
Entity phụ	Parcel, Planning Region, Planning Project
URL Type	Internal / User URL
Mục tiêu SEO	Không SEO; chỉ phục vụ quản lý danh sách theo dõi cá nhân của người dùng
•	Ghi chú
•	Danh sách theo dõi là dữ liệu cá nhân hóa. 
•	Không phải Landing Page. 
•	Không phải nội dung công khai. 
•	Không tham gia Google Search, AI Search, Sitemap hoặc Discover. 
 
•	16.2 URL & Canonical
•	URL chuẩn
Không áp dụng cho SEO.
•	URL NoIndex
/tai-khoan/theo-doi
/quy-hoach/theo-doi
/watchlist
•	Canonical
Không áp dụng.
•	Không tạo URL SEO riêng cho
•	Danh sách theo dõi cá nhân. 
•	Bộ lọc danh sách theo dõi. 
•	Trạng thái đã đọc / chưa đọc. 
•	Nhóm theo dõi. 
•	Sắp xếp theo thời gian. 
•	Đối tượng theo dõi riêng tư. 
 
•	16.3 Metadata
•	Title
Không sinh Metadata SEO riêng.
•	Description
Không sinh Description SEO riêng.
•	H1
Kế thừa từ màn hình tài khoản / chức năng người dùng.
•	Robots
noindex,nofollow
•	Ghi chú
Không đưa tên tài sản, thửa đất hoặc vùng theo dõi của người dùng vào Metadata.
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Không áp dụng
Entity Mapping	Không áp dụng
•	Điều kiện áp dụng
Không áp dụng.
•	Điều kiện loại bỏ
Luôn loại bỏ.
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Không áp dụng
OG Description	Không áp dụng
OG Image	Không áp dụng
•	Share Preview
Không hỗ trợ chia sẻ danh sách theo dõi cá nhân.
 
•	16.6 AI Summary
•	Có AI Summary
No
•	Dữ liệu nguồn
Không áp dụng.
•	Không public
•	Danh sách đối tượng theo dõi. 
•	Thửa đất người dùng quan tâm. 
•	Vùng quy hoạch người dùng theo dõi. 
•	Lịch sử theo dõi. 
•	Tần suất truy cập. 
•	Ghi chú cá nhân. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	Không áp dụng
•	Tham gia khi
Không có trường hợp nào.
•	Loại bỏ khỏi Sitemap khi
Luôn loại bỏ.
 
•	16.8 Internal Link
•	Link đến
•	Trang chi tiết thửa đất. 
•	Trang chi tiết khu quy hoạch. 
•	Trang chi tiết đồ án quy hoạch. 
•	Trang chi tiết cảnh báo. 
•	Link nhận từ
•	Tài khoản người dùng. 
•	Menu cá nhân. 
•	Notification Center. 
•	Ghi chú
Liên kết trong danh sách theo dõi chỉ phục vụ người dùng đã đăng nhập, không phải Internal Link SEO công khai.
 
•	16.9 Security & Visibility
•	Public
Không có dữ liệu public cho SEO.
•	Restricted
•	Thông tin đối tượng người dùng theo dõi. 
•	Trạng thái thông báo. 
•	Private
•	Danh sách theo dõi cá nhân. 
•	Ghi chú cá nhân. 
•	Lịch sử theo dõi. 
•	Hành vi người dùng. 
•	Permission Rule
Chỉ chủ tài khoản hoặc người được phân quyền mới được xem danh sách theo dõi.
Không đưa bất kỳ dữ liệu nào của Watchlist vào:
•	Metadata. 
•	Structured Data. 
•	Open Graph. 
•	AI Summary. 
•	Sitemap. 
 
•	16.10 Cache & Regeneration
Event	Regenerate
watchlist.updated	User Cache
watched_entity.updated	Notification Data
visibility.changed	Permission Cache
•	Ghi chú
Không phát sinh SEO Cache, Metadata, Sitemap hoặc Schema Regeneration.
 
•	16.11 Acceptance Criteria
•	AC-D6.1-001
Danh sách theo dõi luôn NoIndex.
•	AC-D6.1-002
Không sinh Metadata SEO riêng.
•	AC-D6.1-003
Không sinh Structured Data.
•	AC-D6.1-004
Không xuất hiện trong Sitemap.
•	AC-D6.1-005
Không lộ danh sách theo dõi cá nhân.
•	AC-D6.1-006
Các liên kết trong danh sách điều hướng đúng tới Entity đích sau khi kiểm tra quyền.
 
•	16.12 QA Checklist
•	URL
□ URL Watchlist NoIndex
□ Không có Watchlist URL trong Sitemap
□ Không tạo URL SEO theo filter/sort
•	Metadata
□ Không sinh Metadata riêng
□ Robots đúng noindex,nofollow
•	Structured Data
□ Không inject Schema
•	Open Graph
□ Không sinh OG
•	AI Summary
□ Không sinh AI Summary
□ Không lộ danh sách theo dõi
•	Sitemap
□ Không xuất hiện trong Sitemap
•	Security
□ Không lộ dữ liệu cá nhân
□ Không lộ lịch sử theo dõi
□ Permission hoạt động đúng
 
D.6.2 Cảnh báo biến động
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes, có điều kiện
SEO Class	SEO-C / SEO-N
Entity chính	Planning Alert
Entity phụ	Parcel, Planning Region, Planning Change, Legal Document
URL Type	Alert / Notification URL
Mục tiêu SEO	Chỉ SEO cảnh báo biến động công khai có giá trị tham chiếu độc lập; NoIndex cảnh báo cá nhân hóa
•	Ghi chú
•	Phần lớn cảnh báo là dữ liệu cá nhân hóa, mặc định SEO-N. 
•	Chỉ cảnh báo biến động công khai, đã xác minh, có căn cứ pháp lý và có giá trị cộng đồng mới được xem xét SEO. 
•	Cảnh báo gửi riêng cho người dùng không được Index. 
 
•	16.2 URL & Canonical
•	URL chuẩn
Áp dụng cho cảnh báo công khai đủ điều kiện:
/quy-hoach/canh-bao/{alert-slug}
•	URL NoIndex
/tai-khoan/canh-bao/{alert-id}
/quy-hoach/canh-bao?user={user-id}
/quy-hoach/canh-bao?watched={entity-id}
•	Canonical
•	Cảnh báo công khai độc lập: canonical về URL cảnh báo. 
•	Cảnh báo liên quan đến biến động quy hoạch: canonical về URL biến động quy hoạch. 
•	Cảnh báo liên quan đến thửa đất: canonical về URL thửa đất. 
•	Cảnh báo cá nhân hóa: NoIndex, không canonical SEO. 
•	Không tạo URL SEO riêng cho
•	Alert cá nhân. 
•	Notification unread/read. 
•	Alert filter. 
•	Alert session. 
•	Alert theo tài khoản. 
•	Alert theo watchlist riêng. 
 
•	16.3 Metadata
•	Title
Cảnh báo công khai:
Cảnh báo biến động quy hoạch {alert_name} | QH Pro
•	Description
Thông tin cảnh báo biến động quy hoạch liên quan đến {entity_name}, bao gồm nội dung thay đổi, phạm vi ảnh hưởng và văn bản pháp lý liên quan.
•	H1
Cảnh báo biến động quy hoạch {alert_name}
•	Robots
•	Cảnh báo công khai đủ điều kiện: 
index,follow
•	Cảnh báo cá nhân / restricted: 
noindex,nofollow
•	NoIndex khi
•	Alert cá nhân. 
•	Alert chưa xác minh. 
•	Alert không có căn cứ pháp lý. 
•	Alert Restricted. 
•	Alert chỉ phục vụ notification runtime. 
•	Alert không có nội dung công khai độc lập. 
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Event / CreativeWork
Entity Mapping	Planning Alert
•	Schema bổ sung
Schema	Điều kiện
Place	Cảnh báo gắn với vùng quy hoạch / thửa đất
Legislation / CreativeWork	Có văn bản pháp lý liên quan
BreadcrumbList	Có điều hướng đầy đủ
•	Không inject Schema nếu
•	Alert cá nhân. 
•	Alert chưa xác minh. 
•	Restricted. 
•	Missing Legal Document. 
•	Notification runtime. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Map Preview / Change Preview
OG Type	article
•	Share Preview
Hiển thị:
•	Tên cảnh báo. 
•	Đối tượng bị ảnh hưởng. 
•	Tóm tắt biến động. 
•	Map Preview. 
•	Trạng thái xác minh. 
Không hiển thị dữ liệu người dùng theo dõi.
 
•	16.6 AI Summary
•	Có AI Summary
Yes, chỉ với cảnh báo công khai đủ điều kiện.
•	Dữ liệu nguồn
•	Planning Alert. 
•	Planning Change. 
•	Parcel / Planning Region liên quan. 
•	Legal Document. 
•	GIS Layer công khai. 
•	Planning Project. 
•	Nội dung AI Summary
•	Tóm tắt cảnh báo. 
•	Lý do phát sinh cảnh báo. 
•	Đối tượng bị ảnh hưởng. 
•	Căn cứ pháp lý. 
•	Mức độ tác động. 
•	Liên kết đến biến động hoặc Entity nguồn. 
•	Không public
•	Người dùng nhận cảnh báo. 
•	Watchlist cá nhân. 
•	Alert rule cá nhân. 
•	Ghi chú nội bộ. 
•	Dữ liệu chưa xác minh. 
•	Restricted Data. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	planning-alert.xml
•	Tham gia khi
•	Public. 
•	Verified. 
•	Có URL chuẩn. 
•	Có căn cứ pháp lý. 
•	Có nội dung công khai. 
•	Có Entity nguồn hợp lệ. 
•	Loại bỏ khi
•	Alert cá nhân. 
•	Draft. 
•	Unverified. 
•	Restricted. 
•	Deleted. 
•	Missing Legal Document. 
•	Notification runtime. 
 
•	16.8 Internal Link
•	Link đến
•	Planning Change. 
•	Parcel liên quan. 
•	Planning Region liên quan. 
•	Planning Project. 
•	Legal Document. 
•	Report liên quan. 
•	Link nhận từ
•	Chi tiết thửa đất. 
•	Chi tiết khu quy hoạch. 
•	Biến động quy hoạch. 
•	Report. 
•	Snapshot. 
•	Notification Center. 
•	Ghi chú
Internal Link SEO chỉ áp dụng cho cảnh báo công khai.
•	16.9 Security & Visibility
•	Public
•	Tên cảnh báo công khai. 
•	Nội dung tóm tắt. 
•	Đối tượng bị ảnh hưởng. 
•	Văn bản pháp lý công khai. 
•	Map Preview. 
•	Restricted
•	Cảnh báo chỉ dành cho nhóm người dùng. 
•	Dữ liệu đối tác. 
•	Cảnh báo chưa công bố. 
•	Private
•	Người dùng nhận cảnh báo. 
•	Watchlist cá nhân. 
•	Alert preference. 
•	Notification status. 
•	Ghi chú cá nhân. 
•	Permission Rule
Không đưa dữ liệu cá nhân hoặc dữ liệu watchlist vào bất kỳ SEO Output nào.
 
•	16.10 Cache & Regeneration
Event	Regenerate
alert.created	Metadata, Sitemap
alert.updated	Metadata, AI Summary
alert.verified	Robots, Sitemap
alert.visibility.changed	Robots, Sitemap
planning_change.updated	AI Summary
legal_document.updated	AI Summary
map_preview.updated	Open Graph
 
•	16.11 Acceptance Criteria
•	AC-D6.2-001
Alert cá nhân luôn NoIndex.
•	AC-D6.2-002
Alert công khai, verified, có căn cứ pháp lý mới được Index.
•	AC-D6.2-003
Không có dữ liệu Watchlist cá nhân trong Metadata, Schema, OG hoặc AI Summary.
•	AC-D6.2-004
Canonical đúng về Alert hoặc Entity nguồn.
•	AC-D6.2-005
Sitemap chỉ chứa Alert công khai đủ điều kiện.
•	AC-D6.2-006
Không lộ dữ liệu Restricted hoặc Private.
 
•	16.12 QA Checklist
•	URL
□ Alert cá nhân NoIndex
□ Public Alert URL đúng chuẩn
□ Canonical đúng
□ Không Index Notification runtime
•	Metadata
□ Title đúng Alert
□ Description đúng Entity liên quan
□ Robots đúng
•	Structured Data
□ Schema Event / CreativeWork hợp lệ
□ Không inject khi Alert chưa xác minh
•	Open Graph
□ OG Image đúng Preview
□ Không lộ thông tin người nhận cảnh báo
•	AI Summary
□ Tóm tắt đúng cảnh báo
□ Có căn cứ pháp lý nếu Index
□ Không lộ Watchlist cá nhân
•	Sitemap
□ Alert công khai có trong planning-alert.xml
□ Alert cá nhân / unverified bị loại khỏi Sitemap
•	Security
□ Không lộ người dùng theo dõi
□ Không lộ alert preference
□ Không lộ dữ liệu Restricted
 
D.6.3 Chi tiết cảnh báo
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes, có điều kiện
SEO Class	SEO-C / SEO-N
Entity chính	Planning Alert Detail
Entity phụ	Planning Alert, Planning Change, Parcel, Planning Region, Legal Document
URL Type	Detail URL
Mục tiêu SEO	Hiển thị chi tiết cảnh báo; chỉ SEO bản chi tiết công khai, không SEO chi tiết cảnh báo cá nhân
•	Ghi chú
•	Chi tiết cảnh báo cá nhân mặc định SEO-N. 
•	Chi tiết cảnh báo công khai có thể SEO-C nếu có giá trị độc lập. 
•	Nếu cảnh báo chỉ là bản rút gọn của biến động quy hoạch, canonical về URL biến động. 
 
•	16.2 URL & Canonical
•	URL chuẩn
Cảnh báo công khai:
/quy-hoach/canh-bao/{alert-slug}/chi-tiet
•	URL NoIndex
/tai-khoan/canh-bao/{alert-id}/chi-tiet
/quy-hoach/canh-bao/{alert-id}?notification={notification-id}
/quy-hoach/canh-bao/{alert-id}?user={user-id}
•	Canonical
•	Nếu cảnh báo có URL công khai độc lập: canonical về URL chi tiết cảnh báo. 
•	Nếu cảnh báo thuộc Planning Change: canonical về URL biến động quy hoạch. 
•	Nếu cảnh báo thuộc Parcel: canonical về URL Parcel. 
•	Nếu cảnh báo cá nhân: NoIndex. 
•	Không tạo URL SEO riêng cho
•	Notification ID. 
•	User ID. 
•	Read / unread state. 
•	Alert status cá nhân. 
•	Reminder state. 
•	Alert channel. 
 
•	16.3 Metadata
•	Title
Chi tiết cảnh báo quy hoạch {alert_name} | QH Pro
•	Description
Chi tiết cảnh báo quy hoạch liên quan đến {entity_name}, bao gồm nội dung cảnh báo, phạm vi ảnh hưởng, thời điểm phát sinh và căn cứ pháp lý liên quan.
•	H1
Chi tiết cảnh báo {alert_name}
•	Robots
•	Công khai đủ điều kiện: 
index,follow
•	Cá nhân / restricted: 
noindex,nofollow
•	NoIndex khi
•	Alert cá nhân. 
•	Alert chưa xác minh. 
•	Alert thiếu nội dung công khai. 
•	Alert thiếu căn cứ pháp lý. 
•	Alert Restricted. 
•	Detail chỉ phục vụ notification runtime. 
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Event / CreativeWork
Entity Mapping	Planning Alert Detail
•	Schema bổ sung
Schema	Điều kiện
Place	Có đối tượng địa lý liên quan
Legislation / CreativeWork	Có văn bản pháp lý
BreadcrumbList	Có điều hướng đầy đủ
•	Không inject Schema nếu
•	Alert cá nhân. 
•	Missing Legal Reference. 
•	Restricted. 
•	Unverified. 
•	Notification runtime. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Alert Preview / Map Preview
OG Type	article
•	Share Preview
Hiển thị:
•	Tên cảnh báo. 
•	Nội dung chính. 
•	Entity bị ảnh hưởng. 
•	Map Preview. 
•	Căn cứ pháp lý nếu có. 
Không hiển thị người nhận cảnh báo hoặc trạng thái cá nhân.
 
•	16.6 AI Summary
•	Có AI Summary
Yes, với cảnh báo công khai đủ điều kiện.
•	Dữ liệu nguồn
•	Planning Alert Detail. 
•	Planning Alert. 
•	Planning Change. 
•	Parcel. 
•	Planning Region. 
•	Legal Document. 
•	Planning Project. 
•	Nội dung AI Summary
•	Tóm tắt nội dung cảnh báo. 
•	Lý do cảnh báo. 
•	Đối tượng bị ảnh hưởng. 
•	Phạm vi ảnh hưởng. 
•	Căn cứ pháp lý. 
•	Hướng dẫn xem thêm Entity nguồn. 
•	Không public
•	User ID. 
•	Notification ID. 
•	Watchlist cá nhân. 
•	Alert preference. 
•	Ghi chú cá nhân. 
•	Restricted Data. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	planning-alert.xml
•	Tham gia khi
•	Alert Public. 
•	Verified. 
•	Có URL chuẩn. 
•	Có nội dung chi tiết công khai. 
•	Có căn cứ pháp lý hoặc Entity nguồn rõ ràng. 
•	Loại bỏ khi
•	Personal Alert. 
•	Notification runtime. 
•	Unverified. 
•	Restricted. 
•	Private. 
•	Deleted. 
•	Missing Public Content. 
 
•	16.8 Internal Link
•	Link đến
•	Alert tổng quan. 
•	Planning Change. 
•	Parcel liên quan. 
•	Planning Region liên quan. 
•	Planning Project. 
•	Legal Document. 
•	Report liên quan. 
•	Link nhận từ
•	Danh sách cảnh báo. 
•	Notification Center. 
•	Biến động quy hoạch. 
•	Chi tiết thửa đất. 
•	Chi tiết vùng quy hoạch. 
•	Report. 
•	Snapshot. 
•	Ghi chú
Internal Link SEO chỉ áp dụng với cảnh báo công khai.
 
•	16.9 Security & Visibility
•	Public
•	Nội dung cảnh báo công khai. 
•	Đối tượng bị ảnh hưởng. 
•	Map Preview. 
•	Văn bản pháp lý công khai. 
•	Restricted
•	Cảnh báo nhóm người dùng. 
•	Cảnh báo đối tác. 
•	Nội dung chưa công bố. 
•	Private
•	Notification ID. 
•	User ID. 
•	Watchlist. 
•	Read state. 
•	Personal reminder. 
•	Alert preference. 
•	Permission Rule
Không đưa bất kỳ thông tin cá nhân hóa nào vào SEO Output.
 
•	16.10 Cache & Regeneration
Event	Regenerate
alert_detail.updated	Metadata, AI Summary
alert.verified	Robots, Sitemap
alert.visibility.changed	Robots, Sitemap
legal_document.updated	AI Summary
planning_change.updated	AI Summary
map_preview.updated	Open Graph
alert.deleted	Sitemap
 
•	16.11 Acceptance Criteria
•	AC-D6.3-001
Chi tiết cảnh báo cá nhân luôn NoIndex.
•	AC-D6.3-002
Chi tiết cảnh báo public phải không chứa dữ liệu người dùng.
•	AC-D6.3-003
Canonical đúng về URL cảnh báo hoặc Entity nguồn.
•	AC-D6.3-004
Metadata không chứa Notification ID, User ID hoặc dữ liệu cá nhân hóa.
•	AC-D6.3-005
AI Summary không lộ Watchlist, preference hoặc read state.
•	AC-D6.3-006
Sitemap chỉ chứa cảnh báo công khai đủ điều kiện.

•	16.12 QA Checklist
•	URL
□ Personal Alert Detail NoIndex
□ Notification URL NoIndex
□ Public Alert Detail URL đúng chuẩn
□ Canonical đúng
•	Metadata
□ Title đúng Alert
□ Description không chứa dữ liệu cá nhân
□ Robots đúng
•	Structured Data
□ Schema hợp lệ
□ Không inject khi Alert unverified
•	Open Graph
□ OG Preview đúng dữ liệu công khai
□ Không lộ user/watchlist
•	AI Summary
□ Có AI Summary nếu Alert public
□ Không lộ Notification ID
□ Không lộ Watchlist cá nhân
•	Sitemap
□ Public verified alert có trong Sitemap
□ Personal / notification runtime bị loại khỏi Sitemap
•	Security
□ Không lộ User ID
□ Không lộ read state
□ Không lộ alert preference
□ Không lộ dữ liệu Restricted

 
D.7	QH.7 – Snapshot & Report
•	D.7.1 Snapshot
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes, có điều kiện
SEO Class	SEO-D / SEO-C
Entity chính	Snapshot
Entity phụ	Parcel, Planning Region, Planning Project, Report
URL Type	Snapshot URL
Mục tiêu SEO	Lưu và chia sẻ trạng thái dữ liệu quy hoạch tại một thời điểm; chỉ SEO khi Snapshot công khai, ổn định và có giá trị tham chiếu độc lập
•	Ghi chú
•	Snapshot cá nhân mặc định SEO-N / NoIndex. 
•	Snapshot chia sẻ mặc định SEO-D / NoIndex. 
•	Chỉ Snapshot công khai, có nội dung đủ giá trị và không chứa dữ liệu cá nhân mới được xem xét SEO-C. 
 
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/snapshot/{snapshot-slug}
•	URL NoIndex
/quy-hoach/snapshot/{snapshot-id}?token={share-token}
/quy-hoach/snapshot?session={session-id}
/quy-hoach/snapshot?user={user-id}
•	Canonical
•	Snapshot công khai có giá trị độc lập: canonical về URL Snapshot chuẩn. 
•	Snapshot chỉ chụp lại một Entity: canonical về URL Entity nguồn. 
•	Snapshot cá nhân / token share: NoIndex, canonical về Entity nguồn nếu xác định được. 
•	Không tạo URL SEO riêng cho
•	Token share. 
•	Session. 
•	User ID. 
•	Map state. 
•	Layer toggle. 
•	Zoom / center. 
•	Snapshot draft. 
 
•	16.3 Metadata
•	Title
Snapshot quy hoạch {snapshot_name} | QH Pro
•	Description
Ảnh chụp trạng thái dữ liệu quy hoạch tại {snapshot_name}, bao gồm bản đồ, lớp dữ liệu, phạm vi và các thông tin liên quan tại thời điểm tạo Snapshot.
•	H1
Snapshot quy hoạch {snapshot_name}
•	Robots
•	Snapshot công khai đủ điều kiện: 
index,follow
•	Snapshot cá nhân / chia sẻ bằng token: 
noindex,follow
•	Snapshot private / restricted: 
noindex,nofollow
•	NoIndex khi
•	Snapshot cá nhân. 
•	Snapshot chỉ chia sẻ bằng token. 
•	Snapshot thiếu Entity nguồn. 
•	Snapshot chứa dữ liệu Restricted. 
•	Snapshot không có nội dung công khai độc lập. 
•	Snapshot hết hạn. 
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	CreativeWork
Entity Mapping	Snapshot
•	Schema bổ sung
Schema	Điều kiện
Place	Snapshot gắn với Parcel / Planning Region
Dataset	Snapshot gắn với lớp dữ liệu công khai
BreadcrumbList	Có điều hướng đầy đủ
•	Không inject Schema nếu
•	Snapshot private. 
•	Snapshot share-only. 
•	Snapshot expired. 
•	Snapshot thiếu Entity nguồn. 
•	Snapshot chứa dữ liệu Restricted. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Snapshot Preview
OG Type	article
•	Share Preview
Hiển thị:
•	Tên Snapshot. 
•	Địa bàn / Entity nguồn. 
•	Ảnh preview bản đồ. 
•	Thời điểm tạo Snapshot. 
•	Cảnh báo nếu Snapshot chỉ mang tính tham khảo. 
 
•	16.6 AI Summary
•	Có AI Summary
Yes, chỉ với Snapshot công khai đủ điều kiện.
•	Dữ liệu nguồn
•	Snapshot. 
•	Entity nguồn. 
•	Parcel. 
•	Planning Region. 
•	Planning Project. 
•	GIS Layer công khai. 
•	Legal Document liên quan. 
•	Nội dung AI Summary
•	Tóm tắt Snapshot. 
•	Thời điểm dữ liệu được ghi nhận. 
•	Entity nguồn. 
•	Lớp dữ liệu hiển thị. 
•	Nội dung quy hoạch chính tại thời điểm Snapshot. 
•	Không public
•	User ID. 
•	Share Token. 
•	Session. 
•	Ghi chú cá nhân. 
•	Layer Restricted. 
•	Dữ liệu Private. 
•	Dữ liệu đối tác chưa công khai. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	snapshot.xml
•	Tham gia khi
•	Snapshot Public. 
•	Active. 
•	Có URL chuẩn. 
•	Có Entity nguồn. 
•	Có preview công khai. 
•	Có nội dung đủ giá trị độc lập. 
•	Loại bỏ khi
•	Private. 
•	Share-only. 
•	Restricted. 
•	Expired. 
•	Deleted. 
•	Missing Entity Source. 
•	Token-based URL. 
 
•	16.8 Internal Link
•	Link đến
•	Entity nguồn. 
•	Parcel liên quan. 
•	Planning Region. 
•	Planning Project. 
•	Legal Document. 
•	Report liên quan. 
•	Trang bản đồ quy hoạch. 
•	Link nhận từ
•	Trang chi tiết Entity. 
•	Report. 
•	Share Snapshot. 
•	Trang bản đồ. 
•	Màn hình phân tích. 
•	Thư viện quy hoạch. 
 
•	16.9 Security & Visibility
•	Public
•	Tên Snapshot. 
•	Preview công khai. 
•	Entity nguồn. 
•	Thời điểm tạo. 
•	Dữ liệu quy hoạch công khai. 
•	Restricted
•	Snapshot có Layer trả phí. 
•	Snapshot có dữ liệu đối tác. 
•	Snapshot có nội dung chưa công bố. 
•	Private
•	Snapshot cá nhân. 
•	Ghi chú người dùng. 
•	User ID. 
•	Share Token. 
•	Session. 
•	Permission Rule
Không đưa dữ liệu cá nhân, token hoặc Layer Restricted vào Metadata, Schema, OG, AI Summary hoặc Sitemap.
 
•	16.10 Cache & Regeneration
Event	Regenerate
snapshot.created	Metadata, OG
snapshot.updated	Metadata, AI Summary
snapshot.visibility.changed	Robots, Sitemap
snapshot.expired	Robots, Sitemap
source_entity.updated	AI Summary
map_preview.updated	Open Graph
legal_document.updated	AI Summary

•	16.11 Acceptance Criteria
•	AC-D7.1-001
Snapshot cá nhân luôn NoIndex.
•	AC-D7.1-002
Snapshot share bằng token không xuất hiện trong Sitemap.
•	AC-D7.1-003
Snapshot public đủ điều kiện có URL chuẩn.
•	AC-D7.1-004
Canonical đúng về Snapshot hoặc Entity nguồn.
•	AC-D7.1-005
Metadata không chứa User ID, token hoặc dữ liệu cá nhân.
•	AC-D7.1-006
AI Summary chỉ dùng dữ liệu công khai.
•	AC-D7.1-007
Không lộ Layer Restricted hoặc dữ liệu Private.
 
•	16.12 QA Checklist
•	URL
□ Snapshot URL đúng chuẩn
□ Token URL NoIndex
□ Snapshot private NoIndex
□ Canonical đúng
•	Metadata
□ Title tồn tại
□ Description tồn tại
□ H1 tồn tại
□ Robots đúng trạng thái
•	Structured Data
□ Schema CreativeWork hợp lệ
□ Không inject khi Snapshot private/share-only
•	Open Graph
□ OG Image đúng Snapshot Preview
□ Share Preview không lộ dữ liệu riêng tư
•	AI Summary
□ Có AI Summary với Snapshot public
□ Không lộ User ID / token / session
□ Không dùng Layer Restricted
•	Sitemap
□ Snapshot public đủ điều kiện có trong snapshot.xml
□ Snapshot private/share-only/expired bị loại khỏi Sitemap
•	Security
□ Không lộ Share Token
□ Không lộ ghi chú cá nhân
□ Không lộ dữ liệu Restricted
 
•	D.7.2 Share Snapshot
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes, giới hạn
SEO Class	SEO-D
Entity chính	Share Snapshot
Entity phụ	Snapshot, Parcel, Planning Region, Planning Project
URL Type	Share URL
Mục tiêu SEO	Hỗ trợ chia sẻ Snapshot trên mạng xã hội, ứng dụng nhắn tin và Mobile Deep Link; không ưu tiên Index
•	Ghi chú
•	Share Snapshot phục vụ Discover & Share, không phải SEO Ranking. 
•	Mặc định NoIndex. 
•	Ưu tiên Open Graph, preview và bảo mật token. 
 
•	16.2 URL & Canonical
•	URL chuẩn
/share/snapshot/{share-id}
•	URL NoIndex
Toàn bộ Share Snapshot URL mặc định NoIndex:
/share/snapshot/{share-id}?token={token}
/share/snapshot/{share-id}?utm={source}
•	Canonical
•	Nếu Snapshot public: canonical về URL Snapshot chuẩn. 
•	Nếu Snapshot chỉ share bằng token: NoIndex, không đưa vào Sitemap. 
•	Nếu Snapshot gắn với Entity nguồn: canonical về Entity nguồn khi cần. 
•	Không tạo URL SEO riêng cho
•	Share token. 
•	UTM. 
•	Referrer. 
•	Device. 
•	App open state. 
•	Preview session. 
 
•	16.3 Metadata
•	Title
Chia sẻ Snapshot quy hoạch {snapshot_name} | QH Pro
•	Description
Xem Snapshot quy hoạch được chia sẻ từ QH Pro, bao gồm bản đồ preview và thông tin quy hoạch công khai liên quan.
•	H1
Snapshot quy hoạch được chia sẻ
•	Robots
noindex,follow
•	NoIndex khi
Luôn áp dụng cho Share Snapshot URL, trừ trường hợp được chuyển đổi thành Snapshot public chính thức.
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Không áp dụng mặc định
Entity Mapping	Snapshot
•	Điều kiện áp dụng
Chỉ inject Schema tại URL Snapshot public chuẩn, không inject tại Share URL token.
•	Không inject Schema nếu
•	URL chứa token. 
•	Share-only. 
•	Private. 
•	Restricted. 
•	Expired. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Snapshot Preview
OG Type	article
•	Share Preview
Bắt buộc hiển thị:
•	Snapshot Preview. 
•	Tên Snapshot. 
•	Entity nguồn nếu public. 
•	Cảnh báo “Thông tin chỉ có giá trị tham khảo” nếu cần. 
Không hiển thị:
•	User ID. 
•	Token. 
•	Ghi chú cá nhân. 
•	Layer Restricted. 
 
•	16.6 AI Summary
•	Có AI Summary
No, đối với Share URL.
•	Dữ liệu nguồn
Không áp dụng.
•	Ghi chú
AI Summary chỉ áp dụng tại Snapshot public chuẩn hoặc Entity nguồn.
•	Không public
•	Token. 
•	Người chia sẻ. 
•	Người nhận. 
•	Session. 
•	Ghi chú cá nhân. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	Không áp dụng
•	Tham gia khi
Không có trường hợp nào.
•	Loại bỏ khỏi Sitemap khi
Luôn loại bỏ.
 
•	16.8 Internal Link
•	Link đến
•	Snapshot public nếu có. 
•	Entity nguồn. 
•	Trang bản đồ. 
•	Report liên quan. 
•	Link nhận từ
•	Mobile App. 
•	Web Share. 
•	Email. 
•	Zalo / Messenger / Social Share. 
•	Ghi chú
Share URL không phải Internal Link SEO chính thức.
 
•	16.9 Security & Visibility
•	Public
•	Preview an toàn. 
•	Tên Snapshot nếu được phép. 
•	Entity nguồn công khai. 
•	Restricted
•	Snapshot có Layer trả phí. 
•	Snapshot có dữ liệu chưa công bố. 
•	Private
•	Token. 
•	User ID. 
•	Recipient. 
•	Share Log. 
•	Personal Note. 
•	Permission Rule
Share Preview phải dùng Safe Preview. Không hiển thị dữ liệu vượt quá quyền truy cập của người mở link.
 
•	16.10 Cache & Regeneration
Event	Regenerate
share.created	OG Preview
share.expired	Robots
snapshot.visibility.changed	OG Preview, Robots
snapshot.updated	OG Preview
map_preview.updated	Open Graph
 
•	16.11 Acceptance Criteria
•	AC-D7.2-001
Share Snapshot URL luôn NoIndex.
•	AC-D7.2-002
Share Snapshot không xuất hiện trong Sitemap.
•	AC-D7.2-003
Open Graph hiển thị đúng Safe Preview.
•	AC-D7.2-004
Không lộ token, User ID hoặc ghi chú cá nhân.
•	AC-D7.2-005
Canonical đúng về Snapshot public hoặc Entity nguồn nếu có.
 
•	16.12 QA Checklist
•	URL
□ Share URL NoIndex
□ Token URL NoIndex
□ Canonical đúng
•	Metadata
□ Title tồn tại
□ Description tồn tại
□ Robots noindex,follow
•	Structured Data
□ Không inject Schema tại Share URL
•	Open Graph
□ OG Image đúng Safe Preview
□ Không lộ dữ liệu riêng tư
•	Sitemap
□ Không xuất hiện trong Sitemap
•	Security
□ Không lộ token
□ Không lộ người chia sẻ/người nhận
□ Không lộ ghi chú cá nhân
□ Permission kiểm soát đúng preview
 
•	D.7.3 Report
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes, có điều kiện
SEO Class	SEO-A / SEO-C
Entity chính	Report
Entity phụ	Parcel, Planning Region, Planning Project, Legal Document, GIS Layer
URL Type	Report URL
Mục tiêu SEO	SEO các báo cáo công khai có giá trị phân tích, tham chiếu và AI Search
•	Ghi chú
•	Report công khai, có nội dung đầy đủ và dữ liệu nguồn rõ ràng có thể xếp SEO-A. 
•	Report cá nhân, nội bộ, trả phí hoặc chưa công bố mặc định NoIndex. 
•	Report là nhóm nội dung có giá trị SEO cao nếu được chuẩn hóa tốt. 
 
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/bao-cao/{report-slug}
•	URL phụ
/quy-hoach/bao-cao/{administrative-unit-slug}/{report-slug}
•	URL NoIndex
/quy-hoach/bao-cao/{report-id}?draft=true
/quy-hoach/bao-cao/{report-id}?token={share-token}
/quy-hoach/bao-cao/preview/{report-id}
•	Canonical
•	Report public: canonical về URL Report chuẩn. 
•	Report sinh từ Parcel: có Internal Link đến Parcel, không canonical về Parcel nếu Report có giá trị độc lập. 
•	Report chỉ là bản xuất tạm: NoIndex hoặc canonical về Entity nguồn. 
•	Không tạo URL SEO riêng cho
•	Preview. 
•	Draft. 
•	Export format. 
•	Share token. 
•	Filter. 
•	Tab. 
•	Print view. 
 
•	16.3 Metadata
•	Title
Báo cáo quy hoạch {report_name} | QH Pro
•	Description
Báo cáo phân tích quy hoạch {report_name}, bao gồm phạm vi dữ liệu, kết quả phân tích, văn bản pháp lý, lớp dữ liệu liên quan và các cảnh báo chính.
•	H1
Báo cáo quy hoạch {report_name}
•	Robots
•	Report public: 
index,follow
•	Report private/draft/share-token: 
noindex,nofollow
•	NoIndex khi
•	Draft. 
•	Private. 
•	Restricted. 
•	Share-only. 
•	Missing Source Entity. 
•	Missing Public Content. 
•	Report hết hạn. 
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Report / CreativeWork
Entity Mapping	Report
•	Schema bổ sung
Schema	Điều kiện
Dataset	Có dữ liệu phân tích công khai
Place	Report gắn với địa bàn / vùng
Legislation / CreativeWork	Có văn bản pháp lý liên quan
BreadcrumbList	Có điều hướng đầy đủ
•	Không inject Schema nếu
•	Draft. 
•	Restricted. 
•	Private. 
•	Missing Source Entity. 
•	Missing Public Data. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Report Preview / Map Preview
OG Type	article
•	Share Preview
Hiển thị:
•	Tên Report. 
•	Phạm vi báo cáo. 
•	Tóm tắt kết quả chính. 
•	Map Preview / Report Cover. 
•	Cảnh báo dữ liệu nếu có. 
 
•	16.6 AI Summary
•	Có AI Summary
Yes, với Report public.
•	Dữ liệu nguồn
•	Report. 
•	Parcel. 
•	Planning Region. 
•	Planning Project. 
•	Legal Document. 
•	GIS Layer công khai. 
•	Analysis Result. 
•	Nội dung AI Summary
•	Tóm tắt phạm vi Report. 
•	Kết quả phân tích chính. 
•	Căn cứ pháp lý. 
•	Dữ liệu nguồn. 
•	Cảnh báo quan trọng. 
•	Entity liên quan. 
•	Không public
•	Report private. 
•	Report draft. 
•	Dữ liệu trả phí chưa mở. 
•	Internal Note. 
•	Raw Analysis Output. 
•	User Context. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	report.xml
•	Tham gia khi
•	Public. 
•	Active. 
•	Có URL chuẩn. 
•	Có nội dung đủ giá trị. 
•	Có nguồn dữ liệu rõ ràng. 
•	Có Metadata hợp lệ. 
•	Loại bỏ khi
•	Draft. 
•	Private. 
•	Restricted. 
•	Share-only. 
•	Deleted. 
•	Expired. 
•	Missing Public Content. 
 
•	16.8 Internal Link
•	Link đến
•	Entity nguồn. 
•	Administrative Unit. 
•	Parcel. 
•	Planning Region. 
•	Planning Project. 
•	Legal Document. 
•	GIS Layer. 
•	Snapshot liên quan. 
•	Link nhận từ
•	Màn hình phân tích. 
•	Chi tiết thửa đất. 
•	Chi tiết khu quy hoạch. 
•	Snapshot. 
•	Thư viện quy hoạch. 
•	Share Report. 
•	Trang bản đồ. 
 
•	16.9 Security & Visibility
•	Public
•	Report public. 
•	Metadata. 
•	AI Summary. 
•	Report Preview. 
•	Dữ liệu nguồn công khai. 
•	Restricted
•	Report trả phí. 
•	Report đối tác. 
•	Report có lớp dữ liệu hạn chế. 
•	Phân tích chuyên sâu chưa mở. 
•	Private
•	Report cá nhân. 
•	Draft. 
•	User Context. 
•	Ghi chú nội bộ. 
•	Raw Engine Output. 
•	Permission Rule
Không đưa Report Restricted / Private vào SEO Output. Report chỉ được Index khi toàn bộ nội dung public đã được kiểm tra.
 
•	16.10 Cache & Regeneration
Event	Regenerate
report.created	Metadata, Sitemap
report.updated	Metadata, AI Summary
report.visibility.changed	Robots, Sitemap
report.expired	Robots, Sitemap
source_entity.updated	AI Summary
legal_document.updated	AI Summary
report_preview.updated	Open Graph
 
•	16.11 Acceptance Criteria
•	AC-D7.3-001
Report public đủ điều kiện có URL SEO chuẩn.
•	AC-D7.3-002
Report draft/private/share-only luôn NoIndex.
•	AC-D7.3-003
Metadata sinh đầy đủ theo Report.
•	AC-D7.3-004
Structured Data hợp lệ.
•	AC-D7.3-005
AI Summary chỉ sử dụng dữ liệu công khai.
•	AC-D7.3-006
Report đủ điều kiện xuất hiện trong report.xml.
•	AC-D7.3-007
Không lộ raw engine output, user context hoặc dữ liệu Restricted.
 
•	16.12 QA Checklist
•	URL
□ Report URL đúng chuẩn
□ Draft / Preview / Share-token NoIndex
□ Canonical đúng
•	Metadata
□ Title tồn tại
□ Description tồn tại
□ H1 tồn tại
□ Robots đúng trạng thái
•	Structured Data
□ Schema Report / CreativeWork hợp lệ
□ Không inject khi Report private
•	Open Graph
□ OG Image đúng Report Preview
□ Share Preview không lộ dữ liệu hạn chế
•	AI Summary
□ AI Summary đúng nội dung Report
□ Không dùng raw engine output
□ Không lộ user context
•	Sitemap
□ Report public có trong report.xml
□ Report private/draft/share-only bị loại khỏi Sitemap
•	Security
□ Không lộ Report cá nhân
□ Không lộ dữ liệu trả phí
□ Không lộ ghi chú nội bộ
 
•	D.7.4 Share Report
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes, giới hạn
SEO Class	SEO-D
Entity chính	Share Report
Entity phụ	Report, Parcel, Planning Region, Planning Project
URL Type	Share URL
Mục tiêu SEO	Hỗ trợ chia sẻ Report an toàn, có preview tốt, không ưu tiên Index
•	Ghi chú
•	Share Report phục vụ chia sẻ, không phải Landing Page SEO. 
•	Mặc định NoIndex. 
•	Nếu Report công khai có URL chuẩn, Share URL canonical về URL Report. 
 
•	16.2 URL & Canonical
•	URL chuẩn
/share/report/{share-id}
•	URL NoIndex
/share/report/{share-id}?token={token}
/share/report/{share-id}?utm={source}
/share/report/{share-id}?viewer={viewer-id}
•	Canonical
•	Report public: canonical về URL Report chuẩn. 
•	Report private/share-token: NoIndex, không Sitemap. 
•	Report expired: NoIndex. 
•	Không tạo URL SEO riêng cho
•	Token. 
•	UTM. 
•	Viewer. 
•	Device. 
•	App state. 
•	Download state. 
•	Print state. 
 
•	16.3 Metadata
•	Title
Chia sẻ báo cáo quy hoạch {report_name} | QH Pro
•	Description
Xem báo cáo quy hoạch được chia sẻ từ QH Pro, bao gồm tóm tắt phân tích, phạm vi dữ liệu và bản xem trước an toàn.
•	H1
Báo cáo quy hoạch được chia sẻ
•	Robots
noindex,follow
•	NoIndex khi
Áp dụng mặc định cho toàn bộ Share Report URL.
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Không áp dụng mặc định
Entity Mapping	Report
•	Điều kiện áp dụng
Structured Data chỉ inject tại URL Report public chuẩn, không inject tại Share URL token.
•	Không inject Schema nếu
•	Share URL. 
•	Token URL. 
•	Private Report. 
•	Restricted Report. 
•	Expired Share. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Report Preview / Safe Preview
OG Type	article
•	Share Preview
Hiển thị:
•	Tên Report. 
•	Tóm tắt an toàn. 
•	Ảnh preview. 
•	Nguồn QH Pro. 
Không hiển thị:
•	Token. 
•	Người chia sẻ. 
•	Người nhận. 
•	Nội dung Restricted. 
•	Ghi chú nội bộ. 
 
•	16.6 AI Summary
•	Có AI Summary
No, đối với Share URL.
•	Dữ liệu nguồn
Không áp dụng.
•	Ghi chú
AI Summary chỉ áp dụng tại URL Report public chuẩn.
•	Không public
•	Token. 
•	Viewer ID. 
•	Người chia sẻ. 
•	Người nhận. 
•	Share Log. 
•	Report private. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	Không áp dụng
•	Tham gia khi
Không có trường hợp nào.
•	Loại bỏ khỏi Sitemap khi
Luôn loại bỏ.
 
•	16.8 Internal Link
•	Link đến
•	Report public nếu có. 
•	Entity nguồn. 
•	Trang bản đồ. 
•	Snapshot liên quan. 
•	Link nhận từ
•	Web Share. 
•	Mobile App. 
•	Email. 
•	Zalo / Messenger / Social Share. 
•	Ghi chú
Share URL không được tính là Internal Link SEO chính thức.
 
•	16.9 Security & Visibility
•	Public
•	Safe Preview. 
•	Tên Report nếu được phép. 
•	Tóm tắt public. 
•	Restricted
•	Report trả phí. 
•	Report đối tác. 
•	Báo cáo chưa công bố. 
•	Private
•	Token. 
•	Viewer ID. 
•	Share Log. 
•	Người chia sẻ / người nhận. 
•	Ghi chú cá nhân. 
•	Permission Rule
Share Preview phải kiểm tra quyền trước khi hiển thị. Không được preview nội dung vượt quyền truy cập.
 
•	16.10 Cache & Regeneration
Event	Regenerate
share.created	OG Preview
share.expired	Robots
report.visibility.changed	OG Preview, Robots
report.updated	OG Preview
report_preview.updated	Open Graph
 
•	16.11 Acceptance Criteria
•	AC-D7.4-001
Share Report URL luôn NoIndex.
•	AC-D7.4-002
Share Report không xuất hiện trong Sitemap.
•	AC-D7.4-003
Open Graph hiển thị Safe Preview.
•	AC-D7.4-004
Canonical về Report public nếu tồn tại.
•	AC-D7.4-005
Không lộ token, viewer, người chia sẻ hoặc nội dung Restricted.
 
•	16.12 QA Checklist
•	URL
□ Share URL NoIndex
□ Token URL NoIndex
□ Canonical đúng Report public nếu có
•	Metadata
□ Title tồn tại
□ Description tồn tại
□ Robots noindex,follow
•	Structured Data
□ Không inject Schema tại Share URL
•	Open Graph
□ OG Image đúng Safe Preview
□ Không lộ dữ liệu hạn chế
•	Sitemap
□ Không xuất hiện trong Sitemap
•	Security
□ Không lộ token
□ Không lộ viewer ID
□ Không lộ người chia sẻ/người nhận
□ Không lộ Report private/restricted

 
D.8	QH.8 – Tài khoản & Gói dịch vụ
D.8.1 Các màn hình SEO-N
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	No
SEO Class	SEO-N
Entity chính	User Account / Subscription / Payment
Entity phụ	User Profile, Billing, Access Permission
URL Type	Internal / User URL
Mục tiêu SEO	Không SEO; chỉ phục vụ xác thực, tài khoản, gói dịch vụ, thanh toán và phân quyền người dùng
•	Ghi chú
•	QH.8 là nhóm chức năng tài khoản và dịch vụ cá nhân hóa. 
•	Không phải Landing Page SEO. 
•	Không phải nội dung công khai. 
•	Không tham gia Google Search, AI Search, Sitemap hoặc Knowledge Graph công khai. 
•	Nếu sau này có trang giới thiệu bảng giá công khai, landing page gói dịch vụ hoặc trang marketing riêng thì phải tách thành nhóm SEO khác, không thuộc phạm vi D.8.1. 
 
•	16.2 URL & Canonical
•	URL chuẩn
Không áp dụng cho SEO.
•	URL NoIndex
Các nhóm URL sau mặc định NoIndex:
/dang-nhap
/dang-ky
/tai-khoan
/tai-khoan/ho-so
/tai-khoan/goi-dich-vu
/tai-khoan/thanh-toan
/tai-khoan/lich-su-thanh-toan
/tai-khoan/quyen-truy-cap
/tai-khoan/cai-dat
•	Canonical
Không áp dụng cho SEO.
•	Không tạo URL SEO riêng cho
•	Hồ sơ người dùng. 
•	Thông tin tài khoản. 
•	Gói dịch vụ đang sử dụng. 
•	Lịch sử thanh toán. 
•	Hóa đơn. 
•	Phân quyền truy cập. 
•	Cài đặt cá nhân. 
•	Trạng thái đăng nhập. 
•	Session. 
•	Token. 
•	Callback thanh toán. 
•	Duplicate Control
Mọi URL thuộc khu vực tài khoản và thanh toán phải được đánh dấu NoIndex, không tham gia Sitemap và không sinh Canonical SEO độc lập.
 
•	16.3 Metadata
•	Title
Không sinh Metadata SEO riêng.
Có thể dùng Title kỹ thuật phục vụ trải nghiệm người dùng, ví dụ:
Tài khoản của tôi | QH Pro
nhưng không dùng để SEO.
•	Description
Không sinh Description SEO riêng.
•	H1
Theo nội dung màn hình tài khoản, không dùng làm H1 SEO.
•	Robots
noindex,nofollow
•	Ghi chú
Không đưa thông tin cá nhân, gói dịch vụ, lịch sử thanh toán hoặc trạng thái tài khoản vào Metadata.

•	16.4 Structured Data
Thành phần	Giá trị
Schema	Không áp dụng
Entity Mapping	Không áp dụng
•	Điều kiện áp dụng
Không áp dụng.
•	Điều kiện loại bỏ
Luôn loại bỏ đối với toàn bộ màn hình QH.8.
•	Ghi chú
Không inject Schema cho:
•	User Profile. 
•	Subscription. 
•	Payment. 
•	Invoice. 
•	Billing. 
•	Access Permission. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Không áp dụng
OG Description	Không áp dụng
OG Image	Không áp dụng
•	Share Preview
Không hỗ trợ Share Preview cho màn hình tài khoản, gói dịch vụ hoặc thanh toán cá nhân.
•	Ghi chú
Không được tạo Open Graph chứa dữ liệu tài khoản hoặc thông tin giao dịch.
 
•	16.6 AI Summary
•	Có AI Summary
No
•	Dữ liệu nguồn
Không áp dụng.
•	Không public
•	Thông tin tài khoản. 
•	Email / số điện thoại người dùng. 
•	Hồ sơ cá nhân. 
•	Gói dịch vụ đang sử dụng. 
•	Lịch sử thanh toán. 
•	Hóa đơn. 
•	Quyền truy cập. 
•	Token. 
•	Session. 
•	Dữ liệu xác thực. 
•	Ghi chú
Không đưa bất kỳ dữ liệu QH.8 nào vào AI Summary công khai hoặc nội dung đọc hiểu bằng máy phục vụ AI Search.
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	Không áp dụng
•	Tham gia Sitemap khi
Không có trường hợp nào.
•	Loại bỏ khỏi Sitemap khi
Luôn loại bỏ toàn bộ URL thuộc QH.8.
 
•	16.8 Internal Link
•	Link đến
Các màn hình QH.8 có thể điều hướng nội bộ tới:
•	Trang chủ. 
•	Trang bản đồ quy hoạch. 
•	Trang chi tiết Entity đã lưu. 
•	Danh sách theo dõi. 
•	Báo cáo cá nhân. 
•	Snapshot cá nhân. 
•	Link nhận từ
•	Menu tài khoản. 
•	Header user menu. 
•	Notification. 
•	Payment callback. 
•	Subscription flow. 
•	Ghi chú
Các liên kết này là liên kết chức năng, không phải Internal Link SEO công khai.
Không tính các link trong khu vực tài khoản vào kiến trúc Internal Linking SEO.
 
•	16.9 Security & Visibility
•	Public
Không có dữ liệu public SEO trong QH.8.
•	Restricted
•	Thông tin gói dịch vụ. 
•	Trạng thái quyền truy cập. 
•	Trạng thái thanh toán. 
•	Lịch sử giao dịch. 
•	Private
•	Hồ sơ người dùng. 
•	Email. 
•	Số điện thoại. 
•	Token. 
•	Session. 
•	Hóa đơn. 
•	Phương thức thanh toán. 
•	Lịch sử đăng nhập. 
•	Cài đặt cá nhân. 
•	Permission Rule
Toàn bộ màn hình QH.8 yêu cầu kiểm tra đăng nhập và phân quyền.
Không đưa dữ liệu Restricted hoặc Private vào:
•	Metadata. 
•	Structured Data. 
•	Open Graph. 
•	AI Summary. 
•	Sitemap. 
•	Public HTML không cần xác thực. 
 
•	16.10 Cache & Regeneration
Event	Regenerate
user.updated	User Cache
subscription.updated	User Permission Cache
payment.completed	Billing Cache
payment.failed	Billing Cache
permission.changed	Access Cache
session.expired	Auth Cache
•	Ghi chú
Không phát sinh:
•	SEO Metadata Regeneration. 
•	Sitemap Regeneration. 
•	Schema Regeneration. 
•	AI Summary Regeneration. 
•	Open Graph Regeneration. 
 
•	16.11 Acceptance Criteria
•	AC-D8.1-001
Toàn bộ URL tài khoản, gói dịch vụ và thanh toán phải NoIndex.
•	AC-D8.1-002
Không URL nào thuộc QH.8 được xuất hiện trong Sitemap.
•	AC-D8.1-003
Không sinh Structured Data cho User, Subscription, Payment hoặc Billing.
•	AC-D8.1-004
Không sinh AI Summary cho dữ liệu tài khoản.
•	AC-D8.1-005
Không sinh Open Graph chứa dữ liệu người dùng hoặc giao dịch.
•	AC-D8.1-006
Không lộ thông tin cá nhân, token, session, hóa đơn hoặc trạng thái thanh toán trong HTML công khai.
•	AC-D8.1-007
Các màn hình QH.8 chỉ truy cập được sau khi kiểm tra quyền phù hợp.
 
•	16.12 QA Checklist
•	URL
□ URL tài khoản NoIndex
□ URL gói dịch vụ NoIndex
□ URL thanh toán NoIndex
□ Callback thanh toán không Index
□ Không URL QH.8 nào có trong Sitemap
•	Metadata
□ Không sinh Metadata SEO riêng
□ Robots đúng noindex,nofollow
□ Không chứa email / số điện thoại / thông tin thanh toán trong Metadata
•	Structured Data
□ Không inject Schema
□ Không có User Schema
□ Không có Payment / Invoice Schema công khai
•	Open Graph
□ Không sinh OG riêng
□ Không có Share Preview tài khoản
□ Không lộ dữ liệu giao dịch trong OG
•	AI Summary
□ Không sinh AI Summary
□ Không đưa dữ liệu tài khoản vào nội dung public machine-readable
•	Sitemap
□ Không xuất hiện trong bất kỳ Sitemap nào
□ Không có URL thanh toán / billing / profile trong Sitemap
•	Security
□ Không lộ email
□ Không lộ số điện thoại
□ Không lộ token/session
□ Không lộ hóa đơn
□ Không lộ lịch sử thanh toán
□ Không lộ quyền truy cập cá nhân
□ Permission hoạt động đúng trên toàn bộ màn hình QH.8

 
D.9	QH.9 – API & Tích hợp
D.9.1 API Public
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	No đối với API runtime
SEO Class	SEO-N
Entity chính	API Public
Entity phụ	Parcel, Administrative Unit, Planning Region, Planning Project, Legal Document
URL Type	API URL
Mục tiêu SEO	Không SEO endpoint API; chỉ dùng API để cung cấp dữ liệu cho các URL SEO công khai khác
•	Ghi chú
•	API Public không phải Landing Page SEO. 
•	API Public không được Index như một trang nội dung. 
•	SEO chỉ áp dụng cho các trang Web/Share/Report sử dụng dữ liệu từ API Public. 
•	Nếu sau này có Developer Portal / API Documentation công khai, phần đó phải được đặc tả SEO riêng, không gộp với API runtime. 
 
•	16.2 URL & Canonical
•	URL chuẩn
Không áp dụng cho SEO.
•	URL NoIndex
Toàn bộ API endpoint runtime phải NoIndex, ví dụ:
/api/public/parcels
/api/public/planning-regions
/api/public/planning-projects
/api/public/legal-documents
/api/public/gis-layers
•	Canonical
Không áp dụng cho API runtime.
•	Không tạo URL SEO riêng cho
•	API endpoint. 
•	Query API. 
•	Pagination API. 
•	Filter API. 
•	Search API. 
•	GeoJSON response. 
•	JSON response. 
•	API token. 
•	API version runtime. 
•	Duplicate Control
Không để Google hoặc crawler index các endpoint trả về JSON, GeoJSON, file dữ liệu hoặc payload kỹ thuật.
 
•	16.3 Metadata
•	Title
Không áp dụng.
•	Description
Không áp dụng.
•	H1
Không áp dụng.
•	Robots
noindex,nofollow
•	Ghi chú
API runtime không sinh Metadata SEO.
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Không áp dụng
Entity Mapping	Không áp dụng
•	Điều kiện áp dụng
Không áp dụng với API runtime.
•	Điều kiện loại bỏ
Luôn loại bỏ.
•	Ghi chú
Structured Data được sinh tại trang Web SEO, không sinh tại API endpoint.
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Không áp dụng
OG Description	Không áp dụng
OG Image	Không áp dụng
•	Share Preview
Không hỗ trợ chia sẻ trực tiếp API endpoint.
 
•	16.6 AI Summary
•	Có AI Summary
No đối với API endpoint.
•	Dữ liệu nguồn
API Public có thể là nguồn dữ liệu để các màn hình khác sinh AI Summary, nhưng bản thân endpoint không sinh AI Summary.
•	Không public
•	API key. 
•	Token. 
•	Rate limit state. 
•	Request log. 
•	User context. 
•	Partner context. 
•	Internal metadata. 
•	Response debug. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	Không áp dụng
•	Tham gia khi
Không có trường hợp nào.
•	Loại bỏ khỏi Sitemap khi
Luôn loại bỏ toàn bộ API runtime URL.
 
•	16.8 Internal Link
•	Link đến
Không áp dụng cho SEO.
•	Link nhận từ
Không áp dụng cho SEO.
•	Ghi chú
API Public có thể được gọi bởi:
•	Trang chi tiết thửa đất. 
•	Trang chi tiết khu quy hoạch. 
•	Trang chi tiết đồ án. 
•	Report. 
•	Snapshot. 
•	AI Summary. 
Nhưng các quan hệ này là quan hệ dữ liệu, không phải Internal Link SEO.
 
•	16.9 Security & Visibility
•	Public
•	Dữ liệu API được phép công khai. 
•	Payload đã kiểm soát quyền. 
•	Dữ liệu đã lọc theo trạng thái Public. 
•	Restricted
•	Dữ liệu yêu cầu gói dịch vụ. 
•	Dữ liệu yêu cầu API key. 
•	Dữ liệu bị giới hạn theo quota. 
•	Dữ liệu đối tác. 
•	Private
•	API key. 
•	Access token. 
•	Request log. 
•	User ID. 
•	Partner ID. 
•	Debug payload. 
•	Internal permission rule. 
•	Permission Rule
API Public chỉ được trả về dữ liệu đúng phạm vi công khai hoặc đúng quyền truy cập.
Không để API Public trở thành kênh làm lộ dữ liệu Restricted/Private thông qua crawler, bot hoặc request trực tiếp.
 
•	16.10 Cache & Regeneration
Event	Regenerate
public_api_schema.updated	API Cache
entity.visibility.changed	API Permission Cache
parcel.updated	API Data Cache
planning_region.updated	API Data Cache
legal_document.updated	API Data Cache
api_rate_limit.updated	API Runtime Cache
•	Ghi chú
Không phát sinh:
•	SEO Metadata Regeneration. 
•	Sitemap Regeneration. 
•	Schema Regeneration. 
•	Open Graph Regeneration. 
•	AI Summary Regeneration tại endpoint. 
 
•	16.11 Acceptance Criteria
•	AC-D9.1-001
Toàn bộ API Public runtime URL phải NoIndex.
•	AC-D9.1-002
API endpoint không xuất hiện trong Sitemap.
•	AC-D9.1-003
API endpoint không sinh Metadata, Schema, Open Graph hoặc AI Summary.
•	AC-D9.1-004
API Public chỉ trả về dữ liệu đúng quyền công khai hoặc đúng quyền truy cập.
•	AC-D9.1-005
Không lộ API key, token, request log hoặc debug payload.
•	AC-D9.1-006
Các trang SEO sử dụng API Public phải tự sinh SEO Output tại tầng Web, không phụ thuộc vào endpoint API được Index.
 
•	16.12 QA Checklist
•	URL
□ API Public URL NoIndex
□ Không có API URL trong Sitemap
□ Không Index JSON / GeoJSON response
□ Không Index API query URL
•	Metadata
□ Không sinh Metadata SEO
□ Không có Title / Description SEO tại endpoint
•	Structured Data
□ Không inject Schema tại API endpoint
•	Open Graph
□ Không sinh OG tại API endpoint
•	AI Summary
□ Không sinh AI Summary tại endpoint
□ API chỉ là nguồn dữ liệu cho màn hình SEO khác
•	Sitemap
□ Không xuất hiện trong Sitemap
□ Không có API pagination / filter trong Sitemap
•	Security
□ Không lộ API key
□ Không lộ token
□ Không lộ request log
□ Không lộ debug payload
□ Không trả dữ liệu vượt quyền truy cập
 
D.9.2 API Partner
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	No
SEO Class	SEO-N
Entity chính	Partner API
Entity phụ	Partner, Integration, Parcel, Planning Region, Report
URL Type	API URL / Partner Runtime URL
Mục tiêu SEO	Không SEO API Partner; bảo vệ dữ liệu tích hợp, quyền truy cập và phạm vi chia sẻ với đối tác
•	Ghi chú
•	API Partner là kênh tích hợp có kiểm soát. 
•	Không Index endpoint Partner. 
•	Không đưa endpoint Partner vào Sitemap. 
•	Không để crawler truy cập hoặc đọc payload Partner. 
•	Nếu có trang giới thiệu đối tác hoặc tài liệu tích hợp công khai, phần đó phải đặc tả riêng. 
 
•	16.2 URL & Canonical
•	URL chuẩn
Không áp dụng cho SEO.
•	URL NoIndex
Toàn bộ endpoint Partner phải NoIndex, ví dụ:
/api/partner/{partner-code}/parcels
/api/partner/{partner-code}/planning
/api/partner/{partner-code}/reports
/api/partner/{partner-code}/sync
/api/partner/{partner-code}/webhook
•	Canonical
Không áp dụng.
•	Không tạo URL SEO riêng cho
•	Partner endpoint. 
•	Partner webhook. 
•	Sync URL. 
•	Export URL. 
•	Callback URL. 
•	API token URL. 
•	Partner dashboard runtime. 
•	Query / pagination / filter. 
•	Duplicate Control
Không để endpoint Partner bị crawl, index hoặc xuất hiện qua link công khai.
 
•	16.3 Metadata
•	Title
Không áp dụng.
•	Description
Không áp dụng.
•	H1
Không áp dụng.
•	Robots
noindex,nofollow
•	Ghi chú
Không sinh Metadata SEO cho API Partner.
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Không áp dụng
Entity Mapping	Không áp dụng
•	Điều kiện áp dụng
Không áp dụng.
•	Điều kiện loại bỏ
Luôn loại bỏ.
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Không áp dụng
OG Description	Không áp dụng
OG Image	Không áp dụng
•	Share Preview
Không hỗ trợ chia sẻ API Partner.
 
•	16.6 AI Summary
•	Có AI Summary
No.
•	Dữ liệu nguồn
API Partner không sinh nội dung AI Summary công khai.
•	Không public
•	Partner ID. 
•	Partner token. 
•	API key. 
•	Secret. 
•	Webhook payload. 
•	Đồng bộ dữ liệu. 
•	Log tích hợp. 
•	Quyền truy cập đối tác. 
•	Dữ liệu trả riêng cho đối tác. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	Không áp dụng
•	Tham gia khi
Không có trường hợp nào.
•	Loại bỏ khỏi Sitemap khi
Luôn loại bỏ toàn bộ API Partner URL.
 
•	16.8 Internal Link
•	Link đến
Không áp dụng cho SEO.
•	Link nhận từ
Không áp dụng cho SEO.
•	Ghi chú
Không tạo link công khai tới API Partner endpoint.
Nếu có tài liệu tích hợp, chỉ link tới trang tài liệu công khai đã được tách riêng và kiểm soát quyền.
 
•	16.9 Security & Visibility
•	Public
Không có dữ liệu SEO public.
•	Restricted
•	Dữ liệu chia sẻ theo hợp đồng đối tác. 
•	Dữ liệu giới hạn theo quyền Partner. 
•	Gói API đối tác. 
•	Private
•	API key. 
•	Secret. 
•	Token. 
•	Webhook signature. 
•	Partner credential. 
•	Sync log. 
•	Request / response log. 
•	Internal mapping. 
•	Permission Rule
API Partner bắt buộc kiểm tra:
•	Partner identity. 
•	API key / token. 
•	Scope. 
•	Quota. 
•	IP allowlist nếu có. 
•	Permission theo loại dữ liệu. 
Không được expose bất kỳ payload Partner nào vào SEO Output.

•	16.10 Cache & Regeneration
Event	Regenerate
partner.permission.changed	Partner Permission Cache
partner.api_key.rotated	Auth Cache
partner.scope.updated	Access Cache
partner.contract.changed	Access Cache
entity.visibility.changed	Partner Data Cache
webhook.updated	Webhook Runtime Cache
•	Ghi chú
Không phát sinh SEO Cache hoặc SEO Regeneration.
 
•	16.11 Acceptance Criteria
•	AC-D9.2-001
Toàn bộ API Partner endpoint phải NoIndex.
•	AC-D9.2-002
Không có API Partner URL trong Sitemap.
•	AC-D9.2-003
Không sinh Metadata, Schema, Open Graph hoặc AI Summary.
•	AC-D9.2-004
Không endpoint Partner nào được link công khai nếu không có kiểm soát quyền.
•	AC-D9.2-005
Không lộ Partner token, API key, webhook payload hoặc sync log.
•	AC-D9.2-006
Payload Partner chỉ trả dữ liệu đúng scope và quyền truy cập.
 
•	16.12 QA Checklist
•	URL
□ Partner API URL NoIndex
□ Webhook URL NoIndex
□ Sync URL NoIndex
□ Không có Partner API trong Sitemap
•	Metadata
□ Không sinh Metadata SEO
□ Không có Title / Description công khai
•	Structured Data
□ Không inject Schema
•	Open Graph
□ Không sinh OG
•	AI Summary
□ Không sinh AI Summary
□ Không đưa dữ liệu Partner vào machine-readable public content
•	Sitemap
□ Không xuất hiện trong Sitemap
□ Không có URL callback / webhook / export trong Sitemap
•	Security
□ Không lộ API key
□ Không lộ Partner token
□ Không lộ webhook secret
□ Không lộ sync log
□ Scope truy cập đúng
□ Quota / permission hoạt động đúng
 
D.9.3 API Internal
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	No
SEO Class	SEO-N
Entity chính	Internal API
Entity phụ	System Service, Admin Object, Runtime Job
URL Type	Internal API URL
Mục tiêu SEO	Không SEO; API Internal chỉ phục vụ vận hành nội bộ, backend service, admin và runtime system
•	Ghi chú
•	API Internal tuyệt đối không thuộc phạm vi SEO. 
•	Không được Index. 
•	Không được Sitemap. 
•	Không được expose ra public crawler. 
•	Không sinh bất kỳ SEO Output nào. 
 
•	16.2 URL & Canonical
•	URL chuẩn
Không áp dụng.
•	URL NoIndex
Toàn bộ Internal API URL, ví dụ:
/api/internal/*
/internal-api/*
/admin-api/*
/system/*
/runtime/*
/jobs/*
/queue/*
/cache/*
•	Canonical
Không áp dụng.
•	Không tạo URL SEO riêng cho
•	Internal endpoint. 
•	Admin API. 
•	Runtime job. 
•	Queue. 
•	Cache. 
•	Debug endpoint. 
•	Health check nội bộ. 
•	Data pipeline endpoint. 
•	System callback. 
•	Duplicate Control
Internal API không được xuất hiện trong public routing, public sitemap, public link hoặc share preview.
 
•	16.3 Metadata
•	Title
Không áp dụng.
•	Description
Không áp dụng.
•	H1
Không áp dụng.
•	Robots
noindex,nofollow
•	Ghi chú
Internal API không sinh bất kỳ Metadata SEO nào.
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Không áp dụng
Entity Mapping	Không áp dụng
•	Điều kiện áp dụng
Không áp dụng.
•	Điều kiện loại bỏ
Luôn loại bỏ.
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Không áp dụng
OG Description	Không áp dụng
OG Image	Không áp dụng
•	Share Preview
Không hỗ trợ.
 
•	16.6 AI Summary
•	Có AI Summary
No.
•	Dữ liệu nguồn
Không áp dụng.
•	Không public
•	Internal payload. 
•	Debug data. 
•	Runtime state. 
•	Job status. 
•	Queue payload. 
•	Cache key. 
•	Admin data. 
•	System config. 
•	Secret. 
•	Credential. 
•	Log. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	Không áp dụng
•	Tham gia khi
Không có trường hợp nào.
•	Loại bỏ khỏi Sitemap khi
Luôn loại bỏ toàn bộ API Internal.
 
•	16.8 Internal Link
•	Link đến
Không áp dụng cho SEO.
•	Link nhận từ
Không áp dụng cho SEO.
•	Ghi chú
Không có public link tới Internal API.
 
•	16.9 Security & Visibility
•	Public
Không có.
•	Restricted
Không áp dụng cho SEO.
•	Private
Toàn bộ dữ liệu Internal API, bao gồm:
•	System config. 
•	Credential. 
•	Secret. 
•	Runtime state. 
•	Admin payload. 
•	Queue payload. 
•	Cache data. 
•	Internal log. 
•	Error trace. 
•	Permission Rule
Internal API phải được bảo vệ bằng:
•	Authentication. 
•	Authorization. 
•	Network restriction nếu có. 
•	Role-based permission. 
•	Service-to-service permission. 
•	Audit log. 
Không expose bất kỳ Internal API output nào ra public SEO surface.
 
•	16.10 Cache & Regeneration
Event	Regenerate
internal_config.updated	Internal Cache
job.status.changed	Runtime Cache
queue.updated	Queue Runtime
cache.invalidated	Internal Cache
permission.changed	Access Cache
•	Ghi chú
Không có SEO Regeneration.
 
•	16.11 Acceptance Criteria
•	AC-D9.3-001
Toàn bộ Internal API phải NoIndex.
•	AC-D9.3-002
Không có Internal API URL trong Sitemap.
•	AC-D9.3-003
Không sinh Metadata, Schema, OG hoặc AI Summary.
•	AC-D9.3-004
Không endpoint Internal nào được public link.
•	AC-D9.3-005
Không lộ credential, secret, config, log hoặc runtime state.
•	AC-D9.3-006
Crawler không được truy cập Internal API ngoài phạm vi cho phép.
 
•	16.12 QA Checklist
•	URL
□ Internal API NoIndex
□ Không có Internal API trong Sitemap
□ Không public link Internal API
□ Debug URL không Index
•	Metadata
□ Không sinh Metadata SEO
•	Structured Data
□ Không inject Schema
•	Open Graph
□ Không sinh OG
•	AI Summary
□ Không sinh AI Summary
□ Không đưa internal payload vào public content
•	Sitemap
□ Không xuất hiện trong Sitemap
□ Không có /internal/*, /admin-api/*, /runtime/* trong Sitemap
•	Security
□ Không lộ secret
□ Không lộ credential
□ Không lộ config
□ Không lộ log
□ Không lộ queue payload
□ Permission / authentication hoạt động đúng

 
D.10	QH.10 – Quản trị hệ thống
D.10.1 Các màn hình SEO-N
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	No
SEO Class	SEO-N
Entity chính	Admin Object / System Configuration
Entity phụ	User, Role, Permission, Data Job, GIS Dataset, System Log
URL Type	Admin / Internal URL
Mục tiêu SEO	Không SEO; toàn bộ khu vực quản trị chỉ phục vụ vận hành nội bộ, cấu hình hệ thống, quản trị dữ liệu và kiểm soát phân quyền
•	Ghi chú
•	QH.10 là khu vực quản trị nội bộ. 
•	Không phải Landing Page. 
•	Không phải nội dung công khai. 
•	Không tham gia Google Search, AI Search, Sitemap, Open Graph hoặc Discover. 
•	Mọi dữ liệu quản trị, cấu hình, log, job, quyền truy cập và dữ liệu vận hành phải được xem là nội bộ. 
 
•	16.2 URL & Canonical
•	URL chuẩn
Không áp dụng cho SEO.
•	URL NoIndex
Toàn bộ URL Admin mặc định NoIndex, bao gồm:
/admin
/admin/*
/quan-tri
/quan-tri/*
/system/*
/dashboard/admin
•	Canonical
Không áp dụng.
•	Không tạo URL SEO riêng cho
•	Dashboard quản trị. 
•	Quản lý người dùng. 
•	Quản lý phân quyền. 
•	Quản lý dữ liệu GIS. 
•	Quản lý layer. 
•	Quản lý đồ án. 
•	Quản lý văn bản. 
•	Quản lý job. 
•	Quản lý log. 
•	Cấu hình hệ thống. 
•	Lịch sử thao tác. 
•	Import / Export dữ liệu. 
•	Data pipeline. 
•	Runtime operation. 
•	Duplicate Control
Không URL nào thuộc QH.10 được phép xuất hiện như một public URL hoặc SEO URL.
 
•	16.3 Metadata
•	Title
Không sinh Metadata SEO riêng.
Có thể dùng Title kỹ thuật phục vụ vận hành, ví dụ:
Quản trị hệ thống | QH Pro
nhưng không sử dụng cho SEO.
•	Description
Không sinh Description SEO.
•	H1
Theo màn hình quản trị, không dùng làm H1 SEO.
•	Robots
noindex,nofollow
•	Ghi chú
Không đưa thông tin quản trị, dữ liệu nội bộ, trạng thái job, log hoặc cấu hình vào Metadata.
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Không áp dụng
Entity Mapping	Không áp dụng
•	Điều kiện áp dụng
Không áp dụng.
•	Điều kiện loại bỏ
Luôn loại bỏ đối với toàn bộ QH.10.
•	Ghi chú
Không inject Schema cho bất kỳ màn hình Admin nào.
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Không áp dụng
OG Description	Không áp dụng
OG Image	Không áp dụng
•	Share Preview
Không hỗ trợ Share Preview cho khu vực Admin.
•	Ghi chú
Không được tạo Open Graph chứa dữ liệu quản trị, log, cấu hình hoặc dữ liệu nội bộ.
 
•	16.6 AI Summary
•	Có AI Summary
No
•	Dữ liệu nguồn
Không áp dụng.
•	Không public
•	Cấu hình hệ thống. 
•	Quyền người dùng. 
•	Role / Permission. 
•	Admin Log. 
•	System Log. 
•	Data Job. 
•	Import Log. 
•	Validation Log. 
•	Error Log. 
•	Data Pipeline. 
•	Cache Key. 
•	API Key. 
•	Credential. 
•	Secret. 
•	Runtime State. 
•	Ghi chú
Không đưa bất kỳ dữ liệu QH.10 nào vào AI Summary công khai hoặc nội dung machine-readable phục vụ AI Search.

•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	Không áp dụng
•	Tham gia Sitemap khi
Không có trường hợp nào.
•	Loại bỏ khỏi Sitemap khi
Luôn loại bỏ toàn bộ URL thuộc QH.10.
 
•	16.8 Internal Link
•	Link đến
Các màn hình Admin có thể điều hướng nội bộ tới:
•	Quản lý người dùng. 
•	Quản lý phân quyền. 
•	Quản lý dữ liệu. 
•	Quản lý layer. 
•	Quản lý đồ án. 
•	Quản lý văn bản. 
•	Quản lý job. 
•	Quản lý log. 
•	Cấu hình hệ thống. 
•	Link nhận từ
•	Admin dashboard. 
•	Menu quản trị. 
•	Notification nội bộ. 
•	Job monitor. 
•	Data pipeline monitor. 
•	Ghi chú
Các liên kết này là liên kết vận hành nội bộ, không phải Internal Link SEO công khai.
Không tính các link trong Admin vào kiến trúc Internal Linking SEO.
 
•	16.9 Security & Visibility
•	Public
Không có dữ liệu public SEO trong QH.10.
•	Restricted
•	Dữ liệu quản trị theo vai trò. 
•	Dữ liệu vận hành theo quyền. 
•	Dữ liệu cấu hình theo quyền Admin. 
•	Private
•	Credential. 
•	Secret. 
•	API Key. 
•	Access Token. 
•	Refresh Token. 
•	System Config. 
•	Admin Log. 
•	Audit Log. 
•	Runtime Log. 
•	Import Log. 
•	Validation Error. 
•	Pipeline State. 
•	Job Payload. 
•	Cache Data. 
•	Permission Rule
Toàn bộ QH.10 bắt buộc yêu cầu:
•	Đăng nhập. 
•	Phân quyền vai trò. 
•	Kiểm tra quyền truy cập từng màn hình. 
•	Kiểm tra quyền thao tác. 
•	Audit log đối với thao tác quan trọng. 
Không đưa dữ liệu Restricted hoặc Private vào:
•	Metadata. 
•	Structured Data. 
•	Open Graph. 
•	AI Summary. 
•	Sitemap. 
•	Public HTML không cần xác thực. 
 
•	16.10 Cache & Regeneration
Event	Regenerate
admin_config.updated	Admin Cache
user_role.changed	Permission Cache
permission.changed	Access Cache
data_job.updated	Job Runtime Cache
layer_admin.updated	Admin Data Cache
import_job.updated	Pipeline Cache
system_setting.updated	Runtime Config Cache
•	Ghi chú
Không phát sinh:
•	SEO Metadata Regeneration. 
•	Sitemap Regeneration. 
•	Schema Regeneration. 
•	Open Graph Regeneration. 
•	AI Summary Regeneration. 
 
•	16.11 Acceptance Criteria
•	AC-D10.1-001
Toàn bộ URL Admin phải NoIndex.
•	AC-D10.1-002
Không URL nào thuộc QH.10 được xuất hiện trong Sitemap.
•	AC-D10.1-003
Không sinh Structured Data cho màn hình Admin.
•	AC-D10.1-004
Không sinh Open Graph cho màn hình Admin.
•	AC-D10.1-005
Không sinh AI Summary cho dữ liệu Admin.
•	AC-D10.1-006
Không lộ credential, secret, token, API key, system config, log hoặc job payload trong HTML công khai.
•	AC-D10.1-007
Mọi màn hình Admin phải yêu cầu đăng nhập và phân quyền phù hợp.
•	AC-D10.1-008
Crawler không được truy cập nội dung Admin ngoài phạm vi cho phép.
 
•	16.12 QA Checklist
•	URL
□ Admin URL NoIndex
□ /admin/* NoIndex
□ /quan-tri/* NoIndex
□ Không có Admin URL trong Sitemap
□ Không public link Admin URL
•	Metadata
□ Không sinh Metadata SEO riêng
□ Robots đúng noindex,nofollow
□ Không chứa dữ liệu quản trị trong Metadata
•	Structured Data
□ Không inject Schema
□ Không có Entity Mapping công khai
•	Open Graph
□ Không sinh OG
□ Không có Share Preview Admin
•	AI Summary
□ Không sinh AI Summary
□ Không đưa dữ liệu Admin vào machine-readable public content
•	Sitemap
□ Không xuất hiện trong bất kỳ Sitemap nào
□ Không có URL job / log / config / user management trong Sitemap
•	Security
□ Không lộ credential
□ Không lộ secret
□ Không lộ API key
□ Không lộ access token
□ Không lộ system config
□ Không lộ admin log
□ Không lộ import log
□ Không lộ job payload
□ Permission hoạt động đúng trên toàn bộ QH.10

 
D.11	QH.11 – Thư viện số Quy hoạch
D.11.1 Workspace thư viện
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes
SEO Class	SEO-B
Entity chính	Planning Library
Entity phụ	Planning Project, Legal Document, Planning Map, Administrative Unit
URL Type	Landing / List URL
Mục tiêu SEO	Trang điều hướng trung tâm cho thư viện số quy hoạch, hỗ trợ Google Search, AI Search, Internal Linking và Knowledge Graph
•	Ghi chú
•	Workspace thư viện là điểm vào tổng thể của QH.11. 
•	Không phải URL chi tiết cuối cùng. 
•	Có nhiệm vụ điều hướng đến đồ án, văn bản, bản đồ, timeline và quan hệ phiên bản. 
 
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/thu-vien
•	URL phụ
/quy-hoach/thu-vien/{administrative-unit-slug}
•	URL NoIndex
/quy-hoach/thu-vien?filter={filter}
/quy-hoach/thu-vien?sort={sort}
/quy-hoach/thu-vien?view={view-mode}
•	Canonical
•	Workspace tổng: canonical về /quy-hoach/thu-vien. 
•	Workspace theo địa bàn: canonical về URL thư viện địa bàn. 
•	Các URL filter/sort/view mode canonical về URL chuẩn tương ứng. 
•	Không tạo URL SEO riêng cho
•	Bộ lọc. 
•	Sort. 
•	View mode. 
•	Tab. 
•	Search runtime. 
•	Pagination tạm thời. 
 
•	16.3 Metadata
•	Title
Thư viện số Quy hoạch | QH Pro
•	Description
Tra cứu thư viện số quy hoạch gồm đồ án, văn bản pháp lý, bản đồ quy hoạch, timeline pháp lý và quan hệ phiên bản trên hệ thống QH Pro.
•	H1
Thư viện số Quy hoạch
•	Robots
index,follow
•	NoIndex khi
•	Workspace không có dữ liệu công khai. 
•	Trang chỉ hiển thị dữ liệu Restricted. 
•	Filter runtime không tạo URL chuẩn. 
•	Trang lỗi hoặc trạng thái rỗng. 
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	CollectionPage
Entity Mapping	Planning Library
•	Schema bổ sung
Schema	Điều kiện
Dataset	Có tập dữ liệu thư viện công khai
BreadcrumbList	Luôn áp dụng nếu có breadcrumb
ItemList	Có danh sách đồ án / văn bản / bản đồ công khai
•	Không inject Schema nếu
•	Không có dữ liệu công khai. 
•	Trang Restricted. 
•	Dữ liệu đang Draft. 
•	Trang filter runtime. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Library Preview / Map Preview
OG Type	website
•	Share Preview
Hiển thị:
•	Tên thư viện. 
•	Phạm vi địa bàn nếu có. 
•	Nhóm nội dung: đồ án, văn bản, bản đồ. 
•	Preview an toàn. 
 
•	16.6 AI Summary
•	Có AI Summary
Yes
•	Dữ liệu nguồn
•	Planning Library. 
•	Planning Project công khai. 
•	Legal Document công khai. 
•	Planning Map công khai. 
•	Administrative Unit. 
•	Timeline pháp lý công khai. 
•	Nội dung AI Summary
•	Mô tả phạm vi thư viện. 
•	Số nhóm dữ liệu chính. 
•	Các loại quy hoạch có trong thư viện. 
•	Các đồ án / văn bản / bản đồ nổi bật. 
•	Liên kết tới nhóm nội dung trọng tâm. 
•	Không public
•	Hồ sơ Restricted. 
•	Văn bản chưa công bố. 
•	Bản đồ nội bộ. 
•	Metadata vận hành. 
•	Ghi chú kiểm duyệt. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	planning-library.xml
•	Tham gia khi
•	Public. 
•	Active. 
•	Có dữ liệu công khai. 
•	Có URL chuẩn. 
•	Có Metadata hợp lệ. 
•	Loại bỏ khi
•	Empty Library. 
•	Restricted. 
•	Draft. 
•	Deleted. 
•	Chỉ có dữ liệu nội bộ. 
 
•	16.8 Internal Link
•	Link đến
•	Danh sách đồ án. 
•	Chi tiết đồ án. 
•	Danh sách văn bản. 
•	Chi tiết văn bản. 
•	Danh sách bản đồ. 
•	Chi tiết bản đồ. 
•	Timeline pháp lý. 
•	Quan hệ phiên bản. 
•	Đơn vị hành chính liên quan. 
•	Link nhận từ
•	Trang bản đồ quy hoạch. 
•	Trang chi tiết thửa đất. 
•	Trang chi tiết khu quy hoạch. 
•	Report. 
•	Snapshot. 
•	Trang chủ / menu chính. 
 
•	16.9 Security & Visibility
•	Public
•	Danh sách nội dung công khai. 
•	Metadata công khai. 
•	Preview thư viện. 
•	AI Summary công khai. 
•	Restricted
•	Hồ sơ chưa công bố. 
•	Văn bản hạn chế. 
•	Bản đồ trả phí / đối tác. 
•	Tài liệu kiểm duyệt. 
•	Private
•	Ghi chú nội bộ. 
•	Trạng thái xử lý. 
•	Import log. 
•	Audit log. 
•	Quyền truy cập nội bộ. 
•	Permission Rule
Chỉ các nội dung Public mới được đưa vào Metadata, Schema, OG, AI Summary và Sitemap.
 
•	16.10 Cache & Regeneration
Event	Regenerate
planning_library.updated	Metadata, AI Summary
planning_project.created	AI Summary, Sitemap
legal_document.created	AI Summary, Sitemap
planning_map.created	AI Summary, Sitemap
visibility.changed	Robots, Sitemap
administrative_unit.updated	Metadata
 
•	16.11 Acceptance Criteria
•	AC-D11.1-001
Workspace thư viện có URL chuẩn và được Index khi có dữ liệu công khai.
•	AC-D11.1-002
Filter/sort/view mode không tạo URL SEO riêng.
•	AC-D11.1-003
Metadata sinh đầy đủ.
•	AC-D11.1-004
Structured Data dạng CollectionPage hợp lệ.
•	AC-D11.1-005
AI Summary chỉ dùng dữ liệu công khai.
•	AC-D11.1-006
Sitemap đúng điều kiện.
•	AC-D11.1-007
Không lộ dữ liệu Restricted hoặc Private.
 
•	16.12 QA Checklist
•	URL
□ URL Workspace đúng chuẩn
□ Filter/sort NoIndex hoặc canonical đúng
□ Canonical đúng
•	Metadata
□ Title tồn tại
□ Description tồn tại
□ H1 tồn tại
□ Robots đúng
•	Structured Data
□ CollectionPage hợp lệ
□ ItemList đúng nếu có
•	Open Graph
□ OG Title đúng
□ OG Image đúng preview
•	AI Summary
□ Có AI Summary
□ Không dùng dữ liệu Restricted
•	Sitemap
□ Có trong planning-library.xml khi đủ điều kiện
□ Empty / Restricted bị loại khỏi Sitemap
•	Security
□ Không lộ hồ sơ nội bộ
□ Không lộ văn bản chưa công bố
□ Không lộ bản đồ Restricted
 
D.11.2 Danh sách đồ án
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes
SEO Class	SEO-B
Entity chính	Planning Project List
Entity phụ	Planning Project, Administrative Unit, Legal Document
URL Type	List URL
Mục tiêu SEO	Trang danh mục đồ án quy hoạch, hỗ trợ điều hướng đến các URL chi tiết đồ án SEO-A
•	Ghi chú
•	Đây là trang danh sách quan trọng. 
•	Có thể SEO theo địa bàn, loại quy hoạch, cấp quy hoạch nếu có URL chuẩn. 
•	Không SEO các filter runtime rời rạc. 
 
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/do-an
•	URL theo địa bàn
/quy-hoach/{administrative-unit-slug}/do-an
•	URL theo loại quy hoạch
/quy-hoach/do-an/{planning-type-slug}
•	URL NoIndex
/quy-hoach/do-an?filter={filter}
/quy-hoach/do-an?sort={sort}
/quy-hoach/do-an?page={page}
•	Canonical
•	Danh sách tổng: canonical về /quy-hoach/do-an. 
•	Danh sách địa bàn: canonical về URL địa bàn tương ứng. 
•	Filter/sort/pagination tạm thời canonical về URL danh sách chuẩn. 
•	Không tạo URL SEO riêng cho
•	Sort. 
•	Filter runtime. 
•	Search keyword. 
•	View mode. 
•	Pagination không ổn định. 
 
•	16.3 Metadata
•	Title
Danh sách đồ án quy hoạch | QH Pro
•	Title theo địa bàn
Danh sách đồ án quy hoạch {administrative_unit_name} | QH Pro
•	Description
Tra cứu danh sách đồ án quy hoạch, quy hoạch chung, quy hoạch phân khu, quy hoạch chi tiết, quy hoạch sử dụng đất và các văn bản liên quan trên QH Pro.
•	H1
Danh sách đồ án quy hoạch
•	Robots
index,follow
•	NoIndex khi
•	Không có đồ án công khai. 
•	Chỉ có dữ liệu Restricted. 
•	Filter runtime không có URL chuẩn. 
•	Trang tìm kiếm tạm thời. 
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	CollectionPage / ItemList
Entity Mapping	Planning Project List
•	Schema bổ sung
Schema	Điều kiện
BreadcrumbList	Có breadcrumb
Dataset	Danh sách đại diện tập dữ liệu công khai
•	Không inject Schema nếu
•	Danh sách rỗng. 
•	Danh sách Restricted. 
•	Filter runtime. 
•	Không có item công khai. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Planning Project List Preview
OG Type	website
•	Share Preview
Hiển thị:
•	Tên danh sách. 
•	Địa bàn nếu có. 
•	Loại quy hoạch nếu có. 
•	Số lượng đồ án công khai nếu được phép. 
 
•	16.6 AI Summary
•	Có AI Summary
Yes
•	Dữ liệu nguồn
•	Planning Project công khai. 
•	Administrative Unit. 
•	Planning Type. 
•	Legal Document liên quan. 
•	Planning Library. 
•	Nội dung AI Summary
•	Tóm tắt nhóm đồ án. 
•	Phân loại đồ án. 
•	Địa bàn áp dụng. 
•	Đồ án nổi bật. 
•	Liên kết đến chi tiết đồ án. 
•	Không public
•	Đồ án Draft. 
•	Đồ án chưa công bố. 
•	Hồ sơ Restricted. 
•	Ghi chú kiểm duyệt. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	planning-project-list.xml
•	Tham gia khi
•	Public. 
•	Active. 
•	Có tối thiểu một đồ án công khai. 
•	Có URL chuẩn. 
•	Có Metadata hợp lệ. 
•	Loại bỏ khi
•	Empty List. 
•	Restricted. 
•	Draft. 
•	Runtime filter. 
•	Search result tạm thời. 
 
•	16.8 Internal Link
•	Link đến
•	Chi tiết đồ án. 
•	Đơn vị hành chính. 
•	Văn bản pháp lý. 
•	Bản đồ quy hoạch. 
•	Timeline pháp lý. 
•	Quan hệ phiên bản. 
•	Link nhận từ
•	Workspace thư viện. 
•	Trang bản đồ. 
•	Chi tiết khu quy hoạch. 
•	Chi tiết thửa đất. 
•	Report. 
•	Snapshot. 
 
•	16.9 Security & Visibility
•	Public
•	Tên đồ án công khai. 
•	Loại quy hoạch. 
•	Địa bàn. 
•	Trạng thái hiệu lực công khai. 
•	Link chi tiết công khai. 
•	Restricted
•	Đồ án chưa công bố. 
•	Hồ sơ hạn chế. 
•	Dữ liệu kiểm duyệt. 
•	Private
•	Ghi chú nội bộ. 
•	Workflow xử lý. 
•	Import log. 
•	Audit log. 
•	Permission Rule
Danh sách công khai chỉ được hiển thị và SEO các đồ án Public.
 
•	16.10 Cache & Regeneration
Event	Regenerate
planning_project.created	Metadata, Sitemap
planning_project.updated	AI Summary
planning_project.visibility.changed	Robots, Sitemap
administrative_unit.updated	Metadata
legal_document.updated	AI Summary
 
•	16.11 Acceptance Criteria
•	AC-D11.2-001
Danh sách đồ án có URL chuẩn.
•	AC-D11.2-002
Filter/sort/search runtime không tạo URL SEO riêng.
•	AC-D11.2-003
Metadata đúng theo phạm vi danh sách.
•	AC-D11.2-004
Structured Data dạng CollectionPage / ItemList hợp lệ.
•	AC-D11.2-005
Danh sách đủ điều kiện xuất hiện trong Sitemap.
•	AC-D11.2-006
Không hiển thị đồ án Restricted trong danh sách public.
 
•	16.12 QA Checklist
•	URL
□ URL danh sách đúng chuẩn
□ Filter/sort NoIndex hoặc canonical đúng
□ Canonical đúng
•	Metadata
□ Title đúng danh sách
□ Description tồn tại
□ H1 tồn tại
□ Robots đúng
•	Structured Data
□ CollectionPage hợp lệ
□ ItemList chỉ chứa item public
•	Open Graph
□ OG đúng danh sách
□ Preview không lộ dữ liệu Restricted
•	AI Summary
□ Tóm tắt đúng danh sách
□ Không chứa đồ án Restricted
•	Sitemap
□ Có trong planning-project-list.xml khi đủ điều kiện
□ Empty / Restricted bị loại khỏi Sitemap
•	Security
□ Không lộ đồ án chưa công bố
□ Không lộ ghi chú nội bộ

D.11.3 Chi tiết đồ án
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes
SEO Class	SEO-A
Entity chính	Planning Project
Entity phụ	Administrative Unit, Legal Document, Planning Map, GIS Layer, Planning Region
URL Type	Detail URL
Mục tiêu SEO	URL SEO trọng tâm cho đồ án quy hoạch, phục vụ Google Search, AI Search, Knowledge Graph, Historical SEO và Internal Linking
•	Ghi chú
•	Đây là một trong các URL SEO quan trọng nhất của QH Pro. 
•	Mỗi đồ án công khai phải có URL chi tiết chuẩn. 
•	Đồ án cũ, đồ án thay thế, đồ án điều chỉnh vẫn cần hỗ trợ Historical SEO nếu có dữ liệu công khai. 
 
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/do-an/{planning-project-slug}
•	URL theo địa bàn
/quy-hoach/{administrative-unit-slug}/do-an/{planning-project-slug}
•	URL lịch sử
/quy-hoach/do-an/{old-project-slug}
•	URL NoIndex
/quy-hoach/do-an/{project-id}?tab={tab}
/quy-hoach/do-an/{project-id}?version={version-id}
/quy-hoach/do-an/{project-id}?preview=true
•	Canonical
•	Đồ án hiện hành: canonical về URL chuẩn của đồ án. 
•	Đồ án cũ còn giá trị tham khảo: canonical về URL đồ án cũ nếu có giá trị độc lập. 
•	Đồ án bị thay thế: canonical về đồ án thay thế hoặc giữ URL riêng nếu cần Historical SEO. 
•	Tab/version runtime canonical về URL chi tiết đồ án. 
•	Không tạo URL SEO riêng cho
•	Tab. 
•	Filter. 
•	Preview mode. 
•	Runtime version selector. 
•	Layer state. 
•	Sort tài liệu. 
 
•	16.3 Metadata
•	Title
{planning_project_name} | Đồ án quy hoạch | QH Pro
•	Description
Thông tin chi tiết đồ án {planning_project_name}, bao gồm phạm vi áp dụng, loại quy hoạch, trạng thái hiệu lực, văn bản pháp lý, bản đồ quy hoạch và quan hệ phiên bản.
•	H1
{planning_project_name}
•	Robots
index,follow
•	NoIndex khi
•	Draft. 
•	Deleted. 
•	Restricted. 
•	Chưa công bố. 
•	Thiếu tên đồ án. 
•	Không có dữ liệu công khai tối thiểu. 
 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	CreativeWork
Entity Mapping	Planning Project
•	Schema bổ sung
Schema	Điều kiện
Place	Có phạm vi địa lý
Legislation / CreativeWork	Có văn bản pháp lý
Dataset	Có dữ liệu GIS / bản đồ công khai
BreadcrumbList	Có breadcrumb
•	Không inject Schema nếu
•	Draft. 
•	Restricted. 
•	Missing Legal Status. 
•	Missing Public Content. 
•	Deleted. 
 
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Planning Project Preview / Map Preview
OG Type	article
•	Share Preview
Hiển thị:
•	Tên đồ án. 
•	Địa bàn áp dụng. 
•	Loại quy hoạch. 
•	Trạng thái hiệu lực. 
•	Map Preview / Cover Preview. 
 
•	16.6 AI Summary
•	Có AI Summary
Yes
•	Dữ liệu nguồn
•	Planning Project. 
•	Legal Document. 
•	Planning Map. 
•	GIS Layer công khai. 
•	Administrative Unit. 
•	Timeline pháp lý. 
•	Quan hệ phiên bản. 
•	Nội dung AI Summary
•	Tóm tắt đồ án. 
•	Phạm vi áp dụng. 
•	Loại quy hoạch. 
•	Trạng thái pháp lý. 
•	Văn bản phê duyệt. 
•	Bản đồ liên quan. 
•	Phiên bản / điều chỉnh / thay thế. 
•	Cảnh báo dữ liệu nếu có. 
•	Không public
•	Hồ sơ chưa công bố. 
•	Bản thảo. 
•	Ghi chú nội bộ. 
•	Layer Restricted. 
•	Văn bản hạn chế. 
•	Workflow kiểm duyệt. 
 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	planning-project.xml
•	Tham gia khi
•	Public. 
•	Active hoặc có giá trị Historical SEO. 
•	Có URL chuẩn. 
•	Có Metadata hợp lệ. 
•	Có trạng thái pháp lý rõ ràng. 
•	Có nội dung công khai tối thiểu. 
•	Loại bỏ khi
•	Draft. 
•	Restricted. 
•	Deleted. 
•	Chưa công bố. 
•	Không có dữ liệu công khai. 
 
•	16.8 Internal Link
•	Link đến
•	Đơn vị hành chính. 
•	Văn bản pháp lý. 
•	Bản đồ quy hoạch. 
•	Layer GIS. 
•	Timeline pháp lý. 
•	Quan hệ phiên bản. 
•	Khu quy hoạch liên quan. 
•	Report liên quan. 
•	Link nhận từ
•	Workspace thư viện. 
•	Danh sách đồ án. 
•	Chi tiết văn bản. 
•	Chi tiết bản đồ. 
•	Trang bản đồ quy hoạch. 
•	Chi tiết thửa đất. 
•	Chi tiết khu quy hoạch. 
•	Report. 
•	Snapshot. 
 
•	16.9 Security & Visibility
•	Public
•	Tên đồ án. 
•	Loại quy hoạch. 
•	Phạm vi áp dụng. 
•	Trạng thái hiệu lực. 
•	Văn bản công khai. 
•	Bản đồ công khai. 
•	AI Summary. 
•	Restricted
•	Hồ sơ hạn chế. 
•	Bản đồ chưa công bố. 
•	Tài liệu đối tác. 
•	Dữ liệu trả phí. 
•	Private
•	Ghi chú nội bộ. 
•	Workflow xử lý. 
•	Log kiểm duyệt. 
•	File nội bộ. 
•	Import log. 
•	Permission Rule
Không đưa dữ liệu Restricted hoặc Private vào Metadata, Schema, OG, AI Summary hoặc Sitemap.
 
•	16.10 Cache & Regeneration
Event	Regenerate
planning_project.created	Metadata, Sitemap
planning_project.updated	Metadata, AI Summary
planning_project.visibility.changed	Robots, Sitemap
legal_document.updated	AI Summary
planning_map.updated	Open Graph, AI Summary
timeline.updated	AI Summary, Internal Link
version_relation.updated	AI Summary, Internal Link
administrative_unit.updated	Metadata
 
•	16.11 Acceptance Criteria
•	AC-D11.3-001
Mỗi đồ án công khai có đúng một URL chuẩn.
•	AC-D11.3-002
Canonical xử lý đúng đồ án hiện hành, đồ án cũ và đồ án thay thế.
•	AC-D11.3-003
Metadata sinh đầy đủ theo đồ án.
•	AC-D11.3-004
Structured Data hợp lệ.
•	AC-D11.3-005
AI Summary phản ánh đúng dữ liệu công khai.
•	AC-D11.3-006
Đồ án đủ điều kiện xuất hiện trong planning-project.xml.
•	AC-D11.3-007
Không lộ hồ sơ Restricted hoặc dữ liệu nội bộ.
•	AC-D11.3-008
Internal Link đầy đủ tới văn bản, bản đồ, timeline và phiên bản.
 
•	16.12 QA Checklist
•	URL
□ URL chi tiết đồ án đúng chuẩn
□ Canonical đúng
□ Tab/version runtime không tạo URL SEO riêng
□ Legacy URL xử lý đúng
•	Metadata
□ Title đúng tên đồ án
□ Description đủ thông tin pháp lý / phạm vi
□ H1 tồn tại
□ Robots đúng
•	Structured Data
□ Schema CreativeWork hợp lệ
□ Có BreadcrumbList nếu có breadcrumb
□ Không inject khi Restricted
•	Open Graph
□ OG Image đúng Project Preview
□ Share Preview không lộ dữ liệu hạn chế
•	AI Summary
□ Tóm tắt đúng đồ án
□ Có trạng thái pháp lý
□ Có văn bản liên quan
□ Không lộ hồ sơ chưa công bố
•	Sitemap
□ Đồ án public có trong planning-project.xml
□ Draft / Restricted bị loại khỏi Sitemap
•	Security
□ Không lộ bản thảo
□ Không lộ file nội bộ
□ Không lộ workflow kiểm duyệt
□ Permission hoạt động đúng
D.11.4 Danh sách văn bản
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes
SEO Class	SEO-B
Entity chính	Legal Document List
Entity phụ	Legal Document, Planning Project, Administrative Unit
URL Type	List URL
Mục tiêu SEO	Trang danh mục văn bản pháp lý quy hoạch, điều hướng đến các văn bản SEO-A
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/van-ban
•	URL theo địa bàn
/quy-hoach/{administrative-unit-slug}/van-ban
•	URL NoIndex
/quy-hoach/van-ban?filter={filter}
/quy-hoach/van-ban?sort={sort}
/quy-hoach/van-ban?keyword={keyword}
•	Canonical
•	Danh sách tổng canonical về /quy-hoach/van-ban. 
•	Danh sách địa bàn canonical về URL địa bàn tương ứng. 
•	Filter/sort/search runtime canonical về URL danh sách chuẩn. 
•	16.3 Metadata
•	Title
Danh sách văn bản pháp lý quy hoạch | QH Pro
•	Description
Tra cứu danh sách quyết định, nghị quyết, thông báo, văn bản phê duyệt, điều chỉnh và thay thế liên quan đến quy hoạch trên QH Pro.
•	H1
Danh sách văn bản pháp lý quy hoạch
•	Robots
index,follow
•	NoIndex khi
•	Danh sách rỗng. 
•	Chỉ có văn bản Restricted. 
•	Search/filter runtime không có URL chuẩn. 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	CollectionPage / ItemList
Entity Mapping	Legal Document List
Không inject Schema nếu danh sách rỗng, Restricted hoặc chỉ là filter runtime.
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Legal Document List Preview
OG Type	website
•	16.6 AI Summary
•	Có AI Summary
Yes.
•	Dữ liệu nguồn
•	Legal Document công khai. 
•	Planning Project. 
•	Administrative Unit. 
•	Timeline pháp lý. 
•	Nội dung AI Summary
•	Nhóm văn bản hiện có. 
•	Loại văn bản. 
•	Địa bàn áp dụng. 
•	Đồ án liên quan. 
•	Văn bản nổi bật. 
•	Không public
•	Văn bản chưa công bố. 
•	Văn bản Restricted. 
•	Ghi chú nội bộ. 
•	Trạng thái kiểm duyệt. 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	legal-document-list.xml
•	Tham gia khi
•	Public. 
•	Active. 
•	Có ít nhất một văn bản công khai. 
•	Có URL chuẩn. 
•	Loại bỏ khi
•	Empty List. 
•	Restricted. 
•	Runtime filter. 
•	Search result tạm thời. 
•	16.8 Internal Link
•	Link đến
•	Chi tiết văn bản. 
•	Chi tiết đồ án. 
•	Timeline pháp lý. 
•	Quan hệ phiên bản. 
•	Đơn vị hành chính. 
•	Link nhận từ
•	Workspace thư viện. 
•	Chi tiết đồ án. 
•	Chi tiết bản đồ. 
•	Timeline pháp lý. 
•	Report. 
•	16.9 Security & Visibility
•	Public
•	Số hiệu văn bản. 
•	Tên văn bản. 
•	Ngày ban hành. 
•	Cơ quan ban hành. 
•	Trạng thái hiệu lực công khai. 
•	Restricted
•	Văn bản chưa công bố. 
•	Văn bản có điều kiện truy cập. 
•	Tài liệu nội bộ đi kèm. 
•	Private
•	Ghi chú kiểm duyệt. 
•	Workflow xử lý. 
•	Import log. 
•	16.10 Cache & Regeneration
Event	Regenerate
legal_document.created	Metadata, Sitemap
legal_document.updated	Metadata, AI Summary
legal_document.visibility.changed	Robots, Sitemap
administrative_unit.updated	Metadata
planning_project.updated	AI Summary
•	16.11 Acceptance Criteria
•	AC-D11.4-001: Danh sách văn bản có URL chuẩn. 
•	AC-D11.4-002: Filter/sort/search runtime không tạo URL SEO riêng. 
•	AC-D11.4-003: Metadata sinh đúng phạm vi danh sách. 
•	AC-D11.4-004: ItemList chỉ chứa văn bản Public. 
•	AC-D11.4-005: Sitemap chỉ chứa danh sách đủ điều kiện. 
•	AC-D11.4-006: Không lộ văn bản Restricted hoặc ghi chú nội bộ. 
•	16.12 QA Checklist
□ URL danh sách đúng chuẩn
□ Canonical đúng
□ Filter/sort/search NoIndex hoặc canonical đúng
□ Title/Description/H1 tồn tại
□ Robots đúng
□ CollectionPage / ItemList hợp lệ
□ AI Summary không chứa văn bản Restricted
□ Sitemap đúng điều kiện
□ Không lộ ghi chú nội bộ
 
D.11.5 Chi tiết văn bản
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes
SEO Class	SEO-A
Entity chính	Legal Document
Entity phụ	Planning Project, Planning Map, Administrative Unit, Timeline
URL Type	Detail URL
Mục tiêu SEO	URL SEO trọng tâm cho văn bản pháp lý quy hoạch, phục vụ Google Search, AI Search, Legal SEO và Knowledge Graph
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/van-ban/{legal-document-slug}
•	URL phụ
/quy-hoach/do-an/{planning-project-slug}/van-ban/{legal-document-slug}
•	URL NoIndex
/quy-hoach/van-ban/{id}?preview=true
/quy-hoach/van-ban/{id}?download=true
/quy-hoach/van-ban/{id}?tab={tab}
•	Canonical
•	Văn bản công khai canonical về URL chuẩn. 
•	Văn bản thuộc đồ án vẫn canonical về URL văn bản chuẩn. 
•	Preview/download/tab canonical về URL văn bản chuẩn. 
•	Văn bản thay thế hoặc hết hiệu lực vẫn giữ URL riêng nếu có giá trị Historical SEO. 
•	16.3 Metadata
•	Title
{document_number} - {document_title} | QH Pro
•	Description
Thông tin văn bản {document_number}, gồm tên văn bản, ngày ban hành, cơ quan ban hành, trạng thái hiệu lực, đồ án quy hoạch và bản đồ liên quan.
•	H1
{document_title}
•	Robots
index,follow
•	NoIndex khi
•	Draft. 
•	Restricted. 
•	Deleted. 
•	Chưa công bố. 
•	Thiếu số hiệu hoặc tiêu đề văn bản. 
•	Không có trạng thái công khai. 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Legislation / CreativeWork
Entity Mapping	Legal Document
•	Schema bổ sung
Schema	Điều kiện
BreadcrumbList	Có breadcrumb
Place	Có địa bàn áp dụng
CreativeWork	Nếu không đủ điều kiện dùng Legislation
Không inject Schema nếu văn bản Restricted, Draft, thiếu dữ liệu pháp lý tối thiểu hoặc chưa công bố.
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Document Preview / Project Preview
OG Type	article
•	16.6 AI Summary
•	Có AI Summary
Yes.
•	Dữ liệu nguồn
•	Legal Document. 
•	Planning Project. 
•	Planning Map. 
•	Timeline pháp lý. 
•	Administrative Unit. 
•	Version Relation. 
•	Nội dung AI Summary
•	Tóm tắt văn bản. 
•	Số hiệu, ngày ban hành, cơ quan ban hành. 
•	Đối tượng áp dụng. 
•	Đồ án/bản đồ liên quan. 
•	Trạng thái hiệu lực. 
•	Quan hệ thay thế/điều chỉnh nếu có. 
•	Không public
•	Bản thảo. 
•	Văn bản chưa công bố. 
•	Ghi chú nội bộ. 
•	File hạn chế. 
•	Workflow kiểm duyệt. 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	legal-document.xml
•	Tham gia khi
•	Public. 
•	Active hoặc có giá trị Historical SEO. 
•	Có URL chuẩn. 
•	Có số hiệu/tên/ngày ban hành hợp lệ. 
•	Có nội dung công khai tối thiểu. 
•	Loại bỏ khi
•	Draft. 
•	Restricted. 
•	Deleted. 
•	Chưa công bố. 
•	Thiếu thông tin pháp lý tối thiểu. 
•	16.8 Internal Link
•	Link đến
•	Đồ án liên quan. 
•	Bản đồ liên quan. 
•	Timeline pháp lý. 
•	Quan hệ phiên bản. 
•	Địa bàn áp dụng. 
•	Văn bản thay thế / bị thay thế / điều chỉnh. 
•	Link nhận từ
•	Danh sách văn bản. 
•	Chi tiết đồ án. 
•	Chi tiết bản đồ. 
•	Timeline pháp lý. 
•	Quan hệ phiên bản. 
•	Report. 
•	16.9 Security & Visibility
•	Public
•	Số hiệu. 
•	Tên văn bản. 
•	Ngày ban hành. 
•	Cơ quan ban hành. 
•	Trạng thái hiệu lực. 
•	Nội dung công khai. 
•	File công khai. 
•	Restricted
•	File hạn chế. 
•	Văn bản chưa công bố. 
•	Tài liệu đi kèm chưa public. 
•	Private
•	Ghi chú nội bộ. 
•	Luồng xử lý. 
•	Import log. 
•	Kiểm duyệt nội bộ. 
•	16.10 Cache & Regeneration
Event	Regenerate
legal_document.created	Metadata, Sitemap
legal_document.updated	Metadata, AI Summary
legal_document.visibility.changed	Robots, Sitemap
legal_status.changed	Metadata, AI Summary
planning_project.updated	AI Summary, Internal Link
version_relation.updated	AI Summary, Internal Link
•	16.11 Acceptance Criteria
•	AC-D11.5-001: Mỗi văn bản Public có đúng một URL chuẩn. 
•	AC-D11.5-002: Metadata chứa đúng số hiệu, tên văn bản và trạng thái hiệu lực. 
•	AC-D11.5-003: Structured Data hợp lệ. 
•	AC-D11.5-004: AI Summary không dùng dữ liệu chưa công bố. 
•	AC-D11.5-005: Văn bản đủ điều kiện xuất hiện trong legal-document.xml. 
•	AC-D11.5-006: Internal Link đầy đủ tới đồ án, bản đồ, timeline và phiên bản. 
•	AC-D11.5-007: Không lộ file Restricted hoặc workflow nội bộ. 
•	16.12 QA Checklist
□ URL văn bản đúng chuẩn
□ Canonical đúng
□ Preview/download/tab không tạo URL SEO riêng
□ Title có số hiệu/tên văn bản
□ Description có ngày/cơ quan/trạng thái
□ Schema Legislation/CreativeWork hợp lệ
□ OG Preview an toàn
□ AI Summary đúng dữ liệu công khai
□ Có trong legal-document.xml khi đủ điều kiện
□ Không lộ file hạn chế hoặc ghi chú nội bộ
 
D.11.6 Danh sách bản đồ
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes
SEO Class	SEO-B
Entity chính	Planning Map List
Entity phụ	Planning Map, Planning Project, GIS Layer, Administrative Unit
URL Type	List URL
Mục tiêu SEO	Trang danh mục bản đồ quy hoạch, điều hướng đến các bản đồ SEO-A và hỗ trợ truy vấn “bản đồ quy hoạch” theo địa bàn/chủ đề
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/ban-do
•	URL theo địa bàn
/quy-hoach/{administrative-unit-slug}/ban-do
•	URL NoIndex
/quy-hoach/ban-do?filter={filter}
/quy-hoach/ban-do?sort={sort}
/quy-hoach/ban-do?keyword={keyword}
•	Canonical
•	Danh sách tổng canonical về /quy-hoach/ban-do. 
•	Danh sách địa bàn canonical về URL địa bàn. 
•	Filter/sort/search runtime canonical về URL chuẩn. 
•	16.3 Metadata
•	Title
Danh sách bản đồ quy hoạch | QH Pro
•	Description
Tra cứu danh sách bản đồ quy hoạch, bản đồ sử dụng đất, bản đồ giao thông, bản đồ phân khu và các lớp GIS công khai trên QH Pro.
•	H1
Danh sách bản đồ quy hoạch
•	Robots
index,follow
•	NoIndex khi
•	Danh sách rỗng. 
•	Chỉ có bản đồ Restricted. 
•	Search/filter runtime. 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	CollectionPage / ItemList
Entity Mapping	Planning Map List
Có thể bổ sung Dataset nếu danh sách đại diện tập dữ liệu bản đồ công khai.
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Planning Map List Preview
OG Type	website
•	16.6 AI Summary
•	Có AI Summary
Yes.
•	Dữ liệu nguồn
•	Planning Map công khai. 
•	Planning Project. 
•	GIS Layer. 
•	Administrative Unit. 
•	Legal Document. 
•	Nội dung AI Summary
•	Nhóm bản đồ hiện có. 
•	Loại bản đồ. 
•	Địa bàn áp dụng. 
•	Đồ án/văn bản liên quan. 
•	Các bản đồ nổi bật. 
•	Không public
•	Bản đồ chưa công bố. 
•	Bản đồ Restricted. 
•	Layer nội bộ. 
•	Ghi chú kiểm duyệt. 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	planning-map-list.xml
•	Tham gia khi
•	Public. 
•	Có ít nhất một bản đồ công khai. 
•	Có URL chuẩn. 
•	Có Metadata hợp lệ. 
•	Loại bỏ khi
•	Empty List. 
•	Restricted. 
•	Runtime filter/search. 
•	16.8 Internal Link
•	Link đến
•	Chi tiết bản đồ. 
•	Chi tiết đồ án. 
•	Chi tiết văn bản. 
•	Layer GIS. 
•	Đơn vị hành chính. 
•	Link nhận từ
•	Workspace thư viện. 
•	Chi tiết đồ án. 
•	Chi tiết văn bản. 
•	Trang bản đồ quy hoạch. 
•	Report. 
•	16.9 Security & Visibility
•	Public
•	Tên bản đồ. 
•	Loại bản đồ. 
•	Địa bàn. 
•	Preview công khai. 
•	Link chi tiết công khai. 
•	Restricted
•	Bản đồ trả phí. 
•	Bản đồ đối tác. 
•	Bản đồ chưa công bố. 
•	Private
•	File gốc nội bộ. 
•	Import log. 
•	Georeference note. 
•	Workflow kiểm duyệt. 
•	16.10 Cache & Regeneration
Event	Regenerate
planning_map.created	Metadata, Sitemap
planning_map.updated	AI Summary
planning_map.visibility.changed	Robots, Sitemap
map_preview.updated	Open Graph
planning_project.updated	AI Summary
•	16.11 Acceptance Criteria
•	AC-D11.6-001: Danh sách bản đồ có URL chuẩn. 
•	AC-D11.6-002: Filter/sort/search runtime không tạo URL SEO riêng. 
•	AC-D11.6-003: Metadata sinh đúng phạm vi. 
•	AC-D11.6-004: ItemList chỉ chứa bản đồ Public. 
•	AC-D11.6-005: Sitemap đúng điều kiện. 
•	AC-D11.6-006: Không lộ bản đồ Restricted hoặc file nội bộ. 
•	16.12 QA Checklist
□ URL danh sách đúng chuẩn
□ Canonical đúng
□ Filter/sort/search NoIndex hoặc canonical đúng
□ Title/Description/H1 tồn tại
□ ItemList hợp lệ
□ AI Summary không chứa bản đồ Restricted
□ Sitemap đúng điều kiện
□ Không lộ file gốc nội bộ
 
D.11.7 Chi tiết bản đồ
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes
SEO Class	SEO-A
Entity chính	Planning Map
Entity phụ	Planning Project, GIS Layer, Legal Document, Administrative Unit
URL Type	Detail URL
Mục tiêu SEO	URL SEO trọng tâm cho bản đồ quy hoạch, phục vụ truy vấn bản đồ theo địa bàn, loại quy hoạch, đồ án và dữ liệu GIS
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/ban-do/{planning-map-slug}
•	URL phụ
/quy-hoach/do-an/{planning-project-slug}/ban-do/{planning-map-slug}
•	URL NoIndex
/quy-hoach/ban-do/{id}?preview=true
/quy-hoach/ban-do/{id}?layer={layer-id}
/quy-hoach/ban-do/{id}?download=true
•	Canonical
•	Bản đồ công khai canonical về URL chuẩn. 
•	Bản đồ thuộc đồ án vẫn canonical về URL bản đồ chuẩn. 
•	Preview/download/layer state canonical về URL bản đồ chuẩn. 
•	16.3 Metadata
•	Title
{planning_map_name} | Bản đồ quy hoạch | QH Pro
•	Description
Thông tin bản đồ {planning_map_name}, gồm loại bản đồ, phạm vi áp dụng, đồ án liên quan, lớp dữ liệu GIS, văn bản pháp lý và trạng thái công khai.
•	H1
{planning_map_name}
•	Robots
index,follow
•	NoIndex khi
•	Draft. 
•	Restricted. 
•	Deleted. 
•	Chưa công bố. 
•	Thiếu preview hoặc metadata tối thiểu. 
•	Không có phạm vi áp dụng rõ ràng. 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	Dataset / CreativeWork
Entity Mapping	Planning Map
•	Schema bổ sung
Schema	Điều kiện
Place	Có phạm vi địa lý
BreadcrumbList	Có breadcrumb
Legislation / CreativeWork	Có văn bản pháp lý liên quan
Không inject Schema nếu bản đồ Restricted, Draft, thiếu metadata hoặc chưa công bố.
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Map Preview
OG Type	article
•	16.6 AI Summary
•	Có AI Summary
Yes.
•	Dữ liệu nguồn
•	Planning Map. 
•	GIS Layer. 
•	Planning Project. 
•	Legal Document. 
•	Administrative Unit. 
•	Layer Metadata. 
•	Nội dung AI Summary
•	Mô tả bản đồ. 
•	Loại bản đồ. 
•	Phạm vi áp dụng. 
•	Nguồn dữ liệu. 
•	Đồ án/văn bản liên quan. 
•	Trạng thái công khai. 
•	Cảnh báo dữ liệu nếu có. 
•	Không public
•	File gốc Restricted. 
•	Layer nội bộ. 
•	Georeference note nội bộ. 
•	Import log. 
•	Ghi chú kiểm duyệt. 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	planning-map.xml
•	Tham gia khi
•	Public. 
•	Active hoặc có giá trị Historical SEO. 
•	Có URL chuẩn. 
•	Có preview công khai. 
•	Có metadata hợp lệ. 
•	Loại bỏ khi
•	Draft. 
•	Restricted. 
•	Deleted. 
•	Chưa công bố. 
•	Missing Preview. 
•	Missing Metadata. 
•	16.8 Internal Link
•	Link đến
•	Đồ án liên quan. 
•	Văn bản pháp lý. 
•	Layer GIS. 
•	Đơn vị hành chính. 
•	Timeline pháp lý. 
•	Quan hệ phiên bản. 
•	Link nhận từ
•	Danh sách bản đồ. 
•	Chi tiết đồ án. 
•	Chi tiết văn bản. 
•	Trang bản đồ quy hoạch. 
•	Report. 
•	Snapshot. 
•	16.9 Security & Visibility
•	Public
•	Tên bản đồ. 
•	Preview công khai. 
•	Loại bản đồ. 
•	Phạm vi áp dụng. 
•	Metadata công khai. 
•	Restricted
•	File bản đồ trả phí. 
•	File đối tác. 
•	Bản đồ chưa công bố. 
•	Private
•	File gốc nội bộ. 
•	Georeference note. 
•	Import log. 
•	Validation log. 
•	Workflow xử lý. 
•	16.10 Cache & Regeneration
Event	Regenerate
planning_map.created	Metadata, Sitemap
planning_map.updated	Metadata, AI Summary
planning_map.visibility.changed	Robots, Sitemap
map_preview.updated	Open Graph
layer.updated	AI Summary
legal_document.updated	AI Summary
•	16.11 Acceptance Criteria
•	AC-D11.7-001: Mỗi bản đồ Public có URL chuẩn. 
•	AC-D11.7-002: Metadata sinh đầy đủ theo bản đồ. 
•	AC-D11.7-003: Schema Dataset/CreativeWork hợp lệ. 
•	AC-D11.7-004: Open Graph dùng đúng Map Preview. 
•	AC-D11.7-005: AI Summary chỉ dùng dữ liệu công khai. 
•	AC-D11.7-006: Bản đồ đủ điều kiện xuất hiện trong planning-map.xml. 
•	AC-D11.7-007: Không lộ file gốc nội bộ hoặc Layer Restricted. 
•	16.12 QA Checklist
□ URL chi tiết bản đồ đúng chuẩn
□ Canonical đúng
□ Preview/layer/download không tạo URL SEO riêng
□ Title/Description/H1 tồn tại
□ Schema hợp lệ
□ OG Image đúng Map Preview
□ AI Summary đúng metadata công khai
□ Có trong planning-map.xml khi đủ điều kiện
□ Không lộ file gốc / import log / georeference note

D.11.8 Timeline pháp lý
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes, có điều kiện
SEO Class	SEO-C
Entity chính	Legal Timeline
Entity phụ	Planning Project, Legal Document, Planning Map, Version Relation
URL Type	Timeline URL
Mục tiêu SEO	Hỗ trợ Historical SEO, AI Search và Knowledge Graph cho quá trình phê duyệt, điều chỉnh, thay thế, hết hiệu lực của đồ án/văn bản
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/do-an/{planning-project-slug}/timeline
•	URL phụ
/quy-hoach/van-ban/{legal-document-slug}/timeline
•	URL NoIndex
/quy-hoach/timeline?project={id}
/quy-hoach/timeline?event={event-id}
/quy-hoach/timeline?filter={filter}
•	Canonical
•	Timeline có giá trị độc lập canonical về URL Timeline. 
•	Timeline chỉ bổ trợ đồ án canonical về URL chi tiết đồ án. 
•	Event runtime canonical về Timeline hoặc Entity nguồn. 
•	16.3 Metadata
•	Title
Timeline pháp lý {entity_name} | QH Pro
•	Description
Dòng thời gian pháp lý của {entity_name}, bao gồm các mốc phê duyệt, điều chỉnh, thay thế, hết hiệu lực và các văn bản liên quan.
•	H1
Timeline pháp lý {entity_name}
•	Robots
index,follow
•	NoIndex khi
•	Timeline không có mốc công khai. 
•	Thiếu Entity nguồn. 
•	Chỉ là filter runtime. 
•	Dữ liệu Restricted. 
•	Timeline đang Draft. 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	ItemList / CreativeWork
Entity Mapping	Legal Timeline
•	Schema bổ sung
Schema	Điều kiện
Event	Mỗi mốc pháp lý đủ điều kiện
Legislation / CreativeWork	Có văn bản liên quan
BreadcrumbList	Có breadcrumb
Không inject Schema nếu timeline thiếu dữ liệu công khai hoặc có sự kiện Restricted.
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Timeline Preview / Project Preview
OG Type	article
•	16.6 AI Summary
•	Có AI Summary
Yes.
•	Dữ liệu nguồn
•	Legal Timeline. 
•	Planning Project. 
•	Legal Document. 
•	Planning Map. 
•	Version Relation. 
•	Nội dung AI Summary
•	Tóm tắt các mốc pháp lý. 
•	Văn bản phê duyệt. 
•	Văn bản điều chỉnh. 
•	Văn bản thay thế. 
•	Trạng thái hiệu lực hiện tại. 
•	Quan hệ với phiên bản trước/sau. 
•	Không public
•	Sự kiện chưa công bố. 
•	Ghi chú kiểm duyệt. 
•	Workflow nội bộ. 
•	Văn bản Restricted. 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	legal-timeline.xml
•	Tham gia khi
•	Public. 
•	Có Entity nguồn. 
•	Có tối thiểu một mốc công khai. 
•	Có URL chuẩn. 
•	Có Metadata hợp lệ. 
•	Loại bỏ khi
•	Empty Timeline. 
•	Restricted. 
•	Draft. 
•	Missing Entity Source. 
•	Runtime filter. 
•	16.8 Internal Link
•	Link đến
•	Chi tiết đồ án. 
•	Chi tiết văn bản. 
•	Chi tiết bản đồ. 
•	Quan hệ phiên bản. 
•	Văn bản thay thế / điều chỉnh. 
•	Link nhận từ
•	Chi tiết đồ án. 
•	Chi tiết văn bản. 
•	Chi tiết bản đồ. 
•	Danh sách đồ án. 
•	Report. 
•	16.9 Security & Visibility
•	Public
•	Mốc pháp lý công khai. 
•	Văn bản công khai. 
•	Trạng thái hiệu lực công khai. 
•	Restricted
•	Mốc chưa công bố. 
•	Văn bản hạn chế. 
•	Tài liệu nội bộ. 
•	Private
•	Ghi chú kiểm duyệt. 
•	Workflow. 
•	Audit log. 
•	16.10 Cache & Regeneration
Event	Regenerate
timeline.updated	Metadata, AI Summary
timeline_event.created	AI Summary, Sitemap
legal_document.updated	AI Summary
version_relation.updated	AI Summary, Internal Link
visibility.changed	Robots, Sitemap
•	16.11 Acceptance Criteria
•	AC-D11.8-001: Timeline đủ điều kiện có URL chuẩn. 
•	AC-D11.8-002: Runtime event/filter không tạo URL SEO riêng. 
•	AC-D11.8-003: Metadata thể hiện đúng Entity nguồn. 
•	AC-D11.8-004: AI Summary tóm tắt đúng mốc pháp lý công khai. 
•	AC-D11.8-005: Sitemap đúng điều kiện. 
•	AC-D11.8-006: Không lộ sự kiện hoặc văn bản Restricted. 
•	16.12 QA Checklist
□ URL Timeline đúng chuẩn
□ Event/filter runtime canonical đúng
□ Title/Description/H1 tồn tại
□ Schema ItemList/Event hợp lệ
□ AI Summary đúng mốc pháp lý
□ Có trong legal-timeline.xml khi đủ điều kiện
□ Không lộ mốc chưa công bố hoặc ghi chú kiểm duyệt
 
D.11.9 Quan hệ phiên bản
•	16. SEO
•	16.1 SEO Scope
Thuộc tính	Giá trị
SEO Enabled	Yes, có điều kiện
SEO Class	SEO-C
Entity chính	Version Relation
Entity phụ	Planning Project, Legal Document, Planning Map, Legal Timeline
URL Type	Version / Relation URL
Mục tiêu SEO	Hỗ trợ Historical SEO, Canonical Control, AI Search và Knowledge Graph đối với quan hệ kế thừa, điều chỉnh, thay thế và hết hiệu lực
•	16.2 URL & Canonical
•	URL chuẩn
/quy-hoach/do-an/{planning-project-slug}/phien-ban
•	URL quan hệ cụ thể
/quy-hoach/phien-ban/{version-relation-slug}
•	URL NoIndex
/quy-hoach/phien-ban?from={version-a}&to={version-b}
/quy-hoach/phien-ban?relation={relation-id}
/quy-hoach/phien-ban?preview=true
•	Canonical
•	Quan hệ phiên bản có giá trị độc lập canonical về URL version relation. 
•	Quan hệ chỉ bổ trợ đồ án canonical về URL chi tiết đồ án hoặc URL phiên bản của đồ án. 
•	Version selector runtime canonical về Entity nguồn. 
•	16.3 Metadata
•	Title
Quan hệ phiên bản {entity_name} | QH Pro
•	Description
Thông tin quan hệ phiên bản của {entity_name}, bao gồm phiên bản trước, phiên bản sau, nội dung điều chỉnh, thay thế, kế thừa và văn bản pháp lý liên quan.
•	H1
Quan hệ phiên bản {entity_name}
•	Robots
index,follow
•	NoIndex khi
•	Quan hệ chưa xác minh. 
•	Thiếu Entity nguồn. 
•	Thiếu văn bản pháp lý. 
•	Restricted. 
•	Draft. 
•	Preview/runtime. 
•	16.4 Structured Data
Thành phần	Giá trị
Schema	CreativeWork / ItemList
Entity Mapping	Version Relation
•	Schema bổ sung
Schema	Điều kiện
Legislation / CreativeWork	Có văn bản pháp lý
BreadcrumbList	Có breadcrumb
Không inject Schema nếu thiếu version nguồn, version đích hoặc văn bản pháp lý cần thiết.
•	16.5 Open Graph
Thành phần	Giá trị
OG Title	Theo Metadata Title
OG Description	Theo Metadata Description
OG Image	Version Relation Preview
OG Type	article
•	16.6 AI Summary
•	Có AI Summary
Yes.
•	Dữ liệu nguồn
•	Version Relation. 
•	Planning Project. 
•	Legal Document. 
•	Legal Timeline. 
•	Planning Map. 
•	GIS Layer công khai. 
•	Nội dung AI Summary
•	Phiên bản trước. 
•	Phiên bản sau. 
•	Loại quan hệ: điều chỉnh, thay thế, kế thừa, hết hiệu lực. 
•	Văn bản căn cứ. 
•	Tác động chính. 
•	Link tới Entity hiện hành và Entity lịch sử. 
•	Không public
•	Quan hệ chưa xác minh. 
•	Version Draft. 
•	Văn bản chưa công bố. 
•	Ghi chú nội bộ. 
•	Workflow kiểm duyệt. 
•	16.7 Sitemap
Thành phần	Giá trị
Sitemap	version-relation.xml
•	Tham gia khi
•	Public. 
•	Verified. 
•	Có URL chuẩn. 
•	Có Entity nguồn và Entity đích. 
•	Có nội dung mô tả công khai. 
•	Có Metadata hợp lệ. 
•	Loại bỏ khi
•	Draft. 
•	Unverified. 
•	Restricted. 
•	Missing Source Version. 
•	Missing Target Version. 
•	Runtime preview. 
•	16.8 Internal Link
•	Link đến
•	Đồ án hiện hành. 
•	Đồ án cũ. 
•	Văn bản điều chỉnh/thay thế. 
•	Timeline pháp lý. 
•	Bản đồ liên quan. 
•	Biến động quy hoạch. 
•	Link nhận từ
•	Chi tiết đồ án. 
•	Chi tiết văn bản. 
•	Timeline pháp lý. 
•	Biến động quy hoạch. 
•	Report. 
•	Snapshot. 
•	16.9 Security & Visibility
•	Public
•	Quan hệ phiên bản công khai. 
•	Văn bản căn cứ. 
•	Trạng thái hiệu lực. 
•	Entity hiện hành và lịch sử. 
•	Restricted
•	Phiên bản chưa công bố. 
•	Tài liệu hạn chế. 
•	Quan hệ đang kiểm duyệt. 
•	Private
•	Ghi chú nội bộ. 
•	Workflow kiểm duyệt. 
•	Import log. 
•	Audit log. 
•	16.10 Cache & Regeneration
Event	Regenerate
version_relation.created	Metadata, Sitemap
version_relation.updated	Metadata, AI Summary
version_relation.verified	Robots, Sitemap
legal_document.updated	AI Summary
planning_project.updated	AI Summary, Internal Link
visibility.changed	Robots, Sitemap
•	16.11 Acceptance Criteria
•	AC-D11.9-001: Quan hệ phiên bản public/verified có URL chuẩn. 
•	AC-D11.9-002: Runtime selector không tạo URL SEO riêng. 
•	AC-D11.9-003: Canonical đúng về quan hệ phiên bản hoặc Entity nguồn. 
•	AC-D11.9-004: AI Summary nêu đúng phiên bản trước/sau và loại quan hệ. 
•	AC-D11.9-005: Sitemap chỉ chứa quan hệ đủ điều kiện. 
•	AC-D11.9-006: Không lộ version draft hoặc ghi chú kiểm duyệt. 
•	AC-D11.9-007: Internal Link đầy đủ giữa phiên bản cũ, phiên bản mới và văn bản pháp lý. 
•	16.12 QA Checklist
□ URL quan hệ phiên bản đúng chuẩn
□ Runtime selector NoIndex hoặc canonical đúng
□ Title/Description/H1 tồn tại
□ Schema hợp lệ
□ AI Summary đúng quan hệ trước/sau
□ Có link tới đồ án cũ/mới và văn bản liên quan
□ Có trong version-relation.xml khi đủ điều kiện
□ Không lộ version draft / ghi chú nội bộ

 
 
PHẦN E.  CHECKLIST TRIỂN KHAI & NGHIỆM THU
Phần này dùng để kiểm tra việc triển khai SEO sau khi các yêu cầu tại Phần D đã được bổ sung vào từng màn hình SRS và được Dev triển khai trên hệ thống.
Mục tiêu của Phần E là bảo đảm:
•	Không bỏ sót SEO Output. 
•	Không Index nhầm dữ liệu không được phép công khai. 
•	Không tạo URL trùng lặp hoặc URL rác. 
•	Metadata, Canonical, Sitemap, Structured Data và AI Summary được sinh đúng. 
•	Frontend, Backend, Mobile và QA có cùng tiêu chí nghiệm thu. 
 
E.1	Checklist Dev
•	E.1.1 URL & Canonical
□ Mỗi Entity SEO có đúng một URL chuẩn.
□ URL chuẩn đúng định dạng đã đặc tả tại Phần D.
□ URL filter, sort, tab, map state, session, token không được Index.
□ Canonical luôn trỏ về URL chuẩn hoặc Entity nguồn đúng.
□ Không có URL trùng lặp cho cùng một Entity.
□ URL lịch sử, URL trước sáp nhập, URL tên cũ được redirect/canonical đúng.
•	E.1.2 SEO Output
□ Metadata được sinh cho tất cả URL SEO-A, SEO-B, SEO-C đủ điều kiện.
□ Robots đúng theo trạng thái Public / Restricted / Private / Draft / Deleted.
□ Structured Data chỉ inject khi đủ điều kiện.
□ Open Graph có Safe Preview, không lộ dữ liệu hạn chế.
□ AI Summary chỉ sinh từ dữ liệu công khai.
□ Sitemap chỉ chứa URL đủ điều kiện Index.
•	E.1.3 Security
□ Không đưa dữ liệu Restricted vào Metadata.
□ Không đưa dữ liệu Private vào Schema.
□ Không đưa dữ liệu nội bộ vào AI Summary.
□ Không đưa token, session, user id, API key, log vào bất kỳ SEO Output nào.
□ Admin, Account, Internal API, Partner API luôn SEO-N.
 
E.2	Checklist Frontend
•	E.2.1 Head Tags
□ Title render đúng.
□ Description render đúng.
□ Canonical render đúng.
□ Robots render đúng.
□ Open Graph render đúng.
□ Structured Data JSON-LD render đúng điều kiện.
•	E.2.2 Render Mode
□ Nội dung SEO quan trọng có trong HTML đầu tiên hoặc SSR/Hybrid theo yêu cầu.
□ Không phụ thuộc hoàn toàn vào client runtime đối với trang SEO-A.
□ Trang NoIndex vẫn render đúng robots.
□ Trang lỗi, rỗng, restricted không bị Index nhầm.
•	E.2.3 UI hỗ trợ SEO
□ Breadcrumb hiển thị đúng.
□ Internal Link đến Entity liên quan đầy đủ.
□ AI Summary block hiển thị đúng nếu được phép.
□ SEO text block không lộ dữ liệu hạn chế.
□ Share Preview lấy đúng dữ liệu public.
 
E.3	Checklist Backend
•	E.3.1 SEO Payload
□ Backend trả đủ SEO Payload cho Frontend.
□ Payload có title, description, h1, canonical, robots.
□ Payload có schema khi đủ điều kiện.
□ Payload có Open Graph khi cần chia sẻ.
□ Payload có AI Summary khi Entity đủ điều kiện.
□ Payload có trạng thái sitemap eligibility.
•	E.3.2 Sitemap
□ Sitemap tách đúng nhóm: parcel, planning-project, legal-document, planning-map, report, snapshot, alert, timeline.
□ Chỉ URL Public/Active/Valid mới vào Sitemap.
□ Draft, Deleted, Restricted, Private bị loại khỏi Sitemap.
□ Khi visibility thay đổi, Sitemap được cập nhật.
□ Khi Entity bị xóa hoặc chuyển Restricted, URL bị loại khỏi Sitemap.
•	E.3.3 Regeneration
□ Entity update → regenerate Metadata.
□ Geometry update → regenerate Schema / OG Preview.
□ Legal status update → regenerate Metadata / AI Summary.
□ Visibility change → regenerate Robots / Sitemap.
□ Map Preview update → regenerate Open Graph.
 
E.4	Checklist Mobile
•	E.4.1 Deep Link
□ Mobile mở đúng Entity từ URL Web.
□ Share Link mở đúng màn hình tương ứng.
□ Snapshot / Report share mở đúng trạng thái được phép.
□ Token share không làm lộ dữ liệu vượt quyền.
•	E.4.2 Share Preview
□ Share từ Mobile dùng đúng URL chuẩn hoặc Share URL.
□ Không share URL nội bộ, session URL, token không cần thiết.
□ Preview không lộ dữ liệu Private / Restricted.
□ Nếu nội dung không public, Share URL phải NoIndex.
•	E.4.3 Đồng bộ Web / Mobile
□ Entity ID đồng nhất giữa Web và Mobile.
□ URL Web là nguồn SEO chính.
□ Mobile không sinh URL SEO riêng.
□ Deep Link không tạo duplicate URL.

E.5	Checklist QA
•	E.5.1 URL
□ URL đúng chuẩn.
□ URL duy nhất cho mỗi Entity.
□ Canonical đúng.
□ URL filter/sort/tab/session/token không Index.
□ Legacy URL redirect/canonical đúng.
•	E.5.2 Metadata
□ Title tồn tại.
□ Description tồn tại.
□ H1 tồn tại.
□ Robots đúng.
□ Metadata không chứa dữ liệu Restricted / Private.
•	E.5.3 Structured Data
□ JSON-LD hợp lệ.
□ Schema đúng Entity.
□ Không inject Schema khi Draft / Restricted / Missing Data.
□ BreadcrumbList đúng nếu có.
•	E.5.4 Open Graph
□ OG Title đúng.
□ OG Description đúng.
□ OG Image đúng.
□ Share Preview hiển thị an toàn.
□ Không lộ token, user id, dữ liệu cá nhân.
•	E.5.5 AI Summary
□ AI Summary sinh đúng Entity.
□ Chỉ dùng dữ liệu công khai.
□ Không dùng raw layer nội bộ.
□ Không lộ ghi chú nội bộ, workflow, log, dữ liệu chưa công bố.
•	E.5.6 Sitemap
□ URL đủ điều kiện có trong Sitemap.
□ URL không đủ điều kiện bị loại.
□ Sitemap group đúng.
□ Khi trạng thái thay đổi, Sitemap cập nhật đúng.
•	E.5.7 Security
□ Admin luôn NoIndex.
□ Account luôn NoIndex.
□ API luôn NoIndex.
□ Restricted không xuất hiện trong SEO Output.
□ Private không xuất hiện trong SEO Output.
 
E.6	Checklist Go-Live SEO
•	E.6.1 Trước Go-Live
□ Đã rà soát toàn bộ QH.1 → QH.11.
□ Đã gán SEO Class cho toàn bộ màn hình.
□ Đã gán SEO Class cho toàn bộ Entity.
□ Đã gán Index/NoIndex cho toàn bộ URL Type.
□ Đã kiểm tra Sitemap.
□ Đã kiểm tra robots.
□ Đã kiểm tra Canonical.
□ Đã kiểm tra Structured Data.
□ Đã kiểm tra Open Graph.
□ Đã kiểm tra AI Summary.
•	E.6.2 Kiểm tra nhóm SEO-A
□ Administrative Unit.
□ Parcel.
□ Planning Region.
□ Planning Project.
□ Legal Document.
□ Planning Map.
□ Report public.
□ Planning Change public.
•	E.6.3 Kiểm tra nhóm SEO-N
□ Account.
□ Admin.
□ Internal API.
□ Partner API.
□ Runtime search.
□ GPS runtime.
□ Filter / sort / tab / session URL.
□ Token share URL.
•	E.6.4 Điều kiện cho phép Go-Live
Hệ thống chỉ được Go-Live SEO khi:
•	Không có URL Private/Restricted trong Sitemap. 
•	Không có Admin/API/Account URL bị Index. 
•	Tất cả SEO-A có Metadata, Canonical, Robots và Sitemap đúng. 
•	Structured Data không lỗi nghiêm trọng. 
•	Open Graph không lộ dữ liệu hạn chế. 
•	AI Summary chỉ dùng dữ liệu công khai. 
•	QA xác nhận đầy đủ checklist theo từng nhóm màn hình.
