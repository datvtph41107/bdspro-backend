package service

import (
	_errors "common/errors"

	"google.golang.org/grpc/codes"
)

var (
	AuthModeInvalid                       = _errors.MustSpec(210001, "USER_AUTH_MODE_INVALID", "Mode phải là auth hoặc verify", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PhoneInvalid                          = _errors.MustSpec(210002, "USER_PHONE_INVALID", "Số điện thoại không đúng", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	FullNameRequiredForRegistration       = _errors.MustSpec(210003, "USER_FULL_NAME_REQUIRED_FOR_REGISTRATION", "Vui lòng nhập họ tên khi đăng ký", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	AuthIDRequired                        = _errors.MustSpec(210004, "USER_AUTH_ID_REQUIRED", "AuthId không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	AuthIDInvalid                         = _errors.MustSpec(210005, "USER_AUTH_ID_INVALID", "AuthId không đúng", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PhoneOrIDRequired                     = _errors.MustSpec(210006, "USER_PHONE_OR_ID_REQUIRED", "Phone/ID không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	AccountLocked                         = _errors.MustSpec(210007, "USER_ACCOUNT_LOCKED", "Tài khoản của bạn đã bị khóa. Vui lòng liên hệ quản trị viên để được hỗ trợ", codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403))
	LoginRequired                         = _errors.MustSpec(210008, "USER_LOGIN_REQUIRED", "Vui lòng đăng nhập để thực hiện thao tác này", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	AccountInfoNotFound                   = _errors.MustSpec(210009, "USER_ACCOUNT_INFO_NOT_FOUND", "Không tìm thấy thông tin tài khoản", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	CurrentAccountProfileMissing          = _errors.MustSpec(210010, "USER_CURRENT_ACCOUNT_PROFILE_MISSING", "Tài khoản hiện tại chưa có profile", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PhoneNotFound                         = _errors.MustSpec(210012, "USER_PHONE_NOT_FOUND", "Số điện thoại không tồn tại", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	AuthKeyNotInitialized                 = _errors.MustSpec(210013, "USER_AUTH_KEY_NOT_INITIALIZED", "Tài khoản chưa được khởi tạo auth key. Vui lòng đăng ký lại", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	ClientPublicKeyInvalid                = _errors.MustSpec(210015, "USER_CLIENT_PUBLIC_KEY_INVALID", "Client public key không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	AuthKeyInvalid                        = _errors.MustSpec(210016, "USER_AUTH_KEY_INVALID", "Auth key không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	AuthKeyMismatch                       = _errors.MustSpec(210017, "USER_AUTH_KEY_MISMATCH", "Auth key không đúng", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	AccountLockedContactAdmin             = _errors.MustSpec(210018, "USER_ACCOUNT_LOCKED_CONTACT_ADMIN", "Tài khoản của bạn đã bị khóa. Vui lòng liên hệ quản trị viên", codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403))
	UserInfoNotFound                      = _errors.MustSpec(210019, "USER_USER_INFO_NOT_FOUND", "Không tìm thấy thông tin người dùng", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	SessionInvalid                        = _errors.MustSpec(210020, "USER_SESSION_INVALID", "Phiên làm việc không hợp lệ. Vui lòng đăng nhập lại", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	SessionEnded                          = _errors.MustSpec(210021, "USER_SESSION_ENDED", "Phiên làm việc đã kết thúc. Vui lòng đăng nhập lại", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	SessionAuthenticationFailed           = _errors.MustSpec(210022, "USER_SESSION_AUTHENTICATION_FAILED", "Không thể xác thực phiên đăng nhập", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	SessionIDRequiredForDeviceLogout      = _errors.MustSpec(210023, "USER_SESSION_ID_REQUIRED_FOR_DEVICE_LOGOUT", "sessionId là bắt buộc khi logout theo thiết bị", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	OTPLengthInvalid                      = _errors.MustSpec(210024, "USER_OTP_LENGTH_INVALID", "OTP không đủ độ dài", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	AccountNotFound                       = _errors.MustSpec(210025, "USER_ACCOUNT_NOT_FOUND", "Tài khoản không tồn tại", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	UsernamePasswordRequired              = _errors.MustSpec(210030, "USER_USERNAME_PASSWORD_REQUIRED", "Username và password không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	UsernamePasswordInvalid               = _errors.MustSpec(210031, "USER_USERNAME_PASSWORD_INVALID", "Username hoặc password không đúng", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	PasswordConfirmationMismatch          = _errors.MustSpec(210032, "USER_PASSWORD_CONFIRMATION_MISMATCH", "Mật khẩu xác nhận không khớp", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PasswordTooShort                      = _errors.MustSpec(210033, "USER_PASSWORD_TOO_SHORT", "Mật khẩu phải có ít nhất 8 ký tự", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PasswordAlreadySet                    = _errors.MustSpec(210034, "USER_PASSWORD_ALREADY_SET", "Bạn đã có mật khẩu. Vui lòng sử dụng chức năng đổi mật khẩu nếu muốn thay đổi", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PhonePasswordInvalid                  = _errors.MustSpec(210037, "USER_PHONE_PASSWORD_INVALID", "Số điện thoại hoặc mật khẩu không đúng", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	PasswordNotSetUseOTP                  = _errors.MustSpec(210038, "USER_PASSWORD_NOT_SET_USE_OTP", "Bạn chưa thiết lập mật khẩu. Vui lòng đăng nhập bằng OTP", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	UsernameOrPasswordInvalid             = _errors.MustSpec(210039, "USER_USERNAME_OR_PASSWORD_INVALID", "Username hoặc mật khẩu không đúng", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	LoginMethodInvalid                    = _errors.MustSpec(210040, "USER_LOGIN_METHOD_INVALID", "Phương thức đăng nhập không đúng", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	TokenInvalidOrExpired                 = _errors.MustSpec(210041, "USER_TOKEN_INVALID_OR_EXPIRED", "Token không hợp lệ hoặc đã hết hạn", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	RefreshTokenRequired                  = _errors.MustSpec(210042, "USER_REFRESH_TOKEN_REQUIRED", "Token không phải là refresh token", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	SessionExpired                        = _errors.MustSpec(210043, "USER_SESSION_EXPIRED", "Phiên làm việc đã hết hạn, vui lòng đăng nhập lại", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	LoginInfoNotFound                     = _errors.MustSpec(210044, "USER_LOGIN_INFO_NOT_FOUND", "Không tìm thấy thông tin đăng nhập", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	QRSessionInvalid                      = _errors.MustSpec(210045, "USER_QR_SESSION_INVALID", "Thông tin phiên QR không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	LoginSessionIncorrect                 = _errors.MustSpec(210046, "USER_LOGIN_SESSION_INCORRECT", "Phiên đăng nhập không chính xác", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	FullNameRequired                      = _errors.MustSpec(210047, "USER_FULL_NAME_REQUIRED", "Họ tên không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	UserAuthenticationFailed              = _errors.MustSpec(210048, "USER_USER_AUTHENTICATION_FAILED", "Không thể xác thực người dùng", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	LoginSessionNotFound                  = _errors.MustSpec(210049, "USER_LOGIN_SESSION_NOT_FOUND", "Không tìm thấy thông tin phiên đăng nhập", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	UsernameAlreadyExists                 = _errors.MustSpec(210052, "USER_USERNAME_ALREADY_EXISTS", "Username đã tồn tại", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	EmailAlreadyExists                    = _errors.MustSpec(210053, "USER_EMAIL_ALREADY_EXISTS", "Email đã tồn tại", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	AccessDenied                          = _errors.MustSpec(210058, "USER_ACCESS_DENIED", "Không có quyền truy cập", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	AccountNotAdmin                       = _errors.MustSpec(210059, "USER_ACCOUNT_NOT_ADMIN", "Tài khoản không phải admin", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	AdminAccountNotFound                  = _errors.MustSpec(210060, "USER_ADMIN_ACCOUNT_NOT_FOUND", "Không tìm thấy tài khoản admin", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	SelfAccountDeleteDenied               = _errors.MustSpec(210061, "USER_SELF_ACCOUNT_DELETE_DENIED", "Không thể xóa tài khoản của chính mình", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	OldPasswordInvalid                    = _errors.MustSpec(210064, "USER_OLD_PASSWORD_INVALID", "Mật khẩu cũ không đúng", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	UsernameRequired                      = _errors.MustSpec(210069, "USER_USERNAME_REQUIRED", "Username không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PasswordRequired                      = _errors.MustSpec(210070, "USER_PASSWORD_REQUIRED", "Password không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	AdminDataStoreUnavailable             = _errors.MustSpec(210074, "USER_ADMIN_DATA_STORE_UNAVAILABLE", "Kho dữ liệu quản trị chưa sẵn sàng", codes.Unavailable, _errors.LegacyHTTP200(), _errors.LegacyCode(503))
	RootAdminMutationDenied               = _errors.MustSpec(210076, "USER_ROOT_ADMIN_MUTATION_DENIED", "Không thể thay đổi tài khoản quản trị gốc qua API", codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403))
	EmailAlreadyExistsInSystem            = _errors.MustSpec(210078, "USER_EMAIL_ALREADY_EXISTS_IN_SYSTEM", "Email đã tồn tại trong hệ thống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	UsernameAlreadyExistsInSystem         = _errors.MustSpec(210079, "USER_USERNAME_ALREADY_EXISTS_IN_SYSTEM", "Tên đăng nhập đã tồn tại trong hệ thống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PasswordRequiredVI                    = _errors.MustSpec(210080, "USER_PASSWORD_REQUIRED_VI", "Mật khẩu là bắt buộc", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	UserNotFound                          = _errors.MustSpec(210085, "USER_USER_NOT_FOUND", "Không tìm thấy người dùng", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	PINLengthInvalid                      = _errors.MustSpec(210089, "USER_PIN_LENGTH_INVALID", "Mã PIN phải có đúng 6 số", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	OTPInvalid                            = _errors.MustSpec(210090, "USER_OTP_INVALID", "OTP không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	AccountNotFoundGeneric                = _errors.MustSpec(210091, "USER_ACCOUNT_NOT_FOUND_GENERIC", "Không tìm thấy tài khoản", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	PINAlreadyExists                      = _errors.MustSpec(210092, "USER_PIN_ALREADY_EXISTS", "Mã PIN đã tồn tại. Vui lòng sử dụng chức năng cập nhật PIN", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	NewPINLengthInvalid                   = _errors.MustSpec(210098, "USER_NEW_PIN_LENGTH_INVALID", "Mã PIN mới phải có đúng 6 số", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PINNotSet                             = _errors.MustSpec(210099, "USER_PIN_NOT_SET", "Chưa thiết lập mã PIN", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	OldPINInvalid                         = _errors.MustSpec(210100, "USER_OLD_PIN_INVALID", "Mã PIN cũ không đúng", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	PINDisabled                           = _errors.MustSpec(210104, "USER_PIN_DISABLED", "Mã PIN đã bị vô hiệu hóa", codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403))
	UserStatusNotFound                    = _errors.MustSpec(210106, "USER_USER_STATUS_NOT_FOUND", "Không tìm thấy trạng thái người dùng", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	ValidOTPNotFound                      = _errors.MustSpec(210107, "USER_VALID_OTP_NOT_FOUND", "Không tìm thấy OTP hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	RoleNameRequired                      = _errors.MustSpec(210108, "USER_ROLE_NAME_REQUIRED", "name không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PermissionIDsRequired                 = _errors.MustSpec(210109, "USER_PERMISSION_I_DS_REQUIRED", "permissionIds không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	RootAdminRoleBootstrapOnly            = _errors.MustSpec(210110, "USER_ROOT_ADMIN_ROLE_BOOTSTRAP_ONLY", "Vai trò quản trị gốc chỉ được tạo bởi bootstrap operator", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	RoleNotFound                          = _errors.MustSpec(210111, "USER_ROLE_NOT_FOUND", "Role không tồn tại", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	RoleNameRequiredVI                    = _errors.MustSpec(210112, "USER_ROLE_NAME_REQUIRED_VI", "Tên role không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	RoleKeyRequired                       = _errors.MustSpec(210113, "USER_ROLE_KEY_REQUIRED", "Key role không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	RootAdminRoleMutationDenied           = _errors.MustSpec(210114, "USER_ROOT_ADMIN_ROLE_MUTATION_DENIED", "Không thể thay đổi vai trò quản trị gốc", codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403))
	RootAdminRoleDeleteDenied             = _errors.MustSpec(210115, "USER_ROOT_ADMIN_ROLE_DELETE_DENIED", "Không thể xóa vai trò quản trị gốc", codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403))
	RootAdminRolePermissionMutationDenied = _errors.MustSpec(210116, "USER_ROOT_ADMIN_ROLE_PERMISSION_MUTATION_DENIED", "Không thể thay đổi quyền của vai trò quản trị gốc", codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403))
	OrganizationMemberNotFound            = _errors.MustSpec(210117, "USER_ORGANIZATION_MEMBER_NOT_FOUND", "Organization member not found", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	UserIDRequired                        = _errors.MustSpec(210118, "USER_USER_ID_REQUIRED", "UserID không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	RoleIDRequired                        = _errors.MustSpec(210119, "USER_ROLE_ID_REQUIRED", "RoleID không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	RoleByIDNotFound                      = _errors.MustSpec(210120, "USER_ROLE_BY_ID_NOT_FOUND", "Không tìm thấy role với roleId đã cung cấp", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	RoleAPIAssignmentDenied               = _errors.MustSpec(210121, "USER_ROLE_API_ASSIGNMENT_DENIED", "Vai trò này không cho phép gán qua API", codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403))
	ModuleInvalid                         = _errors.MustSpec(210122, "USER_MODULE_INVALID", "Module không đúng", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	IPAddressAccountLimitReached          = _errors.MustSpec(210123, "USER_IP_ADDRESS_ACCOUNT_LIMIT_REACHED", "IP này đã được gắn với tối đa 3 tài khoản", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	AccessGrantNotFound                   = _errors.MustSpec(210127, "USER_ACCESS_GRANT_NOT_FOUND", "Không tìm thấy quyền truy cập", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	UserIDOrIPAddressRequired             = _errors.MustSpec(210132, "USER_USER_ID_OR_IP_ADDRESS_REQUIRED", "Cần cung cấp UserID hoặc IPAddress", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	UserNotFoundProfile                   = _errors.MustSpec(210134, "USER_USER_NOT_FOUND_PROFILE", "Người dùng không tồn tại", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	EmailInvalid                          = _errors.MustSpec(210135, "USER_EMAIL_INVALID", "email không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PhoneAlreadyUsed                      = _errors.MustSpec(210136, "USER_PHONE_ALREADY_USED", "số điện thoại đã được sử dụng", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PayloadInvalid                        = _errors.MustSpec(210137, "USER_PAYLOAD_INVALID", "payload không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	ProfileIDInvalid                      = _errors.MustSpec(210138, "USER_PROFILE_ID_INVALID", "profileId không hợp lệ", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	RealEstateRoleInvalid                 = _errors.MustSpec(210140, "USER_REAL_ESTATE_ROLE_INVALID", "roleRealEstate không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	OnlineStatusHidden                    = _errors.MustSpec(210141, "USER_ONLINE_STATUS_HIDDEN", "Người dùng đã ẩn trạng thái online", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PermissionDenied                      = _errors.MustSpec(210142, "USER_PERMISSION_DENIED", "Bạn không có quyền truy cập", codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403))
	GroupKeyRequired                      = _errors.MustSpec(210143, "USER_GROUP_KEY_REQUIRED", "groupKey là bắt buộc", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	RequestInvalid                        = _errors.MustSpec(210144, "USER_REQUEST_INVALID", "request không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	AdminIDRequired                       = _errors.MustSpec(210145, "USER_ADMIN_ID_REQUIRED", "admin id là bắt buộc", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	AuthIDRequiredLower                   = _errors.MustSpec(210148, "USER_AUTH_ID_REQUIRED_LOWER", "authId là bắt buộc", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	ProfileIDRequired                     = _errors.MustSpec(210149, "USER_PROFILE_ID_REQUIRED", "profileId là bắt buộc", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	TagPayloadRequired                    = _errors.MustSpec(210150, "USER_TAG_PAYLOAD_REQUIRED", "Thiếu thông tin tag", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	TagNotFound                           = _errors.MustSpec(210151, "USER_TAG_NOT_FOUND", "Tag không tồn tại", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	TagTypeInvalid                        = _errors.MustSpec(210152, "USER_TAG_TYPE_INVALID", "Tag type không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	ProfileIDInvalidCaps                  = _errors.MustSpec(210153, "USER_PROFILE_ID_INVALID_CAPS", "Profile ID không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	TagNameRequired                       = _errors.MustSpec(210154, "USER_TAG_NAME_REQUIRED", "Tên tag không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	AccountForProfileNotFound             = _errors.MustSpec(210158, "USER_ACCOUNT_FOR_PROFILE_NOT_FOUND", "Không tìm thấy tài khoản cho profile này", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	DeviceRequestInvalid                  = _errors.MustSpec(210162, "USER_DEVICE_REQUEST_INVALID", "Yêu cầu không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	DeviceIDRequired                      = _errors.MustSpec(210163, "USER_DEVICE_ID_REQUIRED", "deviceId là bắt buộc", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	LastSeenAtInvalid                     = _errors.MustSpec(210164, "USER_LAST_SEEN_AT_INVALID", "lastSeenAt không hợp lệ, định dạng phải là RFC3339", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	RoleGroupNameRequired                 = _errors.MustSpec(210168, "USER_ROLE_GROUP_NAME_REQUIRED", "Tên role group không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	RoleGroupNameTooLong                  = _errors.MustSpec(210169, "USER_ROLE_GROUP_NAME_TOO_LONG", "Tên role group không được quá 20 ký tự", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	RoleGroupKeyRequired                  = _errors.MustSpec(210170, "USER_ROLE_GROUP_KEY_REQUIRED", "Key role group không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	RoleGroupDescriptionTooLong           = _errors.MustSpec(210171, "USER_ROLE_GROUP_DESCRIPTION_TOO_LONG", "Mô tả role group không được quá 500 ký tự", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	RoleGroupIDRequired                   = _errors.MustSpec(210172, "USER_ROLE_GROUP_ID_REQUIRED", "ID role group không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	PermissionNotFound                    = _errors.MustSpec(210173, "USER_PERMISSION_NOT_FOUND", "Permission không tồn tại", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	ModuleRequired                        = _errors.MustSpec(210174, "USER_MODULE_REQUIRED", "Module không được để trống", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	OrganizationIDRequired                = _errors.MustSpec(210175, "USER_ORGANIZATION_ID_REQUIRED", "organizationId là bắt buộc", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	CurrentUserUnknown                    = _errors.MustSpec(210176, "USER_CURRENT_USER_UNKNOWN", "Không xác định được người dùng hiện tại", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	ZaloUserIDUnavailable                 = _errors.MustSpec(210177, "USER_ZALO_USER_ID_UNAVAILABLE", "Không lấy được id từ Zalo", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	UserNotFoundAndCreateFailed           = _errors.MustSpec(210178, "USER_USER_NOT_FOUND_AND_CREATE_FAILED", "Không tìm thấy thông tin người dùng và không thể tạo mới", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	UserNotFoundAndCreateFailedLower      = _errors.MustSpec(210179, "USER_USER_NOT_FOUND_AND_CREATE_FAILED_LOWER", "không tìm thấy thông tin người dùng và không thể tạo mới", codes.FailedPrecondition, _errors.LegacyHTTP200())
	RolesNotFound                         = _errors.MustSpec(210180, "USER_ROLES_NOT_FOUND", "no roles found", codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404))
	AccountLockTypeInvalid                = _errors.MustSpec(210181, "USER_ACCOUNT_LOCK_TYPE_INVALID", "Loại khóa tài khoản không hợp lệ: 20(tạm thời) hoặc 30 (vĩnh viễn)", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	OTPExpired                            = _errors.MustSpec(210182, "USER_OTP_EXPIRED", "Mã OTP đã hết hạn. Nhấn Gửi lại mã.", codes.FailedPrecondition, _errors.LegacyHTTP200())
	OTPActivated                          = _errors.MustSpec(210183, "USER_OTP_ACTIVATED", "OTP đã kích hoạt", codes.FailedPrecondition, _errors.LegacyHTTP200())
	AccountInactive                       = _errors.MustSpec(210184, "USER_ACCOUNT_INACTIVE", "Tài khoản đang bị vô hiệu hóa", codes.FailedPrecondition, _errors.LegacyHTTP200())
	PINLocked                             = _errors.MustSpec(210185, "USER_PIN_LOCKED", "Mã PIN đã bị khóa", codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(423))
	PINAttemptsExceeded                   = _errors.MustSpec(210186, "USER_PIN_ATTEMPTS_EXCEEDED", "Đã nhập sai mã PIN quá số lần cho phép", codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(423))
	PINIncorrect                          = _errors.MustSpec(210187, "USER_PIN_INCORRECT", "Mã PIN không đúng", codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401))
	OTPLocked                             = _errors.MustSpec(210188, "USER_OTP_LOCKED", "Tài khoản tạm khóa do OTP", codes.FailedPrecondition, _errors.LegacyHTTP200())
	OTPSendLimitExceeded                  = _errors.MustSpec(210189, "USER_OTP_SEND_LIMIT_EXCEEDED", "Đã vượt giới hạn gửi OTP", codes.ResourceExhausted, _errors.LegacyHTTP200())
	OTPEnterLimitExceeded                 = _errors.MustSpec(210190, "USER_OTP_ENTER_LIMIT_EXCEEDED", "Đã vượt giới hạn nhập OTP", codes.ResourceExhausted, _errors.LegacyHTTP200())
	OTPIncorrect                          = _errors.MustSpec(210191, "USER_OTP_INCORRECT", "Mã OTP không chính xác", codes.Unauthenticated, _errors.LegacyHTTP200())
	DeviceLimitReached                    = _errors.MustSpec(210192, "USER_DEVICE_LIMIT_REACHED", "Tài khoản đã đạt giới hạn thiết bị", codes.ResourceExhausted, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	RequestValidationFailed               = _errors.MustSpec(210193, "USER_REQUEST_VALIDATION_FAILED", "Dữ liệu không hợp lệ", codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400))
	OTPNextSendLimited                    = _errors.MustSpec(210194, "USER_OTP_NEXT_SEND_LIMITED", "Hãy thử lại sau", codes.ResourceExhausted, _errors.LegacyHTTP200())
	OTPRequestLimited                     = _errors.MustSpec(210195, "USER_OTP_REQUEST_LIMITED", "Bạn đã yêu cầu OTP quá nhiều", codes.ResourceExhausted, _errors.LegacyHTTP200())
	RoleGroupNotFound                     = _errors.MustSpec(210196, "USER_ROLE_GROUP_NOT_FOUND", "role group not found", codes.NotFound, _errors.LegacyProblemCode("iam.role_group.not_found"))
	PlanTierRankMustBePositive            = _errors.MustSpec(210197, "USER_PLAN_TIER_RANK_MUST_BE_POSITIVE", "tier rank must be positive", codes.InvalidArgument, _errors.LegacyProblemCode("catalog.plan_version.tier_rank_positive"))
	PlanVersionIDRequired                 = _errors.MustSpec(210198, "USER_PLAN_VERSION_ID_REQUIRED", "plan version id is required", codes.InvalidArgument, _errors.LegacyProblemCode("catalog.plan_version.id_required"))
	PlanActorIDRequired                   = _errors.MustSpec(210199, "USER_PLAN_ACTOR_ID_REQUIRED", "actor id is required", codes.InvalidArgument, _errors.LegacyProblemCode("catalog.plan_version.actor_id_required"))
	PlanProductDisplayNameRequired        = _errors.MustSpec(210200, "USER_PLAN_PRODUCT_DISPLAY_NAME_REQUIRED", "product display name is required", codes.InvalidArgument, _errors.LegacyProblemCode("catalog.plan_version.product_display_name_required"))
	PlanPageSizeOutOfRange                = _errors.MustSpec(210201, "USER_PLAN_PAGE_SIZE_OUT_OF_RANGE", "page size must be between 1 and 100", codes.InvalidArgument, _errors.LegacyProblemCode("catalog.plan_version.page_size_invalid"))
	PlanStatusInvalid                     = _errors.MustSpec(210202, "USER_PLAN_STATUS_INVALID", "invalid plan status", codes.InvalidArgument, _errors.LegacyProblemCode("catalog.plan_version.status_invalid"))
	PlanTermsInvalid                      = _errors.MustSpec(210203, "USER_PLAN_TERMS_INVALID", "plan terms are invalid", codes.InvalidArgument, _errors.LegacyProblemCode("catalog.plan_version.terms_invalid"))
	PlanSubscriptionTermRequired          = _errors.MustSpec(210204, "USER_PLAN_SUBSCRIPTION_TERM_REQUIRED", "subscription term is required", codes.InvalidArgument, _errors.LegacyProblemCode("catalog.plan_version.subscription_term_required"))
	SelfBlockNotAllowed                   = _errors.MustSpec(
		210205, "USER_SELF_BLOCK_NOT_ALLOWED", "Không thể chặn chính mình",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	SelfUnblockNotAllowed = _errors.MustSpec(
		210206, "USER_SELF_UNBLOCK_NOT_ALLOWED", "Không thể bỏ chặn chính mình",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	UnblockRequired = _errors.MustSpec(
		210207, "USER_UNBLOCK_REQUIRED", "Vui lòng bỏ chặn",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(505),
	)
	SelfFriendNotAllowed = _errors.MustSpec(
		210208, "USER_SELF_FRIEND_NOT_ALLOWED", "Bạn bè là chính mình",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	FriendRequestAlreadySent = _errors.MustSpec(
		210209, "USER_FRIEND_REQUEST_ALREADY_SENT", "Đã gửi yêu cầu",
		codes.AlreadyExists, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	FriendTargetUnavailable = _errors.MustSpec(
		210210, "USER_FRIEND_TARGET_UNAVAILABLE", "Người dùng không tồn tại hoặc chưa đăng ký",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	FriendRequestCannotBePerformed = _errors.MustSpec(
		210211, "USER_FRIEND_REQUEST_CANNOT_BE_PERFORMED", "Không thể thực hiện yêu cầu",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	RecoveryIdentifierRequired = _errors.MustSpec(
		210212, "USER_RECOVERY_IDENTIFIER_REQUIRED", "Vui lòng cung cấp phone, email hoặc username để khôi phục tài khoản",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	DeletedAccountNotFound = _errors.MustSpec(
		210213, "USER_DELETED_ACCOUNT_NOT_FOUND", "Không tìm thấy tài khoản đã bị xóa",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	AccountAlreadyRestored = _errors.MustSpec(
		210214, "USER_ACCOUNT_ALREADY_RESTORED", "Tài khoản đã được khôi phục trước đó",
		codes.FailedPrecondition, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	OTPIPAddressRateLimited = _errors.MustSpec(
		210215, "USER_OTP_IP_ADDRESS_RATE_LIMITED", "IP của bạn đã bị chặn do gửi quá nhiều yêu cầu. Vui lòng thử lại sau.",
		codes.ResourceExhausted, _errors.LegacyHTTP200(), _errors.LegacyCode(401),
	)
	QRSessionCapacityExceeded = _errors.MustSpec(
		210216, "USER_QR_SESSION_CAPACITY_EXCEEDED", "Truy cập đạt giới hạn vui lòng thử lại sau",
		codes.ResourceExhausted, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	OrganizationMembershipAuthenticationFailed = _errors.MustSpec(
		210217, "USER_ORGANIZATION_MEMBERSHIP_AUTHENTICATION_FAILED", "Không thể xác thực thành viên tổ chức",
		codes.Unauthenticated, _errors.LegacyHTTP200(), _errors.LegacyCode(401),
	)
	OrganizationMembershipInactive = _errors.MustSpec(
		210218, "USER_ORGANIZATION_MEMBERSHIP_INACTIVE", "Bạn không phải thành viên đang hoạt động của tổ chức",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403),
	)
	ProfileAccessDenied = _errors.MustSpec(
		210219, "USER_PROFILE_ACCESS_DENIED", "Không có quyền truy cập profile này",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(403),
	)
	ProfileNotFound = _errors.MustSpec(
		210220, "USER_PROFILE_NOT_FOUND", "Không tìm thấy profile",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	AuthenticationInfoNotFound = _errors.MustSpec(
		210221, "USER_AUTHENTICATION_INFO_NOT_FOUND", "Không tìm thấy thông tin xác thực",
		codes.NotFound, _errors.LegacyHTTP200(), _errors.LegacyCode(404),
	)
	SelfFollowNotAllowed = _errors.MustSpec(
		210222, "USER_SELF_FOLLOW_NOT_ALLOWED", "Không thể tự follow chính mình",
		codes.InvalidArgument, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	AlreadyFollowing = _errors.MustSpec(
		210223, "USER_ALREADY_FOLLOWING", "Đã follow người này rồi",
		codes.AlreadyExists, _errors.LegacyHTTP200(), _errors.LegacyCode(400),
	)
	FriendGroupDeleteDenied = _errors.MustSpec(
		210224, "USER_FRIEND_GROUP_DELETE_DENIED", "bạn không có quyền xóa nhóm này",
		codes.PermissionDenied, _errors.LegacyHTTP200(), _errors.LegacyCode(401),
	)
)
