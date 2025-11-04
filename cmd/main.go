package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	pb "github.com/franciscozamorau/osmi-gateway/gen"
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
		runtime.WithRoutingErrorHandler(handleRoutingError),
		runtime.WithErrorHandler(handleGatewayError),
	)

	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := pb.RegisterOsmiServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC-Gateway handler: %v", err)
	}

	handler := withLogging(withValidationMiddleware(mux))

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Gateway running on :%s", port)
	log.Printf("Available endpoints:")
	log.Printf("  POST   /customers")
	log.Printf("  GET    /customers/{id}")
	log.Printf("  POST   /events")
	log.Printf("  GET    /events/{public_id}")
	log.Printf("  GET    /events")
	log.Printf("  POST   /tickets")
	log.Printf("  GET    /users/{user_id}/tickets")
	log.Printf("  POST   /users")

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	log.Fatal(server.ListenAndServe())
}

func validateRequest(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	switch req := resp.(type) {
	case *pb.TicketRequest:
		if !utils.IsValidEventID(req.EventId) {
			return writeErrorResponse(w, "Invalid event_id format", http.StatusBadRequest)
		}
		if !utils.IsValidUserID(req.UserId) {
			return writeErrorResponse(w, "Invalid user_id format", http.StatusBadRequest)
		}
	case *pb.EventLookup:
		if !utils.IsValidEventID(req.PublicId) {
			return writeErrorResponse(w, "Invalid event ID format", http.StatusBadRequest)
		}
	case *pb.UserLookup:
		if !utils.IsValidUserID(req.UserId) {
			return writeErrorResponse(w, "Invalid user ID format", http.StatusBadRequest)
		}
	}
	return nil
}

func handleRoutingError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, httpStatus int) {
	if httpStatus == http.StatusNotFound {
		if hasInvalidPatternInPath(r.URL.Path) {
			writeErrorResponse(w, "Invalid ID format in URL parameter", http.StatusBadRequest)
			return
		}
	}
	runtime.DefaultRoutingErrorHandler(ctx, mux, marshaler, w, r, httpStatus)
}

func handleGatewayError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Gateway error: %v", err)
	writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
}

func withValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := validateURLParameters(r); err != nil {
			writeErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func validateURLParameters(r *http.Request) error {
	path := r.URL.Path

	if strings.Contains(path, "/events/") {
		parts := strings.Split(path, "/events/")
		if len(parts) > 1 {
			eventID := strings.TrimSuffix(parts[1], "/")
			if eventID != "" && !utils.IsValidEventID(eventID) {
				return fmt.Errorf("invalid event ID format")
			}
		}
	}

	if strings.Contains(path, "/users/") && strings.Contains(path, "/tickets") {
		parts := strings.Split(path, "/users/")
		if len(parts) > 1 {
			userParts := strings.Split(parts[1], "/tickets")
			if len(userParts) > 0 {
				userID := userParts[0]
				if userID != "" && !utils.IsValidUserID(userID) {
					return fmt.Errorf("invalid user ID format")
				}
			}
		}
	}

	if strings.Contains(path, "/customers/") {
		parts := strings.Split(path, "/customers/")
		if len(parts) > 1 {
			customerID := strings.TrimSuffix(parts[1], "/")
			if customerID != "" {
				if _, err := fmt.Sscanf(customerID, "%d", new(int)); err != nil {
					return fmt.Errorf("invalid customer ID format: must be a numeric value")
				}
			}
		}
	}

	return nil
}

func hasInvalidPatternInPath(path string) bool {
	invalidPatterns := []string{"/event-", "/user-", "/customer-"}
	for _, pattern := range invalidPatterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}

func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := map[string]interface{}{
		"code":    statusCode,
		"message": message,
		"details": []string{},
	}

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Printf("Error encoding error response: %v", err)
		return err
	}

	return nil
}
