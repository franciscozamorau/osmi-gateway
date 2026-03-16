// internal/server/server.go
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/franciscozamorau/osmi-gateway/internal/config"
	gatewayGrpc "github.com/franciscozamorau/osmi-gateway/internal/grpc" // ALIAS para tu paquete
	"github.com/franciscozamorau/osmi-gateway/internal/handlers/health"
	"github.com/franciscozamorau/osmi-gateway/internal/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc" // ← paquete oficial
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/franciscozamorau/osmi-protobuf/gen/pb"
)

type Server struct {
	config     *config.Config
	grpcConn   *gatewayGrpc.ClientConnection // ← USAMOS EL ALIAS
	httpServer *http.Server
}

func NewServer(cfg *config.Config) (*Server, error) {
	// 1. Crear conexión gRPC (TU paquete, no el oficial)
	grpcConn, err := gatewayGrpc.NewClientConnection(cfg)
	if err != nil {
		return nil, err
	}

	// 2. Crear mux del gateway
	mux := runtime.NewServeMux()

	// 3. Registrar endpoints automáticos (usa el paquete oficial grpc)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err = pb.RegisterOsmiServiceHandlerFromEndpoint(
		context.Background(),
		mux,
		cfg.GRPCServerAddr,
		opts,
	)
	if err != nil {
		grpcConn.Close()
		return nil, err
	}

	// 4. APLICAR MIDDLEWARE MANUALMENTE (SIN CHAIN)
	// Envolvemos cada middleware de forma explícita
	var handler http.Handler = mux

	// ORDEN CORRECTO (de afuera hacia adentro)
	handler = middleware.Recovery(handler)  // 1. Recovery (el más externo)
	handler = middleware.RequestID(handler) // 2. Request ID
	handler = middleware.Logging(handler)   // 3. Logging
	handler = middleware.CORS(handler)      // 4. CORS
	handler = middleware.RateLimit(handler) // 5. Rate Limit
	handler = middleware.Auth(handler)      // 6. Auth (el más interno)

	// 5. Registrar rutas manuales
	mainMux := http.NewServeMux()
	mainMux.Handle("/", handler)                        // Todas las rutas automáticas
	mainMux.HandleFunc("/health", health.HealthHandler) // Ruta manual

	// 6. Crear servidor HTTP
	httpServer := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      mainMux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		config:     cfg,
		grpcConn:   grpcConn,
		httpServer: httpServer,
	}, nil
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop() error {
	if s.grpcConn != nil {
		s.grpcConn.Close()
	}
	return s.httpServer.Close()
}
