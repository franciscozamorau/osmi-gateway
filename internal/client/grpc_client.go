package client

import (
	"fmt"

	"github.com/franciscozamorau/osmi-gateway/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCClient maneja las conexiones gRPC
type GRPCClient struct {
	Conn *grpc.ClientConn
	cfg  *config.Config
}

// NewGRPCClient crea una nueva conexión gRPC
func NewGRPCClient(cfg *config.Config) (*GRPCClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(cfg.GRPCServerAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return &GRPCClient{
		Conn: conn,
		cfg:  cfg,
	}, nil
}

// Close cierra la conexión
func (c *GRPCClient) Close() error {
	return c.Conn.Close()
}
