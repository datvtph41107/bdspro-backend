package config

import "fmt"

type Env string

const (
	Local Env = "local"
	Dev   Env = "dev"
	Prod  Env = "production"
)

func ParseEnv(val string) (Env, error) {
	switch Env(val) {
	case Local, Dev, Prod:
		return Env(val), nil
	default:
		return "", fmt.Errorf("invalid ENV_RUNTIME: %s", val)
	}
}
