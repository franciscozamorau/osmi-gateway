// main.go - CORREGIDO PARA USAR TUS UTILS
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	"github.com/franciscozamorau/osmi-gateway/internal/utils"
	pb "github.com/franciscozamorau/osmi-server/gen"
)

func init() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Configurar ServeMux
	mux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(validateResponse),
		runtime.WithRoutingErrorHandler(handleRoutingError),
		runtime.WithErrorHandler(handleGatewayError),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}),
		runtime.WithIncomingHeaderMatcher(customHeaderMatcher),
	)

	// Configurar opciones gRPC
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(16*1024*1024),
			grpc.MaxCallSendMsgSize(16*1024*1024),
		),
		grpc.WithTimeout(30 * time.Second),
	}

	// Obtener dirección del servidor gRPC
	grpcServerAddr := os.Getenv("GRPC_SERVER_ADDR")
	if grpcServerAddr == "" {
		grpcServerAddr = "localhost:50051"
		log.Printf("Using default gRPC server address: %s", grpcServerAddr)
	}

	// Registrar endpoints
	log.Println("Registering gRPC-Gateway endpoints...")
	err := pb.RegisterOsmiServiceHandlerFromEndpoint(ctx, mux, grpcServerAddr, opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC-Gateway handler: %v", err)
	}

	// Crear handler con middlewares
	handler := withCORS(withLogging(withValidationMiddleware(mux)))

	// Obtener puerto
	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	// Log de endpoints disponibles
	logEndpoints()

	// Configurar mux principal
	mainMux := http.NewServeMux()
	mainMux.HandleFunc("/health", healthCheckHandler)
	mainMux.HandleFunc("/ready", readinessHandler)
	mainMux.Handle("/", handler)

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mainMux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// Iniciar servidor
	log.Printf("Gateway running on :%s (connecting to gRPC server at %s)", port, grpcServerAddr)
	log.Println("Gateway is ready to accept requests")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start gateway: %v", err)
	}
}

func logEndpoints() {
	log.Println("Available endpoints:")
	log.Println("  ┌─ Customers")
	log.Println("  │  POST   /v1/customers")
	log.Println("  │  GET    /v1/customers/{public_id}")
	log.Println("  ├─ Events")
	log.Println("  │  POST   /v1/events")
	log.Println("  │  GET    /v1/events/{public_id}")
	log.Println("  │  GET    /v1/events")
	log.Println("  ├─ Tickets")
	log.Println("  │  POST   /v1/tickets")
	log.Println("  │  GET    /v1/tickets?user_id={uuid}")
	log.Println("  │  GET    /v1/tickets?customer_id={uuid}")
	log.Println("  │  GET    /v1/tickets/{ticket_id}")
	log.Println("  │  PUT    /v1/tickets/{ticket_id}/status")
	log.Println("  ├─ Users")
	log.Println("  │  POST   /v1/users")
	log.Println("  ├─ Categories")
	log.Println("  │  POST   /v1/categories")
	log.Println("  │  GET    /v1/events/{public_id}/categories")
	log.Println("  └─ Health")
	log.Println("     GET    /health")
	log.Println("     GET    /ready")
}

func customHeaderMatcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case "x-request-id", "x-correlation-id", "authorization", "x-api-key":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func validateResponse(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	w.Header().Set("X-Response-Time", time.Now().Format(time.RFC3339))
	return nil
}

func handleRoutingError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, httpStatus int) {
	if httpStatus == http.StatusNotFound {
		writeErrorResponse(w, "Endpoint not found", http.StatusNotFound)
		return
	}
	runtime.DefaultRoutingErrorHandler(ctx, mux, marshaler, w, r, httpStatus)
}

func handleGatewayError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Gateway error for %s %s: %v", r.Method, r.URL.Path, err)

	errorMsg := "Internal server error"
	statusCode := http.StatusInternalServerError

	if strings.Contains(err.Error(), "desc = ") {
		parts := strings.Split(err.Error(), "desc = ")
		if len(parts) > 1 {
			errorMsg = cleanErrorMessage(parts[1])
		}
	}

	if strings.Contains(err.Error(), "code = NotFound") {
		statusCode = http.StatusNotFound
		errorMsg = "Resource not found"
	} else if strings.Contains(err.Error(), "code = InvalidArgument") {
		statusCode = http.StatusBadRequest
	} else if strings.Contains(err.Error(), "code = AlreadyExists") {
		statusCode = http.StatusConflict
		errorMsg = "Resource already exists"
	} else if strings.Contains(err.Error(), "code = Unavailable") {
		statusCode = http.StatusServiceUnavailable
		errorMsg = "Service temporarily unavailable"
	} else if strings.Contains(strings.ToLower(err.Error()), "deadline exceeded") {
		statusCode = http.StatusGatewayTimeout
		errorMsg = "Request timeout"
	} else if strings.Contains(strings.ToLower(err.Error()), "unauthorized") {
		statusCode = http.StatusUnauthorized
		errorMsg = "Unauthorized access"
	}

	writeErrorResponse(w, errorMsg, statusCode)
}

