// cmd/main.go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/franciscozamorau/osmi-gateway/internal/client"
	"github.com/franciscozamorau/osmi-gateway/internal/config"
	"github.com/franciscozamorau/osmi-gateway/internal/routes"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/franciscozamorau/osmi-protobuf/gen/pb"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 1. Crear cliente gRPC para conexión (opcional, para clientes específicos)
	grpcClient, err := client.NewGRPCClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer grpcClient.Close()

	// 2. Crear mux para el gateway
	mux := runtime.NewServeMux()

	// 3. Registrar TODOS los métodos gRPC automáticamente (¡ESTA ES LA LÍNEA CLAVE!)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err = pb.RegisterOsmiServiceHandlerFromEndpoint(ctx, mux, cfg.GRPCServerAddr, opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC gateway: %v", err)
	}

	// 4. Configurar router con middleware y rutas personalizadas
	handler := routes.SetupRouter(mux, grpcClient)

	log.Printf("Gateway iniciado en puerto %s", cfg.HTTPPort)

	if err := http.ListenAndServe(":"+cfg.HTTPPort, handler); err != nil {
		log.Fatalf("Error en servidor: %v", err)
	}
}
