package tracer

type Carrier interface {
	Get(key string) string
	Set(key string, value string)
	Keys() []string
}
