package contract

type Config interface {
	Get(key string) any
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
}
