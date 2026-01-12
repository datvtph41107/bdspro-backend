package _errors

type ErrorCode struct {
	Code    int32
	Message string
}

var defaultHTTPMessages = map[int32]string{
	400: "Bad Request",
	401: "Unauthorized",
	403: "Forbidden",
	404: "Not Found",
	409: "Conflict",
	429: "Too Many Requests",
	500: "Internal Server Error",
}

func throwErrorMessage(httpCode int32, message ...string) error {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = defaultHTTPMessages[httpCode]
	}
	return ReturnError(httpCode, msg)
}

func BadRequestException(message ...string) error {
	return throwErrorMessage(400, message...)
}

func UnauthorizedException(message ...string) error {
	return throwErrorMessage(401, message...)
}

func ForbiddenException(message ...string) error {
	return throwErrorMessage(403, message...)
}

func NotFoundException(message ...string) error {
	return throwErrorMessage(404, message...)
}

func ConflictException(message ...string) error {
	return throwErrorMessage(409, message...)
}

func TooManyRequestsException(message ...string) error {
	return throwErrorMessage(429, message...)
}

func InternalServerException(message ...string) error {
	return throwErrorMessage(500, message...)
}
