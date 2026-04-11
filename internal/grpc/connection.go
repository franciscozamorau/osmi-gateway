// internal/grpc/connection.go
package grpc

import (
	"log"
	"time"

	"github.com/franciscozamorau/osmi-gateway/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type ClientConnection struct {
	conn *grpc.ClientConn
	cfg  *config.Config
}

func NewClientConnection(cfg *config.Config) (*ClientConnection, error) {
	// Configuración keepalive para conexiones estables
	keepaliveParams := keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             5 * time.Second,
		PermitWithoutStream: true,
	}

	// CORREGIDO: El orden correcto es (target, opts...)
	// Antes tenía grpc.NewClient(cfg.GRPCServerAddr, opts...) pero faltaban las opciones
	conn, err := grpc.NewClient(
		cfg.GRPCServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepaliveParams),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024*10), // 10MB
			grpc.MaxCallSendMsgSize(1024*1024*10), // 10MB también para envío
		),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: 5 * time.Second,
		}),
	)
	if err != nil {
		return nil, err
	}

	log.Printf("✅ Conectado a gRPC server en %s", cfg.GRPCServerAddr)

	return &ClientConnection{
		conn: conn,
		cfg:  cfg,
	}, nil
}

func (c *ClientConnection) GetConnection() *grpc.ClientConn {
	return c.conn
}

func (c *ClientConnection) Close() error {
	if c.conn != nil {
		log.Printf("Cerrando conexión gRPC con %s", c.cfg.GRPCServerAddr)
		return c.conn.Close()
	}
	return nil
}

// NOTA: Los clientes específicos (customer_client, event_client) ya NO EXISTEN
// El gateway NO necesita clientes específicos porque usa el gateway automático
