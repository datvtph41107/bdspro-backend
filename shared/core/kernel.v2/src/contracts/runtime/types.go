package contract_runtime

type ServiceType int

const (
	TypeLogger ServiceType = iota
	TypeTracer
	TypeMetrics
	TypeConfig
	TypeDatabase
	TypeCache
	TypeQueue
	TypeHTTPServer
	TypeGRPCServer
	TypeScheduler
	TypeWorker
	TypeCustom
)

func (t ServiceType) String() string {
	switch t {
	case TypeLogger:
		return "logger"
	case TypeTracer:
		return "tracer"
	case TypeMetrics:
		return "metrics"
	case TypeDatabase:
		return "database"
	case TypeHTTPServer:
		return "http-server"
	case TypeGRPCServer:
		return "grpc-server"
	default:
		return "unknown"
	}
}
