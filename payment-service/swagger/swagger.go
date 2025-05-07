package swagger_comment

// @Summary Tạo tài khoản
// @Tags Account
// @Accept json
// @Produce json
// @Param publish body domain.Account true "Thông tin tài khoản"
// @Router /wallet/new [post]
func AccountRegister() {}

// @Summary Cập nhật tài khoản
// @Tags Account
// @Accept json
// @Produce json
// @Param publish body domain.Account true "Thông tin tài khoản"
// @Router /wallet/:id [put]
func AccountUpdate() {}

// @Summary [LOCAL]Tạo giao dịch mới
// @Tags Account
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param publish body domain.Transaction true "Thông tin tài khoản"
// @Router /wallet/transaction [post]
func MakeTransaction() {}

// @Summary Cập nhật tài khoản
// @Tags Account
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /wallet/balance [get]
func Balance() {}
