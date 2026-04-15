package config

type TLSConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	CertFile   string `mapstructure:"cert_file"   validate:"required_if=Enabled true"`
	KeyFile    string `mapstructure:"key_file"    validate:"required_if=Enabled true" sensitive:"true"`
	CAFile     string `mapstructure:"ca_file"`    // cho mTLS
	MutualTLS  bool   `mapstructure:"mutual_tls"` // yêu cầu client cert
	ServerName string `mapstructure:"server_name"`
}
