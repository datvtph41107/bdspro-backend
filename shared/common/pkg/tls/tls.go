// shared/pkg/tls/tls.go
package tlsutil

import (
	"crypto/tls"
	"fmt"

	sharedconfig "common/pkg/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// ServerCredentials tạo TLS credentials cho gRPC server từ config
func ServerCredentials(cfg sharedconfig.TLSConfig) (credentials.TransportCredentials, error) {
	if !cfg.Enabled {
		return insecure.NewCredentials(), nil
	}

	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("load cert/key: %w", err)
	}

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13, // không cho phép TLS 1.0, 1.1, 1.2
	}

	if cfg.MutualTLS && cfg.CAFile != "" {
		// caPool, err := loadCA(cfg.CAFile)
		if err != nil {
			return nil, fmt.Errorf("load CA: %w", err)
		}
		tlsCfg.ClientAuth = tls.RequireAndVerifyClientCert
		// tlsCfg.ClientCAs = caPool
	}

	return credentials.NewTLS(tlsCfg), nil
}

// ClientCredentials cho gateway khi dial upstream (mTLS)
func ClientCredentials(cfg sharedconfig.TLSConfig) (credentials.TransportCredentials, error) {
	if !cfg.Enabled {
		return insecure.NewCredentials(), nil
	}
	// Trong dev: skip verify để không cần cert cho localhost
	if cfg.ServerName == "localhost" || cfg.ServerName == "" {
		return credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true, // chỉ dev
		}), nil
	}

	tlsCfg := &tls.Config{
		ServerName: cfg.ServerName,
		MinVersion: tls.VersionTLS13,
	}

	// mTLS: gateway gửi client cert cho upstream
	if cfg.MutualTLS {
		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, err
		}
		tlsCfg.Certificates = []tls.Certificate{cert}
	}

	return credentials.NewTLS(tlsCfg), nil
}

// Reflection chỉ bật ở dev — cho phép grpcurl, Postman, Evans CLI
func RegisterReflection(s *grpc.Server, env string) {
	if env == "development" || env == "dev" {
		reflection.Register(s)
	}
	// Không bao giờ bật reflection ở prod — lộ toàn bộ API schema
}

// yaml# dev.yaml
// gateway:
//   tls:
//     enabled: false  # dev không cần TLS, kết nối local

// # prod.yaml
// gateway:
//   tls:
//     enabled: true
//     cert_file: "/etc/ssl/certs/gateway.crt"
//     key_file: "/etc/ssl/private/gateway.key"
//     mutual_tls: true
//     ca_file: "/etc/ssl/certs/internal-ca.crt"
//     server_name: "gateway.prod.internal"
