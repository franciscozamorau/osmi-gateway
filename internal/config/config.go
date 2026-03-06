package config

import (
	"os"
	"time"
)

type Config struct {
	GRPCServerAddr string
	HTTPPort       string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

func Load() *Config {
	grpcAddr := os.Getenv("GRPC_SERVER_ADDR")
	if grpcAddr == "" {
		grpcAddr = "localhost:50051"
	}

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8083"
	}

	return &Config{
		GRPCServerAddr: grpcAddr,
		HTTPPort:       port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
}
