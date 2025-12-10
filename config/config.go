package config

import "os"

func GetGRPCServerAddress() string {
	addr := os.Getenv("OSMI_GRPC_ADDR")
	if addr == "" {
		return "localhost:50051"
	}
	return addr
}

/* expandirlo para más configuraciones futuro
package config

import (
	"os"
	"strconv"
)

type Config struct {
	GRPCServerAddress string
	GatewayPort       string
	DatabaseURL       string
	EnableTLS         bool
}

func Load() *Config {
	return &Config{
		GRPCServerAddress: getEnv("OSMI_GRPC_ADDR", "localhost:50051"),
		GatewayPort:       getEnv("GATEWAY_PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		EnableTLS:         getEnvBool("ENABLE_TLS", false),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}*/
