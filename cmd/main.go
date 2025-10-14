package main

import (
	"context"
	"log"
	"net/http"

	osmi "osmi-gateway/gen"
	"osmi-gateway/internal/utils"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(validateRequest),
	)

	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := osmi.RegisterOsmiServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to register handler: %v", err)
	}

	log.Println("Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func validateRequest(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	if req, ok := resp.(*osmi.TicketRequest); ok {
		if !utils.IsValidEventID(req.EventId) || !utils.IsValidUserID(req.UserId) {
			http.Error(w, "Invalid event_id or user_id", http.StatusBadRequest)
			return nil
		}
	}
	return nil
}
