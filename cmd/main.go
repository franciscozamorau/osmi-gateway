package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	pb "osmi-gateway/internal/pb"
	"osmi-gateway/internal/utils"
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

	// Este es el handler generado por protoc-gen-grpc-gateway
	err := pb.RegisterOsmiServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC-Gateway handler: %v", err)
	}

	log.Println("Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
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
