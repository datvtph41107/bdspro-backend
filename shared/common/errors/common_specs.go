package _errors

import "google.golang.org/grpc/codes"

var (
	AuthorizationHeaderMissing = MustSpec(
		100001,
		"AUTH_AUTHORIZATION_HEADER_MISSING",
		"Authorization header is missing",
		codes.Unauthenticated,
		LegacyProblemCode("auth.authorization_header_missing"),
	)
	TokenMissing = MustSpec(
		100002,
		"AUTH_TOKEN_MISSING",
		"Missing token",
		codes.Unauthenticated,
		LegacyProblemCode("auth.token_missing"),
	)
	TokenInvalidOrExpired = MustSpec(
		100003,
		"AUTH_TOKEN_INVALID_OR_EXPIRED",
		"Invalid or expired token",
		codes.Unauthenticated,
		LegacyProblemCode("auth.token_invalid_or_expired"),
	)
	TempTokenRouteForbidden = MustSpec(
		100004,
		"AUTH_TEMP_TOKEN_ROUTE_FORBIDDEN",
		"Token không đúng",
		codes.PermissionDenied,
		LegacyProblemCode("auth.temp_token_route_forbidden"),
	)
	APIKeyHeaderInvalid = MustSpec(
		100005,
		"AUTH_API_KEY_HEADER_INVALID",
		"Invalid API key header",
		codes.InvalidArgument,
		LegacyProblemCode("auth.api_key_header_invalid"),
	)
	APIKeyInvalid = MustSpec(
		100006,
		"AUTH_API_KEY_INVALID",
		"API key verification failed",
		codes.Unauthenticated,
		LegacyProblemCode("auth.api_key_invalid"),
	)
	APIKeyVerificationUnavailable = MustSpec(
		100007,
		"AUTH_API_KEY_VERIFICATION_UNAVAILABLE",
		"API key verification unavailable",
		codes.Unavailable,
		LegacyProblemCode("auth.api_key_verification_unavailable"),
	)
	AuthenticationRequired = MustSpec(
		100008,
		"AUTH_AUTHENTICATION_REQUIRED",
		"Authorization header missing",
		codes.Unauthenticated,
		LegacyProblemCode("auth.authentication_required"),
	)
	CallerClassificationInvalid = MustSpec(
		100009,
		"AUTH_CALLER_CLASSIFICATION_INVALID",
		"Invalid caller classification",
		codes.Unauthenticated,
		LegacyProblemCode("auth.caller_classification_invalid"),
	)
	RequestValidationFailed = MustSpec(
		100010,
		"COMMON_REQUEST_VALIDATION_FAILED",
		"Dữ liệu không hợp lệ",
		codes.InvalidArgument,
		LegacyHTTP200(),
		LegacyCode(400),
	)
	ResourceIDInvalid = MustSpec(
		100011,
		"COMMON_RESOURCE_ID_INVALID",
		"ID không hợp lệ",
		codes.InvalidArgument,
		LegacyHTTP200(),
		LegacyCode(400),
	)
	DataNotFound = MustSpec(
		100012,
		"COMMON_DATA_NOT_FOUND",
		"Không tìm thấy dữ liệu",
		codes.NotFound,
		LegacyHTTP200(),
		LegacyCode(404),
	)
)
