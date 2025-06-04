package main

// // @Summary Đăng tin sản phẩm
// // @Description API để đăng tin sản phẩm lên hệ thống
// // @Tags User/Product-search
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param publish body data.StatusRequest true "Thông tin đăng tin"
// // @Router /v1/bdspro/user/product/internal [get]
// func ProductInternal() {}

// // @Summary Đăng tin sản phẩm
// // @Tags User/Post
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param post body u_post.PostSaveRequest true "Thông tin đăng tin"
// // @Router /v1/bdspro/user/post [post]
// func PostCreate() {}

// // @Summary Cập nhật tin
// // @Tags User/Post
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int64 true "ID bài đăng"
// // @Param post body u_post.PostUpdateRequest true "Thông tin đăng tin"
// // @Router /v1/bdspro/user/post/{id} [put]
// func PostUpdate() {}

// // @Summary Gia hạn tin
// // @Description API để gia hạn
// // @Tags User/Post
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param post body u_post.PostExpiredRequest true "Thông tin gia hạn"
// // @Router /v1/bdspro/user/post/new-expired [post]
// func PostMakeNewExpired() {}

// // @Summary Xem chi tiết
// // @Tags User/Post
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int64 true "ID bài đăng"
// // @Router /v1/bdspro/user/post/{id} [get]
// func PostDetail() {}

// // @Summary Xóa
// // @Tags User/Post
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int64 true "ID bài đăng"
// // @Router /v1/bdspro/user/post/{id} [delete]
// func PostDelete() {}

// // @Summary Danh sách tin của tôi
// // @Tags User/Post
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param dto query u_post.SearchRequest true "Tham số"
// // @Router /v1/bdspro/user/post/personal [get]
// func PostPerson() {}

// // @Summary Danh sách tin của tôi
// // @Tags User/Post
// // @Accept json
// // @Produce json
// // @Param dto query u_post.SearchRequest true "Tham số"
// // @Router /v1/bdspro/user/post/global [get]
// func PostGlobal() {}

// // @Summary Thêm sự kiện
// // @Tags User/Action
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param dto body u_action.ActionRequest true "Tham số"
// // @Router /v1/bdspro/user/action [post]
// func PostAction() {}

// // @Summary Đếm sự kiện của bài viết
// // @Tags User/Action
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path uint64 true "ID"
// // @Router /v1/bdspro/user/action/dashboard/post/{id} [get]
// func PostDashboard() {}

// // @Summary Đếm sự kiện của bài viết
// // @Tags User/Action
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param query query u_action.PostBarRequest true "Query"
// // @Router /v1/bdspro/user/action/dashboard/post/bar [get]
// func PostDashboardBar() {}

// // @Summary Lấy danh sách tài sản
// // @Description Trả về danh sách tất cả tài sản hiện có trong hệ thống
// // @Tags Assets
// // @Produce json
// // @Param page query int true "Page"
// // @Param size query int true "Size"
// // @Security BearerAuth
// // @Router /v2/bdspro/v2/assets [get]
// func GetAssets() {}

// // @Summary Lưu trữ tài sản
// // @Tags Assets
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path uint64 true "ID"
// // @Param body body pb_asset.ArchivedRequest false "Body"
// // @Router /v2/bdspro/v2/assets/archived/{id} [put]
// func FilterAssets() {}

// // // @Summary Thao tác hàng loạt với tài sản
// // // @Description Thực hiện hành động hàng loạt như xóa, cập nhật trạng thái cho nhiều tài sản
// // // @Tags Assets
// // // @Accept json
// // // @Produce json
// // // @Security BearerAuth
// // // @Param bulk body dto.BulkAssetActionRequest true "Thông tin thao tác hàng loạt"
// // // @Success 200 {object} data.StatusResponse
// // // @Router /v2/bdspro/v2/assets/bulk-action [post]
// // func BulkActionAssets() {}

// // @Summary Tạo tài sản mới
// // @Description API dùng để tạo mới một tài sản trong hệ thống
// // @Tags Assets
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param asset body pb_asset.CreateAssetRequest true "Thông tin tài sản"
// // @Router /v2/bdspro/v2/assets/create [post]
// func CreateAsset(pb_asset.CreateAssetRequest) {}

// // @Summary Cập nhật tài sản
// // @Description Cập nhật thông tin một tài sản theo ID
// // @Tags Assets
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path string true "ID tài sản"
// // @Param asset body pb_asset.UpdateAssetRequest true "Thông tin cập nhật"
// // @Router /v2/bdspro/v2/assets/{id} [put]
// func UpdateAsset() {}

// // @Summary Xóa tài sản
// // @Description Xóa một tài sản khỏi hệ thống
// // @Tags Assets
// // @Produce json
// // @Security BearerAuth
// // @Param id path string true "ID tài sản"
// // @Success 200 {object} data.StatusResponse
// // @Router /v2/bdspro/v2/assets/{id} [delete]
// func DeleteAsset() {}

// // @Summary Lấy chi tiết tài sản
// // @Description Lấy thông tin chi tiết của một tài sản theo ID
// // @Tags Assets
// // @Produce json
// // @Security BearerAuth
// // @Param id path string true "ID tài sản"
// // @Success 200 {object} dto.AssetDetailResponse
// // @Router /v2/bdspro/v2/assets/{id}/details [get]
// func GetAssetDetails() {}

// // @Summary Chia sẻ tài sản
// // @Description Chia sẻ quyền quản lý tài sản cho người khác
// // @Tags Assets
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path string true "ID tài sản"
// // @Param share body pb_asset.ShareAssetRequest true "Thông tin chia sẻ"
// // @Success 200 {object} data.StatusResponse
// // @Router /v2/bdspro/v2/assets/{id}/share [post]
// func ShareAsset() {}

// // @Summary Tách tài sản
// // @Description Tách tài sản hiện tại thành nhiều tài sản nhỏ hơn
// // @Tags Assets
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path string true "ID tài sản"
// // @Param split body pb_asset.SplitAssetRequest true "Thông tin tách tài sản"
// // @Success 200 {object} data.StatusResponse
// // @Router /v2/bdspro/v2/assets/{id}/split [post]
// func SplitAsset() {}

// // @Summary Gộp tài sản
// // @Description Gộp nhiều tài sản thành một tài sản chung
// // @Tags Assets
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path string true "ID tài sản chính"
// // @Param merge body pb_asset.MergeAssetRequest true "Thông tin gộp tài sản"
// // @Success 200 {object} data.StatusResponse
// // @Router /v2/bdspro/v2/assets/{id}/merge [post]
// func MergeAsset() {}

// // @Summary Lịch sử tài sản
// // @Description Lấy lịch sử thay đổi của một tài sản
// // @Tags Assets
// // @Produce json
// // @Security BearerAuth
// // @Param id path string true "ID tài sản"
// // @Router /v2/bdspro/v2/assets/{id}/history [get]
// func GetAssetHistory() {}
