package config

import "os"

func GetGRPCServerAddress() string {
	addr := os.Getenv("OSMI_GRPC_ADDR")
	if addr == "" {
		return "localhost:50051"
	}
	return addr
}
