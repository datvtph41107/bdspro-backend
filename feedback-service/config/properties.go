package config

type Properties struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Grpc     GrpcConfig     `mapstructure:"grpc"`
}

type ServerConfig struct {
	Name    string `mapstructure:"name"`
	TCPPort int    `mapstructure:"tcp_port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DSN      string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Pass string `mapstructure:"pass"`
}

type GrpcConfig struct {
	User         GrpcServiceConfig `mapstructure:"user"`
	Organization GrpcServiceConfig `mapstructure:"organization"`
	Social       GrpcServiceConfig `mapstructure:"social"`
}

type GrpcServiceConfig struct {
	Host string `mapstructure:"host"`
}

var AppProperties *Properties
