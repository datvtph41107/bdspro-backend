package env

type DBType string

var (
	USER_CONTEXT = "user_id"
	SERVICE_NAME = "marketing-service"

	DBTypePostgres DBType = "postgres"

	DefaultDBDNS    = "host=localhost user=postgres password=postgres dbname=marketing port=5432 sslmode=disable"
	DefaultGRPCPort = "8209"
)

var (
	GRPC_PORT              = "8214"
	PRIMARY_DB_TYPE DBType = ""
	DB_DNS                 = ""

	PAYMENT_SERVICE_HOST = "localhost"
	PAYMENT_SERVICE_PORT = "8205"

	NOTIFICATION_SERVICE_HOST = "localhost"
	NOTIFICATION_SERVICE_PORT = "8204"

	ENV_RUNTIME = ""
)
