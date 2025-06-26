package env

type DBType string

var (
	USER_CONTEXT             = "user_id"
	ADMIN_ROLE_KEY           = "admin"
	ADMIN_ROLE_NAME          = "Admin"
	MAX_ORGANIZATION_MEMBERS = 10

	// Organization Log Types
	LOG_TYPE_ORGANIZATION_CREATE = "ORGANIZATION_CREATE"
	LOG_TYPE_ORGANIZATION_UPDATE = "ORGANIZATION_UPDATE"
	LOG_TYPE_MEMBER_CREATE       = "MEMBER_CREATE"
	LOG_TYPE_MEMBER_UPDATE       = "MEMBER_UPDATE"
	LOG_TYPE_MEMBER_DELETE       = "MEMBER_DELETE"
	LOG_TYPE_ROLE_CREATE         = "ROLE_CREATE"
	LOG_TYPE_ROLE_UPDATE         = "ROLE_UPDATE"
	LOG_TYPE_ROLE_DELETE         = "ROLE_DELETE"
	LOG_TYPE_PERMISSION_CREATE   = "PERMISSION_CREATE"
	LOG_TYPE_PERMISSION_UPDATE   = "PERMISSION_UPDATE"
	LOG_TYPE_PERMISSION_DELETE   = "PERMISSION_DELETE"
	SERVICE_NAME                 = "organization-service"

	DBTypePostgres DBType = "postgres"

	DefaultDBDNS    = "host=localhost user=postgres password=postgres dbname=organization port=5432 sslmode=disable"
	DefaultGRPCPort = "8207"
)

var (
	GRPC_PORT                    = "8207"
	LOG_WORKER_QUEUE_SIZE        = 1000
	PRIMARY_DB_TYPE       DBType = ""
	DB_DNS                       = ""

	CHAT_SERVICE_HOST = "localhost"
	CHAT_SERVICE_PORT = "8208"

	USER_SERVICE_HOST = "localhost"
	USER_SERVICE_PORT = "8201"
	ENV_RUNTIME       = ""
)
