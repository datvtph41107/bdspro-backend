package config

type Options struct {
	ConfigName  string   // tên file config
	ConfigPaths []string // các đường dẫn tìm
	ConfigType  string   // loại file (yaml, json)
	EnvVar      string   // biến môi trường xác định env
	EnvPrefix   string   // prefix cho env
	ReadEnv     bool     // có đọc env không
	Validate    bool     // có validate không
}

type Option func(*Options)

// tên của file config viper.SetConfigName()
// .WithConfigName("config") tìm file config.yaml, config.json, config.toml
func WithConfigName(name string) Option {
	return func(o *Options) { o.ConfigName = name }
}

func WithConfigPath(path string) Option {
	return func(o *Options) { o.ConfigPaths = append(o.ConfigPaths, path) }
}

func WithConfigType(typ string) Option {
	return func(o *Options) { o.ConfigType = typ }
}

func WithEnvVar(envVar string) Option {
	return func(o *Options) { o.EnvVar = envVar }
}

func WithEnvPrefix(prefix string) Option {
	return func(o *Options) {
		o.ReadEnv = true
		o.EnvPrefix = prefix
	}
}

func WithValidation(enabled bool) Option {
	return func(o *Options) { o.Validate = enabled }
}

func DefaultOptions() *Options {
	return &Options{
		ConfigName:  "config",
		ConfigPaths: []string{".", "./config"},
		ConfigType:  "yaml",
		EnvVar:      "ENV_RUNTIME",
		ReadEnv:     false,
		Validate:    true,
	}
}
