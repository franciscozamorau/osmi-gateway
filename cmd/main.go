package main

import (
	"context"
	"log"
	"net/http"

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

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	log.Printf("Conectando a gRPC server: %s", cfg.GRPCServerAddr)

	// Registrar TODOS los métodos gRPC automáticamente
	err := pb.RegisterOsmiServiceHandlerFromEndpoint(ctx, mux, cfg.GRPCServerAddr, opts)
	if err != nil {
		log.Fatalf("Error registrando gateway: %v", err)
	}

	// Configurar router con middleware y rutas personalizadas
	handler := routes.SetupRouter(mux)

	log.Printf("🚀 Gateway iniciado en puerto %s", cfg.HTTPPort)
	log.Println("📡 Endpoints REST disponibles:")

	if err := http.ListenAndServe(":"+cfg.HTTPPort, handler); err != nil {
		log.Fatalf("Error en servidor: %v", err)
	}
}
