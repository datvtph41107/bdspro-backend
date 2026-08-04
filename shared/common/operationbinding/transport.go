package operationbinding

import (
	"fmt"
	"strings"
	"unicode"
)

/** TransportKind identifies one technical entrypoint type. */
type TransportKind string

const (
	TransportGRPCUnary        TransportKind = "grpc_unary"
	TransportGRPCServerStream TransportKind = "grpc_server_stream"
	TransportHTTPRoute        TransportKind = "http_route"
)

/** Transport identifies one RPC or route template. */
type Transport struct {
	Kind     TransportKind
	Method   string
	Template string
}

/** GRPCUnary creates one unary gRPC transport key. */
func GRPCUnary(fullMethod string) Transport {
	return Transport{Kind: TransportGRPCUnary, Method: fullMethod}
}

/** GRPCServerStream creates one streaming gRPC transport key. */
func GRPCServerStream(fullMethod string) Transport {
	return Transport{Kind: TransportGRPCServerStream, Method: fullMethod}
}

/** HTTPRoute creates one HTTP route-template key. */
func HTTPRoute(method, template string) Transport {
	return Transport{Kind: TransportHTTPRoute, Method: strings.ToUpper(method), Template: template}
}

/** CanonicalKey returns a stable transport lookup key. */
func (t Transport) CanonicalKey() string {
	switch t.Kind {
	case TransportGRPCUnary, TransportGRPCServerStream:
		return string(t.Kind) + ":" + t.Method
	case TransportHTTPRoute:
		return string(t.Kind) + ":" + strings.ToUpper(t.Method) + ":" + t.Template
	default:
		return string(t.Kind) + ":" + t.Method + ":" + t.Template
	}
}

/** Validate checks one transport key. */
func (t Transport) Validate() error {
	switch t.Kind {
	case TransportGRPCUnary, TransportGRPCServerStream:
		if t.Template != "" {
			return fmt.Errorf("gRPC transport cannot declare route template")
		}
		return validateFullMethod(t.Method)
	case TransportHTTPRoute:
		if !validHTTPMethod(t.Method) {
			return fmt.Errorf("unsupported HTTP method %q", t.Method)
		}
		return validateRouteTemplate(t.Template)
	default:
		return fmt.Errorf("unsupported transport kind %q", t.Kind)
	}
}

func validateFullMethod(value string) error {
	if value == "" || strings.TrimSpace(value) != value || !strings.HasPrefix(value, "/") {
		return fmt.Errorf("invalid gRPC full method %q", value)
	}
	parts := strings.Split(strings.TrimPrefix(value, "/"), "/")
	if len(parts) != 2 || !validDottedIdentifier(parts[0]) || !validIdentifier(parts[1]) {
		return fmt.Errorf("invalid gRPC full method %q", value)
	}
	return nil
}

func validateRouteTemplate(value string) error {
	if value == "" || strings.TrimSpace(value) != value || !strings.HasPrefix(value, "/") {
		return fmt.Errorf("invalid HTTP route template %q", value)
	}
	if strings.ContainsAny(value, "?#") {
		return fmt.Errorf("HTTP route template cannot contain query or fragment")
	}
	for _, r := range value {
		if unicode.IsControl(r) || unicode.IsSpace(r) {
			return fmt.Errorf("HTTP route template contains unsafe character")
		}
	}
	return nil
}

func validHTTPMethod(value string) bool {
	switch strings.ToUpper(value) {
	case "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS":
		return true
	default:
		return false
	}
}

func validDottedIdentifier(value string) bool {
	parts := strings.Split(value, ".")
	if len(parts) < 2 {
		return false
	}
	for _, part := range parts {
		if !validIdentifier(part) {
			return false
		}
	}
	return true
}

func validIdentifier(value string) bool {
	if value == "" {
		return false
	}
	for index, r := range value {
		if index == 0 {
			if !isASCIIAlpha(r) && r != '_' {
				return false
			}
			continue
		}
		if !isASCIIAlpha(r) && (r < '0' || r > '9') && r != '_' {
			return false
		}
	}
	return true
}

func isASCIIAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
