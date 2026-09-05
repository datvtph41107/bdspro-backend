package connect

// DefaultPorts map tên service → gRPC port (khớp gateway-service/config).
var DefaultPorts = map[string]string{
	"user":         "8201",
	"bdspro":       "8202",
	"file":         "8203",
	"notification": "8204",
	"payment":      "8205",
	"crm":          "8206",
	"organization": "8207",
	"chat":         "8208",
	"social":       "8209",
	"map":          "8210",
	"search":       "8211",
	"appointment":  "8213",
	"marketing":    "8214",
	"transaction":  "8215",
	"auth":         "8216",
	"feedback":     "8217",
	"assistant":    "8218",
	"tqd":          "8219",
	"hub":          "8280",
}

// DefaultPort trả port mặc định; rỗng nếu không biết.
func DefaultPort(name string) string {
	return DefaultPorts[name]
}
