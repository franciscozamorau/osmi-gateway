package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/franciscozamorau/osmi-protobuf/gen/pb"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	grpcServerAddr := os.Getenv("GRPC_SERVER_ADDR")
	if grpcServerAddr == "" {
		grpcServerAddr = "localhost:50051"
	}

	log.Printf("Connecting to gRPC at: %s", grpcServerAddr)

	err := pb.RegisterOsmiServiceHandlerFromEndpoint(ctx, mux, grpcServerAddr, opts)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	mux.HandlePath("GET", "/health", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"osmi-gateway"}`))
	})

	// Usar puerto 8083 para evitar conflicto
	port := "8083"

	log.Printf("Starting gRPC Gateway on port %s", port)
	log.Println("REST Endpoints available:")
	log.Println("  GET  /health")
	log.Println("  POST /tickets")
	log.Println("  GET  /users/{id}/tickets")
	log.Println("  POST /customers")
	log.Println("  GET  /customers/{id}")
	log.Println("  POST /users")
	log.Println("  POST /events")
	log.Println("  GET  /events/{id}")
	log.Println("  GET  /events")
	log.Println("  POST /categories")
	log.Println("  GET  /events/{id}/categories")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
