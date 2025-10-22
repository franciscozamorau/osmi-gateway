package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	pb "github.com/franciscozamorau/osmi-gateway/gen"
	"github.com/franciscozamorau/osmi-gateway/internal/handlers"
	"github.com/franciscozamorau/osmi-gateway/internal/utils"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment")
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(validateRequest),
	)

	opts := []grpc.DialOption{grpc.WithInsecure()}

	// Gateway generado por protoc-gen-grpc-gateway
	err := pb.RegisterOsmiServiceHandlerFromEndpoint(ctx, mux, "osmi-server:50051", opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC-Gateway handler: %v", err)
	}

	// Handlers manuales conectados al router
	http.HandleFunc("/customers", handlers.CreateCustomerHandler)
	http.HandleFunc("/tickets", handlers.CreateTicketHandler)
	http.HandleFunc("/events", handlers.CreateEventHandler)
	http.HandleFunc("/users", handlers.CreateUserHandler)

	// Gateway multiplexado en /api/
	http.Handle("/api/", http.StripPrefix("/api", mux))

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Gateway running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func validateRequest(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	if req, ok := resp.(*pb.TicketRequest); ok {
		if !utils.IsValidEventID(req.EventId) || !utils.IsValidUserID(req.UserId) {
			http.Error(w, "Invalid event_id or user_id", http.StatusBadRequest)
			return nil
		}
	}
	return nil
}