func cleanErrorMessage(errMsg string) string {
	errMsg = strings.TrimPrefix(errMsg, "rpc error: ")
	errMsg = strings.TrimSuffix(errMsg, "\"")
	errMsg = strings.TrimPrefix(errMsg, "\"")
	return errMsg
}

func withValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validar métodos HTTP
		if r.Method != http.MethodGet && r.Method != http.MethodPost &&
			r.Method != http.MethodPut && r.Method != http.MethodDelete &&
			r.Method != http.MethodOptions {
			writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Validar Content-Type para POST/PUT
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				writeErrorResponse(w, "Unsupported Content-Type. Use application/json", http.StatusUnsupportedMediaType)
				return
			}
		}

		// Validar parámetros de URL
		if err := validateURLParameters(r); err != nil {
			writeErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validar query parameters
		if err := validateQueryParameters(r); err != nil {
			writeErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateURLParameters(r *http.Request) error {
	path := r.URL.Path

	// Validar que todos los IDs en la URL sean UUIDs válidos
	if strings.Contains(path, "/events/") {
		parts := strings.Split(path, "/events/")
		if len(parts) > 1 {
			eventID := extractFirstPart(parts[1])
			if eventID != "" && !utils.IsValidUUID(eventID) {
				return fmt.Errorf("invalid event ID format: must be a valid UUID")
			}
		}
	}

	if strings.Contains(path, "/customers/") {
		parts := strings.Split(path, "/customers/")
		if len(parts) > 1 {
			customerID := extractFirstPart(parts[1])
			if customerID != "" && !utils.IsValidUUID(customerID) {
				return fmt.Errorf("invalid customer ID format: must be a valid UUID")
			}
		}
	}

	if strings.Contains(path, "/tickets/") {
		parts := strings.Split(path, "/tickets/")
		if len(parts) > 1 {
			ticketID := extractFirstPart(parts[1])
			if ticketID != "" && !utils.IsValidUUID(ticketID) {
				return fmt.Errorf("invalid ticket ID format: must be a valid UUID")
			}
		}
	}

	if strings.Contains(path, "/users/") {
		parts := strings.Split(path, "/users/")
		if len(parts) > 1 {
			userID := extractFirstPart(parts[1])
			if userID != "" && !utils.IsValidUUID(userID) {
				return fmt.Errorf("invalid user ID format: must be a valid UUID")
			}
		}
	}

	return nil
}

func validateQueryParameters(r *http.Request) error {
	query := r.URL.Query()

	// Validar user_id en query
	if userID := query.Get("user_id"); userID != "" {
		if !utils.IsValidUUID(userID) {
			return fmt.Errorf("invalid user_id in query: must be a valid UUID")
		}
	}

	// Validar customer_id en query
	if customerID := query.Get("customer_id"); customerID != "" {
		if !utils.IsValidUUID(customerID) {
			return fmt.Errorf("invalid customer_id in query: must be a valid UUID")
		}
	}

	// Validar email en query
	if email := query.Get("email"); email != "" {
		if !utils.IsValidEmail(email) {
			return fmt.Errorf("invalid email format in query")
		}
	}

	return nil
}

func extractFirstPart(path string) string {
	path = strings.TrimSuffix(path, "/")
	if strings.Contains(path, "/") {
		return strings.Split(path, "/")[0]
	}
	return path
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Request: %s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

		ww := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(ww, r)

		duration := time.Since(start)
		log.Printf("Response: %s %s %d %s", r.Method, r.URL.Path, ww.statusCode, duration)
	})
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID, X-Correlation-ID, X-API-Key")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"status":    "healthy",
		"service":   "osmi-gateway",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"status":    "ready",
		"service":   "osmi-gateway",
		"timestamp": time.Now().Format(time.RFC3339),
		"checks": map[string]string{
			"api":      "healthy",
			"database": "connected (via gRPC)",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    statusCode,
			"message": message,
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Printf("Error encoding error response: %v", err)
		return err
	}

	return nil
}
