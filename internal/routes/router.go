// internal/routes/router.go
package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/franciscozamorau/osmi-gateway/internal/client"
	"github.com/franciscozamorau/osmi-gateway/internal/middleware"
	pb "github.com/franciscozamorau/osmi-protobuf/gen/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func SetupRouter(gwmux *runtime.ServeMux, grpcClient *client.GRPCClient) http.Handler {
	// Cliente para customers
	customerClient := client.NewCustomerClient(grpcClient.Conn)

	// ============ RUTAS PERSONALIZADAS ============

	// GET /customers/{public_id} - Usando public_id en lugar de id numérico
	gwmux.HandlePath("GET", "/customers/{public_id}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		publicID := pathParams["public_id"]
		if publicID == "" {
			http.Error(w, "public_id is required", http.StatusBadRequest)
			return
		}

		// Construir lookup con public_id
		lookup := &pb.CustomerLookup{
			Lookup: &pb.CustomerLookup_PublicId{
				PublicId: publicID,
			},
		}

		// Llamar al cliente gRPC
		resp, err := customerClient.GetCustomer(r.Context(), lookup)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Devolver respuesta
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	// ============ RUTAS PÚBLICAS ============
	gwmux.HandlePath("GET", "/health", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"osmi-gateway"}`))
	})

	// ============ MIDDLEWARE ============
	handler := middleware.Logging(gwmux)
	handler = middleware.CORS(handler)

	rateLimiter := middleware.NewRateLimiter(10, 20, 3*time.Minute)
	handler = rateLimiter.Limit(handler)

	return handler
}
